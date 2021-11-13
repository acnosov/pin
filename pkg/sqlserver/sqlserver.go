package sqlserver

import (
	"github.com/aibotsoft/pin/pkg/config"
	_ "github.com/denisenkom/go-mssqldb"
	"github.com/jmoiron/sqlx"
	"net"
	"net/url"
)

func BuildConnString(cfg *config.Config) string {
	query := url.Values{}
	query.Add("database", cfg.MssqlDatabase)
	query.Add("app name", cfg.ServiceName)
	//in seconds; 0 to disable (default is 30) //my def 1440
	query.Add("keepAlive", cfg.KeepAlive)
	//in bytes; 512 to 32767 (default is 4096)
	//Encrypted connections have a maximum packet size of 16383 bytes
	//query.Add("packet size", cfg.PacketSize)
	//logging flags (default 0/no logging, 63 for full logging)
	query.Add("log", cfg.Log)
	u := &url.URL{
		Scheme:   "sqlserver",
		User:     url.UserPassword(cfg.MssqlUser, cfg.MssqlPassword),
		Host:     net.JoinHostPort(cfg.MssqlHost, cfg.MssqlPort),
		RawQuery: query.Encode(),
	}
	return u.String()
}

func MustConnect(cfg *config.Config) *sqlx.DB {
	connString := BuildConnString(cfg)
	db, err := sqlx.Open("sqlserver", connString)
	if err != nil {
		panic(err)
	}
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	db.MapperFunc(func(s string) string { return s })
	return db
}
