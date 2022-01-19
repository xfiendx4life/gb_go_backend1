package deliver

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/gb_go_backend1/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/pkg/url"
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
	sh, err := u.usecase.Add(ectx.Request().Context(), temUrl.Raw, temUrl.UserId, u.z)
	if err != nil {
		u.z.Errorf("can't add url: %s", err)
		return echo.ErrBadRequest
	}
	return ectx.JSON(200, struct {
		Shortened string `json:"shortened"`
	}{Shortened: sh})
}

func (u *urlDeliver) Get(ectx echo.Context) (*models.Url, error) {
	return nil, nil
}

func (u *urlDeliver) List(ectx echo.Context) ([]models.Url, error) {
	return nil, nil
}
