package usecase_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_backend1/internal/logger"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/redirects/usecase"
)

type mockStorage struct {
	err error
}

func (mc *mockStorage) AddRedirect(ctx context.Context, redirect *models.Redirects) error {
	return mc.err
}

func (mc *mockStorage) GetRedirects(ctx context.Context, shortened string) ([]models.Redirects, error) {
	if mc.err != nil {
		return nil, mc.err
	}
	return []models.Redirects{
		{Id: 1, UrlId: 1, Date: time.Now().Add(time.Duration(time.Hour * (-72)))},
		{Id: 2, UrlId: 1, Date: time.Now().Add(time.Duration(time.Hour * (-2)))},
		{Id: 3, UrlId: 1, Date: time.Now().Add(time.Duration(time.Hour * (-179)))},
		{Id: 4, UrlId: 1, Date: time.Now().Add(time.Duration(time.Hour * (-1)))},
	}, nil
}

var (
	lgr        = logger.InitLogger(-1, "")
	ctx        = context.Background()
	correctGet = models.Summary{
		Week:  3,
		Today: 2,
		Month: 4,
	}
)

func TestAdd(t *testing.T) {
	uc := usecase.New(&mockStorage{}, lgr)
	err := uc.Add(ctx, 1)
	assert.NoError(t, err)
}

func TestAddError(t *testing.T) {
	uc := usecase.New(&mockStorage{err: fmt.Errorf("testerror")}, lgr)
	err := uc.Add(ctx, 1)
	assert.Error(t, err)
}

func TestGet(t *testing.T) {
	uc := usecase.New(&mockStorage{}, lgr)
	sum, err := uc.Get(ctx, "test")
	assert.NoError(t, err)
	assert.Equal(t, correctGet, sum)
}

func TestGetError(t *testing.T) {
	uc := usecase.New(&mockStorage{err: fmt.Errorf("testerror")}, lgr)
	_, err := uc.Get(ctx, "test")
	assert.Error(t, err)
}
