package deliver

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/url"
	"go.uber.org/zap"
)

type urlDeliver struct {
	usecase url.UseCase
	z       *zap.SugaredLogger
}

func New(u url.UseCase, z *zap.SugaredLogger) url.Deliver {
	return &urlDeliver{
		usecase: u,
		z:       z,
	}
}

// TODO: TEST IT!
func (u *urlDeliver) Save(ectx echo.Context) error {
	temUrl := struct {
		Raw    string `json:"rawurl"`
		UserId int    `json:"userid"`
	}{}
	err := json.NewDecoder(ectx.Request().Body).Decode(&temUrl)
	if err != nil {
		u.z.Errorf("can't parse json: %s", err)
		return echo.ErrBadRequest
	}
	sh, err := u.usecase.Add(ectx.Request().Context(), temUrl.Raw, temUrl.UserId)
	if err != nil {
		u.z.Errorf("can't add url: %s", err)
		return echo.ErrBadRequest
	}
	return ectx.JSON(http.StatusCreated, struct {
		Shortened string `json:"shortened"`
	}{Shortened: sh})
}

func (u *urlDeliver) Get(ectx echo.Context) error {
	shortened := ectx.Param("shortened")
	u.z.Infof("get short url %s", shortened)
	url, err := u.usecase.Get(ectx.Request().Context(), shortened)
	if err != nil {
		u.z.Errorf("Can't get url: %s", err)
		return fmt.Errorf("can't get url: %s", err)
	}
	return ectx.Redirect(http.StatusSeeOther, url)
}

func (u *urlDeliver) List(ectx echo.Context) error {
	id, err := strconv.Atoi(ectx.QueryParam("id"))
	u.z.Infof("listing all the urls for user: %d", id)
	if err != nil {
		u.z.Errorf("can't parse id param to string %s", err)
		return fmt.Errorf("can't parse id param to string %s", err)
	}
	ms, err := u.usecase.List(ectx.Request().Context(), id)
	if err != nil {
		u.z.Errorf("can't get list: %s", err)
		return fmt.Errorf("can't get list: %s", err)
	}
	return ectx.JSON(http.StatusOK, ms)
}
