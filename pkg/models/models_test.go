package models

import (
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewUrl(t *testing.T) {
	u1 := NewUrl("teststring", 1)
	u2 := NewUrl("teststring", 2)
	log.Print(u1.Shortened)
	assert.Equal(t, u1.Shortened, u2.Shortened)
}

func TestNewUrlNotEqual(t *testing.T) {
	u1 := NewUrl("teststrin", 1)
	u2 := NewUrl("teststring", 2)
	log.Print(u1.Shortened)
	assert.NotEqual(t, u1.Shortened, u2.Shortened)
}
