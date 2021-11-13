package store

import (
	"context"
	"database/sql"
	"fmt"
	pb "github.com/aibotsoft/gen/fortedpb"
	api "github.com/aibotsoft/gen/pinapi"
	"github.com/aibotsoft/pin/pkg/config"
	mssql "github.com/denisenkom/go-mssqldb"
	"github.com/dgraph-io/ristretto"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"strconv"
	"time"
)

type Store struct {
	cfg   *config.Config
	log   *zap.SugaredLogger
	db    *sqlx.DB
	Cache *ristretto.Cache
}

const getAccount = "select top 1 Id, AccountType, CurrencyCode, ServiceName, Username, Password, Commission, Share from Account"

func (s *Store) SetVerifyWithTTL(key string, value interface{}, ttl time.Duration) bool {
	s.Cache.SetWithTTL(key, value, 1, ttl)
	for i := 0; i < 100; i++ {
		got, b := s.Cache.Get(key)
		if b {
			if got == value {
				return true
			} else {
				s.log.Info("got != value:", got, value)
				return false
			}
		}
		time.Sleep(time.Microsecond * 5)
	}
	return false
}
func (s *Store) IsInCache(ctx context.Context, sport api.Sport) bool {
	if sport.Id == nil {
		return false
	}
	key := "Sport:" + strconv.Itoa(int(sport.GetId()))
	got, ok := s.Cache.Get(key)
	if !ok {
		s.Cache.Set(key, sport, 1)
		return false
	}
	if sport == got {
		s.log.Info("got from Cache", got)
		return true
	}
	return false
}

func (s *Store) SaveSports(ctx context.Context, sports []api.Sport) error {
	newSports := sports[:0]
	for _, x := range sports {
		if !s.IsInCache(ctx, x) {
			newSports = append(newSports, x)
		}
	}
	if len(newSports) == 0 {
		return nil
	}
	tvpType := mssql.TVP{TypeName: "SportType", Value: newSports}
	_, err := s.db.ExecContext(ctx, "uspCreateSportTVP", tvpType)

	if err != nil {
		return errors.Wrap(err, "uspCreateSportTVP error")
	}
	return nil
}

func (s *Store) SaveLeagues(ctx context.Context, sportId int, leagues []api.League) error {
	tvpType := mssql.TVP{TypeName: "LeagueType", Value: leagues}
	_, err := s.db.ExecContext(ctx, "uspSaveLeagues", sportId, tvpType)
	if err != nil {
		return errors.Wrap(err, "uspSaveLeagues error")
	}
	return nil
}

func (s *Store) SavePeriods(ctx context.Context, sportId int, got []api.Period) error {
	//s.log.Infow("periods", "periods", got[0])
	tvp := mssql.TVP{TypeName: "PeriodType", Value: got}
	_, err := s.db.ExecContext(ctx, "uspSavePeriods", sportId, tvp)
	if err != nil {
		return errors.Wrap(err, "uspSavePeriods error")
	}
	return nil

}

func (s *Store) SaveEvents(ctx context.Context, leagueId int, events []api.Fixture) error {
	tvpType := mssql.TVP{TypeName: "EventType", Value: events}
	_, err := s.db.ExecContext(ctx, "uspSaveEvents", leagueId, tvpType)
	if err != nil {
		return errors.Wrap(err, "uspSaveEvents error")
	}
	return nil
}

type SpecialType struct {
	Id         *int64
	BetType    *string
	Name       *string
	Date       *time.Time
	Cutoff     *time.Time
	Category   *string
	Units      *string
	Status     *string
	LiveStatus *int
}
type SpecialEventType struct {
	Id           *int
	PeriodNumber *int
	Home         *string
	Away         *string
	SpecialId    *int64
}

type SpecialContestantType struct {
	Id        *int64
	Name      *string
	RotNum    *int
	SpecialId *int64
}

