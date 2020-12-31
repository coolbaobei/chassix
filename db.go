package chassis

import (
	"c6x.io/chassis/config"
	"c6x.io/chassis/logx"
	"errors"
	"sync"
	"time"

	"gorm.io/driver/mysql"
	pg "gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type MultiDBSource struct {
	lock sync.RWMutex
	dbs  []*gorm.DB
}

const (
	DriverMysql    = "mysql"
	DriverPostgres = "postgres"
)

var (
	ErrNoDatabaseConfiguration = errors.New("there isn't any database setting in the configuration file")
)

var (
	multiDBSource *MultiDBSource
	initOnce      sync.Once
)

func initMultiDBSource() {
	initOnce.Do(func() {
		multiCfg := config.Databases()
		multiDBSource = new(MultiDBSource)
		multiDBSource.lock.Lock()
		defer multiDBSource.lock.Unlock()
		for _, v := range multiCfg {
			multiDBSource.dbs = append(multiDBSource.dbs, mustConnectDB(v))
		}
	})
}

func mustConnectDB(dbCfg *config.DatabaseConfig) *gorm.DB {
	log := logx.New().Service("chassis").Category("gorm")
	dialect := dbCfg.Dialect
	var db *gorm.DB
	var err error
	if "" == dialect || DriverMysql == dialect {
		db, err = gorm.Open(mysql.New(mysql.Config{
			DSN:                       dbCfg.DSN, // data source name, refer https://github.com/go-sql-driver/mysql#dsn-data-source-name
			DefaultStringSize:         256,       // add default size for string fields, by default, will use db type `longtext` for fields without size, not a primary key, no index defined and don't have default values
			DisableDatetimePrecision:  true,      // disable datetime precision support, which not supported before MySQL 5.6
			DontSupportRenameIndex:    true,      // drop & create index when rename index, rename index not supported before MySQL 5.7, MariaDB
			DontSupportRenameColumn:   true,      // use change when rename column, rename rename not supported before MySQL 8, MariaDB
			SkipInitializeWithVersion: false,     // smart configure based on used version
		}), &gorm.Config{Logger: DefaultLogger(&dbCfg.Logger)})
		if err == nil {
			if sqlDB, err := db.DB(); err == nil {
				if dbCfg.MaxIdle > 0 {
					sqlDB.SetMaxIdleConns(dbCfg.MaxIdle)
				}
				if dbCfg.MaxOpen > 0 && dbCfg.MaxOpen > dbCfg.MaxIdle {
					sqlDB.SetMaxOpenConns(100)
				}
				if dbCfg.MaxLifetime > 0 {
					sqlDB.SetConnMaxLifetime(time.Duration(dbCfg.MaxLifetime) * time.Second)
				}
				return db
			}
		} else {
			log.Errorf("connect mysql db failed: error=%s", err.Error())
			log.Fatalln(err)
			return nil
		}

	}
	if DriverPostgres == dialect {
		if db, err := gorm.Open(pg.New(pg.Config{DSN: dbCfg.DSN}), &gorm.Config{
			Logger: DefaultLogger(&dbCfg.Logger),
		}); err == nil {
			if sqlDB, err := db.DB(); err == nil {
				if dbCfg.MaxIdle > 0 {
					sqlDB.SetMaxIdleConns(dbCfg.MaxIdle)
				}
				if dbCfg.MaxOpen > 0 && dbCfg.MaxOpen > dbCfg.MaxIdle {
					sqlDB.SetMaxOpenConns(100)
				}
				if dbCfg.MaxLifetime > 0 {
					sqlDB.SetConnMaxLifetime(time.Duration(dbCfg.MaxLifetime) * time.Second)
				}
				return db
			} else {
				return nil
			}
		} else {
			log.Errorf("connect db failed: error=%s", err.Error())
			log.Fatalln(err)
		}
		return nil
	}

	return nil
}

//DB get the default(first) *Db connection
func DB() (*gorm.DB, error) {
	if dbs, err := DBs(); nil != err {
		return nil, err
	} else {
		return dbs[0], nil
	}
}

//DBs get all database connections
func DBs() ([]*gorm.DB, error) {
	if initMultiDBSource(); 0 == multiDBSource.Size() {
		return nil, ErrNoDatabaseConfiguration
	}
	return multiDBSource.dbs, nil
}

//Close close all db connection
func CloseAllDB() error {
	if 0 == multiDBSource.Size() {
		return ErrNoDatabaseConfiguration
	}
	for _, v := range multiDBSource.dbs {
		if db, err := v.DB(); err == nil {
			if err := db.Close(); nil != err {
				return err
			}
		}
	}
	return nil
}

//Size get db connection size
func (s MultiDBSource) Size() int {
	s.lock.RLock()
	defer s.lock.RUnlock()
	return len(s.dbs)
}
