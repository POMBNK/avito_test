package config

import (
	"github.com/POMBNK/avito_test_task/pkg/logger"
	"github.com/ilyakaznacheev/cleanenv"

	"sync"
)

type Config struct {
	Storage struct {
		Postgresql Postgresql `json:"postgresql"`
	} `yaml:"storage"`
}

type Postgresql struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Database string `json:"database"`
	User     string `json:"user"`
	Password string `json:"password"`
}

var once sync.Once
var cfg = &Config{}

func GetCfg() *Config {
	once.Do(func() {
		logs := logger.GetLogger()
		logs.Info("Reading config...")
		err := cleanenv.ReadConfig("config.yml", cfg)
		if err != nil {
			help, _ := cleanenv.GetDescription(cfg, nil)
			logs.Info(help)
			logs.Fatal(err)
		}
	})
	return cfg
}
