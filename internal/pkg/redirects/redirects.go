package redirects

import (
	"context"

	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
)

type Repository interface {
	AddRedirect(ctx context.Context, redirect *models.Redirects) error
	GetRedirects(ctx context.Context, shortened string) ([]models.Redirects, error)
}

type UseCase interface {
	Add(ctx context.Context, id int) error
	Get(ctx context.Context, shortened string) (models.Summary, error)
}

type Deliver interface {
	GetSummary(ectx echo.Context) error
}
