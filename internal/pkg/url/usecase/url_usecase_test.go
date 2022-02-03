//go:build unit
// +build unit

package usecase_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_backend1/internal/config"
	"github.com/xfiendx4life/gb_go_backend1/internal/logger"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/url/usecase"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	return &models.User{
		Name:     "testname",
		Password: "correctPassword",
	}, mc.err
}

func (mc *mockStorage) AddUrl(ctx context.Context, url *models.Url, z *zap.SugaredLogger) error {
	url.Id = 1
	return mc.err
}

func (mc *mockStorage) GetUrl(ctx context.Context, id int, z *zap.SugaredLogger) (*models.Url, error) {
	return &models.Url{
		Raw:       "RawTestUrl",
		Shortened: "shortenedTestUrl",
	}, mc.err
}

func (mc *mockStorage) GetUrls(ctx context.Context, userID int, z *zap.SugaredLogger) ([]models.Url, error) {
	return nil, nil
}

func (mc *mockStorage) GetUrlByShortened(ctx context.Context, shortened string, z *zap.SugaredLogger) (*models.Url, error) {
	return &models.Url{
		Raw:       "RawTestUrl",
		Shortened: "shortenedTestUrl",
	}, mc.err
}

func (mc *mockStorage) AddRedirect(ctx context.Context, r *models.Redirects, z *zap.SugaredLogger) error {
	return nil
}

func (mc *mockStorage) GetRedirects(ctx context.Context, urlId int, z *zap.SugaredLogger) ([]models.Redirects, error) {
	return []models.Redirects{}, nil
}

var (
	lgr = logger.InitLogger(zapcore.DebugLevel, "")
	ctx = context.Background()
)

func TestAddUrl(t *testing.T) {
	st := mockStorage{err: nil}
	uc := usecase.New(&st)
	short, err := uc.Add(ctx, "testurl", 1, lgr)
	assert.NoError(t, err)
	assert.NotEmpty(t, short)
}

func TestAddStorageError(t *testing.T) {
	st := mockStorage{err: fmt.Errorf("some storage error")}
	uc := usecase.New(&st)
	short, err := uc.Add(ctx, "testurl", 1, lgr)
	assert.Error(t, err)
	assert.Empty(t, short)
}

func TestNewUser(t *testing.T) {
	set := make(map[string]struct{})
	var err error
	var i int
	var u *models.Url
	for i = 0; i < 10000; i++ {
		u = usecase.NewUrl("TestNewUser", 1, lgr)
		if _, ok := set[u.Shortened]; ok {
			err = fmt.Errorf("already exists")
			break
		}
		set[u.Shortened] = struct{}{}
	}
	assert.NoError(t, err, fmt.Sprintf("Url: %s, i: %d", u.Shortened, i))

}

func TestGetNoErr(t *testing.T) {
	uc := usecase.New(&mockStorage{})
	raw, err := uc.Get(ctx, "shortenedTestUrl", lgr)
	assert.NoError(t, err)
	assert.Equal(t, "RawTestUrl", raw)
}

func TestGetErr(t *testing.T) {
	uc := usecase.New(&mockStorage{err: fmt.Errorf("can't add error")})
	raw, err := uc.Get(ctx, "SomeNotValidURL", lgr)
	assert.Error(t, err)
	assert.Equal(t, "", raw)
}
