package usecase

import (
	"context"
	"fmt"

	"github.com/xfiendx4life/gb_go_backend1/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/pkg/url"
	"github.com/xfiendx4life/gb_go_backend1/storage"
	"go.uber.org/zap"
)

type gres struct {
	store storage.Storage
}

func New(st storage.Storage) url.UseCase {
	return &gres{store: st}
}

func (g *gres) Add(ctx context.Context, raw string, userId int, z *zap.SugaredLogger) (shortened string, err error) {
	url := models.NewUrl(raw, userId, z)
	err = g.store.AddUrl(ctx, url, z)
	if err != nil {
		z.Errorf("can't add url to storage: %s", err)
		return "", fmt.Errorf("can't add url to storage: %s", err)
	}
	z.Infof("Url %s was successfully added to storage and shortened to -> %s", url.Raw, url.Shortened)
	return url.Shortened, nil
}

func (g *gres) Get(ctx context.Context, shortened string, z *zap.SugaredLogger) (raw string, err error) {
	return "", nil
}

func (g *gres) List(ctx context.Context, userId int, z *zap.SugaredLogger) ([]models.Url, error) {
	return []models.Url{}, nil
}
