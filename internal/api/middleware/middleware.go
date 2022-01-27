package middleware

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

func RecoverMiddleware(z *zap.SugaredLogger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		defer func() {
			if err := recover(); err != nil {
				z.Panicf("Panic recovered with error: %s", err)
			}
		}()
		return func(ectx echo.Context) error {
			return next(ectx)
		}
	}
}
