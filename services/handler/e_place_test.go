package handler

import (
	"context"
	pb "github.com/aibotsoft/gen/fortedpb"
	"github.com/aibotsoft/micro/util"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHandler_EPlaceBet(t *testing.T) {
	ctx := context.Background()
	line := &pb.SurebetSide{
		ServiceName: "Pinnacle",
		SportName:   "E Sports",
		LeagueName:  "Dota 2 - SEA Dota Invitational",
		Home:        "Impunity",
		Away:        "LGD.int",
		//MarketName:  "П1",
		MarketName: "Ф1(1,5)",
		Url:        "https://members.pinnacle.com/Sportsbook/Mobile/ru-RU/Enhanced/Regular/SportsBookAll/35/Curacao/Odds/E Sports-12/Market/2/208947/1122534935",
		Check:      &pb.Check{},
		CheckCalc: &pb.CheckCalc{
			Status:   "Ok",
			MaxStake: 1,
			MinStake: 1,
			MaxWin:   2,
			Stake:    1,
			Win:      1,
			IsFirst:  true,
		},
		BetConfig: &pb.BetConfig{MaxStake: 100, RoundValue: 0.01},
		Bet:       &pb.Bet{},
		ToBet:     &pb.ToBet{Id: util.UnixMsNow()},
	}
	sb := &pb.Surebet{Members: []*pb.SurebetSide{line}, Currency: currencyList}

	err := h.CheckLine(ctx, sb)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, line.SportId)
		t.Log(line.SportId)
		err := h.PlaceBet(ctx, sb)
		assert.NoError(t, err)
	}
}
