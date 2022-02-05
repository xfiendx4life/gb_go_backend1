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
	"go.uber.org/zap"
)

type gres struct {
	store url.Repository
	z     *zap.SugaredLogger
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

func New(st url.Repository, z *zap.SugaredLogger) url.UseCase {
	return &gres{store: st, z: z}
}

func (g *gres) Add(ctx context.Context, raw string, userId int) (shortened string, err error) {
	url := NewUrl(raw, userId, g.z)
	err = g.store.AddUrl(ctx, url)
	if err != nil {
		g.z.Errorf("can't add url to storage: %s", err)
		return "", fmt.Errorf("can't add url to storage: %s", err)
	}
	g.z.Infof("Url %s was successfully added to storage and shortened to -> %s", url.Raw, url.Shortened)
	return url.Shortened, nil
}

// TODO: test this
// ? Maybe move it to Redirects usecase
func (g *gres) AddStats(ctx context.Context, urlID int) error {
	rdr := models.Redirects{
		UrlId: urlID,
		Date:  time.Now(),
	}
	err := g.store.AddRedirect(ctx, &rdr)
	if err != nil {
		g.z.Errorf("can't add stats of redirect: %s", err)
		return fmt.Errorf("can't add stats of redirect: %s", err)
	}
	return nil
}

func (g *gres) Get(ctx context.Context, shortened string) (raw string, err error) {
	url, err := g.store.GetUrlByShortened(ctx, shortened)
	if err != nil {
		return "", echo.ErrBadRequest
	}
	err = g.AddStats(ctx, url.Id)
	if err != nil {
		g.z.Errorf("can't add stats to storage %s", err)
		return "", echo.ErrInternalServerError
	}
	return url.Raw, nil
}

func (g *gres) List(ctx context.Context, userId int) ([]models.Url, error) {
	urls, err := g.store.GetUrls(ctx, userId)
	if err != nil {
		g.z.Errorf("can't get urls for user_id %d, error: %s", userId, err)
		return nil, fmt.Errorf("can't get urls for user_id %d, error: %s", userId, err)
	}

	return urls, nil
}
