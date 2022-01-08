package models

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/xfiendx4life/gb_go_backend1/internal/logger"
	"go.uber.org/zap/zapcore"
)

var lgr = logger.InitLogger(zapcore.DebugLevel, "")

func TestNewUrl(t *testing.T) {
	u1 := NewUrl("teststring", 1, lgr)
	u2 := NewUrl("teststring", 2, lgr)
	log.Print(u1.Shortened)
	assert.Equal(t, u1.Shortened, u2.Shortened)
}

func TestNewUrlNotEqual(t *testing.T) {
	u1 := NewUrl("teststrin", 1, lgr)
	u2 := NewUrl("teststring", 2, lgr)
	log.Print(u1.Shortened)
	assert.NotEqual(t, u1.Shortened, u2.Shortened)
}
