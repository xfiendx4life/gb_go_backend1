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

// func FormToJsonMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
// 	return func (ectx echo.Context) error {
// 		data := struct {
// 			name string `json:"name"`
// 			password string `josn:"password"`
// 			email string `json:"password"`
// 		}{
// 			name: ectx.FormValue("name")
// 			password: ectx.FormValue("password"),
// 			email: ectx.FormValue("email"),
// 		}
// 		buf := new(bytes.Buffer)
// 		err := json.NewEncoder(buf).Encode(data)
// 		if err != nil {
// 			return echo.ErrBadRequest
// 		}
// 		c := echo.New
// 	}
// }
