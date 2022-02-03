package url

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"go.uber.org/zap"
)

type UseCase interface {
	Add(ctx context.Context, raw string, userId int, z *zap.SugaredLogger) (shortened string, err error)
	// * Get method changes stats
	Get(ctx context.Context, shortened string, z *zap.SugaredLogger) (raw string, err error)
	List(ctx context.Context, userId int, z *zap.SugaredLogger) ([]models.Url, error)
}

type Deliver interface {
	Save(ectx echo.Context) error
	Get(ectx echo.Context) (string, error)
	List(ectx echo.Context) ([]models.Url, error)
}

type Repository interface {
	AddUrl(ctx context.Context, url *models.Url, z *zap.SugaredLogger) error
	GetUrls(ctx context.Context, userID int, z *zap.SugaredLogger) ([]models.Url, error)
	GetUrlByShortened(ctx context.Context, shortened string, z *zap.SugaredLogger) (*models.Url, error)
	// ! Here for a while but have to move it to repository pkg
	AddRedirect(ctx context.Context, redirect *models.Redirects, z *zap.SugaredLogger) error
}
