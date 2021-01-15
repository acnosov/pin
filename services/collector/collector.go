package collector

import (
	"context"
	api "github.com/aibotsoft/gen/pinapi"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/pin/pkg/client"
	"github.com/aibotsoft/pin/pkg/store"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"time"
)

const sportMinPeriodSeconds = 120
const leagueMinPeriodSeconds = 120
const eventMinPeriodSeconds = 60
const periodMinPeriodSeconds = 60 * 60 * 24 * 7
const betResultMinPeriodSeconds = 5 * 60
const collectMinPeriod = 2 * time.Second

var stopSportId = map[int]bool{12: true, 58: true, 57: true, 24: true}
var leagueMinPeriod = make(map[int]time.Time)
var eventMinPeriod = make(map[int]time.Time)
var periodMinPeriod = make(map[int]time.Time)

var betResultLast time.Time
var eventLast = make(map[int]int64)
var specialLast = make(map[int]int64)
var lineLast = make(map[int]int64)
var specialLineLast = make(map[int]int64)

type Collector struct {
	cfg       *config.Config
	log       *zap.SugaredLogger
	client    *client.Client
	store     *store.Store
	sports    []api.Sport
	sportLast time.Time
	account   store.Account
}

func (c *Collector) AccId() int {
	return c.account.Id
}
func (c *Collector) CollectJob() {
	for {
		start := time.Now()
		//ctx, cancel := context.WithTimeout(context.Background(), 6*time.Minute)
		err := c.CollectRound(context.Background())
		//cancel()
		if err != nil {
			c.log.Info(err)
			time.Sleep(time.Minute)
		} else {
			c.log.Infow("collect_job_done", "time", time.Since(start))
			time.Sleep(time.Second * 10)
		}
	}
}

func (c *Collector) CollectRound(ctx context.Context) error {
	sports, err := c.Sports(ctx)
	if err != nil {
		return err
	}
	err = c.BetList(ctx)
	c.errLogAndSleep(err)

	for _, sport := range sports {
		if stopSportId[sport.GetId()] {
			continue
		}
		if sport.GetEventCount()+sport.GetEventSpecialsCount() == 0 {
			continue
		}
		err := c.Leagues(ctx, sport)
		c.errLogAndSleep(err)

		err = c.Periods(ctx, sport)
		c.errLogAndSleep(err)

		err = c.Events(ctx, sport)
		c.errLogAndSleep(err)

		//err = c.Specials(ctx, sport)
		//c.errLogAndSleep(err)

		err = c.Lines(ctx, sport)
		c.errLogAndSleep(err)

		//err = c.SpecialLines(ctx, sport)
		//c.errLogAndSleep(err)
	}
	return nil
}

func (c *Collector) errLogAndSleep(err error) {
	if err != nil {
		c.log.Info(err)
	}
	time.Sleep(collectMinPeriod)
}

func (c *Collector) Sports(ctx context.Context) ([]api.Sport, error) {
	var err error
	if time.Since(c.sportLast).Seconds() < sportMinPeriodSeconds {
		return c.sports, nil
	}
	c.sports, err = c.client.GetSports(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "client.GetSports error")
	}
	err = c.store.SaveSports(ctx, c.sports)
	if err != nil {
		return nil, errors.Wrap(err, "store.SaveSports error")
	}
	c.sportLast = time.Now()
	return c.sports, nil
}

func (c *Collector) Leagues(ctx context.Context, sport api.Sport) error {
	if time.Since(leagueMinPeriod[sport.GetId()]).Seconds() < leagueMinPeriodSeconds {
		return nil
	}
	leagues, err := c.client.GetLeagues(ctx, sport.GetId())
	if err != nil {
		return err
	}
	//start := time.Now()
	err = c.store.SaveLeagues(ctx, sport.GetId(), leagues)
	//c.log.Infow("saveLeague", "time", time.Since(start).Milliseconds(), "sportId", sport.GetId(), "count", len(leagues))
	if err != nil {
		return err
	}
	leagueMinPeriod[sport.GetId()] = time.Now()
	return nil
}
func (c *Collector) Periods(ctx context.Context, sport api.Sport) error {
	if time.Since(periodMinPeriod[sport.GetId()]).Seconds() < periodMinPeriodSeconds {
		return nil
	}
	got, err := c.client.GetPeriods(ctx, sport.GetId())
	if err != nil {
		return err
	}
	//start := time.Now()
	err = c.store.SavePeriods(ctx, sport.GetId(), got)
	//c.log.Infow("saveLeague", "time", time.Since(start).Milliseconds(), "sportId", sport.GetId(), "count", len(leagues))
	if err != nil {
		return err
	}
	periodMinPeriod[sport.GetId()] = time.Now()
	return nil
}

