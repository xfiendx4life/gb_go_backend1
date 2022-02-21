package usecase_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/xfiendx4life/gb_go_backend1/internal/logger"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	rdrUse "github.com/xfiendx4life/gb_go_backend1/internal/pkg/redirects/usecase"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/url/usecase"
	"github.com/xfiendx4life/gb_go_backend1/storage"
	"go.uber.org/zap/zapcore"
)

type mockStorage struct {
	err error
}

func (mc *mockStorage) AddUrl(ctx context.Context, url *models.Url) error {
	url.Id = 1
	return mc.err
}

func (mc *mockStorage) GetUrl(ctx context.Context, id int) (*models.Url, error) {
	return &models.Url{
		Raw:       "RawTestUrl",
		Shortened: "shortenedTestUrl",
	}, mc.err
}

func (mc *mockStorage) GetUrls(ctx context.Context, userID int) ([]models.Url, error) {
	return nil, nil
}

func (mc *mockStorage) GetUrlByShortened(ctx context.Context, shortened string) (*models.Url, error) {
	return &models.Url{
		Raw:       "RawTestUrl",
		Shortened: "shortenedTestUrl",
	}, mc.err
}

func (mc *mockStorage) AddRedirect(ctx context.Context, r *models.Redirects) error {
	return nil
}

func (mc *mockStorage) GetStorage() storage.Storage {
	return storage.New() // ! shit happens here because we shouldn't return real storage
}

type mockRdrStorage struct {
	err error
}

func (mr *mockRdrStorage) AddRedirect(ctx context.Context, redirect *models.Redirects) error {
	return mr.err
}

func (mr *mockRdrStorage) GetRedirects(ctx context.Context, shortened string) ([]models.Redirects, error) {
	if mr.err != nil {
		return nil, mr.err
	}
	return []models.Redirects{
		{Id: 1, UrlId: 1, Date: time.Now().Add(time.Duration(time.Hour * (-72)))},
		{Id: 2, UrlId: 1, Date: time.Now().Add(time.Duration(time.Hour * (-2)))},
		{Id: 3, UrlId: 1, Date: time.Now().Add(time.Duration(time.Hour * (-179)))},
		{Id: 4, UrlId: 1, Date: time.Now().Add(time.Duration(time.Hour * (-1)))},
	}, nil
}

var (
	lgr   = logger.InitLogger(zapcore.DebugLevel, "")
	ctx   = context.Background()
	rdrUc = rdrUse.New(&mockRdrStorage{}, lgr)
)

func TestAddUrl(t *testing.T) {
	st := mockStorage{err: nil}
	uc := usecase.New(&st, rdrUc, lgr)
	short, err := uc.Add(ctx, "testurl", 1)
	assert.NoError(t, err)
	assert.NotEmpty(t, short)
}

func TestAddStorageError(t *testing.T) {
	st := mockStorage{err: fmt.Errorf("some storage error")}
	uc := usecase.New(&st, rdrUc, lgr)
	short, err := uc.Add(ctx, "testurl", 1)
	assert.Error(t, err)
	assert.Empty(t, short)
}

func isClean(s string) bool {
	et := "!*'();:@&=+$,/\\?%#[]"
	for _, chr := range s {
		for _, sym := range et {
			if sym == chr {
				return false
			}
		}
	}
	return true
}

func TestNewUrl(t *testing.T) {
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
		require.True(t, isClean(u.Shortened), fmt.Errorf("not clean %s", u.Shortened))
		set[u.Shortened] = struct{}{}
	}
	assert.NoError(t, err, fmt.Sprintf("Url: %s, i: %d", u.Shortened, i))

}

func TestGetNoErr(t *testing.T) {
	uc := usecase.New(&mockStorage{}, rdrUc, lgr)
	raw, err := uc.Get(ctx, "shortenedTestUrl")
	assert.NoError(t, err)
	assert.Equal(t, "RawTestUrl", raw)
}

func TestGetErr(t *testing.T) {
	uc := usecase.New(&mockStorage{err: fmt.Errorf("can't add error")}, rdrUc, lgr)
	raw, err := uc.Get(ctx, "SomeNotValidURL")
	assert.Error(t, err)
	assert.Equal(t, "", raw)
}
