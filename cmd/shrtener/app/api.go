package app

import (
	"context"
	"log"
	"os"

	"github.com/labstack/echo/v4"
	userDel "github.com/xfiendx4life/gb_go_backend1/pkg/user/deliver"
	"github.com/xfiendx4life/gb_go_backend1/pkg/user/usecase"
	"github.com/xfiendx4life/gb_go_backend1/storage"
	"go.uber.org/zap"
)

const port = ":8000"

func App(z *zap.SugaredLogger) {
	ctx := context.Background() // TODO Change for context with Timeout
	server := echo.New()
	store := storage.New()
	os.Setenv("MAX_CONS", "10") // TODO Change for config
	os.Setenv("MIN_CONS", "5")  // TODO Change for config
	err := store.InitNewStorage(ctx, "postgres://xfiendx4life:123456@172.17.0.2:5432/shortener", z)
	if err != nil {
		log.Fatalf("can't connect to storage")
	}
	user := usecase.New(store)
	userDeliver := userDel.New(user, z)
	server.POST("/user/create", userDeliver.Create)
	log.Fatal(server.Start(port))
}
