//go:build integration
// +build integration

package config_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_backend1/internal/config"
	"github.com/xfiendx4life/gb_go_backend1/internal/logger"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func createTestFile(data []byte, z *zap.SugaredLogger) {
	f, err := os.Create("test.yaml")
	if err != nil {
		z.Fatalf("%s", err)
		panic(err)
	}
	defer func() {
		err := f.Close()
		if err != nil {
			z.Fatalf("can't close file")
		}
	}()
	f.Write(data)
}

func TestReadFromFile(t *testing.T) {
	l := logger.InitLogger(zapcore.DebugLevel, "")
	data := []byte(`timeout: 2
loglevel: debug 
logfile: access.txt 
uri: postgres://xfiendx4life:123456@172.17.0.2:5432/shortener
maxcons: 10
mincons: 5
secretkey: somesecret
ttl: 60`)
	createTestFile(data, l)
	res, err := config.ReadFromFile("test.yaml", l)
	os.Remove("test.yaml")
	assert.Nil(t, err)
	assert.Equal(t, data, res)
}

func TestReadFromFileError(t *testing.T) {
	l := logger.InitLogger(zapcore.DebugLevel, "")
	_, err := config.ReadFromFile("test.yaml", l)
	os.Remove("test.yaml")
	assert.NotNil(t, err)
}
