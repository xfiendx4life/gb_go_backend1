package storage

import (
	"context"
	"errors"
	"fmt"
	"net"
	"time"

	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/xfiendx4life/gb_go_backend1/internal/config"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"go.uber.org/zap"
)

func New() Storage {
	return &PG{}
}

// TODO: change this function go get all configuration from config object
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

func (pg *PG) AddUser(ctx context.Context, user *models.User, z *zap.SugaredLogger) error {
	select {
	case <-ctx.Done():
		z.Error("done with context")
		return fmt.Errorf("done with context")
	default:
		q := `INSERT INTO users (name, password, email) VALUES ($1, $2, $3) RETURNING id`
		row := pg.dbPool.QueryRow(ctx, q, user.Name, user.Password, user.Email)
		err := row.Scan(&user.Id)
		z.Info(user.Id)
		if err != nil {
			z.Errorf("error while inserting to db: %s", err)
			return fmt.Errorf("error while inserting to db: %s", err)
		}
		return nil
	}

}

func (pg *PG) GetUserByLogin(ctx context.Context, login string, z *zap.SugaredLogger) (*models.User, error) {
	select {
	case <-ctx.Done():
		z.Error("done with context")
		return nil, fmt.Errorf("done with context")
	default:
		q := `SELECT * FROM users WHERE name=$1;`
		u := models.User{}
		err := pg.dbPool.QueryRow(ctx, q, login).Scan(&u.Id, &u.Name, &u.Password, &u.Email)
		if err != nil {
			z.Errorf("can't get user by name: %s", err)
			if errors.Is(err, pgx.ErrNoRows) {
				return &models.User{}, nil
			}
			return nil, fmt.Errorf("can't get user by name: %s", err)
		}
		return &u, nil
	}
}

func (pg *PG) AddUrl(ctx context.Context, url *models.Url, z *zap.SugaredLogger) error {
	select {
	case <-ctx.Done():
		z.Error("done with context")
		return fmt.Errorf("done with context")
	default:
		q := `INSERT INTO urls (raw, shortened, user_id) VALUES ($1, $2, $3) RETURNING id`
		row := pg.dbPool.QueryRow(ctx, q, url.Raw, url.Shortened, url.UserId)
		err := row.Scan(&url.Id)
		z.Infof("Added url with id: %d", url.Id)
		if err != nil {
			z.Errorf("error while inserting url to db: %s", err)
			return fmt.Errorf("error while inserting url to db: %s", err)
		}
		return nil
	}
}

func (pg *PG) GetUrls(ctx context.Context, userID int, z *zap.SugaredLogger) ([]models.Url, error) {
	select {
	case <-ctx.Done():
		z.Error("done with context")
		return nil, fmt.Errorf("done with context")
	default:
		q := `SELECT id, raw, shortened, user_id FROM urls WHERE user_id = $1`
		rows, err := pg.dbPool.Query(ctx, q, userID)
		if err != nil {
			z.Errorf("can't get urls: %s", err)
			return nil, fmt.Errorf("can't get urls: %s", err)
		}
		defer rows.Close()
		urls := make([]models.Url, 0)
		for rows.Next() {
			url := models.Url{}
			err := rows.Scan(&url.Id, &url.Raw, &url.Shortened, &url.UserId)
			//	// &url.RedirectsNum.Month, &url.RedirectsNum.Week, &url.RedirectsNum.Today)
			if err != nil {
				z.Errorf("can't parse urls: %s", err)
				return nil, fmt.Errorf("can't parse urls: %s", err)
			}
			urls = append(urls, url)
		}
		return urls, nil
	}
}

func (pg *PG) GetUrlByShortened(ctx context.Context, shortened string, z *zap.SugaredLogger) (*models.Url, error) {
	select {
	case <-ctx.Done():
		z.Error("done with context")
		return nil, fmt.Errorf("done with context")
	default:
		q := `SELECT id, raw, shortened, user_id FROM urls WHERE shortened=$1;`
		u := models.Url{}
		err := pg.dbPool.QueryRow(ctx, q, shortened).Scan(&u.Id, &u.Raw, &u.Shortened, &u.UserId) ////, &u.RedirectsNum.Month, &u.RedirectsNum.Week, &u.RedirectsNum.Today)
		if err != nil {
			z.Errorf("can't get url by shortened: %s", err)
			return nil, fmt.Errorf("can't get url by shortened: %s", err)
		}
		return &u, nil
	}
}

func (pg *PG) AddRedirect(ctx context.Context, redirect *models.Redirects, z *zap.SugaredLogger) error {
	select {
	case <-ctx.Done():
		z.Error("done with context")
		return fmt.Errorf("done with context")
	default:
		q := `INSERT INTO redirects (url_id, date_of_usage) VALUES ($1, $2) RETURNING id`
		err := pg.dbPool.QueryRow(ctx, q, redirect.UrlId, redirect.Date).Scan(&redirect.Id)
		if err != nil {
			z.Errorf("can't add to database: %s", err)
			return fmt.Errorf("can't add to database: %s", err)
		}
		return nil
	}
}

func (pg *PG) GetRedirects(ctx context.Context, urlId int, z *zap.SugaredLogger) ([]models.Redirects, error) {
	select {
	case <-ctx.Done():
		z.Error("done with context")
		return nil, fmt.Errorf("done with context")
	default:
		q := `SELECT * FROM redirects WHERE url_id = $1`
		row, err := pg.dbPool.Query(ctx, q, urlId)
		if err != nil {
			z.Error("can't get rows %s", err)
			return nil, fmt.Errorf("can't get rows %s", err)
		}
		defer row.Close()
		res := make([]models.Redirects, 0)
		for row.Next() {
			red := models.Redirects{}
			err = row.Scan(&red.Id, &red.UrlId, &red.Date)
			if err != nil {
				z.Error("can't scan row %s", err)
				return nil, fmt.Errorf("can't scan row %s", err)
			}
			z.Infof("Got redirects: %v", red)
			res = append(res, red)
		}
		return res, nil
	}
}
