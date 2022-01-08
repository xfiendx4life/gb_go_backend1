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
	if u.Password == password {
		return false, nil
	}
	return true, nil
}

func (g *gres) Add(ctx context.Context, name, password, email string, z *zap.SugaredLogger) (*models.User, error) {
	return &models.User{}, nil
}

func New(repo storage.Storage) user.UseCase {
	return &gres{repo: repo}
}
