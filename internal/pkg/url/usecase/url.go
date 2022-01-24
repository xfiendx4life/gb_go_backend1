package usecase

import (
	"context"
	"time"

	crypto_rand "crypto/rand"
	"encoding/binary"
	"fmt"
	math_rand "math/rand"

	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/url"
	"github.com/xfiendx4life/gb_go_backend1/storage"
	"go.uber.org/zap"
)

type gres struct {
	store storage.Storage
}

func getSeedNumber() (int64, error) {
	var b [8]byte
	_, err := crypto_rand.Read(b[:])
	if err != nil {
		return 0, fmt.Errorf("cannot seed math/rand package with cryptographically secure random number generator")
	}
	return int64(binary.LittleEndian.Uint64(b[:])), nil
}

// takes []byte of raw url and then performes random
// permutation and takes 15 or less symbols for short link
func NewUrl(raw string, userId int, z *zap.SugaredLogger) *models.Url {
	bts := []byte(raw)
	n, err := getSeedNumber()
	if err != nil {
		z.Error(err)
	}
	math_rand.Seed(n)
	var i int
	for i = 0; i < 15 && i < len(bts); i++ {
		toChange := math_rand.Intn(len(bts) - i - 1 + i)
		bts[i], bts[toChange] = bts[toChange], bts[i]
		bts[i] -= byte(math_rand.Intn(10))
	}
	return &models.Url{
		Raw:       raw,
		Shortened: string(bts[:i]),
		UserId:    userId,
	}
}

func New(st storage.Storage) url.UseCase {
	return &gres{store: st}
}

func (g *gres) Add(ctx context.Context, raw string, userId int, z *zap.SugaredLogger) (shortened string, err error) {
	url := NewUrl(raw, userId, z)
	err = g.store.AddUrl(ctx, url, z)
	if err != nil {
		z.Errorf("can't add url to storage: %s", err)
		return "", fmt.Errorf("can't add url to storage: %s", err)
	}
	z.Infof("Url %s was successfully added to storage and shortened to -> %s", url.Raw, url.Shortened)
	return url.Shortened, nil
}

// TODO: test this
// ? Maybe move it to Redirects usecase
func (g *gres) AddStats(ctx context.Context, urlID int, z *zap.SugaredLogger) error {
	rdr := models.Redirects{
		UrlId: urlID,
		Date:  time.Now(),
	}
	err := g.store.AddRedirect(ctx, &rdr, z)
	if err != nil {
		z.Errorf("can't add stats of redirect: %s", err)
		return fmt.Errorf("can't add stats of redirect: %s", err)
	}
	return nil
}

func (g *gres) Get(ctx context.Context, shortened string, z *zap.SugaredLogger) (raw string, err error) {
	url, err := g.store.GetUrlByShortened(ctx, shortened, z)
	if err != nil {
		return "", echo.ErrBadRequest
	}
	err = g.AddStats(ctx, url.Id, z)
	if err != nil {
		z.Errorf("can't add stats to storage %s", err)
		return "", echo.ErrInternalServerError
	}
	return url.Raw, nil
}

func (g *gres) List(ctx context.Context, userId int, z *zap.SugaredLogger) ([]models.Url, error) {
	urls, err := g.store.GetUrls(ctx, userId, z)
	if err != nil {
		z.Errorf("can't get urls for user_id %d, error: %s", userId, err)
		return nil, fmt.Errorf("can't get urls for user_id %d, error: %s", userId, err)
	}

	return urls, nil
}
