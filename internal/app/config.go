package app

import (
	"errors"
	"fmt"
	"github.com/hyuti/api-blueprint/config"
)

var ErrCfgEmpty = errors.New("config must not be empty")

func WithCfg() (*config.Config, error) {
	cfg, err := config.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("cannot init config: %w", err)
	}
	return cfg, nil
}
func Cfg() *config.Config {
	mutex.Lock()
	defer mutex.Unlock()
	return app.cfg
}
