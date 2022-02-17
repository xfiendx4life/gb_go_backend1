package user

import (
	"context"

	//// "github.com/labstack/echo"
	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
)

type Deliver interface {
	Login(ectx echo.Context) error
	Create(ectx echo.Context) error
	CreateFromForm(ectx echo.Context) error
}

type UseCase interface {
	Validate(ctx context.Context, name, password string) (bool, error)
	Add(ctx context.Context, user *models.User) error
}

type Repository interface {
	AddUser(ctx context.Context, user *models.User) error
	GetUserByLogin(ctx context.Context, login string) (*models.User, error)
}