func (s *Store) SaveSpecials(ctx context.Context, leagueId int, specials []api.SpecialFixture) error {
	var specialTypes []SpecialType
	var specialEventTypes []SpecialEventType
	var specialContestantTypes []SpecialContestantType
	for _, sp := range specials {
		st := SpecialType{
			Id:         sp.Id,
			BetType:    sp.BetType,
			Name:       sp.Name,
			Date:       sp.Date,
			Cutoff:     sp.Cutoff,
			Category:   sp.Category,
			Units:      sp.Units,
			Status:     sp.Status,
			LiveStatus: sp.LiveStatus,
		}
		specialTypes = append(specialTypes, st)
		if se, ok := sp.GetEventOk(); ok {
			specialEventTypes = append(specialEventTypes, SpecialEventType{
				Id:           se.Id,
				PeriodNumber: se.PeriodNumber,
				Home:         se.Home,
				Away:         se.Away,
				SpecialId:    st.Id,
			})
		}
		for _, co := range sp.GetContestants() {
			specialContestantTypes = append(specialContestantTypes, SpecialContestantType{
				Id:        co.Id,
				Name:      co.Name,
				RotNum:    co.RotNum,
				SpecialId: st.Id,
			})
		}
	}
	if len(specialTypes) > 0 {
		tvpType := mssql.TVP{TypeName: "SpecialType", Value: specialTypes}
		_, err := s.db.ExecContext(ctx, "uspSaveSpecials", leagueId, tvpType)
		if err != nil {
			s.log.Info(errors.Wrap(err, "uspSaveSpecials error"))
		}
	}

	if len(specialEventTypes) > 0 {
		tvpSpecialEventType := mssql.TVP{TypeName: "SpecialEventType", Value: specialEventTypes}
		_, err := s.db.ExecContext(ctx, "uspSaveSpecialEvents", tvpSpecialEventType)
		if err != nil {
			s.log.Info(errors.Wrap(err, "uspSaveSpecialEvents error"))
			for i, se := range specialEventTypes {
				set := fmt.Sprintf("Id: %v PeriodNumber: %v Home: %v Away: %v SpecialId: %v", *se.Id, *se.PeriodNumber, *se.Home, *se.Away, *se.SpecialId)
				s.log.Info(i, set)
			}
		}
	}
	if len(specialContestantTypes) > 0 {
		tvpSpecialContestantType := mssql.TVP{TypeName: "SpecialContestantType", Value: specialContestantTypes}
		_, err := s.db.ExecContext(ctx, "uspSaveSpecialContestants", tvpSpecialContestantType)
		if err != nil {
			s.log.Info(errors.Wrap(err, "uspSaveSpecialContestants error"))
		}
	}
	return nil
}

type LineType struct {
	LineId       *int64
	EventId      *int64
	Number       *int
	Cutoff       *time.Time
	Status       *int
	MaxSpread    *float64
	MaxMoneyline *float64
	MaxTotal     *float64
	MaxTeamTotal *float64
}
type MoneylineType struct {
	LineId *int64
	Home   *float64
	Away   *float64
	Draw   *float64
}
type TotalType struct {
	LineId    *int64
	AltLineId *int64
	Points    *float64
	Over      *float64
	Under     *float64
}

type SpreadType struct {
	LineId    *int64
	AltLineId *int64
	Hdp       *float64
	Home      *float64
	Away      *float64
}

