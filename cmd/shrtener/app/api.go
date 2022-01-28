package app

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xfiendx4life/gb_go_backend1/internal/api/middlewares"
	"github.com/xfiendx4life/gb_go_backend1/internal/config"
	urlDel "github.com/xfiendx4life/gb_go_backend1/internal/pkg/url/deliver"
	urlCase "github.com/xfiendx4life/gb_go_backend1/internal/pkg/url/usecase"
	userDel "github.com/xfiendx4life/gb_go_backend1/internal/pkg/user/deliver"
	userCase "github.com/xfiendx4life/gb_go_backend1/internal/pkg/user/usecase"
	"github.com/xfiendx4life/gb_go_backend1/storage"
	"go.uber.org/zap"
)

const port = ":8000"

func App(z *zap.SugaredLogger) {
	conf := config.New()
	path, err := os.Getwd()
	if err != nil {
		z.Fatalf("can't find path: %s", err)
		return
	}
	confFile, err := config.ReadFromFile(path+"/configs/config.yml", z)
	if err != nil {
		z.Fatalf("can't read config: %s", err)
		return
	}
	err = conf.ReadConfig(confFile, z)
	if err != nil {
		z.Fatalf("can't read config: %s", err)
		return //
	}
	ctx := context.Background() // TODO Change for context with Timeout
	server := echo.New()
	// server.Use(middlewares.RecoverMiddleware(z))
	server.Use(middleware.Recover())
	store := storage.New()
	err = store.InitNewStorage(ctx, z, conf.GetConfStorage())
	if err != nil {
		log.Fatalf("can't connect to storage")
	}
	user := userCase.New(store)
	ttl := conf.GetConfAuth().GetTtl()
	dur := time.Duration(ttl) * time.Minute
	z.Infof("expiry time %v", dur)
	z.Infof("expiry time %d", time.Now().Add(dur).Unix())
	userDeliver := userDel.New(user,
		time.Now().Add(dur).Unix(),
		conf.GetConfAuth().GetSecretKey(),
		z)
	url := urlCase.New(store)
	urlDeliver := urlDel.New(url, z)
	server.POST("/user/create", userDeliver.Create)
	server.GET("/user/login", userDeliver.Login)
	server.POST("/url", urlDeliver.Save, middlewares.JWTAuthMiddleware(conf.GetConfAuth().GetSecretKey()))
	log.Fatal(server.Start(port))
	// middleware.REcoverWithConfig()
}
