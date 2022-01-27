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
	GetConfAuth() Auth
	GetConfStorage() Storage
}

type ConfYML struct {
	Timeout  time.Duration `yaml:"timeout"`
	LogLevel zapcore.Level `yaml:"loglevel"`
	LogFile  string        `yaml:"logfile"`
	ConfStorage
	ConfAuth
}

type Storage interface {
	GetURI() string
	GetMaxCons() int
	GetMinCons() int
}

type Auth interface {
	GetSecretKey() string
	GetTtl() int64
}

type ConfStorage struct {
	Uri     string `yaml:"uri"`
	MaxCons int    `yaml:"maxcons"`
	MinCons int    `yaml:"mincons"`
}

type ConfAuth struct {
	SecretKey string `yaml:"secretkey"`
	Ttl       int64  `yaml:"ttl"`
}
