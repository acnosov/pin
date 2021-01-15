package collector

import (
	"context"
	"github.com/aibotsoft/gen/pinapi"
	api "github.com/aibotsoft/gen/pinapi"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/sqlserver"
	"github.com/aibotsoft/pin/pkg/store"
	"github.com/stretchr/testify/assert"
	"testing"
)

var c *Collector

func TestMain(m *testing.M) {
	cfg := config.New()
	log := logger.New()
	db := sqlserver.MustConnectX(cfg)
	sto := store.NewStore(cfg, log, db)
	c = NewCollector(cfg, log, sto)
	m.Run()
}

func TestCollector_CollectJob(t *testing.T) {
	c.CollectJob()
	//time.Sleep(time.Second)
	//c.CollectJob()
	//assert.NoError(t, err)
}

func TestCollector_Periods(t *testing.T) {
	s := pinapi.Sport{}
	s.SetId(29)
	err := c.Periods(context.Background(), s)
	assert.NoError(t, err)

}

func TestCollector_Balance(t *testing.T) {
	got, err := c.Balance(context.Background())
	if assert.NoError(t, err) {
		assert.NotEmpty(t, got)
		t.Log(got)
	}
}

func TestCollector_Bets(t *testing.T) {
	err := c.BetList(context.Background())
	assert.NoError(t, err)
}

func TestCollector_Lines(t *testing.T) {
	err := c.Lines(context.Background(), api.Sport{
		Id: api.PtrInt(4),
	})
	assert.NoError(t, err)
}
