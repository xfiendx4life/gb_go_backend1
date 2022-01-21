package app

import (
	"context"
	"log"
	"time"

	"github.com/labstack/echo/v4"
	urlDel "github.com/xfiendx4life/gb_go_backend1/pkg/url/deliver"
	urlCase "github.com/xfiendx4life/gb_go_backend1/pkg/url/usecase"
	userDel "github.com/xfiendx4life/gb_go_backend1/pkg/user/deliver"
	userCase "github.com/xfiendx4life/gb_go_backend1/pkg/user/usecase"
	"github.com/xfiendx4life/gb_go_backend1/storage"
	"go.uber.org/zap"
)

const port = ":8000"

func App(z *zap.SugaredLogger) {
	ctx := context.Background() // TODO Change for context with Timeout
	server := echo.New()
	store := storage.New()
	err := store.InitNewStorage(ctx, "postgres://xfiendx4life:123456@172.17.0.2:5432/shortener", z)
	if err != nil {
		log.Fatalf("can't connect to storage")
	}
	user := userCase.New(store)
	userDeliver := userDel.New(user,
		time.Now().Add(time.Hour).Unix(), // ! read from congig
		"somesecret",                     // ! Read from config
		z)
	url := urlCase.New(store)
	urlDeliver := urlDel.New(url, z)
	server.POST("/user/create", userDeliver.Create)
	server.GET("/user/login", userDeliver.Login)
	server.POST("/url", urlDeliver.Save)
	log.Fatal(server.Start(port))
}
