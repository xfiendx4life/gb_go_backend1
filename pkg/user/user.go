package user

import (
	"context"

	//// "github.com/labstack/echo"
	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/gb_go_backend1/pkg/models"
	"go.uber.org/zap"
)

type Delivery interface {
	Login(ectx echo.Context) error
	Create(ectx echo.Context) error
}

type UseCase interface {
	Validate(ctx context.Context, name, password string, z *zap.SugaredLogger) (bool, error)
	Add(ctx context.Context, name, password, email string, z *zap.SugaredLogger) (*models.User, error)
}

// type Repository interface {
// 	Get(ctx context.Context, login string, z *zap.SugaredLogger) (*models.User, error)
// 	Add(ctx context.Context, name, password, email string, z *zap.SugaredLogger) (*models.User, error)
// }
