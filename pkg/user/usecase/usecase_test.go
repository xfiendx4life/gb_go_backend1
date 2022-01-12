package usecase_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_backend1/internal/logger"
	"github.com/xfiendx4life/gb_go_backend1/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/pkg/user/usecase"
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

func (mc *mockStorage) InitNewStorage(ctx context.Context, connection string, z *zap.SugaredLogger) error {
	return nil
}

func (mc *mockStorage) AddUser(ctx context.Context, user *models.User, z *zap.SugaredLogger) error {
	user.Id = 1
	return mc.err
}

func (mc *mockStorage) GetUserByLogin(ctx context.Context, login string, z *zap.SugaredLogger) (*models.User, error) {
	return &models.User{
		Name:     "testname",
		Password: "correctPassword",
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
	u, err := uc.Add(ctx, "testname", "password", "email", lgr)
	assert.NoError(t, err)
	assert.Equal(t, models.User{
		Id:       1,
		Name:     "testname",
		Password: "password",
		Email:    "email"}, *u)
}

func TestAddUserError(t *testing.T) {
	st := &mockStorage{err: errors.New("some error")}
	uc := usecase.New(st)
	u, err := uc.Add(ctx, "testname", "password", "email", lgr)
	assert.Error(t, err)
	assert.Nil(t, u)
}
