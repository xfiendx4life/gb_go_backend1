package usecase

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"

	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/user"
	"go.uber.org/zap"
)

type gres struct {
	repo user.Repository
	z    *zap.SugaredLogger
}

func (g *gres) Validate(ctx context.Context, name, password string) (bool, error) {
	u, err := g.repo.GetUserByLogin(ctx, name)
	if err != nil {
		g.z.Errorf("can't validate user: %s", err)
		return false, fmt.Errorf("can't validate user: %s", err)
	}
	if u.Password != hashPassword(password) {
		g.z.Info("No such user or wrong password")
		return false, nil
	}
	return true, nil
}

func hashPassword(rawPass string) string {
	s := md5.Sum([]byte(rawPass))
	return hex.EncodeToString(s[:])
}

func (g *gres) Add(ctx context.Context, u *models.User) error {
	g.z.Infof("created new user %v", *u)
	u.Password = hashPassword(u.Password)
	err := g.repo.AddUser(ctx, u)
	if err != nil {
		g.z.Errorf("can't add new user to storage: %s", err)
		return fmt.Errorf("can't add new user to storage: %s", err)
	}
	g.z.Infof("user %v added to storage", u)
	return nil
}

func (g *gres) Get(ctx context.Context, name string) (id int, err error) {
	u, err := g.repo.GetUserByLogin(ctx, name)
	if err != nil {
		g.z.Errorf("can't get user with name %s, err: %s", name, err)
		return 0, fmt.Errorf("can't get user with name %s, err: %s", name, err)
	}
	return u.Id, nil
}

func New(repo user.Repository, z *zap.SugaredLogger) user.UseCase {
	return &gres{repo: repo, z: z}
}