func (s *Store) SaveLines(ctx context.Context, eventId int64, periods []api.OddsPeriod) error {
	var lineTypes []LineType
	var mlTypes []MoneylineType
	var totalTypes []TotalType
	var spreadTypes []SpreadType

	for _, p := range periods {
		lineTypes = append(lineTypes, LineType{
			LineId:       p.LineId,
			EventId:      &eventId,
			Number:       p.Number,
			Cutoff:       p.Cutoff,
			Status:       p.Status,
			MaxSpread:    p.MaxSpread,
			MaxMoneyline: p.MaxMoneyline,
			MaxTotal:     p.MaxTotal,
			MaxTeamTotal: p.MaxTeamTotal,
		})
		if p.HasMoneyline() {
			ml := p.GetMoneyline()
			mlTypes = append(mlTypes, MoneylineType{
				LineId: p.LineId,
				Home:   ml.Home,
				Away:   ml.Away,
				Draw:   ml.Draw,
			})
		}

		for _, t := range p.GetTotals() {
			totalTypes = append(totalTypes, TotalType{
				LineId:    p.LineId,
				AltLineId: t.AltLineId,
				Points:    t.Points,
				Over:      t.Over,
				Under:     t.Under,
			})
		}
		for _, sp := range p.GetSpreads() {
			//s.log.Info(*p.LineId, sp.GetAltLineId(), sp.GetHdp(), sp.GetHome(), sp.GetAway())
			spreadTypes = append(spreadTypes, SpreadType{
				LineId:    p.LineId,
				AltLineId: sp.AltLineId,
				Hdp:       sp.Hdp,
				Home:      sp.Home,
				Away:      sp.Away,
			})
		}
	}

	//s.log.Info("tvpTotalType", tvpTotalType)
	if len(lineTypes) > 0 {
		tvpLineType := mssql.TVP{TypeName: "LineType", Value: lineTypes}
		_, err := s.db.ExecContext(ctx, "uspSaveLines", tvpLineType)
		if err != nil {
			s.log.Error(errors.Wrap(err, "uspSaveLines error"))
		}
	}
	if len(mlTypes) > 0 {
		tvpMlType := mssql.TVP{TypeName: "MoneylineType", Value: mlTypes}
		_, err := s.db.ExecContext(ctx, "uspSaveMoneylines", tvpMlType)
		if err != nil {
			s.log.Error(errors.Wrap(err, "uspSaveMoneylines error"))
		}
	}

	if len(totalTypes) > 0 {
		tvpTotalType := mssql.TVP{TypeName: "TotalType", Value: totalTypes}
		_, err := s.db.ExecContext(ctx, "uspSaveTotals", tvpTotalType)
		if err != nil {
			s.log.Error(errors.Wrap(err, "uspSaveTotals error"), "totalTypes: ", totalTypes)
		}
	}
	if len(spreadTypes) > 0 {
		tvpSpreadType := mssql.TVP{TypeName: "SpreadType", Value: spreadTypes}
		_, err := s.db.ExecContext(ctx, "uspSaveSpreads", tvpSpreadType)
		if err != nil {
			s.log.Error(errors.Wrap(err, "uspSaveSpreads error"),
				" spreadTypes: ", spreadTypes)
		}
	}
	return nil
}

type SpecialLineType struct {
	Id        *int64
	SpecialId *int64
	LineId    *int64
	MaxBet    *float64
	Price     *float64
	Handicap  *float64
}

func (s *Store) SaveSpecialLines(ctx context.Context, leagueId int, specials []api.SpecialOddsSpecial) error {
	var specialLineTypes []SpecialLineType

	for _, special := range specials {
		for _, line := range special.GetContestantLines() {
			//s.log.Info(*special.Id, *line.Id, *line.LineId)
			specialLineTypes = append(specialLineTypes, SpecialLineType{
				Id:        line.Id,
				SpecialId: special.Id,
				LineId:    line.LineId,
				MaxBet:    special.MaxBet,
				Price:     line.Price,
				Handicap:  line.Handicap,
			})
		}
	}
	if len(specialLineTypes) > 0 {
		tvpSpecialLineType := mssql.TVP{TypeName: "SpecialLineType", Value: specialLineTypes}
		_, err := s.db.ExecContext(ctx, "uspSaveSpecialLines", tvpSpecialLineType)
		if err != nil {
			return errors.Wrap(err, "uspSaveSpecialLines error")
		}
	}
	return nil
}

func (s *Store) CreateSport(ctx context.Context, sport api.Sport) error {
	_, err := s.db.ExecContext(ctx, "uspCreateSport", sport.Id, sport.Name, sport.HasOfferings,
		sport.LeagueSpecialsCount, sport.EventSpecialsCount, sport.EventCount)

	if err != nil {
		return errors.Wrap(err, "")
	}
	//s.log.Info(res)
	return nil
}

