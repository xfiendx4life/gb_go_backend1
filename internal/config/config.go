package config

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type Config interface {
	// read config from file with path
	ReadConfig(data []byte, z *zap.SugaredLogger) (err error)
	GetLogLevel() zapcore.Level
	GetLogFile() string
	GetTimeOut() time.Duration
}

type ConfYML struct {
	Timeout  time.Duration `yaml:"timeout"`
	LogLevel zapcore.Level `yaml:"loglevel"`
	LogFile  string        `yaml:"logfile"`
}
