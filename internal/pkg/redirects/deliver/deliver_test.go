package deliver_test

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_backend1/internal/logger"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/redirects/deliver"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/redirects/usecase"
)

type mockStorage struct {
	err error
}

func (mc *mockStorage) AddRedirect(ctx context.Context, redirect *models.Redirects) error {
	return mc.err
}

func (mc *mockStorage) GetRedirects(ctx context.Context, shortened string) ([]models.Redirects, error) {
	if mc.err != nil {
		return nil, mc.err
	}
	return []models.Redirects{
		{Id: 1, UrlId: 1, Date: time.Now().Add(time.Duration(time.Hour * (-72)))},
		{Id: 2, UrlId: 1, Date: time.Now().Add(time.Duration(time.Hour * (-2)))},
		{Id: 3, UrlId: 1, Date: time.Now().Add(time.Duration(time.Hour * (-179)))},
		{Id: 4, UrlId: 1, Date: time.Now().Add(time.Duration(time.Hour * (-1)))},
	}, nil
}

var (
	lgr        = logger.InitLogger(-1, "")
	correctGet = models.Summary{
		Week:  3,
		Today: 2,
		Month: 4,
	}
)

func TestGet(t *testing.T) {
	req := httptest.NewRequest("GET", "/redirects", nil)
	resp := httptest.NewRecorder()
	ectx := echo.New().NewContext(req, resp)
	ectx.SetPath("/:shortened")
	ectx.SetParamNames("shortened")
	ectx.SetParamValues("testurl")
	c := usecase.New(&mockStorage{}, lgr)
	err := deliver.New(c, lgr).GetSummary(ectx)
	assert.NoError(t, err)
	res := models.Summary{}
	json.NewDecoder(strings.NewReader(resp.Body.String())).Decode(&res)
	assert.Equal(t, correctGet, res)

}

func TestGetError(t *testing.T) {
	req := httptest.NewRequest("GET", "/redirects", nil)
	resp := httptest.NewRecorder()
	ectx := echo.New().NewContext(req, resp)
	ectx.SetPath("/:shortened")
	ectx.SetParamNames("shortened")
	ectx.SetParamValues("testurl")
	c := usecase.New(&mockStorage{err: fmt.Errorf("testerror")}, lgr)
	err := deliver.New(c, lgr).GetSummary(ectx)
	assert.Error(t, err)
}
