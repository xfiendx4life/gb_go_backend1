//go:build unit
// +build unit

package usecase_test

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_backend1/internal/logger"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/user/usecase"
	"go.uber.org/zap/zapcore"
)

var (
	lgr = logger.InitLogger(zapcore.DebugLevel, "")
	ctx = context.Background()
)

type mockStorage struct {
	err error
}

func (mc *mockStorage) AddUser(ctx context.Context, user *models.User) error {
	user.Id = 1
	return mc.err
}

func (mc *mockStorage) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	s := md5.Sum([]byte("correctPassword"))
	return &models.User{
		Name:     "testname",
		Password: hex.EncodeToString(s[:]),
	}, mc.err
}

func TestValidateCorrect(t *testing.T) {
	uc := usecase.New(&mockStorage{}, lgr)
	res, err := uc.Validate(ctx, "testname", "correctPassword")
	assert.NoError(t, err)
	assert.True(t, res)
}

func TestValidateIncorrect(t *testing.T) {
	uc := usecase.New(&mockStorage{}, lgr)
	res, err := uc.Validate(ctx, "testname", "incorrectPassword")
	assert.NoError(t, err)
	assert.False(t, res)
}

func TestValidateError(t *testing.T) {
	st := &mockStorage{err: errors.New("some error")}
	uc := usecase.New(st, lgr)
	res, err := uc.Validate(ctx, "testname", "incorrectPassword")
	assert.Error(t, err)
	assert.False(t, res)
}

func TestAddUser(t *testing.T) {
	uc := usecase.New(&mockStorage{}, lgr)
	u := &models.User{Name: "testname", Password: "password", Email: "email"}
	err := uc.Add(ctx, u)
	assert.NoError(t, err)
	assert.Equal(t, models.User{
		Id:       1,
		Name:     "testname",
		Password: "5f4dcc3b5aa765d61d8327deb882cf99",
		Email:    "email"}, *u)
}

func TestAddUserError(t *testing.T) {
	st := &mockStorage{err: errors.New("some error")}
	uc := usecase.New(st, lgr)
	u := &models.User{Name: "testname", Password: "password", Email: "email"}
	err := uc.Add(ctx, u)
	assert.Error(t, err)
}
