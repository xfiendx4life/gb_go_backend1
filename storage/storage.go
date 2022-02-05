package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/xfiendx4life/gb_go_backend1/internal/config"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"go.uber.org/zap"
)

type Storage interface {
	InitNewStorage(ctx context.Context, z *zap.SugaredLogger, config config.Storage) error
	GetRedirects(ctx context.Context, urlId int, z *zap.SugaredLogger) ([]models.Redirects, error)
	GetDbPool() *pgxpool.Pool
}

type PG struct {
	dbPool *pgxpool.Pool
}
