package deliver

import (
	"encoding/json"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/gb_go_backend1/pkg/models"
	"github.com/xfiendx4life/gb_go_backend1/pkg/user"
	"go.uber.org/zap"
)

type EchoDeliver struct {
	User user.UseCase
	z    *zap.SugaredLogger
}

func New(useCase user.UseCase, lgr *zap.SugaredLogger) user.Deliver {
	return &EchoDeliver{
		User: useCase,
		z:    lgr,
	}
}

func (e *EchoDeliver) Login(ectx echo.Context) error {
	return nil
}

func (e *EchoDeliver) Create(ectx echo.Context) (err error) {
	u := &models.User{}
	// //if err = ectx.Bind(u); err != nil {
	// 	//return echo.ErrBadRequest
	// //}
	err = json.NewDecoder(ectx.Request().Body).Decode(u)
	if err != nil {
		return echo.ErrInternalServerError
	}
	if err = e.User.Add(ectx.Request().Context(), u, e.z); err != nil {
		return echo.ErrInternalServerError
	}
	return ectx.JSON(http.StatusCreated, u)
}
