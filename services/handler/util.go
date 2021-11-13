package handler

import (
	"context"
	pb "github.com/aibotsoft/gen/fortedpb"
	api "github.com/aibotsoft/gen/pinapi"
	"github.com/aibotsoft/micro/util"
	"github.com/pkg/errors"
	"regexp"
	"sort"
	"strconv"
	"time"
)

//var urlRe = regexp.MustCompile(`Odds\/(.*?)-(\d+)\/Market\/(\d+)\/(\d+)\/(\d+)`)
var urlRe = regexp.MustCompile(`Sports\/(\d+)\/Leagues\/(\d+)\/Events\/(\d+)`)

func ParseUrl(side *pb.SurebetSide) (err error) {
	u := urlRe.FindStringSubmatch(side.Url)
	if len(u) < 4 {
		return errors.Errorf("parse url error %s", side.Url)
	}
	//if u[1] != side.SportName {
	//	return errors.Errorf("sport name %q does not match in url %q", side.SportName, u[2])
	//}
	side.SportId, err = strconv.ParseInt(u[1], 10, 64)
	if err != nil {
		return errors.Wrapf(err, "sportId error %v", u[2])
	}
	side.LeagueId, err = strconv.ParseInt(u[2], 10, 64)
	if err != nil {
		return errors.Wrapf(err, "leagueId error %v", u[4])
	}
	side.EventId = u[3]
	return nil
}

func (h *Handler) GetCurrency(sb *pb.Surebet) float64 {
	for _, currency := range sb.Currency {
		if currency.Code == h.account.CurrencyCode {
			return currency.Value
		}
	}
	return 0
}

type Ticket struct {
	Market
	api.LineResponse
	SportId int64
	EventId int64
}

func (h *Handler) SetTicket(id int64, t Ticket) {
	h.store.Cache.SetWithTTL(id, t, 1, time.Minute)
}

func (h *Handler) GetTicket(id int64) (Ticket, error) {
	get, b := h.store.Cache.Get(id)
	if !b {
		return Ticket{}, errors.Errorf("not found ticket in cache with id: %v", id)
	}
	return get.(Ticket), nil
}

func (h *Handler) CalcTotalMiddleMargin(ctx context.Context, eventId string, period int) (middleMarginAvg float64, err error) {
	var margin, marginSum, middleMarginSum float64
	spreads, err := h.store.SelectTotals(ctx, eventId, period)
	if err != nil {
		return
	}
	spCount := len(spreads)
	switch spCount {
	case 0:
		h.log.Info("no_spreads_for_event: ", eventId)
	case 1:
		//margin = 1/spreads[0].Home + 1/spreads[0].Away - 1
	default:
		var maxMargin float64

		for a := range spreads {
			m := 1/spreads[a].Over + 1/spreads[a].Under - 1
			marginSum += m
			if m > maxMargin {
				maxMargin = m
			}
		}
		//margin = marginSum / float64(spCount)
		margin = maxMargin

		var marginSlice []float64
		for i := 1; i < spCount; i++ {
			pointDiff := spreads[i].Points - spreads[i-1].Points
			mm := (1/spreads[i-1].Over + 1/spreads[i].Under - 1 - margin) / pointDiff
			marginSlice = append(marginSlice, mm)
			//h.log.Info(pointDiff, " ", util.TruncateFloat(mm, 2))
		}
		sort.Float64s(marginSlice)
		if len(marginSlice) > 2 {
			for i := 0; i < len(marginSlice)-2; i++ {
				middleMarginSum += marginSlice[i]
			}
			middleMarginAvg = util.TruncateFloat(middleMarginSum/float64(len(marginSlice)-2)*100, 2)
		} else {
			for i := 0; i < len(marginSlice); i++ {
				middleMarginSum += marginSlice[i]
			}
			middleMarginAvg = util.TruncateFloat(middleMarginSum/float64(len(marginSlice))*100, 2)
		}
		h.log.Debug("total, margin:", util.TruncateFloat(margin*100, 2), " len:", spCount, " middleMarginAvg:", middleMarginAvg)
	}
	return
}
func (h *Handler) CalcHandicapMiddleMargin(ctx context.Context, eventId string, period int) (middleMarginAvg float64, err error) {
	var margin, marginSum, middleMarginSum float64

	spreads, err := h.store.SelectSpreads(ctx, eventId, period)
	if err != nil {
		return
	}
	spCount := len(spreads)
	switch spCount {
	case 0:
		h.log.Info("no_spreads_for_event: ", eventId)
	case 1:
		margin = 1/spreads[0].Home + 1/spreads[0].Away - 1
	default:
		var maxMargin float64
		for a := range spreads {
			m := 1/spreads[a].Home + 1/spreads[a].Away - 1
			if m > maxMargin {
				maxMargin = m
			}
			marginSum += m
			//h.log.Debug("margin: ", m, "max: ", maxMargin)
		}
		//margin = marginSum / float64(spCount)
		margin = maxMargin
		var marginSlice []float64
		for i := 1; i < spCount; i++ {
			pointDiff := spreads[i].Hdp - spreads[i-1].Hdp
			mm := (1/spreads[i-1].Away + 1/spreads[i].Home - 1 - margin) / pointDiff
			marginSlice = append(marginSlice, mm)
			//h.log.Info(pointDiff, " ", util.TruncateFloat(mm, 2), spreads[i-1].Away, spreads[i].Home, spreads[i].Hdp)
		}
		sort.Float64s(marginSlice)
		if len(marginSlice) > 2 {
			//h.log.Infof("len(marginSlice) > 2 ")
			for i := 0; i < len(marginSlice)-2; i++ {
				middleMarginSum += marginSlice[i]
				//h.log.Info(marginSlice[i])
			}
			middleMarginAvg = util.TruncateFloat(middleMarginSum/float64(len(marginSlice)-2)*100, 2)
		} else {
			if len(marginSlice) > 1 {
				middleMarginAvg = util.TruncateFloat(marginSlice[0]*100, 2)
			}
		}
		h.log.Debug("handicap, margin:", util.TruncateFloat(margin*100, 2), " len:", spCount, " middleMarginAvg:", middleMarginAvg)
	}
	return
}

//if spreads[a].AltLineId == nil {
//	margin = 1/spreads[a].Home + 1/spreads[a].Away - 1
//	//if a < spCount {
//	//	b = a + 1
//	//} else {
//	//	b = a - 1
//	//}
//	b = a + 1
//	diff :=spreads[b].Hdp - spreads[a].Hdp
//	fuck := ((1/spreads[a].Away + 1/spreads[b].Home - 1) - margin) / diff
//	//x:=1/1.869 + 1/1.383 -1
//	h.log.Infow("", "count", spCount, "margin", margin, "diff", diff, "fuck", fuck)
//	if fuck > 0.3 {
//		h.log.Infow("", "a", spreads[a], "margin", margin, "b", spreads[b], "fuck", fuck, "diff", diff)
//
//	}
//}
