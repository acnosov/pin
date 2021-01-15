package handler

import (
	"context"
	"fmt"
	"github.com/aibotsoft/gen/epinapi"
	pb "github.com/aibotsoft/gen/fortedpb"
	api "github.com/aibotsoft/gen/pinapi"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/config_client"
	"github.com/aibotsoft/micro/status"
	"github.com/aibotsoft/micro/util"
	"github.com/aibotsoft/pin/pkg/client"
	"github.com/aibotsoft/pin/pkg/store"
	"github.com/aibotsoft/pin/services/auth"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"net/http"
	"strconv"
	"time"
)

type Handler struct {
	cfg     *config.Config
	log     *zap.SugaredLogger
	client  *client.Client
	store   *store.Store
	balance Balance
	account store.Account
	eClient *epinapi.APIClient
	auth    *auth.Auth
	Conf    *config_client.ConfClient
}

const LockTimeOut = time.Second * 28

func NewHandler(cfg *config.Config, log *zap.SugaredLogger, store *store.Store, auth *auth.Auth, conf *config_client.ConfClient) *Handler {
	ctx := context.Background()
	account, err := store.GetAccount(ctx)
	if err != nil {
		log.Panicw("get account error", "err", err)
	}
	cli := client.NewClient(cfg, log, account.Username, account.Password)

	clientConfig := epinapi.NewConfiguration()
	tr := &http.Transport{TLSHandshakeTimeout: 0 * time.Second, IdleConnTimeout: 0 * time.Second}
	clientConfig.HTTPClient = &http.Client{Transport: tr}
	clientConfig.Debug = cfg.Service.Debug
	eClient := epinapi.NewAPIClient(clientConfig)

	h := &Handler{cfg: cfg, log: log, client: cli, store: store, balance: Balance{},
		account: account, eClient: eClient, auth: auth, Conf: conf}
	return h
}

func (h *Handler) AccId() int {
	return h.account.Id
}
func (h *Handler) FindSportIdByName(ctx context.Context, line *pb.SurebetSide) error {
	err := h.store.FindSportIdByName(ctx, line)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) FindLeagueIdByName(ctx context.Context, line *pb.SurebetSide) error {
	err := h.store.FindLeagueIdByName(ctx, line)
	if err != nil {
		return errors.Wrap(err, "store.FindLeagueIdByName error")
	}
	return nil
}

func (h *Handler) FindEvent(ctx context.Context, line *pb.SurebetSide) error {
	err := h.store.FindEvent(ctx, line)
	if err != nil {
		return err
	}
	return nil
}

func (h *Handler) GetLock(sb *pb.Surebet) bool {
	side := sb.Members[0]
	//key := side.MarketName + side.Url
	key := side.Market.BetType + side.Url
	for i := 0; i < 40; i++ {
		got, b := h.store.Cache.Get(key)
		if b && got.(int64) != sb.SurebetId {
			time.Sleep(time.Millisecond * 50)
		} else {
			return h.store.SetVerifyWithTTL(key, sb.SurebetId, LockTimeOut)
		}
	}
	return false
}
func (h *Handler) ReleaseCheck(ctx context.Context, sb *pb.Surebet) {
	side := sb.Members[0]
	key := side.Market.BetType + side.Url
	got, b := h.store.Cache.Get(key)
	if b && got.(int64) == sb.SurebetId {
		h.store.Cache.Del(key)
	}
}

