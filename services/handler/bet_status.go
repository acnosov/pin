package handler

import (
	"context"
	"time"
)

var BettingStatus bool

func (h *Handler) BetStatusRound(ctx context.Context) {
	resp, err := h.client.GetBettingStatus(ctx)
	if err != nil {
		h.log.Info(err)
		BettingStatus = false
	} else {
		switch resp.GetStatus() {
		case "ALL_BETTING_ENABLED":
			BettingStatus = true
		case "ALL_BETTING_CLOSED":
			BettingStatus = false
			h.log.Infow("got_betting_status", "resp", resp, "betting_status", BettingStatus)
		}
	}
}
func (h *Handler) BetStatusJob() {
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		h.BetStatusRound(ctx)
		cancel()
		time.Sleep(time.Minute)
	}
}