func NewStore(cfg *config.Config, log *zap.SugaredLogger, db *sqlx.DB) *Store {
	cache, err := ristretto.NewCache(&ristretto.Config{
		NumCounters: 10000000,
		MaxCost:     100000000,
		BufferItems: 64,
		Metrics:     false,
	})
	if err != nil {
		panic(err)
	}
	return &Store{log: log, db: db, Cache: cache}
}

func (s *Store) Close() {
	s.log.Info("closing db and Cache")
	err := s.db.Close()
	if err != nil {
		s.log.Fatal(err)
	}
	s.Cache.Close()
}

func (s *Store) FindSportIdByName(ctx context.Context, line *pb.SurebetSide) error {
	err := s.db.QueryRowContext(ctx, "select dbo.fnFindSportIdByName(@SportName)", sql.Named("SportName", line.SportName)).Scan(&line.SportId)
	if err != nil {
		return err
	}
	return nil
}

func (s *Store) FindLeagueIdByName(ctx context.Context, line *pb.SurebetSide) error {
	err := s.db.QueryRowContext(ctx, "select dbo.fnFindLeagueIdByName(@Name, @SportId)",
		sql.Named("Name", line.LeagueName),
		sql.Named("SportId", line.SportId)).Scan(&line.LeagueId)
	if err != nil {
		return errors.Wrap(err, "dbo.fnFindLeagueIdByName error")
	}
	return nil
}

func (s *Store) FindEvent(ctx context.Context, line *pb.SurebetSide) error {
	var id *int64
	err := s.db.QueryRowContext(ctx, "select dbo.fnFindEvent(@Starts, @Home, @Away, @LeagueId)",
		//sql.Named("Starts", line.Starts),
		sql.Named("Home", line.Home),
		sql.Named("Away", line.Away),
		sql.Named("LeagueId", line.LeagueId)).Scan(&id)
	s.log.Info(id)
	if err != nil {
		return errors.Wrap(err, "dbo.fnFindEvent error")
	}
	if id == nil {

	}
	return nil
}

func (s *Store) SaveBalance(ctx context.Context, b api.ClientBalanceResponse, accountId int) error {
	_, err := s.db.ExecContext(ctx, "uspSaveBalance", accountId, b.AvailableBalance, b.OutstandingTransactions,
		b.GivenCredit, b.Currency)
	if err != nil {
		return errors.Wrap(err, "uspSaveBalance error")
	}
	return nil

}

type Account struct {
	Id           int
	AccountType  string
	CurrencyCode string
	ServiceName  string
	Username     string
	Password     string
	Commission   float64
	Share        float64
}

func (s *Store) GetAccount(ctx context.Context) (Account, error) {
	a := Account{}
	err := s.db.QueryRowContext(ctx, getAccount).Scan(&a.Id, &a.AccountType, &a.CurrencyCode, &a.ServiceName, &a.Username, &a.Password, &a.Commission, &a.Share)
	if err != nil {
		return Account{}, errors.Wrap(err, "get account error")
	}
	return a, nil
}

