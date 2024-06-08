package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	App      `yaml:"APP"`
	GRPC     `yaml:"GRPC"`
	Telegram `yaml:"TELEGRAM"`
	Elastic  `yaml:"ELASTIC"`
	Gin      `yaml:"GIN"`
	Postgres `yaml:"POSTGRES"`
}
type Postgres struct {
	Host         string `yaml:"HOST"`
	Port         string `yaml:"PORT"`
	User         string `yaml:"USER"`
	Password     string `yaml:"PASSWORD"`
	DatabaseName string `yaml:"DATABASE_NAME"`
}

type App struct {
	Name       string `yaml:"NAME"`
	Debug      bool   `yaml:"DEBUG"`
	PathLogger string `yaml:"PATH_LOGGER"`
}
type GRPC struct {
	Port string `yaml:"PORT"`
}
type Gin struct {
	Port     string `yaml:"PORT"`
	BasePath string `yaml:"BASE_PATH"`
}
type Telegram struct {
	Token        string `yaml:"TOKEN"`
	ChatID       int64  `yaml:"CHAT_ID"`
	FailSilently bool   `yaml:"FAIL_SILENTLY"`
}
type Elastic struct {
	URL      string `yaml:"URL"`
	Username string `yaml:"USERNAME"`
	Password string `yaml:"PASSWORD"`
	Timeout  int    `yaml:"TIMEOUT_SECOND"`
}

func GetConfig() (*Config, error) {
	cfg := new(Config)
	err := cleanenv.ReadConfig("./config/config.yaml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}
	return cfg, nil
}
