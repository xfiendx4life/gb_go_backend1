package config

import (
	"fmt"
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/yaml.v3"
)

type fakestruct struct {
	Timeout     int    `yaml:"timeout"`
	LogLevel    string `yaml:"loglevel"`
	LogFile     string `yaml:"logfile"`
	Port        string `yaml:"port"`
	ConfStorage `yaml:",inline"`
	ConfAuth    `yaml:",inline"`
}

func New() Config {
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

func (st *ConfStorage) GetURI() string {
	return st.Uri
}

func (st *ConfStorage) GetMaxCons() int {
	return st.MaxCons
}

func (st *ConfStorage) GetMinCons() int {
	return st.MinCons
}

func (a *ConfAuth) GetSecretKey() string {
	return a.SecretKey
}

func (a *ConfAuth) GetTtl() int64 {
	return a.Ttl
}

func (c *ConfYML) GetPort() string {
	return c.Port
}

func ReadFromFile(path string, z *zap.SugaredLogger) (data []byte, err error) {
	data, err = os.ReadFile(path)
	if err != nil {
		z.Errorf("can't open config file: %s", err)
		return nil, fmt.Errorf("can't open config file: %s", err)
	}
	z.Infof("Read from file: %s", path)
	return data, err
}

func (c *ConfYML) GetConfAuth() Auth {
	return &c.ConfAuth
}

func (c *ConfYML) GetConfStorage() Storage {
	return &c.ConfStorage
}

func ReadFromEnv() []byte {
	timeout := os.Getenv("TIMEOUT")
	loglevel := os.Getenv("LOGLEVEL")
	logfile := os.Getenv("LOGFILE")
	uri := os.Getenv("DATABASE_URL")
	maxcons := os.Getenv("MAXCONS")
	mincons := os.Getenv("MINCONS")
	secretkey := os.Getenv("SECRETKEY")
	ttl := os.Getenv("TTL")
	port := os.Getenv("PORT")
	return []byte(fmt.Sprintf(`timeout: %s
loglevel: %s
logfile: %s
uri: %s
maxcons: %s
mincons: %s
secretkey: %s
port: %s
ttl: %s`, timeout, loglevel, logfile, uri, maxcons, mincons, secretkey, port, ttl))
}

func (conf *ConfYML) ReadConfig(data []byte, z *zap.SugaredLogger) (err error) {
	fake := fakestruct{}
	z.Infof("Config data: %s", string(data))
	err = yaml.Unmarshal(data, &fake)
	if err != nil {
		z.Errorf("can't unmarshall data: %s", err)
		return fmt.Errorf("can't unmarshall data: %s", err)
	}
	var levels = map[string]zapcore.Level{
		"debug":   zap.DebugLevel,
		"info":    zap.InfoLevel,
		"warning": zap.WarnLevel,
		"error":   zap.ErrorLevel,
		"panic":   zap.PanicLevel,
		"fatal":   zap.FatalLevel,
	}
	conf.LogFile = fake.LogFile
	conf.Timeout = time.Duration(fake.Timeout) * time.Second
	conf.LogLevel = zapcore.Level(levels[fake.LogLevel])
	conf.ConfStorage = fake.ConfStorage
	conf.ConfAuth = fake.ConfAuth
	conf.Port = fake.Port
	z.Infof("Config file: %v", fake)
	return nil
}
