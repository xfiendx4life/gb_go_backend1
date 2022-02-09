package deliver

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/redirects"
	"go.uber.org/zap"
)

type del struct {
	use redirects.UseCase
	z   zap.SugaredLogger
}

func (d *del) GetSummary(ectx echo.Context) error {
	shortened := ectx.Param("shortened")
	summary, err := d.use.Get(ectx.Request().Context(), shortened)
	if err != nil {
		d.z.Errorf("can't get summary by %s: %s", shortened, err)
		return echo.ErrBadRequest
	}
	return ectx.JSON(http.StatusOK, summary)
}
