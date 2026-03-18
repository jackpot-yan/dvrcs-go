package config

type Config struct {
	HttpConfig  HttpConfig
	AgentConfig Agent
}

type HttpConfig struct {
	Ip   string
	Port int64
}

type Agent struct {
	Address string
}
