package middlewares

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xfiendx4life/gb_go_backend1/internal/pkg/user/deliver"
)

func JWTAuthMiddleware(secret string) echo.MiddlewareFunc {
	return middleware.JWTWithConfig(middleware.JWTConfig{
		SuccessHandler: nil,
		SigningKey:     []byte(secret),
		Claims:         &deliver.Payload{},
	})
}
