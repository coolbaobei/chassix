package chassis

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/coolbaobei/chassix/config"
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
