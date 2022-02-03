//go:build integration
// +build integration

package repository_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_backend1/internal/logger"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/user/repository"
	"github.com/xfiendx4life/gb_go_backend1/storage"
	"go.uber.org/zap/zapcore"
)

var lgr = logger.InitLogger(zapcore.DebugLevel, "")
var pg storage.Storage
var mc = mockConfStorage{}

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

func TestMain(m *testing.M) {
	setUp()
	m.Run()
	tearDown()
}

func setUp() {
	os.Setenv("MAX_CONS", "10")
	os.Setenv("MIN_CONS", "5")
	pg = storage.New()
	pg.InitNewStorage(context.Background(), lgr, &mc)
}

func tearDown() {
	// os.Setenv("MAX_CONS", "0")
	// os.Setenv("MIN_CONS", "0")
	pg.GetDbPool().Exec(context.Background(), "DELETE FROM redirects;")
	pg.GetDbPool().Exec(context.Background(), "DELETE FROM urls;")
	pg.GetDbPool().Exec(context.Background(), "DELETE FROM users;")
	pg.GetDbPool().Exec(context.Background(), "ALTER SEQUENCE urls_id_seq RESTART WITH 1;")
	pg.GetDbPool().Exec(context.Background(), "ALTER SEQUENCE users_id_seq RESTART WITH 1;")
	time.Sleep(time.Millisecond)
}

func TestAddUser(t *testing.T) {
	ctx := context.Background()
	repo := repository.New(pg)
	u := models.User{
		Name:     "TestAddUser",
		Password: "03212345",
		Email:    "somemail@fnd.ru",
	}
	err := repo.AddUser(ctx, &u, lgr)
	assert.NoError(t, err)
	tearDown()
}

func TestAddUserError(t *testing.T) {
	ctx := context.Background()
	repo := repository.New(pg)
	q := `INSERT INTO users (name, password, email) VALUES ($1, $2, $3) RETURNING id`
	_ = pg.GetDbPool().QueryRow(ctx, q, "TestAddUserError", "7892345", "somemail@fnd.ru")
	u := models.User{
		Name:     "TestAddUserError",
		Password: "2345",
		Email:    "somemail@fnd.ru",
	}
	err := repo.AddUser(ctx, &u, lgr)
	assert.Error(t, err)
	tearDown()
}

func TestGetUser(t *testing.T) {
	ctx := context.Background()
	repo := repository.New(pg)
	q := `INSERT INTO users (name, password, email) VALUES ($1, $2, $3) RETURNING id`
	expected := models.User{
		Id:       1,
		Name:     "TestGetUser",
		Password: "2345",
		Email:    "somemail@fnd.ru",
	}
	_ = pg.GetDbPool().QueryRow(ctx, q, expected.Name, expected.Password, expected.Email).Scan(&expected.Id)
	u, err := repo.GetUserByLogin(ctx, "TestGetUser", lgr)
	assert.NoError(t, err)
	assert.Equal(t, expected, *u)
	tearDown()
}

func TestGetUserError(t *testing.T) {
	ctx := context.Background()
	repo := repository.New(pg)
	u, err := repo.GetUserByLogin(ctx, "TestGetUserError", lgr)
	assert.NoError(t, err)
	assert.Equal(t, models.User{}, *u)
	tearDown()
}
