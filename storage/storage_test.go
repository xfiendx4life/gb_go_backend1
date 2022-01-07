//go:build integration
// +build integration

package storage

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
	pg.dbPool.Exec(context.Background(), "DELETE FROM urls;")
	pg.dbPool.Exec(context.Background(), "DELETE FROM users;")
	pg.dbPool.Exec(context.Background(), "ALTER SEQUENCE urls_id_seq RESTART WITH 1;")
	pg.dbPool.Exec(context.Background(), "ALTER SEQUENCE users_id_seq RESTART WITH 1;")
	time.Sleep(time.Millisecond)
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
		Name:     "TestAddUser",
		Password: "03212345",
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
	_ = pg.dbPool.QueryRow(ctx, q, "TestAddUserError", "7892345", "somemail@fnd.ru")
	u := models.User{
		Name:     "TestAddUserError",
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
	expected := models.User{
		Id:       1,
		Name:     "TestGetUser",
		Password: "2345",
		Email:    "somemail@fnd.ru",
	}
	_ = pg.dbPool.QueryRow(ctx, q, expected.Name, expected.Password, expected.Email).Scan(&expected.Id)
	u, err := pg.GetUserByLogin(ctx, "TestGetUser", lgr)
	assert.NoError(t, err)
	assert.Equal(t, expected, *u)
}

func TestGetUserError(t *testing.T) {
	ctx := context.Background()
	err := pg.InitNewStorage(ctx, "postgres://xfiendx4life:123456@172.17.0.2:5432/shortener", lgr)
	assert.NoError(t, err)
	_, err = pg.GetUserByLogin(ctx, "TestGetUserError", lgr)
	assert.Error(t, err)

}

func TestAddUrl(t *testing.T) {
	ctx := context.Background()
	err := pg.InitNewStorage(ctx, "postgres://xfiendx4life:123456@172.17.0.2:5432/shortener", lgr)
	assert.NoError(t, err)
	q := `INSERT INTO users (name, password, email) VALUES ($1, $2, $3) RETURNING id`
	_ = pg.dbPool.QueryRow(ctx, q, "TestAddUrl", "2300045", "somemail@fnd.ru")
	url := models.Url{
		Raw:       "https://google.com",
		Shortened: "shorturl.at/huNP1",
		UserId:    1,
	}
	err = pg.AddUrl(ctx, &url, lgr)
	assert.NoError(t, err)
	assert.Equal(t, 1, url.Id)
}

func TestGetUrl(t *testing.T) {
	ctx := context.Background()
	err := pg.InitNewStorage(ctx, "postgres://xfiendx4life:123456@172.17.0.2:5432/shortener", lgr)
	assert.NoError(t, err)
	q := `INSERT INTO users (name, password, email) VALUES ($1, $2, $3) RETURNING id`
	var userId int
	err = pg.dbPool.QueryRow(ctx, q, "TestGetUrl", "TestGetUrl", "somemail@fnd.ru").Scan(&userId)
	assert.NoError(t, err)
	url := models.Url{
		Raw:       "https://google.com",
		Shortened: "shorturl.at/huNP1",
		UserId:    userId,
	}
	q = `INSERT INTO urls (raw, shortened, user_id) VALUES ($1, $2, $3) RETURNING id`
	err = pg.dbPool.QueryRow(ctx, q, url.Raw, url.Shortened, url.UserId).Scan(&url.Id)
	assert.NoError(t, err)
	res, err := pg.GetUrl(ctx, userId, lgr)
	require.NoError(t, err)
	assert.Equal(t, url, *res)
}

func TestGetUrls(t *testing.T) {
	ctx := context.Background()
	err := pg.InitNewStorage(ctx, "postgres://xfiendx4life:123456@172.17.0.2:5432/shortener", lgr)
	assert.NoError(t, err)
	q := `INSERT INTO users (name, password, email) VALUES ($1, $2, $3) RETURNING id`
	var userId int
	err = pg.dbPool.QueryRow(ctx, q, "TestGetUrls", "TestGetUrl", "somemail@fnd.ru").Scan(&userId)
	assert.NoError(t, err)
	for _, i := range "0123456789" {
		url := models.Url{
			Raw:       "https://google.com" + string(i),
			Shortened: "shorturl.at/huNP1",
			UserId:    userId,
		}
		q = `INSERT INTO urls (raw, shortened, user_id) VALUES ($1, $2, $3) RETURNING id`
		err = pg.dbPool.QueryRow(ctx, q, url.Raw, url.Shortened, url.UserId).Scan(&url.Id)
		assert.NoError(t, err)
	}
	urls, err := pg.GetUrls(ctx, userId, lgr)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(urls))
	assert.Equal(t, "0", string(urls[0].Raw[len(urls[0].Raw)-1]))

}

func TestGetUrlByShortened(t *testing.T) {
	ctx := context.Background()
	err := pg.InitNewStorage(ctx, "postgres://xfiendx4life:123456@172.17.0.2:5432/shortener", lgr)
	assert.NoError(t, err)
	q := `INSERT INTO users (name, password, email) VALUES ($1, $2, $3) RETURNING id`
	var userId int
	err = pg.dbPool.QueryRow(ctx, q, "TestGetUrlBySH", "TestGetUrlBySH", "somemail@fnd.ru").Scan(&userId)
	assert.NoError(t, err)
	url := models.Url{
		Raw:       "https://google.com",
		Shortened: "TestGetUrlByShortened",
		UserId:    userId,
	}
	q = `INSERT INTO urls (raw, shortened, user_id) VALUES ($1, $2, $3) RETURNING id`
	err = pg.dbPool.QueryRow(ctx, q, url.Raw, url.Shortened, url.UserId).Scan(&url.Id)
	assert.NoError(t, err)
	res, err := pg.GetUrlByShortened(ctx, "TestGetUrlByShortened", lgr)
	require.NoError(t, err)
	assert.Equal(t, url, *res)
}
