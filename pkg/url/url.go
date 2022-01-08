package url

import (
	"context"

	"github.com/labstack/echo"
	"github.com/xfiendx4life/gb_go_backend1/pkg/models"
	"go.uber.org/zap"
)

type UseCase interface {
	Add(ctx context.Context, raw string, z *zap.SugaredLogger) (shortened string, err error)
	// * Get method changes stats
	Get(ctx context.Context, shortened string, z *zap.SugaredLogger) (raw string, err error)
	List(ctx context.Context, userId int, z *zap.SugaredLogger) ([]models.Url, error)
}

type Deliver interface {
	Save(ectx echo.Context) error
	Get(ectx echo.Context) (*models.Url, error)
	List(ectx echo.Context) ([]models.Url, error)
}
