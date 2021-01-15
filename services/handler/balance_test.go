package handler

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestBalance_Get(t *testing.T) {
	got, b := h.balance.Get()
	assert.False(t, b)
	assert.Empty(t, got)

	h.balance.Set(1)
	got, b = h.balance.Get()
	assert.True(t, b)
	assert.Equal(t, float64(1), got)

}

func TestHandler_GetBalance(t *testing.T) {
	got := h.GetBalance()
	go h.CheckBalance(false)
	assert.Empty(t, got)
	time.Sleep(time.Second)
	secondTime := h.GetBalance()
	assert.NotEmpty(t, secondTime)
}

func TestHandler_CheckBalance(t *testing.T) {
	h.CheckBalance(false)
}
