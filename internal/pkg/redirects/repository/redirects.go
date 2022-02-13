package repository

import (
	"context"
	"fmt"

	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/redirects"
	"github.com/xfiendx4life/gb_go_backend1/storage"
	"go.uber.org/zap"
)

type repo struct {
	store storage.Storage
	z     *zap.SugaredLogger
}

func New(store storage.Storage, z *zap.SugaredLogger) redirects.Repository {
	return &repo{store: store, z: z}
}

func (r *repo) AddRedirect(ctx context.Context, redirect *models.Redirects) error {
	select {
	case <-ctx.Done():
		r.z.Error("done with context")
		return fmt.Errorf("done with context")
	default:
		q := `INSERT INTO redirects (url_id, date_of_usage) VALUES ($1, $2) RETURNING id`
		err := r.store.GetDbPool().QueryRow(ctx, q, redirect.UrlId, redirect.Date).Scan(&redirect.Id)
		if err != nil {
			r.z.Errorf("can't add to database: %s", err)
			return fmt.Errorf("can't add to database: %s", err)
		}
		return nil
	}
}

func (r *repo) GetRedirects(ctx context.Context, shortened string) ([]models.Redirects, error) {
	select {
	case <-ctx.Done():
		r.z.Error("done with context")
		return nil, fmt.Errorf("done with context")
	default:
		q := `SELECT * FROM redirects WHERE url_id = (SELECT id FROM urls WHERE shortened=$1);`
		row, err := r.store.GetDbPool().Query(ctx, q, shortened)
		if err != nil {
			r.z.Error("can't get rows %s", err)
			return nil, fmt.Errorf("can't get rows %s", err)
		}
		defer row.Close()
		res := make([]models.Redirects, 0)
		for row.Next() {
			red := models.Redirects{}
			err = row.Scan(&red.Id, &red.UrlId, &red.Date)
			if err != nil {
				r.z.Error("can't scan row %s", err)
				return nil, fmt.Errorf("can't scan row %s", err)
			}
			r.z.Infof("Got redirects: %v", red)
			res = append(res, red)
		}
		return res, nil
	}
}
