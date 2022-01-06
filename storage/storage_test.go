package storage

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_backend1/internal/logger"
	"github.com/xfiendx4life/gb_go_backend1/pkg/models"
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
	os.Setenv("MAX_CONS", "0")
	os.Setenv("MIN_CONS", "0")
	pg.dbPool.Exec(context.Background(), "DELETE FROM urls")
	pg.dbPool.Exec(context.Background(), "DELETE FROM users")
	pg.dbPool.Exec(context.Background(), "ALTER SEQUENCE urls_id_seq RESTART WITH 1")
	pg.dbPool.Exec(context.Background(), "ALTER SEQUENCE users_id_seq RESTART WITH 1")
}

func TestConection(t *testing.T) {
	err := pg.InitNewStorage(context.Background(), "postgres://xfiendx4life:123456@172.17.0.2:5432/shortener", lgr)
	assert.NoError(t, err)
}

func TestAddUser(t *testing.T) {
	ctx := context.Background()
	err := pg.InitNewStorage(ctx, "postgres://xfiendx4life:123456@172.17.0.2:5432/shortener", lgr)
	assert.NoError(t, err)
	u := models.User{
		Name:     "sk*nk",
		Password: "2345",
		Email:    "somemail@fnd.ru",
	}
	err = pg.AddUser(ctx, &u, lgr)
	assert.NoError(t, err)
}

func TestAddUserError(t *testing.T) {
	ctx := context.Background()
	err := pg.InitNewStorage(ctx, "postgres://xfiendx4life:123456@172.17.0.2:5432/shortener", lgr)
	assert.NoError(t, err)
	q := `INSERT INTO users (name, password, email) VALUES ($1, $2, $3) RETURNING id`
	_ = pg.dbPool.QueryRow(ctx, q, "sk*nk", "2345", "somemail@fnd.ru")
	u := models.User{
		Name:     "sk*nk",
		Password: "2345",
		Email:    "somemail@fnd.ru",
	}
	err = pg.AddUser(ctx, &u, lgr)
	assert.Error(t, err)
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	err := pg.InitNewStorage(ctx, "postgres://xfiendx4life:123456@172.17.0.2:5432/shortener", lgr)
	assert.NoError(t, err)
	q := `INSERT INTO users (name, password, email) VALUES ($1, $2, $3) RETURNING id`
	_ = pg.dbPool.QueryRow(ctx, q, "sk*nk", "2345", "somemail@fnd.ru")
	u, err := pg.GetUserByLogin(ctx, "sk*nk", lgr)
	assert.NoError(t, err)
	expected := models.User{
		Id:       1,
		Name:     "sk*nk",
		Password: "2345",
		Email:    "somemail@fnd.ru",
	}
	assert.Equal(t, expected, *u)
}

func TestGetUserError(t *testing.T) {
	ctx := context.Background()
	err := pg.InitNewStorage(ctx, "postgres://xfiendx4life:123456@172.17.0.2:5432/shortener", lgr)
	assert.NoError(t, err)
	_, err = pg.GetUserByLogin(ctx, "sknk", lgr)
	assert.Error(t, err)

}

func TestAddUrl(t *testing.T) {
	ctx := context.Background()
	err := pg.InitNewStorage(ctx, "postgres://xfiendx4life:123456@172.17.0.2:5432/shortener", lgr)
	assert.NoError(t, err)
	q := `INSERT INTO users (name, password, email) VALUES ($1, $2, $3) RETURNING id`
	_ = pg.dbPool.QueryRow(ctx, q, "sk*nk", "2345", "somemail@fnd.ru")
	url := models.Url{
		Raw:       "https://google.com",
		Shortened: "shorturl.at/huNP1",
		UserId:    1,
	}
	err = pg.AddUrl(ctx, &url, lgr)
	assert.NoError(t, err)
	assert.Equal(t, 1, url.Id)
}
