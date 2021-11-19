package main

import (
	"fmt"
	"github.com/aibotsoft/pin/pkg/config"
	"github.com/aibotsoft/pin/pkg/sqlserver"
	"github.com/aibotsoft/pin/pkg/store"
	"github.com/aibotsoft/pin/services/auth"
	"github.com/aibotsoft/pin/services/handler"
	"github.com/aibotsoft/pin/services/server"
	"github.com/vrischmann/envconfig"
	"go.uber.org/zap"
	"os"
	"os/signal"
	"syscall"
)

const Version = 0.1

func main() {
	os.Setenv("SERVICE_NAME", "pin-service")
	os.Setenv("MSSQL_DATABASE", "PinServiceDev")
	os.Setenv("GRPC_PORT", "50053")

	logger, err := zap.NewDevelopment()
	if err != nil {
		panic(err)
	}
	log := logger.Sugar()

	cfg := &config.Config{}
	err = envconfig.Init(cfg)
	if err != nil {
		panic(err)
	}

	log.Infow("Begin service", "version", Version, "config", cfg)
	db := sqlserver.MustConnect(cfg)
	sto := store.NewStore(cfg, log, db)

	au := auth.New(cfg, log, sto)
	go au.AuthJob()
	h := handler.NewHandler(cfg, log, sto, au)
	//go h.BalanceJob()
	//go h.BetStatusJob()

	s := server.NewServer(cfg, log, h)
	// Инициализируем Close
	errc := make(chan error)
	go func() {
		c := make(chan os.Signal)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	//c := collector.NewCollector(cfg, log, sto)
	//go c.CollectJob()

	go func() {
		errc <- s.Serve()
	}()
	defer func() {
		s.Close()
	}()
	log.Info("exit: ", <-errc)
}
