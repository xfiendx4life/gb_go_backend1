//go:build integration
// +build integration

package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_backend1/internal/logger"
	"go.uber.org/zap/zapcore"
)

var lgr = logger.InitLogger(zapcore.DebugLevel, "")
var pg PG

func TestMain(m *testing.M) {
	setUp()
	m.Run()
	tearDown()
}

func setUp() {
	os.Setenv("MAX_CONS", "10")
	os.Setenv("MIN_CONS", "5")
	pg = PG{}
}

func tearDown() {
	// os.Setenv("MAX_CONS", "0")
	// os.Setenv("MIN_CONS", "0")
	pg.dbPool.Exec(context.Background(), "DELETE FROM redirects;")
	pg.dbPool.Exec(context.Background(), "DELETE FROM urls;")
	pg.dbPool.Exec(context.Background(), "DELETE FROM users;")
	pg.dbPool.Exec(context.Background(), "ALTER SEQUENCE urls_id_seq RESTART WITH 1;")
	pg.dbPool.Exec(context.Background(), "ALTER SEQUENCE users_id_seq RESTART WITH 1;")
	time.Sleep(time.Millisecond)
}

type mockConfStorage struct {
}

func (m *mockConfStorage) GetURI() string {
	return "postgres://xfiendx4life:123456@172.17.0.2:5432/shortener"
}

func (m *mockConfStorage) GetMaxCons() int {
	return 10
}

func (m *mockConfStorage) GetMinCons() int {
	return 5
}

var mc = mockConfStorage{}

func TestConection(t *testing.T) {
	err := pg.InitNewStorage(context.Background(), lgr, &mc)
	assert.NoError(t, err)
}
