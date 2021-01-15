package client

import (
	"context"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/sqlserver"
	"github.com/aibotsoft/pin/pkg/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

var c *Client

func TestMain(m *testing.M) {
	cfg := config.New()
	log := logger.New()
	db := sqlserver.MustConnectX(cfg)
	sto := store.NewStore(cfg, log, db)
	ctx := context.Background()
	account, err := sto.GetAccount(ctx)
	if err != nil {
		log.Panicw("get account error", "err", err)
	}
	c = NewClient(cfg, log, account.Username, account.Password)
	m.Run()
}

func TestClient_GetSports(t *testing.T) {
	sports, err := c.GetSports(context.Background())
	if assert.NoError(t, err) {
		assert.NotEmpty(t, sports)
	}
	t.Log(sports)
}

//https://members.pinnacle.com/Sportsbook/Mobile/ru-RU/Enhanced/Regular/SportsBookAll/35/Curacao/Odds/Soccer-29/Market/2/6416/1119900162

func TestClient_GetBalance(t *testing.T) {
	got, err := c.GetBalance(context.Background())
	if assert.NoError(t, err) {
		assert.NotEmpty(t, got)
		t.Log(got)
	}
}

func TestClient_GetBettingStatus(t *testing.T) {
	got, err := c.GetBettingStatus(context.Background())
	if assert.NoError(t, err) {
		assert.NotEmpty(t, got)
		t.Log(got.GetStatus())
	}

}
