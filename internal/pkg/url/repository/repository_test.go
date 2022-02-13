//go:build integration
// +build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfiendx4life/gb_go_backend1/internal/logger"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/url/repository"
	"github.com/xfiendx4life/gb_go_backend1/storage"
	"go.uber.org/zap/zapcore"
)

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

var (
	mc  = mockConfStorage{}
	lgr = logger.InitLogger(zapcore.DebugLevel, "")
	pg  storage.Storage
)

func TestMain(m *testing.M) {
	setUp()
	m.Run()
	tearDown()
}

func setUp() {
	pg = storage.New()
	err := pg.InitNewStorage(context.Background(), lgr, &mc)
	if err != nil {
		lgr.Fatalf("can't intialize storage %s", err)
	}
}

func tearDown() {
	pg.GetDbPool().Exec(context.Background(), "DELETE FROM redirects;")
	pg.GetDbPool().Exec(context.Background(), "DELETE FROM urls;")
	pg.GetDbPool().Exec(context.Background(), "DELETE FROM users;")
	pg.GetDbPool().Exec(context.Background(), "ALTER SEQUENCE urls_id_seq RESTART WITH 1;")
	pg.GetDbPool().Exec(context.Background(), "ALTER SEQUENCE users_id_seq RESTART WITH 1;")
	time.Sleep(time.Millisecond)
}

func TestAddUrl(t *testing.T) {
	ctx := context.Background()
	repo := repository.New(pg, lgr)
	q := `INSERT INTO users (name, password, email) VALUES ($1, $2, $3) RETURNING id`
	_ = pg.GetDbPool().QueryRow(ctx, q, "TestAddUrl", "2300045", "somemail@fnd.ru")
	url := models.Url{
		Raw:       "https://google.com",
		Shortened: "shorturl.at/huNP1",
		UserId:    1,
	}
	err := repo.AddUrl(ctx, &url)
	assert.NoError(t, err)
	assert.Equal(t, 1, url.Id)
	tearDown()
}

func TestGetUrls(t *testing.T) {
	ctx := context.Background()
	repo := repository.New(pg, lgr)
	q := `INSERT INTO users (name, password, email) VALUES ($1, $2, $3) RETURNING id`
	var userId int
	err := pg.GetDbPool().QueryRow(ctx, q, "TestGetUrls", "TestGetUrl", "somemail@fnd.ru").Scan(&userId)
	assert.NoError(t, err)
	for _, i := range "0123456789" {
		url := models.Url{
			Raw:       "https://google.com" + string(i),
			Shortened: "shorturl.at/huNP1",
			UserId:    userId,
		}
		q = `INSERT INTO urls (raw, shortened, user_id) VALUES ($1, $2, $3) RETURNING id`
		err = pg.GetDbPool().QueryRow(ctx, q, url.Raw, url.Shortened, url.UserId).Scan(&url.Id)
		assert.NoError(t, err)
	}
	urls, err := repo.GetUrls(ctx, userId)
	assert.NoError(t, err)
	assert.Equal(t, 10, len(urls))
	assert.Equal(t, "0", string(urls[0].Raw[len(urls[0].Raw)-1]))
	tearDown()
}

func TestGetUrlByShortened(t *testing.T) {
	ctx := context.Background()
	repo := repository.New(pg, lgr)
	q := `INSERT INTO users (name, password, email) VALUES ($1, $2, $3) RETURNING id`
	var userId int
	err := pg.GetDbPool().QueryRow(ctx, q, "TestGetUrlBySH", "TestGetUrlBySH", "somemail@fnd.ru").Scan(&userId)
	assert.NoError(t, err)
	url := models.Url{
		Raw:       "https://google.com",
		Shortened: "TestGetUrlByShortened",
		UserId:    userId,
	}
	q = `INSERT INTO urls (raw, shortened, user_id) VALUES ($1, $2, $3) RETURNING id`
	err = pg.GetDbPool().QueryRow(ctx, q, url.Raw, url.Shortened, url.UserId).Scan(&url.Id)
	assert.NoError(t, err)
	res, err := repo.GetUrlByShortened(ctx, "TestGetUrlByShortened")
	require.NoError(t, err)
	assert.Equal(t, url, *res)
	tearDown()
}