func (h *Handler) CheckLine(ctx context.Context, sb *pb.Surebet) error {
	side := sb.Members[0]
	side.Check.AccountId = int64(h.account.Id)
	side.Check.AccountLogin = h.account.Username
	side.Check.Currency = h.GetCurrency(sb)
	lockStart := time.Now()
	ok := h.GetLock(sb)
	if !ok {
		side.Check.Status = status.ServiceBusy
		side.Check.StatusInfo = "service_already_check_this_market"
		h.log.Infow(status.ServiceBusy, "h", side.Home, "a", side.Away, "bt", side.Market.BetType, "lock_time", time.Since(lockStart))
		//h.log.Debugw(status.ServiceBusy, "my_surebetId", sb.SurebetId)
		return nil
	}
	lockTime := time.Since(lockStart)
	balance, _ := h.balance.Get()
	if balance < float64(side.BetConfig.MaxStake) {
		h.CheckBalance(true)
	} else {
		go h.CheckBalance(false)
	}

	m, err := Convert(side.MarketName)
	if err != nil {
		h.log.Error(err)
		return nil
	}
	err = ParseUrl(side)
	if err != nil {
		h.log.Error(err)
		return nil
	}
	//h.log.Infow("parse_url", "eventId", side.EventId)
	if side.SportId == 12 {
		h.ECheck(ctx, sb, m)
	} else {
		event, err := h.store.GetEvent(side.EventId)
		if err != nil {
			side.Check.Status = status.StatusNotFound
			side.Check.StatusInfo = "not_found_event_in_db"
			h.log.Infow(side.Check.StatusInfo, "sport", side.SportName, "league", side.LeagueName, "home", side.Home, "away", side.Away, "eventId", side.EventId)
			return nil
		}
		if side.SportName == "Baseball" && m.BetType != MONEYLINE {
			h.log.Infow("got Baseball event from db", "event", event, "home_pitcher", event.GetHomePitcher(), "away_pitcher", event.GetAwayPitcher())
			if event.HasHomePitcher() || event.HasAwayPitcher() {
				side.Check.Status = status.PitchersRequired
				side.Check.StatusInfo = "pitchers in Baseball event"
				return nil
			}
		}
		side.Starts = event.GetStarts().Format(time.RFC3339)

		//h.log.Infow("begin CheckLine", "market", m, "sportId", side.SportId, "leagueId", side.LeagueId, "eventId", side.EventId)
		eventId, _ := strconv.ParseInt(side.EventId, 10, 64)

		line, err := h.client.CheckLine(ctx, int(side.SportId), int(side.LeagueId), eventId, m.Side, m.Handicap, m.PeriodNumber, m.BetType, m.Team)
		if err != nil {
			h.log.Info(err)
			h.log.Infow("get_line_error", "sport", side.SportId, "league", side.LeagueId, "event", eventId, "side", m.Side,
				"handicap", m.Handicap, "period", m.PeriodNumber, "bet_type", m.BetType, "team", m.Team)
			side.Check.StatusInfo = "check_line_request_error"
			return nil
		}
		if line.GetStatus() == "SUCCESS" {
			t := Ticket{Market: m, LineResponse: line, SportId: side.SportId, EventId: eventId}
			h.SetTicket(side.Check.Id, t)
			side.Check.Status = status.StatusOk
			side.Check.StatusInfo = line.GetStatus()
			side.Check.Price = line.GetPrice()

			side.Check.MinBet = util.ToUSD(line.GetMinRiskStake(), side.Check.Currency)
			side.Check.MaxBet = util.ToUSD(line.GetMaxRiskStake(), side.Check.Currency)
			side.Check.Balance = util.ToUSDInt(h.GetBalance(), side.Check.Currency)
			side.Check.FillFactor = h.balance.CalcFillFactor()
		} else if line.GetStatus() == "NOT_EXISTS" {
			side.Check.Status = status.StatusNotFound
			side.Check.StatusInfo = line.GetStatus()
			h.log.Infow("resp", "got", line, "check", side.Check)
		} else if line.GetStatus() == "ALL_BETTING_CLOSED" {
			side.Check.Status = status.BadBettingStatus
			side.Check.StatusInfo = line.GetStatus()
			h.log.Infow("resp", "got", line, "check", side.Check)
		}
	}
	err = h.store.GetStat(side)
	if err != nil {
		h.log.Error(err)
	} else {
		//h.log.Infow("stat", "CountLine", side.Check.CountLine, "AmountLine", side.Check.AmountLine,
		//	"CountEvent", side.Check.CountEvent,"AmountEvent", side.Check.AmountEvent,)
	}

	if sb.Calc.MiddleDiff > 0 {
		fromCache := false
		key := fmt.Sprintf("MiddleMargin:%v:%v:%v", side.EventId, m.PeriodNumber, m.BetType)
		got, b := h.store.Cache.Get(key)
		if b {
			fromCache = true
			side.Check.MiddleMargin = got.(float64)
		} else if m.BetType == SPREAD || m.BetType == MONEYLINE {
			side.Check.MiddleMargin, err = h.CalcHandicapMiddleMargin(ctx, side.EventId, m.PeriodNumber)
			if err != nil {
				h.log.Error(err)
			}
		} else if m.BetType == TOTAL_POINTS || m.BetType == TEAM_TOTAL_POINTS {
			side.Check.MiddleMargin, err = h.CalcTotalMiddleMargin(ctx, side.EventId, m.PeriodNumber)
			if err != nil {
				h.log.Error(err)
			}
		}
		if side.Check.MiddleMargin > 0 && !fromCache {
			h.store.Cache.SetWithTTL(key, side.Check.MiddleMargin, 1, time.Minute)
		}
	}
	if !BettingStatus {
		side.Check.Status = status.BadBettingStatus
		side.Check.StatusInfo = "ALL_BETTING_CLOSED"
	}
	side.Check.SubService = "pin"

	if side.Check.Status == status.StatusOk {
		msg := "check_ok"
		if side.Check.Price < side.Price-0.01 {
			msg = "p_lower"
		}
		h.log.Infow(msg, "p", side.Check.Price, "fp", side.Price, "g", fmt.Sprintf("%v-%v-%v:%v", side.SportName, side.Home, side.Away, side.MarketName),
			"bt", side.Market.BetType,
			"lt", lockTime,
			"sub", side.Check.SubService,
			"diff", sb.Calc.MiddleDiff,
		)
	} else {
		h.log.Infow("check_not_ok", "status", side.Check.Status, "fp", side.Price, "g", fmt.Sprintf("%v-%v-%v:%v", side.SportName, side.Home, side.Away, side.MarketName),
			"bt", side.Market.BetType,
			"lt", lockTime,
			"sub", side.Check.SubService,
		)
	}
	return nil
}

