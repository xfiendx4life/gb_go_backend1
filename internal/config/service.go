package config

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

func InitConfig() Config {
	return &ConfYML{}
}

func (c *ConfYML) GetLogLevel() zapcore.Level {
	return c.LogLevel
}

func (c *ConfYML) GetLogFile() string {
	return c.LogFile
}

func (c *ConfYML) GetTimeOut() time.Duration {
	return c.Timeout
}

func ReadFromFile(path string, z *zap.SugaredLogger) (data []byte, err error) {
	data, err = os.ReadFile(path)
	if err != nil {
		z.Errorf("can't open config file: %s", err)
		return nil, fmt.Errorf("can't open config file: %s", err)
	}
	return data, err
}

func (conf *ConfYML) ReadConfig(data []byte, z *zap.SugaredLogger) (err error) {
	fakestruct := struct {
		Timeout  int    `yaml:"timeout"`
		LogLevel uint8  `yaml:"loglevel"`
		LogFile  string `yaml:"logfile"`
	}{}
	err = yaml.Unmarshal(data, &fakestruct)
	if err != nil {
		z.Errorf("can't unmarshall data: %s", err)
		return fmt.Errorf("can't unmarshall data: %s", err)
	}
	conf.LogFile = fakestruct.LogFile
	conf.Timeout = time.Duration(fakestruct.Timeout) * time.Second
	conf.LogLevel = zapcore.Level(fakestruct.LogLevel)
	return nil
}