func (s *Store) SaveCheck(sb *pb.Surebet) error {
	side := sb.Members[0]
	_, err := s.db.Exec("dbo.uspSaveSide",
		sql.Named("Id", sb.SurebetId),
		sql.Named("SideIndex", side.Num-1),

		sql.Named("SportName", side.SportName),
		sql.Named("SportId", side.SportId),
		sql.Named("LeagueName", side.LeagueName),
		sql.Named("LeagueId", side.LeagueId),
		sql.Named("Home", side.Home),
		sql.Named("HomeId", side.HomeId),
		sql.Named("Away", side.Away),
		sql.Named("AwayId", side.AwayId),
		sql.Named("MarketName", side.MarketName),
		sql.Named("MarketId", side.MarketId),
		sql.Named("Price", side.Price),
		sql.Named("Initiator", side.Initiator),
		sql.Named("Starts", sb.Starts),
		sql.Named("EventId", side.EventId),

		sql.Named("CheckId", side.GetCheck().GetId()),
		sql.Named("AccountId", side.GetCheck().GetAccountId()),
		sql.Named("CheckPrice", side.GetCheck().GetPrice()),
		sql.Named("CheckStatus", side.GetCheck().GetStatus()),
		sql.Named("CountLine", side.GetCheck().GetCountLine()),
		sql.Named("CountEvent", side.GetCheck().GetCountEvent()),
		sql.Named("AmountEvent", side.GetCheck().GetAmountEvent()),
		sql.Named("MinBet", side.GetCheck().GetMinBet()),
		sql.Named("MaxBet", side.GetCheck().GetMaxBet()),
		sql.Named("Balance", side.GetCheck().GetBalance()),
		sql.Named("Currency", side.GetCheck().GetCurrency()),
		sql.Named("CheckDone", side.GetCheck().GetDone()),

		sql.Named("CalcStatus", side.GetCheckCalc().GetStatus()),
		sql.Named("MaxStake", side.GetCheckCalc().GetMaxStake()),
		sql.Named("MinStake", side.GetCheckCalc().GetMinStake()),
		sql.Named("MaxWin", side.GetCheckCalc().GetMaxWin()),
		sql.Named("Stake", side.GetCheckCalc().GetStake()),
		sql.Named("Win", side.GetCheckCalc().GetWin()),
		sql.Named("IsFirst", side.GetCheckCalc().GetIsFirst()),
	)
	if err != nil {
		return errors.Wrapf(err, "uspSaveSide error")
	}
	return nil
}

func (s *Store) SaveBet(sb *pb.Surebet) error {
	side := sb.Members[0]
	_, err := s.db.Exec("dbo.uspSaveBet",
		sql.Named("SurebetId", sb.SurebetId),
		sql.Named("SideIndex", side.Num-1),

		sql.Named("BetId", side.ToBet.Id),
		sql.Named("TryCount", side.GetToBet().GetTryCount()),
		sql.Named("Status", side.GetBet().GetStatus()),
		sql.Named("StatusInfo", side.GetBet().GetStatusInfo()),
		sql.Named("Start", side.GetBet().GetStart()),
		sql.Named("Done", side.GetBet().GetDone()),
		sql.Named("Price", side.GetBet().GetPrice()),
		sql.Named("Stake", side.GetBet().GetStake()),
		sql.Named("ApiBetId", side.GetBet().GetApiBetId()),
	)
	if err != nil {
		return errors.Wrap(err, "uspSaveBet error")
	}
	return nil
}

type Stat struct {
	MarketName  string
	CountEvent  int64
	CountLine   int64
	AmountEvent int64
	AmountLine  int64
}

func (s *Store) GetStat(side *pb.SurebetSide) error {
	var stat []Stat
	err := s.db.Select(&stat, "dbo.uspCalcStat", sql.Named("EventId", side.EventId))
	if err == sql.ErrNoRows {
		return nil
	} else if err != nil {
		return errors.Wrap(err, "uspCalcStat error")
	} else {
		for i := range stat {
			side.Check.AmountEvent = stat[i].AmountEvent
			side.Check.CountEvent = stat[i].CountEvent
			if stat[i].MarketName == side.MarketName {
				side.Check.CountLine = stat[i].CountLine
				side.Check.AmountLine = stat[i].AmountLine
				return nil
			}
		}
	}
	return nil
}