func (h *Handler) PlaceBet(ctx context.Context, sb *pb.Surebet) error {
	side := sb.Members[0]
	//h.log.Info(side)
	//go func() {
	//	err := h.store.SaveCheck(sb)
	//	if err != nil {
	//		h.log.Error(err)
	//	}
	//}()
	err := h.store.SaveCheck(sb)
	if err != nil {
		h.log.Error(err)
	}

	if side.SportId == 12 {
		h.EPlaceBet(ctx, sb)
		return nil
	}
	ticket, err := h.GetTicket(side.Check.Id)
	if err != nil {
		h.log.Info(err)
		return nil
	}
	h.log.Infow("ticket", "t", ticket)
	uniqueRequestId := uuid.New().String()
	stake := util.AdaptStake(side.CheckCalc.Stake, side.Check.Currency, side.BetConfig.RoundValue)

	placeCtx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	var res api.PlaceBetResponseV2
	res, err = h.client.PlaceBet(placeCtx,
		ticket.SportId,
		ticket.EventId,
		ticket.Side,
		ticket.PeriodNumber,
		ticket.BetType,
		ticket.Team,
		*ticket.LineId,
		ticket.AltLineId,
		uniqueRequestId,
		stake)
	if err != nil {
		h.log.Info(err)
		res, err = h.client.CheckAccepted(uniqueRequestId)
		if err != nil {
			h.log.Error(err)
			return nil
		}
	}
	h.log.Infow("placeResult", "placeResult", res)

	if res.GetStatus() == api.ACCEPTED {
		b := res.GetStraightBet()
		side.Bet.Status = status.StatusOk
		side.Bet.StatusInfo = b.GetBetStatus()

		side.Bet.Price = b.GetPrice()
		side.Bet.Stake = util.ToUSD(b.GetRisk(), side.Check.Currency)

		side.Bet.ApiBetId = strconv.FormatInt(b.GetBetId(), 10)

		h.balance.Sub(b.GetRisk())
	} else if res.GetStatus() == api.PROCESSED_WITH_ERROR {
		if res.GetErrorCode() == "ALL_BETTING_CLOSED" {
			BettingStatus = false
		}
		side.Bet.Status = status.StatusNotAccepted
		side.Bet.StatusInfo = string(res.GetErrorCode())
	}
	side.Bet.Done = util.UnixMsNow()
	err = h.store.SaveBet(sb)
	if err != nil {
		h.log.Error(err)
	}
	return nil
}

func (h *Handler) GetResults(ctx context.Context) ([]pb.BetResult, error) {
	return h.store.GetResults(ctx)
}

func (h *Handler) Close() {
	h.store.Close()
}
