package main

import (
	"github.com/xfiendx4life/gb_go_backend1/cmd/shrtener/app"
	"github.com/xfiendx4life/gb_go_backend1/internal/logger"
	"go.uber.org/zap/zapcore"
)

func main() {
	app.App(logger.InitLogger(zapcore.DebugLevel, ""))
}
