package storage

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/xfiendx4life/gb_go_backend1/internal/config"
	"go.uber.org/zap"
)

type Storage interface {
	InitNewStorage(ctx context.Context, z *zap.SugaredLogger, config config.Storage) error
	GetDbPool() *pgxpool.Pool
}

type PG struct {
	dbPool *pgxpool.Pool
}
