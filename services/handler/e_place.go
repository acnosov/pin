package handler

import (
	"context"
	"github.com/aibotsoft/gen/epinapi"
	pb "github.com/aibotsoft/gen/fortedpb"
	"github.com/aibotsoft/micro/status"
	"github.com/aibotsoft/micro/util"
	"strconv"
)

func (h *Handler) EPlaceBet(ctx context.Context, sb *pb.Surebet) {
	side := sb.Members[0]
	stake := util.AdaptStake(side.CheckCalc.Stake, side.Check.Currency, side.BetConfig.RoundValue)
	got, ok := ETicketMap.Load(side.Check.Id)
	if !ok {
		side.Bet.Status = status.StatusError
		side.Bet.StatusInfo = "load_ticket_error"
		return
	}
	ticket := got.(epinapi.SelectionItem)
	h.log.Infow("ticket", "", ticket)
	resp, _, err := h.eClient.BetApi.PlaceBet(h.auth.Auth(ctx)).PlaceBetRequest(epinapi.PlaceBetRequest{
		OddsFormat:         "decimal",
		AcceptBetterPrice:  epinapi.PtrBool(true),
		AcceptBetterPrices: epinapi.PtrBool(true),
		Class:              epinapi.PtrString("Straight"),
		Stake:              &stake,
		Selections:         []epinapi.SelectionItem{ticket},
	}).Execute()
	if err != nil {
		apiErr := err.(epinapi.GenericOpenAPIError)
		switch v := apiErr.Model().(type) {
		case epinapi.BadRequestError:
			side.Bet.Status = status.StatusNotAccepted
			side.Bet.StatusInfo = v.GetTitle()
		default:
			side.Bet.Status = status.StatusError
			side.Bet.StatusInfo = "request_error"
			h.log.Error(err)
		}
	} else {
		if resp.GetId() > 0 {
			side.Bet.Status = status.StatusOk
			//side.Bet.StatusInfo = resp.Get
			side.Bet.Price = resp.GetPrice()
			side.Bet.Stake = util.ToUSD(resp.GetStake(), side.Check.Currency)
			side.Bet.ApiBetId = strconv.FormatInt(resp.GetId(), 10)
			//h.balance.Sub(resp.GetStake())
		} else {
			side.Bet.Status = status.StatusNotAccepted
			side.Bet.StatusInfo = "api_bet_id_is_0"
			h.log.Infow("place_bet_response", "resp", resp)
		}
	}
	h.log.Infow("side_bet", "bet", side.Bet)
	side.Bet.Done = util.UnixMsNow()
	err = h.store.SaveBet(sb)
	if err != nil {
		h.log.Error(err)
	}
}