func (c *Collector) Events(ctx context.Context, sport api.Sport) error {
	if time.Since(eventMinPeriod[sport.GetId()]).Seconds() < eventMinPeriodSeconds {
		return nil
	}
	if sport.GetEventCount() == 0 {
		c.log.Infow("sport has no events", "sportId", sport.GetId(), "sport.GetEventCount", sport.GetEventCount())
		return nil
	}
	events, err := c.client.GetEvents(ctx, sport.GetId(), eventLast[sport.GetId()])
	if err != nil {
		return errors.Wrap(err, "GetEvents error")
	}
	eventLast[sport.GetId()] = events.GetLast()
	for _, league := range events.GetLeague() {
		//c.log.Info("sportId:", sport.GetId(), " real event count: ", len(league.GetEvents()))
		err := c.store.SaveEvents(ctx, league.GetId(), league.GetEvents())
		if err != nil {
			c.log.Error(err)
			continue
		}
	}
	eventMinPeriod[sport.GetId()] = time.Now()
	return nil
}

func (c *Collector) Specials(ctx context.Context, sport api.Sport) error {
	if sport.GetEventSpecialsCount() == 0 {
		//c.log.Infow("sport has no special", "sportId", sport.GetId(), "EventSpecialsCount", sport.GetEventSpecialsCount())
		return nil
	}
	specials, err := c.client.GetSpecials(ctx, sport.GetId(), specialLast[sport.GetId()])
	if err != nil {
		return errors.Wrap(err, "GetSpecials error")
	}
	specialLast[sport.GetId()] = specials.GetLast()

	for _, league := range specials.GetLeagues() {
		err := c.store.SaveSpecials(ctx, league.GetId(), league.GetSpecials())
		if err != nil {
			c.log.Error(err)
			continue
		}
	}
	return nil
}

func (c *Collector) Lines(ctx context.Context, sport api.Sport) error {
	lines, err := c.client.GetLines(ctx, sport.GetId(), lineLast[sport.GetId()])
	if err != nil {
		return errors.Wrap(err, "GetLines error")
	}
	lineLast[sport.GetId()] = lines.GetLast()
	for _, league := range lines.GetLeagues() {
		for _, event := range league.GetEvents() {
			//c.log.Infow("event", "", event)
			err := c.store.SaveLines(ctx, event.GetId(), event.GetPeriods())
			if err != nil {
				c.log.Error(err)
				continue
			}
		}
	}
	return nil
}

func (c *Collector) SpecialLines(ctx context.Context, sport api.Sport) error {
	if sport.GetEventSpecialsCount() == 0 {
		//c.log.Infow("sport has no special", "sportId", sport.GetId(), "EventSpecialsCount", sport.GetEventSpecialsCount())
		return nil
	}
	lines, err := c.client.GetSpecialLines(ctx, sport.GetId(), specialLineLast[sport.GetId()])
	if err != nil {
		return errors.Wrap(err, "GetSpecialLines error")
	}
	specialLineLast[sport.GetId()] = lines.GetLast()
	for _, league := range lines.GetLeagues() {
		err := c.store.SaveSpecialLines(ctx, league.GetId(), league.GetSpecials())
		if err != nil {
			c.log.Error(err)
			continue
		}
	}
	return nil
}
func (c *Collector) BetList(ctx context.Context) error {
	if time.Since(betResultLast).Seconds() < betResultMinPeriodSeconds {
		return nil
	}

	bets, err := c.client.GetBets(ctx)
	if err != nil {
		return errors.Wrap(err, "Bets error")
	}
	err = c.store.SaveBetList(bets.GetStraightBets())
	if err != nil {
		return errors.Wrap(err, "SaveBetList error")
	}
	betResultLast = time.Now()
	return nil
}

func NewCollector(cfg *config.Config, log *zap.SugaredLogger, store *store.Store) *Collector {
	ctx := context.Background()
	account, err := store.GetAccount(ctx)
	if err != nil {
		log.Panicw("get account error", "err", err)
	}
	cli := client.NewClient(cfg, log, account.Username, account.Password)
	return &Collector{cfg: cfg, log: log, client: cli, store: store, account: account}
}

//func (c *Collector) Balance(ctx context.Context) (api.ClientBalanceResponse, error) {
//	var err error
//	if time.Since(balanceLast).Seconds() < balanceMinPeriodSeconds {
//		return api.ClientBalanceResponse{}, nil
//	}
//	balance, err := c.client.GetBalance(ctx)
//	if err != nil {
//		return api.ClientBalanceResponse{}, errors.Wrap(err, "client.GetBalance error")
//	}
//	err = c.store.SaveBalance(ctx, balance, c.AccId())
//	if err != nil {
//		return balance, errors.Wrap(err, "store.SaveBalance error")
//	}
//	balanceLast = time.Now()
//	return balance, nil
//}
