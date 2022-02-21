package repository

import (
	"context"
	"fmt"

	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/url"
	"github.com/xfiendx4life/gb_go_backend1/storage"
	"go.uber.org/zap"
)

type PG struct {
	store storage.Storage
	z     *zap.SugaredLogger
}

func New(store storage.Storage, z *zap.SugaredLogger) url.Repository {
	return &PG{store: store, z: z}
}

func (pg *PG) AddUrl(ctx context.Context, url *models.Url) error {
	select {
	case <-ctx.Done():
		pg.z.Error("done with context")
		return fmt.Errorf("done with context")
	default:
		q := `INSERT INTO urls (raw, shortened, user_id) VALUES ($1, $2, $3) RETURNING id`
		row := pg.store.GetDbPool().QueryRow(ctx, q, url.Raw, url.Shortened, url.UserId)
		err := row.Scan(&url.Id)
		pg.z.Infof("Added url with id: %d", url.Id)
		if err != nil {
			pg.z.Errorf("error while inserting url to db: %s", err)
			return fmt.Errorf("error while inserting url to db: %s", err)
		}
		return nil
	}
}

func (pg *PG) GetUrls(ctx context.Context, userID int) ([]models.Url, error) {
	select {
	case <-ctx.Done():
		pg.z.Error("done with context")
		return nil, fmt.Errorf("done with context")
	default:
		q := `SELECT id, raw, shortened, user_id FROM urls WHERE user_id = $1`
		rows, err := pg.store.GetDbPool().Query(ctx, q, userID)
		if err != nil {
			pg.z.Errorf("can't get urls: %s", err)
			return nil, fmt.Errorf("can't get urls: %s", err)
		}
		defer rows.Close()
		urls := make([]models.Url, 0)
		for rows.Next() {
			url := models.Url{}
			err := rows.Scan(&url.Id, &url.Raw, &url.Shortened, &url.UserId)
			//	// &url.RedirectsNum.Month, &url.RedirectsNum.Week, &url.RedirectsNum.Today)
			if err != nil {
				pg.z.Errorf("can't parse urls: %s", err)
				return nil, fmt.Errorf("can't parse urls: %s", err)
			}
			urls = append(urls, url)
		}
		return urls, nil
	}
}

func (pg *PG) GetUrlByShortened(ctx context.Context, shortened string) (*models.Url, error) {
	select {
	case <-ctx.Done():
		pg.z.Error("done with context")
		return nil, fmt.Errorf("done with context")
	default:
		q := `SELECT id, raw, shortened, user_id FROM urls WHERE shortened=$1;`
		u := models.Url{}
		err := pg.store.GetDbPool().QueryRow(ctx, q, shortened).Scan(&u.Id, &u.Raw, &u.Shortened, &u.UserId) ////, &u.RedirectsNum.Month, &u.RedirectsNum.Week, &u.RedirectsNum.Today)
		if err != nil {
			pg.z.Errorf("can't get url by shortened: %s", err)
			return nil, fmt.Errorf("can't get url by shortened: %s", err)
		}
		return &u, nil
	}
}

func (pg *PG) GetStorage() storage.Storage {
	return pg.store
}
