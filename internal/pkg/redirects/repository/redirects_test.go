//go:build integration
// +build integration

package repository_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_backend1/internal/logger"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/redirects/repository"
	"github.com/xfiendx4life/gb_go_backend1/storage"
)

var (
	lgr = logger.InitLogger(-1, "")
	pg  = storage.New()
	mc  = mockConfStorage{}
)

func TestMain(m *testing.M) {
	m.Run()
	tearDown()
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

func tearDown() {
	pg.GetDbPool().Exec(context.Background(), "DELETE FROM redirects;")
	pg.GetDbPool().Exec(context.Background(), "DELETE FROM urls;")
	pg.GetDbPool().Exec(context.Background(), "DELETE FROM users;")
	pg.GetDbPool().Exec(context.Background(), "ALTER SEQUENCE urls_id_seq RESTART WITH 1;")
	pg.GetDbPool().Exec(context.Background(), "ALTER SEQUENCE users_id_seq RESTART WITH 1;")
	time.Sleep(time.Millisecond)
}

func TestAddRedirect(t *testing.T) {
	ctx := context.Background()
	err := pg.InitNewStorage(ctx, lgr, &mc)
	r := repository.New(pg, lgr)
	assert.NoError(t, err)
	q := `INSERT INTO users (name, password, email) VALUES ($1, $2, $3) RETURNING id`
	var userId int
	err = pg.GetDbPool().QueryRow(ctx, q, "TestAddRedirect", "TestAddRedirect", "somemail@fnd.ru").Scan(&userId)
	assert.NoError(t, err)
	url := models.Url{
		Raw:       "https://google.com",
		Shortened: "TestAddRedirect",
		UserId:    userId,
	}
	q = `INSERT INTO urls (raw, shortened, user_id) VALUES ($1, $2, $3) RETURNING id`
	err = pg.GetDbPool().QueryRow(ctx, q, url.Raw, url.Shortened, url.UserId).Scan(&url.Id)
	assert.NoError(t, err)

	rdr := models.Redirects{
		UrlId: url.Id,
		Date:  time.Now(),
	}
	err = r.AddRedirect(ctx, &rdr)
	assert.NoError(t, err)
	assert.NotEqual(t, 0, rdr.Id)
	tearDown()
}

func TestGetRedirects(t *testing.T) {
	ctx := context.Background()
	err := pg.InitNewStorage(ctx, lgr, &mc)
	assert.NoError(t, err)
	r := repository.New(pg, lgr)
	q := `INSERT INTO users (name, password, email) VALUES ($1, $2, $3) RETURNING id`
	var userId int
	err = pg.GetDbPool().QueryRow(ctx, q, "TestGetRedirects", "TestGetRedirects", "somemail@fnd.ru").Scan(&userId)
	assert.NoError(t, err)
	url := models.Url{
		Raw:       "https://google.com",
		Shortened: "TestGetRedirects",
		UserId:    userId,
	}
	q = `INSERT INTO urls (raw, shortened, user_id) VALUES ($1, $2, $3) RETURNING id`
	err = pg.GetDbPool().QueryRow(ctx, q, url.Raw, url.Shortened, url.UserId).Scan(&url.Id)
	assert.NoError(t, err)
	rdr := models.Redirects{
		UrlId: url.Id,
		Date:  time.Now().Add(time.Duration(-2) * time.Hour),
	}
	q = `INSERT INTO redirects (url_id, date_of_usage) VALUES ($1, $2) RETURNING id`
	err = pg.GetDbPool().QueryRow(ctx, q, rdr.UrlId, rdr.Date).Scan(&rdr.Id)
	assert.NoError(t, err)
	rdr = models.Redirects{
		UrlId: url.Id,
		Date:  time.Now(),
	}
	q = `INSERT INTO redirects (url_id, date_of_usage) VALUES ($1, $2) RETURNING id`
	err = pg.GetDbPool().QueryRow(ctx, q, rdr.UrlId, rdr.Date).Scan(&rdr.Id)
	assert.NoError(t, err)
	redirects, err := r.GetRedirects(ctx, url.Shortened)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(redirects))
	tearDown()
}
