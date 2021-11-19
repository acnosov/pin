package handler

import (
	"context"
	pb "github.com/aibotsoft/gen/fortedpb"
	"github.com/aibotsoft/pin/pkg/config"
	"github.com/aibotsoft/pin/pkg/sqlserver"
	"github.com/aibotsoft/pin/pkg/store"
	"github.com/aibotsoft/pin/services/auth"
	"github.com/stretchr/testify/assert"
	"github.com/vrischmann/envconfig"
	"go.uber.org/zap"
	"testing"
)

var h *Handler

func TestMain(m *testing.M) {
	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	log := logger.Sugar()

	cfg := &config.Config{}
	err = envconfig.Init(cfg)
	if err != nil {
		panic(err)
	}
	log.Infow("Begin service", "config", cfg)

	db := sqlserver.MustConnect(cfg)
	sto := store.NewStore(cfg, log, db)
	a := auth.New(cfg, log, sto)
	h = NewHandler(cfg, log, sto, a)
	m.Run()
	h.Close()
}

//2. *Pinnacle; ставка ТБ(10,5); коэф. 2.34 (Baseball - Korea Professional Baseball // Hanwha Eagles - Kiwoom Heroes)
//https://members.pinnacle.com/Sportsbook/Mobile/ru-RU/Enhanced/Regular/SportsBookAll/35/Curacao/Odds/Baseball-3/Market/2/6227/1122239544
var currencyList = []pb.Currency{{Code: "USD", Value: 1}, {Code: "EUR", Value: 0.93}}

//П1	229
//П2	105
//1/2 П1	40
//1/2 П2	7
//2/2 П1	15
//2/2 П2	68
//3/2 П1	76
//3/2 П2	7
//4/2 П2	44
//5/2 П2	37
//(карты) ТБ(3,5)	219
//(карты) ТМ(2,5)	3
//(карты) ТМ(3,5)	9
//(карты) Ф1(1,5)	137
//(карты) Ф2(-1,5)	99
//(карты) Ф1(-1,5)	100
//(карты) Ф2(1,5)	9
//(карты) Ф2(2,5)	28
//(карты) Ф1(2,5)	9
//(карты) Ф1(-2,5)	6

var MarketNames = []string{
	"1",
	"2",
	"Х",
	"X",
	//"П1",
	//"П2",
	//"1/2 П1",
	//"1/2 П2",
	//"2/2 П1",
	//"2/2 П2",
	//"3/2 П1",
	//"3/2 П2",
	//"(карты) ТМ(2,5)",
	//"(карты) ТБ(2,5)",
	//"(карты) ТМ(3,5)",
	//"(карты) ТБ(3,5)",
	//"(карты) Ф1(1,5)",
	//"(карты) Ф2(-1,5)",
	//"(карты) Ф1(-1,5)",
	//"(карты) Ф2(1,5)",
	//"1/2 ИТ1М(14,5)",
}

func TestHandler_CheckLine(t *testing.T) {
	ctx := context.Background()
	//h.BetStatusRound(ctx)

	//events, _ := h.store.SelectCurrentEvents(ctx, 29, 1, "Regular")
	//e := events[0]
	side := &pb.SurebetSide{
		ServiceName: "Pinnacle",
		SportName:   "E Sports",
		LeagueName:  "Test",
		Home:        "Test",
		Away:        "Test",
		MarketName:  "1",
		//Url:       "https://members.pinnacle.com/Sportsbook/Mobile/ru-RU/Enhanced/Regular/SportsBookAll/35/Curacao/Odds/E Sports-12/Market/2/208947/1122534935",
		//https://beta.pinnacle.com/en/Sports/29/Leagues/2368/Events/1433894850
		Url:       "https://beta.pinnacle.com/en/Sports/29/Leagues/2368/Events/1433894850",
		Price:     2,
		Check:     &pb.Check{},
		Market:    &pb.Market{},
		BetConfig: &pb.BetConfig{MaxStake: 10},
	}

	sb := &pb.Surebet{Members: []*pb.SurebetSide{side}, Currency: currencyList}
	for _, name := range MarketNames {
		side.MarketName = name
		err := h.CheckLine(ctx, sb)
		if assert.NoError(t, err) {
			//h.log.Infow("side", "sb", sb, "event", e)
			//t.Log(side.Check)
		}
	}

}
