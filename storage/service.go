package storage

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/xfiendx4life/gb_go_backend1/internal/config"
	"go.uber.org/zap"
)

func New() Storage {
	return &PG{}
}

func configurePool(conf *pgxpool.Config, z *zap.SugaredLogger, config config.Storage) (err error) {
	// add cofiguration
	conf.MaxConns = int32(config.GetMaxCons())
	conf.MinConns = int32(config.GetMinCons())

	conf.HealthCheckPeriod = 1 * time.Minute
	conf.MaxConnLifetime = 24 * time.Hour
	conf.MaxConnIdleTime = 30 * time.Minute
	conf.ConnConfig.ConnectTimeout = 1 * time.Second
	conf.ConnConfig.DialFunc = (&net.Dialer{
		KeepAlive: conf.HealthCheckPeriod,
		Timeout:   conf.ConnConfig.ConnectTimeout,
	}).DialContext
	return nil
}

func (pg *PG) InitNewStorage(ctx context.Context, z *zap.SugaredLogger, config config.Storage) error {

	conf, err := pgxpool.ParseConfig(config.GetURI())
	if err != nil {
		z.Errorf("can't init storage: %s", err)
		return fmt.Errorf("can't init storage: %s", err)
	}
	err = configurePool(conf, z, config)
	if err != nil {
		z.Errorf("can't configure pool %s", err)
		return fmt.Errorf("can't configure pool %s", err)
	}

	dbPool, err := pgxpool.ConnectConfig(ctx, conf)
	if err != nil {
		z.Errorf("can't create pool %s", err)
		return fmt.Errorf("can't create pool %s", err)
	}
	pg.dbPool = dbPool
	return nil
}

func (pg *PG) GetDbPool() *pgxpool.Pool {
	return pg.dbPool
}
