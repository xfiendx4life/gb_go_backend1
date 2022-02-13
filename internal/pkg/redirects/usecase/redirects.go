package usecase

import (
	"context"
	"fmt"
	"time"

	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/redirects"
	"go.uber.org/zap"
)

type Uc struct {
	store redirects.Repository
	z     *zap.SugaredLogger
}

func New(store redirects.Repository, z *zap.SugaredLogger) redirects.UseCase {
	return &Uc{store: store, z: z}
}

func (uc *Uc) Add(ctx context.Context, urlId int) error {
	rdr := models.Redirects{
		UrlId: urlId,
		Date:  time.Now(),
	}
	err := uc.store.AddRedirect(ctx, &rdr)
	if err != nil {
		uc.z.Errorf("can't add stats of redirect: %s", err)
		return fmt.Errorf("can't add stats of redirect: %s", err)
	}
	uc.z.Infof("Redirect info added")
	return nil
}

func (uc *Uc) Get(ctx context.Context, shortened string) (models.Summary, error) {
	rdrs, err := uc.store.GetRedirects(ctx, shortened)
	if err != nil {
		uc.z.Errorf("can't get redirects: %s", err)
		return models.Summary{}, fmt.Errorf("can't get redirects: %s", err)
	}
	now := time.Now()
	summary := models.Summary{}
	for _, rdr := range rdrs {
		diffInHours := now.Sub(rdr.Date).Hours()
		if diffInHours < 730 {
			summary.Month++
		}
		if diffInHours < 168 {
			summary.Week++
		}
		if diffInHours < 24 {
			summary.Today++
		}
	}
	return summary, nil
}
