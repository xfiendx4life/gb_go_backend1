package user

import (
	"context"

	//// "github.com/labstack/echo"
	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/gb_go_backend1/pkg/models"
	"go.uber.org/zap"
)

type Deliver interface {
	Login(ectx echo.Context) error
	Create(ectx echo.Context) error
}

type UseCase interface {
	Validate(ctx context.Context, name, password string, z *zap.SugaredLogger) (bool, error)
	Add(ctx context.Context, name, password, email string, z *zap.SugaredLogger) (*models.User, error)
}
