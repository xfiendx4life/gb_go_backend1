//go:build unit
// +build unit

package deliver_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_backend1/internal/logger"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	rdrUse "github.com/xfiendx4life/gb_go_backend1/internal/pkg/redirects/usecase"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/url/deliver"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/url/usecase"
	"github.com/xfiendx4life/gb_go_backend1/storage"
	"go.uber.org/zap/zapcore"
)

type mockStorage struct {
	err error
}

func (mc *mockStorage) AddUrl(ctx context.Context, url *models.Url) error {
	url.Id = 1
	return mc.err
}

func (mc *mockStorage) GetUrl(ctx context.Context, id int) (*models.Url, error) {
	return &models.Url{
		Raw:       "RawTestUrl",
		Shortened: "shortenedTestUrl",
	}, mc.err
}

func (mc *mockStorage) GetUrls(ctx context.Context, userID int) ([]models.Url, error) {
	return []models.Url{
		{
			UserId:    1,
			Raw:       "RawTestUrl",
			Shortened: "shortenedTestUrl",
		},
		{
			UserId:    1,
			Raw:       "RawTestUrl1",
			Shortened: "shortenedTestUrl1",
		},
		{
			UserId:    1,
			Raw:       "RawTestUrl2",
			Shortened: "shortenedTestUrl2",
		},
	}, nil
}

func (m *mockStorage) GetStorage() storage.Storage {
	return storage.New() // ! shit happens here because we shouldn't return real storage
}

func (mc *mockStorage) GetUrlByShortened(ctx context.Context, shortened string) (*models.Url, error) {
	return &models.Url{
		Raw:       "RawTestUrl",
		Shortened: "shortenedTestUrl",
	}, mc.err
}

type mockRdrStorage struct {
	err error
}

func (mr *mockRdrStorage) AddRedirect(ctx context.Context, redirect *models.Redirects) error {
	return mr.err
}

func (mr *mockRdrStorage) GetRedirects(ctx context.Context, shortened string) ([]models.Redirects, error) {
	if mr.err != nil {
		return nil, mr.err
	}
	return []models.Redirects{
		{Id: 1, UrlId: 1, Date: time.Now().Add(time.Duration(time.Hour * (-72)))},
		{Id: 2, UrlId: 1, Date: time.Now().Add(time.Duration(time.Hour * (-2)))},
		{Id: 3, UrlId: 1, Date: time.Now().Add(time.Duration(time.Hour * (-179)))},
		{Id: 4, UrlId: 1, Date: time.Now().Add(time.Duration(time.Hour * (-1)))},
	}, nil
}

var (
	lgr = logger.InitLogger(zapcore.DebugLevel, "")
	// ctx = context.Background()
	rdrUc = rdrUse.New(&mockRdrStorage{}, lgr)
)

func TestAddUrl(t *testing.T) {
	data := `{
		"rawurl": "someVeryLongEtsUrl",    
		"userid": 1
	  }`
	req := httptest.NewRequest("POST", "/url", strings.NewReader(data))
	resp := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	mc := mockStorage{}
	uc := usecase.New(&mc, rdrUc, lgr)
	del := deliver.New(uc, lgr)
	c := echo.New().NewContext(req, resp)
	err := del.Save(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestAddUrlErr(t *testing.T) {
	data := `{
		"rawurl": "someVeryLongEtsUrl",    
		"userid": 1
	  }`
	req := httptest.NewRequest("POST", "/url", strings.NewReader(data))
	resp := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	mc := mockStorage{err: fmt.Errorf("Some error")}
	uc := usecase.New(&mc, rdrUc, lgr)
	del := deliver.New(uc, lgr)
	c := echo.New().NewContext(req, resp)
	err := del.Save(c)
	assert.Error(t, err)
	assert.Equal(t, echo.ErrBadRequest, err)
}

func TestAddUrlJsonErr(t *testing.T) {
	data := `lgf`
	req := httptest.NewRequest("POST", "/url", strings.NewReader(data))
	resp := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	mc := mockStorage{err: fmt.Errorf("Some error")}
	uc := usecase.New(&mc, rdrUc, lgr)
	del := deliver.New(uc, lgr)
	c := echo.New().NewContext(req, resp)
	err := del.Save(c)
	assert.Error(t, err)
	assert.Equal(t, echo.ErrBadRequest, err)
}

func TestGetUrl(t *testing.T) {
	req := httptest.NewRequest("GET", "/", nil)
	resp := httptest.NewRecorder()
	ectx := echo.New().NewContext(req, resp)
	ectx.SetPath("/:shortened")
	ectx.SetParamNames("shortened")
	ectx.SetParamValues("shortenedTestUrl")
	mc := &mockStorage{}
	uc := usecase.New(mc, rdrUc, lgr)
	del := deliver.New(uc, lgr)
	raw, err := del.Get(ectx)
	assert.NoError(t, err)
	assert.Equal(t, "RawTestUrl", raw)
}

func TestGetList(t *testing.T) {
	q := make(url.Values)
	q.Set("id", "1")
	req := httptest.NewRequest("GET", "/urls?"+q.Encode(), nil)
	resp := httptest.NewRecorder()
	ectx := echo.New().NewContext(req, resp)
	mc := &mockStorage{}
	uc := usecase.New(mc, rdrUc, lgr)
	del := deliver.New(uc, lgr)
	urls, err := del.List(ectx)
	assert.NoError(t, err)
	targetRes := []models.Url{
		{
			UserId:    1,
			Raw:       "RawTestUrl",
			Shortened: "shortenedTestUrl",
		},
		{
			UserId:    1,
			Raw:       "RawTestUrl1",
			Shortened: "shortenedTestUrl1",
		},
		{
			UserId:    1,
			Raw:       "RawTestUrl2",
			Shortened: "shortenedTestUrl2",
		},
	}
	assert.Equal(t, targetRes, urls)
}
