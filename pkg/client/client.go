package client

import (
	"context"
	api "github.com/aibotsoft/gen/pinapi"
	"github.com/aibotsoft/micro/config"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

type Client struct {
	cfg *config.Config
	log *zap.SugaredLogger
	*api.APIClient
	basicAuth api.BasicAuth
}

const CheckAcceptedTryMaxCount = 10

func (c *Client) CheckAccepted(uniqueRequestId string) (api.PlaceBetResponseV2, error) {
	//todo: эту функцию можно адаптировать к ожиданию принятия лайф ставок
	var res api.GetBetsByTypeResponseV3
	var err error
	for i := 0; i < CheckAcceptedTryMaxCount; i++ {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		res, err = c.GetBetByUUID(ctx, uniqueRequestId)
		cancel()
		if err != nil {
			break
		}
	}
	if err != nil {
		return api.PlaceBetResponseV2{}, errors.Wrapf(err, "CheckAccepted error after %v try", CheckAcceptedTryMaxCount)
	}
	if len(res.GetStraightBets()) == 0 {
		return api.PlaceBetResponseV2{Status: api.PROCESSED_WITH_ERROR, UniqueRequestId: &uniqueRequestId}, nil
	}
	if len(res.GetStraightBets()) > 1 {
		c.log.Warn("more then one bet", "bets", res.GetStraightBets())
	}
	return api.PlaceBetResponseV2{Status: api.ACCEPTED, UniqueRequestId: &uniqueRequestId, StraightBet: &res.GetStraightBets()[0]}, nil
}

func (c *Client) GetBets(ctx context.Context) (api.GetBetsByTypeResponseV3, error) {
	fromDate := time.Now().Add(-time.Minute * 1440).UTC().Format(time.RFC3339)
	toDate := time.Now().UTC().Format(time.RFC3339)
	req := c.GetBetsApi.BetsGetBetsByTypeV3(c.auth(ctx)).FromDate(fromDate).ToDate(toDate).Betlist("ALL")
	got, _, err := req.Execute()
	if err != nil {
		return got, errors.Wrap(err, "GetBets error")
	}
	return got, nil
}

func (c *Client) GetBetByUUID(ctx context.Context, uniqueRequestId string) (api.GetBetsByTypeResponseV3, error) {
	req := c.GetBetsApi.BetsGetBetsByTypeV3(c.auth(ctx)).UniqueRequestIds([]string{uniqueRequestId})
	res, _, err := req.Execute()
	if err != nil {
		return res, errors.Wrap(err, "GetBetByUUID error")
	}
	return res, nil
}

func (c *Client) PlaceBet(ctx context.Context,
	sportId int64,
	eventId int64,
	side string,
	periodNumber int,
	betType string,
	team string,
	lineId int64,
	altLineId *int64,
	uniqueRequestId string,
	stake float64,

) (api.PlaceBetResponseV2, error) {

	r := api.PlaceBetRequest{
		OddsFormat:       api.DECIMAL,
		UniqueRequestId:  uniqueRequestId,
		AcceptBetterLine: true,

		Stake:             stake,
		WinRiskStake:      api.RISK,
		LineId:            lineId,
		Pitcher1MustStart: nil,
		Pitcher2MustStart: nil,
		FillType:          api.FILLANDKILL,
		SportId:           sportId,
		EventId:           eventId,
		PeriodNumber:      periodNumber,
		BetType:           betType,
	}
	if altLineId != nil {
		r.SetAltLineId(*altLineId)
	}
	if side != "" {
		r.SetSide(side)
	}
	if team != "" {
		r.SetTeam(team)
	}

	req := c.PlaceBetsApi.BetsStraightV2(c.auth(ctx)).Request(r)

	res, _, err := req.Execute()
	if err != nil {
		return res, errors.Wrap(err, "place bet error")
	}
	return res, nil
}

//Side: This is needed only for TOTAL_POINTS and TEAM_TOTAL_POINTS  Available values : OVER, UNDER \n
//Handicap: This is needed for SPREAD, TOTAL_POINTS and TEAM_TOTAL_POINTS bet types
//PeriodNumber: For example, for soccer we have 0 (Game), 1 (1st Half) & 2 (2nd Half)
//BetType: Available values : SPREAD, MONEYLINE, TOTAL_POINTS, TEAM_TOTAL_POINTS
//Team: This is needed only for SPREAD, MONEYLINE and TEAM_TOTAL_POINTS bet types  Available values : Team1, Team2, Draw
func (c *Client) CheckLine(ctx context.Context,
	sportId int,
	leagueId int,
	eventId int64,
	side string,
	handicap *float64, periodNumber int, betType string, team string) (api.LineResponse, error) {
	req := c.LineApi.LineStraightV1Get(c.auth(ctx)).SportId(sportId).LeagueId(leagueId).EventId(eventId).OddsFormat("Decimal").BetType(betType).PeriodNumber(periodNumber)
	if handicap != nil {
		req = req.Handicap(*handicap)
	}
	if side != "" {
		req = req.Side(side)
	}
	if team != "" {
		req = req.Team(team)
	}
	res, _, err := req.Execute()
	if err != nil {
		return res, errors.Wrap(err, "get line error")
	}
	return res, nil
}
func (c *Client) GetSports(ctx context.Context) ([]api.Sport, error) {
	sports, _, err := c.OthersApi.SportsV2Get(c.auth(ctx)).Execute()
	if err != nil {
		return nil, errors.Wrap(err, "get sports error")
	}
	return sports.GetSports(), nil
}
func (c *Client) GetLeagues(ctx context.Context, sportId int) ([]api.League, error) {
	leagues, _, err := c.OthersApi.LeaguesV2Get(c.auth(ctx)).SportId(sportId).Execute()
	if err != nil {
		return nil, errors.Wrap(err, "get leagues error")
	}
	return leagues.GetLeagues(), nil
}
func (c *Client) GetPeriods(ctx context.Context, sportId int) ([]api.Period, error) {
	got, _, err := c.OthersApi.PeriodsV1Get(c.auth(ctx)).SportId(sportId).Execute()
	if err != nil {
		return nil, errors.Wrap(err, "get periods error")
	}
	return got.GetPeriods(), nil
}
func (c *Client) GetEvents(ctx context.Context, sportId int, since int64) (api.FixturesResponse, error) {
	res, _, err := c.FixturesApi.FixturesV1Get(c.auth(ctx)).SportId(sportId).Since(since).Execute()
	if err != nil {
		return api.FixturesResponse{}, errors.Wrap(err, "get fixtures error")
	}
	return res, nil
}

func (c *Client) GetSpecials(ctx context.Context, SportId int, since int64) (api.SpecialsFixturesResponse, error) {
	res, r, err := c.FixturesApi.FixturesSpecialV1Get(c.auth(ctx)).SportId(SportId).Since(since).Execute()
	if err != nil {
		if r != nil && r.StatusCode == 304 {
			//c.log.Info("304")
			return api.SpecialsFixturesResponse{}, nil
		}
		return res, errors.Wrap(err, "get specials error")
	}
	return res, nil
}

func (c *Client) GetLines(ctx context.Context, sportId int, since int64) (api.OddsResponse, error) {
	odds, _, err := c.OddsApi.OddsStraightV1Get(c.auth(ctx)).SportId(sportId).Since(since).OddsFormat("Decimal").Execute()
	if err != nil {
		return api.OddsResponse{}, errors.Wrap(err, "get odds error")
	}
	return odds, nil
}

func (c *Client) GetSpecialLines(ctx context.Context, sportId int, since int64) (api.SpecialOddsResponse, error) {
	odds, r, err := c.OddsApi.OddsSpecialV1Get(c.auth(ctx)).SportId(sportId).Since(since).Execute()
	if err != nil {
		if r != nil && r.StatusCode == 304 {
			//c.log.Info("304")
			return api.SpecialOddsResponse{}, nil
		}
		return odds, errors.Wrap(err, "get special odds error")
	}
	return odds, nil
}

func (c *Client) GetBalance(ctx context.Context) (api.ClientBalanceResponse, error) {
	got, _, err := c.ClientBalanceApi.ClientBalanceV1Get(c.auth(ctx)).Execute()
	if err != nil {
		return got, errors.Wrap(err, "get balance error")
	}
	return got, nil
}
func (c *Client) GetBettingStatus(ctx context.Context) (api.BettingStatusResponse, error) {
	got, _, err := c.BettingStatusApi.BetsGetBettingStatus(c.auth(ctx)).Execute()
	if err != nil {
		return got, errors.Wrap(err, "get_betting_status_error")
	}
	return got, nil
}
func (c *Client) auth(ctx context.Context) context.Context {
	return context.WithValue(ctx, api.ContextBasicAuth, c.basicAuth)
}

func NewClient(cfg *config.Config, log *zap.SugaredLogger, username string, password string) *Client {
	//auth := context.WithValue(context.Background(), api.ContextBasicAuth, api.BasicAuth{UserName: username, Password: password})
	ba := api.BasicAuth{UserName: username, Password: password}
	clientConfig := api.NewConfiguration()
	clientConfig.Debug = cfg.Service.Debug
	return &Client{cfg: cfg, log: log, APIClient: api.NewAPIClient(clientConfig), basicAuth: ba}
}
