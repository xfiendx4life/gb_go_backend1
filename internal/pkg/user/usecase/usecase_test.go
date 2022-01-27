package usecase_test

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_backend1/internal/config"
	"github.com/xfiendx4life/gb_go_backend1/internal/logger"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/user/usecase"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	lgr = logger.InitLogger(zapcore.DebugLevel, "")
	ctx = context.Background()
)

type mockStorage struct {
	err error
}

func (mc *mockStorage) InitNewStorage(ctx context.Context, z *zap.SugaredLogger, conf config.Storage) error {
	return nil
}

func (mc *mockStorage) AddUser(ctx context.Context, user *models.User, z *zap.SugaredLogger) error {
	user.Id = 1
	return mc.err
}

func (mc *mockStorage) GetUserByLogin(ctx context.Context, login string, z *zap.SugaredLogger) (*models.User, error) {
	s := md5.Sum([]byte("correctPassword"))
	return &models.User{
		Name:     "testname",
		Password: hex.EncodeToString(s[:]),
	}, mc.err
}

func (mc *mockStorage) AddUrl(ctx context.Context, url *models.Url, z *zap.SugaredLogger) error {
	return nil
}
func (mc *mockStorage) GetUrl(ctx context.Context, id int, z *zap.SugaredLogger) (*models.Url, error) {
	return nil, nil
}
func (mc *mockStorage) GetUrls(ctx context.Context, userID int, z *zap.SugaredLogger) ([]models.Url, error) {
	return nil, nil
}

func (mc *mockStorage) GetUrlByShortened(ctx context.Context, shortened string, z *zap.SugaredLogger) (*models.Url, error) {
	return &models.Url{}, nil
}

func (mc *mockStorage) AddRedirect(ctx context.Context, r *models.Redirects, z *zap.SugaredLogger) error {
	return nil
}

func (mc *mockStorage) GetRedirects(ctx context.Context, urlId int, z *zap.SugaredLogger) ([]models.Redirects, error) {
	return []models.Redirects{}, nil
}

func TestValidateCorrect(t *testing.T) {
	uc := usecase.New(&mockStorage{})
	res, err := uc.Validate(ctx, "testname", "correctPassword", lgr)
	assert.NoError(t, err)
	assert.True(t, res)
}

func TestValidateIncorrect(t *testing.T) {
	uc := usecase.New(&mockStorage{})
	res, err := uc.Validate(ctx, "testname", "incorrectPassword", lgr)
	assert.NoError(t, err)
	assert.False(t, res)
}

func TestValidateError(t *testing.T) {
	st := &mockStorage{err: errors.New("some error")}
	uc := usecase.New(st)
	res, err := uc.Validate(ctx, "testname", "incorrectPassword", lgr)
	assert.Error(t, err)
	assert.False(t, res)
}

func TestAddUser(t *testing.T) {
	uc := usecase.New(&mockStorage{})
	u := &models.User{Name: "testname", Password: "password", Email: "email"}
	err := uc.Add(ctx, u, lgr)
	assert.NoError(t, err)
	assert.Equal(t, models.User{
		Id:       1,
		Name:     "testname",
		Password: "5f4dcc3b5aa765d61d8327deb882cf99",
		Email:    "email"}, *u)
}

func TestAddUserError(t *testing.T) {
	st := &mockStorage{err: errors.New("some error")}
	uc := usecase.New(st)
	u := &models.User{Name: "testname", Password: "password", Email: "email"}
	err := uc.Add(ctx, u, lgr)
	assert.Error(t, err)
}
