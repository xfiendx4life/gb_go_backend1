package usecase

import (
	"context"
	"fmt"

	"github.com/xfiendx4life/gb_go_backend1/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/pkg/user"
	"go.uber.org/zap"

	"github.com/xfiendx4life/gb_go_backend1/storage"
)

type gres struct {
	repo storage.Storage
}

func (g *gres) Validate(ctx context.Context, name, password string, z *zap.SugaredLogger) (bool, error) {
	u, err := g.repo.GetUserByLogin(ctx, name, z)
	if err != nil {
		z.Errorf("can't validate user: %s", err)
		return false, fmt.Errorf("can't validate user: %s", err)
	}
	if u.Password != password {
		z.Info("No such user or wrong password")
		return false, nil
	}
	return true, nil
}

func (g *gres) Add(ctx context.Context, u *models.User, z *zap.SugaredLogger) error {
	z.Infof("created new user %v", *u)
	err := g.repo.AddUser(ctx, u, z)
	if err != nil {
		z.Errorf("can't add new user to storage: %s", err)
		return fmt.Errorf("can't add new user to storage: %s", err)
	}
	z.Infof("user %v added to storage", u)
	return nil
}

func New(repo storage.Storage) user.UseCase {
	return &gres{repo: repo}
}
