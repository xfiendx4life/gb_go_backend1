package deliver_test

import (
	"context"
	"crypto/md5"
	"encoding/hex"
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
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/user/deliver"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/user/usecase"
	"go.uber.org/zap/zapcore"
)

type mockStorage struct {
	err error
}

func (mc *mockStorage) AddUser(ctx context.Context, user *models.User) error {
	user.Id = 1
	return mc.err
}

func (mc *mockStorage) GetUserByLogin(ctx context.Context, login string) (*models.User, error) {
	s := md5.Sum([]byte("correctpassword"))
	password := hex.EncodeToString(s[:])
	return &models.User{
		Name:     "testname",
		Password: password,
	}, mc.err
}

var (
	lgr = logger.InitLogger(zapcore.DebugLevel, "")
	ttl = time.Now().Add(time.Hour).Unix()
	// ctx = context.Background()
)

func TestCreate(t *testing.T) {
	data := `{
"name": "punk",
"password": "123",
"email": "punk@mail.ru"
}`
	req := httptest.NewRequest("POST", "/user/create", strings.NewReader(data))
	resp := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	mc := mockStorage{}
	uc := usecase.New(&mc, lgr)
	del := deliver.New(uc, ttl, "secret", lgr)
	c := echo.New().NewContext(req, resp)
	err := del.Create(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.Code)
}

func TestCreateError(t *testing.T) {
	data := ""
	req := httptest.NewRequest("POST", "/user/create", strings.NewReader(data))
	resp := httptest.NewRecorder()
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	mc := mockStorage{}
	uc := usecase.New(&mc, lgr)
	del := deliver.New(uc, ttl, "secret", lgr)
	c := echo.New().NewContext(req, resp)
	err := del.Create(c)
	lgr.Error(err)
	assert.Error(t, err)
}

func TestLogin(t *testing.T) {
	q := make(url.Values)
	q.Set("name", "testname")
	q.Set("password", "correctpassword")
	req := httptest.NewRequest("GET", "/user/login?"+q.Encode(), nil)
	resp := httptest.NewRecorder()
	mc := mockStorage{}
	uc := usecase.New(&mc, lgr)
	del := deliver.New(uc, ttl, "secret", lgr)
	c := echo.New().NewContext(req, resp)
	err := del.Login(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
}

func TestGet(t *testing.T) {
	q := make(url.Values)
	q.Set("name", "testname")
	req := httptest.NewRequest("GET", "/user?"+q.Encode(), nil)
	resp := httptest.NewRecorder()
	mc := mockStorage{}
	uc := usecase.New(&mc, lgr)
	del := deliver.New(uc, ttl, "secret", lgr)
	c := echo.New().NewContext(req, resp)
	err := del.Get(c)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.Code)
}
