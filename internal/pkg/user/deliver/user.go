package deliver

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/user"
	"go.uber.org/zap"
)

type EchoDeliver struct {
	User      user.UseCase
	z         *zap.SugaredLogger
	ttl       int64
	secretKey []byte
}

func New(useCase user.UseCase, ttl int64, secret string, lgr *zap.SugaredLogger) user.Deliver {
	return &EchoDeliver{
		User:      useCase,
		z:         lgr,
		ttl:       ttl,
		secretKey: []byte(secret),
	}
}

type Payload struct {
	jwt.StandardClaims
	Name string
}

func (e *EchoDeliver) Login(ectx echo.Context) error {
	name := ectx.QueryParam("name")
	password := ectx.QueryParam("password")
	e.z.Infof("attempt to login with name: %s and password: %s", name, password)
	if name == "" || password == "" {
		e.z.Errorf("Empty login or password")
		return echo.ErrUnauthorized
	}
	if ok, err := e.User.Validate(ectx.Request().Context(), name, password); err != nil || !ok {
		e.z.Errorf("can't validate password %s for user %s -> %s ", name, password, err)
		return echo.ErrUnauthorized
	}
	e.z.Infof("Validation succeded")
	payload := Payload{
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: e.ttl,
		},
		Name: name,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &payload)

	signedToken, err := token.SignedString(e.secretKey)
	if err != nil {
		return err
	}

	ectx.Response().Header().Set("X-Expires-After", time.Unix(e.ttl, 0).String())
	return ectx.JSON(http.StatusOK, signedToken)
}

func (e *EchoDeliver) Create(ectx echo.Context) (err error) {
	u := &models.User{}
	err = json.NewDecoder(ectx.Request().Body).Decode(u)
	if err != nil {
		return echo.ErrInternalServerError
	}
	if err = e.User.Add(ectx.Request().Context(), u); err != nil {
		return echo.ErrInternalServerError
	}
	return ectx.JSON(http.StatusCreated, u)
}

func (e *EchoDeliver) Get(ectx echo.Context) error {
	name := ectx.QueryParam("name")
	id, err := e.User.Get(ectx.Request().Context(), name)
	if err != nil {
		return echo.ErrBadRequest
	}
	return ectx.JSON(http.StatusOK, id)

}
