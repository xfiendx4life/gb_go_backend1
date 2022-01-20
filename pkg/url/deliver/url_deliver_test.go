package deliver_test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_backend1/internal/logger"
	"github.com/xfiendx4life/gb_go_backend1/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/pkg/url/deliver"
	"github.com/xfiendx4life/gb_go_backend1/pkg/url/usecase"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type mockStorage struct {
	err error
}

func (mc *mockStorage) InitNewStorage(ctx context.Context, connection string, z *zap.SugaredLogger) error {
	return nil
}

func (mc *mockStorage) AddUser(ctx context.Context, user *models.User, z *zap.SugaredLogger) error {
	user.Id = 1
	return mc.err
}

func (mc *mockStorage) GetUserByLogin(ctx context.Context, login string, z *zap.SugaredLogger) (*models.User, error) {
	return &models.User{
		Name:     "testname",
		Password: "correctPassword",
	}, mc.err
}

func (mc *mockStorage) AddUrl(ctx context.Context, url *models.Url, z *zap.SugaredLogger) error {
	url.Id = 1
	return mc.err
}

func (mc *mockStorage) GetUrl(ctx context.Context, id int, z *zap.SugaredLogger) (*models.Url, error) {
	return &models.Url{
		Raw:       "RawTestUrl",
		Shortened: "shortenedTestUrl",
	}, mc.err
}

func (mc *mockStorage) GetUrls(ctx context.Context, userID int, z *zap.SugaredLogger) ([]models.Url, error) {
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

func (mc *mockStorage) GetUrlByShortened(ctx context.Context, shortened string, z *zap.SugaredLogger) (*models.Url, error) {
	return &models.Url{
		Raw:       "RawTestUrl",
		Shortened: "shortenedTestUrl",
	}, mc.err
}

func (mc *mockStorage) AddRedirect(ctx context.Context, r *models.Redirects, z *zap.SugaredLogger) error {
	return nil
}

func (mc *mockStorage) GetRedirects(ctx context.Context, urlId int, z *zap.SugaredLogger) ([]models.Redirects, error) {
	return []models.Redirects{}, nil
}

var (
	lgr = logger.InitLogger(zapcore.DebugLevel, "")
	// ctx = context.Background()
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
	uc := usecase.New(&mc)
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
	uc := usecase.New(&mc)
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
	uc := usecase.New(&mc)
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
	uc := usecase.New(mc)
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
	uc := usecase.New(mc)
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
