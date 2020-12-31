package chassis

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"c6x.io/chassis/config"
)

func TestDBs(t *testing.T) {
	//defer CloseAllDB()
	// given
	config.LoadFromEnvFile()
	dbCfg := config.Databases()
	assert.NotEmpty(t, dbCfg)
	// when
	dbs, _ := DBs()
	assert.NotEmpty(t, dbs)
	assert.NotNil(t, dbs[0])
}
