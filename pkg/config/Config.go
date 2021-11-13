package config

type Config struct {
	ServiceName   string
	KeepAlive     string `envconfig:"default=1440"`
	Log           string `envconfig:"default=0"`
	MssqlUser     string `envconfig:"default=sa"`
	MssqlHost     string `envconfig:"default=localhost"`
	MssqlPassword string `envconfig:"default=sa"`
	MssqlDatabase string
	MssqlPort     string `envconfig:"default=1433"`
	//BrokerUrl     string `envconfig:"default=nats://localhost:4222"`
	GrpcHost string `envconfig:"default=0.0.0.0"`
	GrpcPort string `envconfig:"default=50051"`
	//CurrencyCode  string `envconfig:"default=USD"`
	//DafUser       string
	//DafPass       string
	PinEsportHost string
	//DafIncapValue string
	PinDebug bool `envconfig:"default=false"`
}
