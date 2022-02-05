package app

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xfiendx4life/gb_go_backend1/internal/api/middlewares"
	"github.com/xfiendx4life/gb_go_backend1/internal/config"
	urlDel "github.com/xfiendx4life/gb_go_backend1/internal/pkg/url/deliver"
	urlRepo "github.com/xfiendx4life/gb_go_backend1/internal/pkg/url/repository"
	urlCase "github.com/xfiendx4life/gb_go_backend1/internal/pkg/url/usecase"
	userDel "github.com/xfiendx4life/gb_go_backend1/internal/pkg/user/deliver"
	userRepo "github.com/xfiendx4life/gb_go_backend1/internal/pkg/user/repository"
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.GetTimeOut())*time.Second) // TODO hange for context with Timeout
	defer cancel()
	// ctx := context.Background()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	server := echo.New()

	server.Use(middleware.Recover())
	store := storage.New()
	err = store.InitNewStorage(ctx, z, conf.GetConfStorage())
	if err != nil {
		log.Fatalf("can't connect to storage")
	}
	user := userCase.New(userRepo.New(store, z), z)
	ttl := conf.GetConfAuth().GetTtl()
	dur := time.Duration(ttl) * time.Minute
	z.Infof("expiry time %v", dur)
	z.Infof("expiry time %d", time.Now().Add(dur).Unix())
	userDeliver := userDel.New(user,
		time.Now().Add(dur).Unix(),
		conf.GetConfAuth().GetSecretKey(),
		z)
	url := urlCase.New(urlRepo.New(store, z), z)
	urlDeliver := urlDel.New(url, z)
	server.POST("/user/create", userDeliver.Create)
	server.GET("/user/login", userDeliver.Login)
	server.POST("/url", urlDeliver.Save, middlewares.JWTAuthMiddleware(conf.GetConfAuth().GetSecretKey()))

	go func() {
		log.Fatal(server.Start(port))
	}()
	<-sigs
	z.Errorf(("done with syscall"))
	// ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.GetTimeOut().Seconds())) // TODO hange for context with Timeout
	// defer cancel()
	if err := server.Shutdown(ctx); err != nil {
		z.Fatalf("can't shutdown")
	}
}
