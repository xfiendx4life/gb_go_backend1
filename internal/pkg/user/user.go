package user

import (
	"context"

	//// "github.com/labstack/echo"
	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"go.uber.org/zap"
)

type Deliver interface {
	Login(ectx echo.Context) error
	Create(ectx echo.Context) error
}

type UseCase interface {
	Validate(ctx context.Context, name, password string, z *zap.SugaredLogger) (bool, error)
	Add(ctx context.Context, user *models.User, z *zap.SugaredLogger) error
}

type Repository interface {
	AddUser(ctx context.Context, user *models.User, z *zap.SugaredLogger) error
	GetUserByLogin(ctx context.Context, login string, z *zap.SugaredLogger) (*models.User, error)
}