//type StraightBetV3 struct {
//	BetId int64 `json:"betId"`
//	WagerNumber int `json:"wagerNumber"`
//	PlacedAt time.Time `json:"placedAt"`
//	BetStatus string `json:"betStatus"`
//	BetType string `json:"betType"`
//	Win float64 `json:"win"`
//	Risk float64 `json:"risk"`
//	WinLoss *float64 `json:"winLoss,omitempty"`
//	OddsFormat OddsFormat `json:"oddsFormat"`
//	CustomerCommission *float64 `json:"customerCommission,omitempty"`
//	UpdateSequence int64 `json:"updateSequence"`
//	SportId *int `json:"sportId,omitempty"`
//	LeagueId *int `json:"leagueId,omitempty"`
//	EventId *int64 `json:"eventId,omitempty"`
//	Handicap *float64 `json:"handicap,omitempty"`
//	Price *float64 `json:"price,omitempty"`
//	TeamName *string `json:"teamName,omitempty"`
//	Side *string `json:"side,omitempty"`
//	Pitcher1 *string `json:"pitcher1,omitempty"`
//	Pitcher2 *string `json:"pitcher2,omitempty"`
//	Pitcher1MustStart *bool `json:"pitcher1MustStart,omitempty"`
//	Pitcher2MustStart *bool `json:"pitcher2MustStart,omitempty"`
//	Team1 *string `json:"team1,omitempty"`
//	Team2 *string `json:"team2,omitempty"`
//	PeriodNumber *int `json:"periodNumber,omitempty"`
//	Team1Score *float64 `json:"team1Score,omitempty"`
//	Team2Score *float64 `json:"team2Score,omitempty"`
//	FtTeam1Score *float64 `json:"ftTeam1Score,omitempty"`
//	FtTeam2Score *float64 `json:"ftTeam2Score,omitempty"`
//	PTeam1Score *float64 `json:"pTeam1Score,omitempty"`
//	PTeam2Score *float64 `json:"pTeam2Score,omitempty"`
//	EventStartTime *time.Time `json:"eventStartTime,omitempty"`
//}

func (s *Store) SaveBetList(bets []api.StraightBetV3) error {
	if len(bets) == 0 {
		return nil
	}
	tvp := mssql.TVP{TypeName: "StraightBetListType", Value: bets}
	_, err := s.db.Exec("uspStraightBetList", tvp)
	if err != nil {
		return errors.Wrap(err, "uspStraightBetList error")
	}
	return nil
}

func (s *Store) GetResults(ctx context.Context) ([]pb.BetResult, error) {
	var res []pb.BetResult
	rows, err := s.db.QueryxContext(ctx, "select * from dbo.GetResults")
	if err != nil {
		return nil, errors.Wrap(err, "get bet results error")
	}
	for rows.Next() {
		var r pb.BetResult
		var Price, Stake, WinLoss *float64
		var ApiBetId, ApiBetStatus *string
		err := rows.Scan(&r.SurebetId, &r.SideIndex, &r.BetId, &ApiBetId, &ApiBetStatus, &Price, &Stake, &WinLoss)
		if err != nil {
			s.log.Error(err)
			continue
		}
		if ApiBetId != nil {
			r.ApiBetId = *ApiBetId
		}
		if ApiBetStatus != nil {
			r.ApiBetStatus = *ApiBetStatus
		}
		if Price != nil {
			r.Price = *Price
		}
		if Stake != nil {
			r.Stake = *Stake
		}
		if WinLoss != nil {
			r.WinLoss = *WinLoss
		}
		res = append(res, r)
	}
	return res, nil
}

func (s *Store) GetEvent(eventId string) (event api.Fixture, err error) {
	err = s.db.Get(&event, "select Id, ParentId, Starts, Home, Away, RotNum, LiveStatus, HomePitcher, AwayPitcher, ResultingUnit from dbo.Event where Id = @Id", sql.Named("Id", eventId))
	return
}

type Token struct {
	Session   string
	ApiKey    string
	Device    string
	TrustCode string
	Id        int64
}

func (s *Store) LoadToken(ctx context.Context) (Token, error) {
	var t Token
	err := s.db.GetContext(ctx, &t, "select Session, ApiKey, Device,TrustCode,Id from dbo.Auth where Id=1")
	return t, err
}

func (s *Store) SaveToken(ctx context.Context, token Token) error {
	_, err := s.db.ExecContext(ctx, "UPDATE dbo.Auth SET Session=@Session, ApiKey=@ApiKey, Device=@Device, TrustCode=@TrustCode where Id=1",
		sql.Named("Session", token.Session),
		sql.Named("ApiKey", token.ApiKey),
		sql.Named("Device", token.Device),
		sql.Named("TrustCode", token.TrustCode),
	)
	return err
}

