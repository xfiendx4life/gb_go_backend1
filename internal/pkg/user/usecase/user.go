package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/user"
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
	if u.Password != hashPassword(password) {
		z.Info("No such user or wrong password")
		return false, nil
	}
	return true, nil
}

func hashPassword(rawPass string) string {
	s := md5.Sum([]byte(rawPass))
	return hex.EncodeToString(s[:])
}

func (g *gres) Add(ctx context.Context, u *models.User, z *zap.SugaredLogger) error {
	z.Infof("created new user %v", *u)
	u.Password = hashPassword(u.Password)
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
