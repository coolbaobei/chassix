package logx

import (
	"github.com/coolbaobei/chassix/config"
	"testing"
)

func Test_Logger(t *testing.T) {
	config.LoadFromEnvFile()

	New().Component("log").Category("test").Info("test log")
}