const spreadQ = `
with cte as (
    select top 1 l.LineId
    from Event e
             join Line L on e.Id = L.EventId
             join Spread S on l.LineId = S.LineId
    where e.Id = @p1
      and l.Number = @p2
    group by l.LineId
    order by l.LineId desc
)
select s.LineId, s.AltLineId, Hdp, Home, Away
from Spread s
         join cte on cte.LineId = s.LineId
`

type Spread struct {
	LineId    int64
	AltLineId *int64
	Hdp       float64
	Home      float64
	Away      float64
}

func (s *Store) SelectSpreads(ctx context.Context, eventId string, period int) (spreads []Spread, err error) {
	err = s.db.SelectContext(ctx, &spreads, spreadQ, eventId, period)
	return
}

type DemoEvent struct {
	Id            string
	Home          string
	Away          string
	SportName     string
	SportId       string
	LeagueName    string
	LeagueId      string
	ResultingUnit string
	Starts        string
}

const CurrentEventsQ = `
select top (@p2)
               e.Id, 
                 e.Home, 
                 e.Away, 
                 e.ResultingUnit, 
                 e.Starts,
                 s.Name SportName,
                 s.Id SportId,
                 l.Name LeagueName,
                 l.Id LeagueId
from Event e
join League l on e.LeagueId = l.Id
join Sport s on l.SportId=s.Id
where l.EventCount > 0
  and l.SportId = @p1
  and e.LiveStatus in (0, 2)
and e.Starts > sysdatetimeoffset()
and e.ResultingUnit = @p3
and LiveStatus=0
and DATEDIFF(minute, sysdatetimeoffset(), e.Starts) < 900
order by (select count(Line.LineId) from Line where Line.EventId=e.Id) desc
`

func (s *Store) SelectCurrentEvents(ctx context.Context, sportId int64, count int64, eventType string) (events []DemoEvent, err error) {
	err = s.db.SelectContext(ctx, &events, CurrentEventsQ, sportId, count, eventType)
	return
}

type Total struct {
	LineId    int64
	AltLineId *int64
	Points    float64
	Over      float64
	Under     float64
}

const totalQ = `
with cte as (
    select top 1 l.LineId
    from Event e
             join Line L on e.Id = L.EventId
             join Total t on l.LineId = t.LineId
    where e.Id = @p1
      and l.Number = @p2
    group by l.LineId
    order by l.LineId desc
)
select t.LineId, t.AltLineId, Points, [Over], Under
from Total t
         join cte on cte.LineId = t.LineId
`

func (s *Store) SelectTotals(ctx context.Context, eventId string, period int) (totals []Total, err error) {
	err = s.db.SelectContext(ctx, &totals, totalQ, eventId, period)
	return
}

const DeleteOldLinesSQL = `
delete l
from Line l
         join Event e on l.EventId = e.Id
where DATEDIFF(day, Starts, sysdatetimeoffset()) > 7
`

func (s *Store) DeleteOldLines(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, DeleteOldLinesSQL)
	return err
}

const DeleteOldMoneylinesSQL = `
delete m 
from Moneyline m
where DATEDIFF(day, UpdatedAt, sysdatetimeoffset()) > 14
`

func (s *Store) DeleteOldMoneylines(ctx context.Context) error {
	_, err := s.db.ExecContext(ctx, DeleteOldMoneylinesSQL)
	return err
}

//func (s *Store) CreateSport(ctx context.Context, sportName string, serviceId int) (int, error) {
//	key := s.FormKey("Sport", sportName, strconv.Itoa(serviceId))
//	id, b := s.CheckInCache(ctx, key)
//	if b {
//		return id, nil
//	}
//	err := s.db.QueryRowContext(ctx, "uspCreateSport", &sportName, &serviceId).Scan(&id)
//	if err != nil {
//		return 0, errors.Wrap(err, "uspCreateSport error")
//	}
//	s.Cache.Set(key, id, 1)
//	return id, nil
//}
