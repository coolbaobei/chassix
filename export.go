package chassis
//
//import (
//	"gorm.io/gorm"
//)
//
//const (
//	DriverMysql    = "mysql"
//	DriverPostgres = "postgres"
//)
//
//var driver map[string]DatabaseProvider = make(map[string]DatabaseProvider)
//
//func RegisterDriver(dType string, provider DatabaseProvider) {
//	driver[dType] = provider
//}
//
////DB get the default(first) *Db connection
//func DB() *gorm.DB {
//	if dbs := DBs(); dbs == nil || len(dbs) == 0 {
//		return nil
//	} else {
//		return dbs[0]
//	}
//}
//
////DBs get all database connections
//func DBs() []*gorm.DB {
//	if initMultiDBSource(); 0 == multiDBSource.Size() {
//		return nil
//	}
//	return multiDBSource.DBs
//}
//func SetDB(index int, db *gorm.DB) {
//	if len(multiDBSource.DBs) == 0 {
//		initMultiDBSource()
//	}
//	multiDBSource.Lock.Lock()
//	defer multiDBSource.Lock.Unlock()
//	multiDBSource.DBs[index] = db
//}
//
////Close close all db connection
//func CloseAllDB() error {
//	if 0 == multiDBSource.Size() {
//		return ErrNoDatabaseConfiguration
//	}
//	for _, v := range multiDBSource.DBs {
//		if db, err := v.DB(); err == nil {
//			if err := db.Close(); nil != err {
//				return err
//			}
//		}
//	}
//	return nil
//}
////GetDatasource get datasource
//func GetDatasource() *Datasource {
//	return datasource
//}
