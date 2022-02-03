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
	// AddUser(ctx context.Context, user *models.User, z *zap.SugaredLogger) error
	// GetUserByLogin(ctx context.Context, login string, z *zap.SugaredLogger) (*models.User, error)
	// AddUrl(ctx context.Context, url *models.Url, z *zap.SugaredLogger) error
	// GetUrls(ctx context.Context, userID int, z *zap.SugaredLogger) ([]models.Url, error)
	// GetUrlByShortened(ctx context.Context, shortened string, z *zap.SugaredLogger) (*models.Url, error)
	// AddRedirect(ctx context.Context, redirect *models.Redirects, z *zap.SugaredLogger) error
	GetRedirects(ctx context.Context, urlId int, z *zap.SugaredLogger) ([]models.Redirects, error)
	GetDbPool() *pgxpool.Pool
}

type PG struct {
	dbPool *pgxpool.Pool
}
