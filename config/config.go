package config

import (
	logger "dcrcs-go/utils"
	"os"
	"path/filepath"
	"runtime"

	"github.com/goccy/go-yaml"
)

type Config struct {
	HttpConfig  HttpConfig `json:"httpconfig"`
	AgentConfig Agent      `json:"agentconfig"`
}

type HttpConfig struct {
	Ip   string `json:"ip"`
	Port int64  `json:"port"`
}

type Agent struct {
	Address string `json:"address"`
}

func InitConfig() *Config {
	_, fileName, _, ok := runtime.Caller(0)
	if !ok {
		logger.Error("read config error")
		panic("read config error")
	}
	configPath := filepath.Join(filepath.Dir(fileName), "config.yaml")
	configByte, err := os.ReadFile(configPath)
	if err != nil {
		logger.Error("read confg error")
		panic("read config error")
	}
	config := Config{}
	err = yaml.Unmarshal(configByte, &config)
	return &config
}
