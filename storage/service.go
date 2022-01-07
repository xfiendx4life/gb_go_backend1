package storage

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/xfiendx4life/gb_go_backend1/pkg/models"
	"go.uber.org/zap"
)

func New() Storage {
	return &PG{}
}

// TODO: change this function go get all configuration from config object
func configurePool(conf *pgxpool.Config, z *zap.SugaredLogger) (err error) {
	// add cofiguration
	Maxc, err := strconv.Atoi(os.Getenv("MAX_CONS"))
	if err != nil {
		z.Errorf("wrong format of env var: %s", err)
		return fmt.Errorf("wrong format of env var %s", err)
	}
	Minc, err := strconv.Atoi(os.Getenv("MIN_CONS"))
	if err != nil {
		z.Errorf("wrong format of env var: %s", err)
		return fmt.Errorf("wrong format of env var %s", err)
	}
	conf.MaxConns = int32(Maxc) // 10
	conf.MinConns = int32(Minc) // 5

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

func (pg *PG) InitNewStorage(ctx context.Context, connection string, z *zap.SugaredLogger) error {
	conf, err := pgxpool.ParseConfig(connection)
	if err != nil {
		z.Errorf("can't init storage: %s", err)
		return fmt.Errorf("can't init storage: %s", err)
	}
	err = configurePool(conf, z)
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

// TODO: Test it!
func (pg *PG) GetUrl(ctx context.Context, userId int, z *zap.SugaredLogger) (*models.Url, error) {
	select {
	case <-ctx.Done():
		z.Error("done with context")
		return nil, fmt.Errorf("done with context")
	default:
		q := `SELECT * FROM urls WHERE user_id=$1;`
		u := models.Url{}
		err := pg.dbPool.QueryRow(ctx, q, userId).Scan(&u.Id, &u.Raw, &u.Shortened, &u.UserId, &u.RedirectsNum.Month, &u.RedirectsNum.Week, &u.RedirectsNum.Today)
		if err != nil {
			z.Errorf("can't get url by id: %s", err)
			return nil, fmt.Errorf("can't get url by id: %s", err)
		}
		return &u, nil
	}
}

func (pg *PG) GetUrls(ctx context.Context, userID int, z *zap.SugaredLogger) ([]models.Url, error) {
	select {
	case <-ctx.Done():
		z.Error("done with context")
		return nil, fmt.Errorf("done with context")
	default:
		q := `SELECT * FROM urls WHERE user_id = $1`
		rows, err := pg.dbPool.Query(ctx, q, userID)
		if err != nil {
			z.Errorf("can't get urls: %s", err)
			return nil, fmt.Errorf("can't get urls: %s", err)
		}
		defer rows.Close()
		urls := make([]models.Url, 0)
		for rows.Next() {
			url := models.Url{}
			err := rows.Scan(&url.Id, &url.Raw, &url.Shortened, &url.UserId,
				&url.RedirectsNum.Month, &url.RedirectsNum.Week, &url.RedirectsNum.Today)
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
		q := `SELECT * FROM urls WHERE shortened=$1;`
		u := models.Url{}
		err := pg.dbPool.QueryRow(ctx, q, shortened).Scan(&u.Id, &u.Raw, &u.Shortened, &u.UserId, &u.RedirectsNum.Month, &u.RedirectsNum.Week, &u.RedirectsNum.Today)
		if err != nil {
			z.Errorf("can't get url by shortened: %s", err)
			return nil, fmt.Errorf("can't get url by shortened: %s", err)
		}
		return &u, nil
	}
}
