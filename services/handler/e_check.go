package handler

import (
	"context"
	"fmt"
	"github.com/aibotsoft/gen/epinapi"
	pb "github.com/aibotsoft/gen/fortedpb"
	"github.com/aibotsoft/micro/status"
	"github.com/aibotsoft/micro/util"
	"strconv"
	"strings"
	"sync"
)

var EMarketMap = map[string]string{MONEYLINE: "m", SPREAD: "s", TOTAL_POINTS: "ou", TEAM_TOTAL_POINTS: "tt"}
var ETeamMap = map[string]string{Team1: "home", Team2: "away"}
var ETicketMap sync.Map

func calcDesignation(m Market) string {
	if m.BetType == SPREAD || m.BetType == MONEYLINE {
		return ETeamMap[m.Team]
	}
	return strings.ToLower(m.Side)
}
func calcHandicap(m Market) string {
	if m.Handicap == nil {
		return ""
	}
	if m.BetType == SPREAD && m.Team == Team2 {
		p := *m.Handicap * -1
		return strconv.FormatFloat(p, 'g', -1, 64)
	}
	return strconv.FormatFloat(*m.Handicap, 'g', -1, 64)
}

func ConvertEMarket(m Market) string {
	marketKey := fmt.Sprintf("s;%d;%s", m.PeriodNumber, EMarketMap[m.BetType])
	if m.BetType == SPREAD || m.BetType == TOTAL_POINTS {
		marketKey = marketKey + ";" + calcHandicap(m)
	} else if m.BetType == TEAM_TOTAL_POINTS {
		marketKey = marketKey + ";" + calcHandicap(m) + ";" + ETeamMap[m.Team]
	}
	return marketKey
}
func (h *Handler) ECheck(ctx context.Context, sb *pb.Surebet, m Market) {
	side := sb.Members[0]
	marketKey := ConvertEMarket(m)
	//h.log.Infow(marketKey, "team", m.Team, "hand", m.Handicap, "calc", calcHandicap(m))
	eventId, _ := strconv.ParseInt(side.EventId, 10, 64)
	resp, r, err := h.eClient.BetApi.GetLine(h.auth.Auth(ctx)).PlaceBetRequest(epinapi.PlaceBetRequest{
		OddsFormat: "decimal",
		Selections: []epinapi.SelectionItem{{
			MatchupId:   eventId,
			MarketKey:   marketKey,
			Designation: calcDesignation(m),
			Price:       side.Price,
			Points:      m.Handicap,
		}},
	}).Execute()
	if err != nil {
		if r == nil {
			side.Check.Status = status.StatusError
			side.Check.StatusInfo = "request_error"
			return
		}
		switch r.StatusCode {
		case 410:
			side.Check.Status = status.StatusNotFound
			side.Check.StatusInfo = "event_gone"
		case 404:
			side.Check.Status = status.StatusNotFound
			side.Check.StatusInfo = "event_not_found_in_service"
			apiError := err.(epinapi.GenericOpenAPIError)
			body := apiError.Body()
			h.log.Info(string(body))
		case 400:
			side.Check.Status = status.StatusNotFound
			apiError := err.(epinapi.GenericOpenAPIError)
			body := apiError.Body()
			h.log.Info(string(body))

		case 401:
			side.Check.Status = status.BadBettingStatus
			side.Check.StatusInfo = "service_401_error"
		default:
			h.log.Error(err)
			side.Check.Status = status.StatusNotFound
			side.Check.StatusInfo = err.Error()
		}
		h.log.Infow("check_no_ok", "check", side.Check)
		return
	}
	selections := resp.GetSelections()
	if len(selections) != 1 {
		side.Check.Status = status.StatusError
		side.Check.StatusInfo = "no_selections_in_response"
		return
	}
	ticket := selections[0]
	ETicketMap.Store(side.Check.Id, ticket)

	side.Check.Status = status.StatusOk
	side.Check.Balance = util.ToUSDInt(h.GetBalance(), side.Check.Currency)
	side.Check.FillFactor = h.balance.CalcFillFactor()

	side.Check.Price = ticket.GetPrice()
	side.Starts = sb.Starts

	for _, limit := range resp.GetLimits() {
		if limit.GetType() == "minRiskStake" {
			side.Check.MinBet = util.ToUSD(limit.GetAmount(), side.Check.Currency)
		} else if limit.GetType() == "maxRiskStake" {
			side.Check.MaxBet = util.ToUSD(limit.GetAmount(), side.Check.Currency)
		}
	}
	//h.log.Infow(marketKey, "team", m.Team, "hand", m.Handicap, "calc", calcHandicap(m), "price", side.Check.Price)
	h.log.Infow("check", "", side.Check)
}
