package app

import (
	"context"
	"flag"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/xfiendx4life/gb_go_backend1/internal/api/middlewares"
	"github.com/xfiendx4life/gb_go_backend1/internal/config"
	rdrDel "github.com/xfiendx4life/gb_go_backend1/internal/pkg/redirects/deliver"
	rdrRepo "github.com/xfiendx4life/gb_go_backend1/internal/pkg/redirects/repository"
	rdrUsecase "github.com/xfiendx4life/gb_go_backend1/internal/pkg/redirects/usecase"
	urlDel "github.com/xfiendx4life/gb_go_backend1/internal/pkg/url/deliver"
	urlRepo "github.com/xfiendx4life/gb_go_backend1/internal/pkg/url/repository"
	urlCase "github.com/xfiendx4life/gb_go_backend1/internal/pkg/url/usecase"
	userDel "github.com/xfiendx4life/gb_go_backend1/internal/pkg/user/deliver"
	userRepo "github.com/xfiendx4life/gb_go_backend1/internal/pkg/user/repository"
	userCase "github.com/xfiendx4life/gb_go_backend1/internal/pkg/user/usecase"
	"github.com/xfiendx4life/gb_go_backend1/storage"
	"go.uber.org/zap"
)

const port = ":8080"

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func setTemplates(server *echo.Echo) {
	t := &Template{
		templates: template.Must(template.ParseGlob("web/templates/*.html")),
	}
	server.Renderer = t
}

func readConfig(confSource string, z *zap.SugaredLogger) config.Config {
	conf := config.New()
	path, err := os.Getwd()
	if err != nil {
		z.Fatalf("can't find path: %s", err)
		return nil
	}
	var confFile []byte
	if confSource != "" {
		confFile, err = config.ReadFromFile(path+confSource, z)
		if err != nil {
			z.Fatalf("can't read config: %s", err)
			return nil
		}
	} else {
		confFile = config.ReadFromEnv()
	}
	err = conf.ReadConfig(confFile, z)
	if err != nil {
		z.Fatalf("can't read config: %s", err)
		return nil
	}
	return conf
}

func App(z *zap.SugaredLogger) {
	flag.Parse()
	confSource := flag.String("config", "", "Use flag to choose config source, env if empty")
	conf := readConfig(*confSource, z)
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(conf.GetTimeOut())*time.Second) // TODO hange for context with Timeout
	defer cancel()
	// ctx := context.Background()
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	server := echo.New()
	server.Use(middleware.Recover())
	setTemplates(server)
	// * Storage init
	store := storage.New()
	err := store.InitNewStorage(ctx, z, conf.GetConfStorage())
	if err != nil {
		log.Fatalf("can't connect to storage")
	}
	// * User init
	user := userCase.New(userRepo.New(store, z), z)
	ttl := conf.GetConfAuth().GetTtl()
	dur := time.Duration(ttl) * time.Minute
	z.Infof("expiry time %v", dur)
	z.Infof("expiry time %d", time.Now().Add(dur).Unix())
	userDeliver := userDel.New(user,
		time.Now().Add(dur).Unix(),
		conf.GetConfAuth().GetSecretKey(),
		z)
	// *Redirect init
	rdr := rdrUsecase.New(rdrRepo.New(store, z), z)
	rDel := rdrDel.New(rdr, z)

	// *URL init
	url := urlCase.New(urlRepo.New(store, z), rdr, z)
	urlDeliver := urlDel.New(url, z)

	server.Static("static", "./web/static")
	// * Handlers
	server.POST("/user/create", userDeliver.Create)
	server.GET("/user/login", userDeliver.Login)
	server.POST("/url", urlDeliver.Save, middlewares.JWTAuthMiddleware(conf.GetConfAuth().GetSecretKey()))
	server.GET("/:shortened", urlDeliver.Get)
	server.GET("/redirects/:shortened", rDel.GetSummary, middlewares.JWTAuthMiddleware(conf.GetConfAuth().GetSecretKey()))
	server.GET("/web/generate", func(ectx echo.Context) error {
		data := make(map[string]string)
		z.Info("In render handler")
		return ectx.Render(http.StatusOK, "generate", data)
	})

	go func() {
		z.Fatal(server.Start(port))
	}()

	<-sigs
	z.Errorf(("done with syscall"))
	if err := server.Shutdown(ctx); err != nil {
		z.Fatalf("can't shutdown")
	}
}
