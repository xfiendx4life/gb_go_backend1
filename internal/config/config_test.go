package config_test

import (
	"os"
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
targetfile: target.csv
port: :8080`
	lgr := newLogger()
	err := c.ReadConfig([]byte(data), lgr)
	assert.Nil(t, err)
	assert.Equal(t, time.Duration(2)*time.Second, c.GetTimeOut())
	assert.Equal(t, zapcore.FatalLevel, c.GetLogLevel())
	assert.Equal(t, "access.txt", c.GetLogFile())
	assert.Equal(t, ":8080", c.GetPort())
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
port: :8080
`
	lgr := newLogger()
	err := c.ReadConfig([]byte(data), lgr)
	assert.NoError(t, err)
	assert.Equal(t, "postgres://xfiendx4life:123456@172.17.0.2:5432/shortener", c.GetConfStorage().GetURI())
	assert.Equal(t, 10, c.GetConfStorage().GetMaxCons())
	assert.Equal(t, 5, c.GetConfStorage().GetMinCons())
	assert.Equal(t, "somesecret", c.GetConfAuth().GetSecretKey())
	assert.Equal(t, int64(60), c.GetConfAuth().GetTtl())
	assert.Equal(t, ":8080", c.GetPort())

}

func TestReadConfigError(t *testing.T) {
	c := config.New()
	data := `some text`
	lgr := newLogger()
	err := c.ReadConfig([]byte(data), lgr)
	assert.NotNil(t, err)
}

func TestGetFromEnv(t *testing.T) {
	os.Setenv("TIMEOUT", "2")
	os.Setenv("LOGLEVEL", "debug")
	os.Setenv("LOGFILE", "access.txt")
	os.Setenv("URI", "postgres://xfiendx4life:123456@172.17.0.2:5432/shortener")
	os.Setenv("MAXCONS", "10")
	os.Setenv("MINCONS", "5")
	os.Setenv("SECRETKEY", "somesecret")
	os.Setenv("TTL", "60")
	os.Setenv("PORT", ":8080")
	data := config.ReadFromEnv()
	testData := `timeout: 2
loglevel: debug
logfile: access.txt
uri: postgres://xfiendx4life:123456@172.17.0.2:5432/shortener
maxcons: 10
mincons: 5
secretkey: somesecret
port: :8080
ttl: 60`
	assert.Equal(t, testData, string(data))
}
