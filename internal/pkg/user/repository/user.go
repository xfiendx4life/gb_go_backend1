package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v4"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/user"
	"github.com/xfiendx4life/gb_go_backend1/storage"
	"go.uber.org/zap"
)

type PG struct {
	store storage.Storage
	z     *zap.SugaredLogger
}

func New(store storage.Storage, z *zap.SugaredLogger) user.Repository {
	return &PG{store: store, z: z}
}

func (pg *PG) AddUser(ctx context.Context, user *models.User) error {
	select {
	case <-ctx.Done():
		pg.z.Error("done with context")
		return fmt.Errorf("done with context")
	default:
		q := `INSERT INTO users (name, password, email) VALUES ($1, $2, $3) RETURNING id`
		row := pg.store.GetDbPool().QueryRow(ctx, q, user.Name, user.Password, user.Email)
		err := row.Scan(&user.Id)
		pg.z.Info(user.Id)
		if err != nil {
			pg.z.Errorf("error while inserting to db: %s", err)
			return fmt.Errorf("error while inserting to db: %s", err)
		}
		return nil
	}

}

func (pg *PG) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	select {
	case <-ctx.Done():
		pg.z.Error("done with context")
		return nil, fmt.Errorf("done with context")
	default:
		q := `SELECT * FROM users WHERE name=$1;`
		u := models.User{}
		err := pg.store.GetDbPool().QueryRow(ctx, q, login).Scan(&u.Id, &u.Name, &u.Password, &u.Email)
		if err != nil {
			pg.z.Errorf("can't get user by name: %s", err)
			if errors.Is(err, pgx.ErrNoRows) {
				return &models.User{}, nil
			}
			return nil, fmt.Errorf("can't get user by name: %s", err)
		}
		return &u, nil
	}
}
