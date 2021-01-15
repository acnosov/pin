package store

import (
	"context"
	pb "github.com/aibotsoft/gen/fortedpb"
	api "github.com/aibotsoft/gen/pinapi"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/sqlserver"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var s *Store

func TestMain(m *testing.M) {
	cfg := config.New()
	log := logger.New()
	db := sqlserver.MustConnectX(cfg)
	s = NewStore(cfg, log, db)
	m.Run()
	s.Close()
}

func SportHelper(t *testing.T) api.Sport {
	t.Helper()
	sport := api.Sport{}
	sport.SetId(666)
	sport.SetName("TestSport")
	sport.SetHasOfferings(true)
	sport.SetLeagueSpecialsCount(66)
	sport.SetEventSpecialsCount(66)
	sport.SetEventCount(66)
	return sport
}
func SportListHelper(count int) []api.Sport {
	sports := make([]api.Sport, count)
	for i := 0; i < count; i++ {
		s := api.Sport{}
		s.SetId(i)
		s.SetName("TestSport")
		s.SetHasOfferings(true)
		s.SetLeagueSpecialsCount(i)
		s.SetEventSpecialsCount(i)
		s.SetEventCount(i)
		sports[i] = s
	}
	return sports
}

func BenchmarkStore_CreateSportTVP(b *testing.B) {
	b.ReportAllocs()
	sl := SportListHelper(10000)
	_, err := s.db.Exec("TRUNCATE TABLE dbo.Sport")
	b.Log("table truncated", "create: ", len(sl))
	assert.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		err := s.SaveSports(context.Background(), sl)
		if err != nil {
			b.Log(err)
		}
	}
}

func BenchmarkStore_CreateSport(b *testing.B) {
	b.ReportAllocs()
	sl := SportListHelper(1000)
	_, err := s.db.Exec("TRUNCATE TABLE dbo.Sport")
	b.Log("table truncated", "create: ", len(sl))
	assert.NoError(b, err)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for _, sport := range sl {
			err := s.CreateSport(context.Background(), sport)
			if err != nil {
				b.Log(err)
			}
		}
	}
}

func TestStore_CreateSport(t *testing.T) {
	sports := SportListHelper(10)
	t.Log(sports)
	err := s.SaveSports(context.Background(), sports)
	assert.NoError(t, err)
	time.Sleep(time.Millisecond)
	err = s.SaveSports(context.Background(), sports)
	assert.NoError(t, err)
}

func TestStore_IsInCache(t *testing.T) {
	sport := SportHelper(t)
	got := s.IsInCache(context.Background(), sport)
	assert.False(t, got, got)
	time.Sleep(time.Millisecond)
	fromCache := s.IsInCache(context.Background(), sport)
	assert.True(t, fromCache)
}

func TestStore_FindSportIdByName(t *testing.T) {
	l := &pb.SurebetSide{
		SportName: "Badminton",
	}
	err := s.FindSportIdByName(context.Background(), l)
	if assert.NoError(t, err) {
		assert.NotEmpty(t, l.SportId)
	}
}

func TestStore_GetAccount(t *testing.T) {
	got, err := s.GetAccount(context.Background())
	if assert.NoError(t, err) {
		assert.NotEmpty(t, got)
	}
}

func TestStore_GetStat(t *testing.T) {
	side := &pb.SurebetSide{
		EventId:    "1167772745",
		MarketName: "ТБ(2,6)",
		Check:      &pb.Check{},
	}

	err := s.GetStat(side)
	if assert.NoError(t, err) {
		t.Log(side.Check)
	}
}

func TestStore_GetResults(t *testing.T) {
	got, err := s.GetResults(context.Background())
	if assert.NoError(t, err) {
		assert.NotEmpty(t, got)
		t.Log(got)
	}
}

func TestStore_GetEvent(t *testing.T) {
	got, err := s.GetEvent("1122239544")
	if assert.NoError(t, err) {
		assert.NotEmpty(t, got)
		t.Log(got)
	}
}

func TestStore_LoadToken(t *testing.T) {
	got, err := s.LoadToken(context.Background())
	if assert.NoError(t, err) {
		assert.NotEmpty(t, got)
		t.Log(got)
	}
}

func TestStore_DeleteOldLines(t *testing.T) {
	err := s.DeleteOldLines(context.Background())
	assert.NoError(t, err)
}
