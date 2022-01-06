package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/xfiendx4life/gb_go_backend1/pkg/models"
	"go.uber.org/zap"
)

type Storage interface {
	InitNewStorage(ctx context.Context, connection string, z *zap.SugaredLogger) error
	AddUser(ctx context.Context, user *models.User, z *zap.SugaredLogger) error
	GetUserByLogin(ctx context.Context, login string, z *zap.SugaredLogger) (*models.User, error)
	AddUrl(ctx context.Context, url *models.Url, z *zap.SugaredLogger) error
	GetUrl(ctx context.Context, id int, z *zap.SugaredLogger) (*models.Url, error)
	GetUrls(ctx context.Context, userID int, z *zap.SugaredLogger) ([]models.Url, error)
}

type PG struct {
	dbPool *pgxpool.Pool
}
