package main

import (
	"github.com/aibotsoft/micro/config"
	"github.com/aibotsoft/micro/logger"
	"github.com/aibotsoft/micro/sqlserver"
	"github.com/aibotsoft/pin/services/server"
)

func main() {
	cfg := config.New()
	log := logger.New()
	log.Infow("Begin service", "config", cfg)
	db := sqlserver.MustConnect(cfg)

	s := server.NewServer(cfg, log, db)
	log.Info(s)
	//auth := context.WithValue(context.Background(), api.ContextBasicAuth, api.BasicAuth{
	//	UserName: "ON868133",
	//	Password: "@RjE5d4Q",
	//})
	//clientConfig := api.NewConfiguration()
	////cfg.AddDefaultHeader("testheader", "testvalue")
	//////cfg.Host = testHost
	//////cfg.Scheme = testScheme
	//client := api.NewAPIClient(clientConfig)
	////log.Println(auth)
	////log.Println(client)
	//sports, response, err := client.OthersApi.SportsV2Get(auth).Execute()
	//if err != nil {
	//	log.Fatal(err)
	//}
	//log.Info(sports.GetSports())
	//log.Info(response)
}
