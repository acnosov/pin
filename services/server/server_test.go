package server

import (
	"context"
	"github.com/aibotsoft/gen/surebetpb"
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/sqlserver"
	"github.com/stretchr/testify/assert"
	"testing"
)

func InitServerHelper(t *testing.T) *Server {
	t.Helper()
	cfg := config.New()
	log := logger.New()
	db := sqlserver.MustConnect(cfg)
	s := NewServer(cfg, log, db)
	return s
}

func TestServer_CheckLine(t *testing.T) {
	s := InitServerHelper(t)
	defer s.Close()
	ctx := context.Background()
	line := &surebetpb.Line{
		Num:         1,
		ServiceName: "Pinnacle",
		SportName:   "Table Tennis",
		LeagueName:  "Setka Cup",
		Home:        "Vasil Smyk",
		Away:        "Vladimir Melnikov",
		MarketName:  "П1",
		Price:       2.1,
		Url:         "https://members.pinnacle.com/Sportsbook/Mobile/ru-RU/Enhanced/Regular/SportsBookAll/35/Curacao/Odds/Table Tennis-32/Market/2/208761/1119354436",
		Initiator:   false,
	}
	checkLine, err := s.CheckLine(ctx, &surebetpb.CheckLineRequest{
		Line: line,
	})
	if assert.NoError(t, err) {
		assert.Equal(t, 2, checkLine.Price)
	}

}
