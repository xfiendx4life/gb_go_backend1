//go:build unit
// +build unit

package config_test

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_backend1/internal/config"
	"github.com/xfiendx4life/gb_go_backend1/internal/logger"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func newLogger() *zap.SugaredLogger {
	level := zapcore.DebugLevel
	return logger.InitLogger(level, "")
}

func TestReadConfig(t *testing.T) {
	c := config.New()
	data := `timeout: 2
loglevel: fatal
logfile: access.txt
targetfile: target.csv`
	lgr := newLogger()
	err := c.ReadConfig([]byte(data), lgr)
	assert.Nil(t, err)
	assert.Equal(t, time.Duration(2)*time.Second, c.GetTimeOut())
	assert.Equal(t, zapcore.FatalLevel, c.GetLogLevel())
	assert.Equal(t, "access.txt", c.GetLogFile())
}

func TestReadFullConfig(t *testing.T) {
	c := config.New()
	data := `timeout: 2
loglevel: debug 
logfile: access.txt 
uri: postgres://xfiendx4life:123456@172.17.0.2:5432/shortener
maxcons: 10
mincons: 5
secretkey: somesecret
ttl: 60
`
	lgr := newLogger()
	err := c.ReadConfig([]byte(data), lgr)
	assert.NoError(t, err)
	assert.Equal(t, "postgres://xfiendx4life:123456@172.17.0.2:5432/shortener", c.GetConfStorage().GetURI())
	assert.Equal(t, 10, c.GetConfStorage().GetMaxCons())
	assert.Equal(t, 5, c.GetConfStorage().GetMinCons())
	assert.Equal(t, "somesecret", c.GetConfAuth().GetSecretKey())
	assert.Equal(t, int64(60), c.GetConfAuth().GetTtl())

}

func TestReadConfigError(t *testing.T) {
	c := config.New()
	data := `some text`
	lgr := newLogger()
	err := c.ReadConfig([]byte(data), lgr)
	assert.NotNil(t, err)
}
