package app

import (
	"fmt"
	"github.com/hyuti/api-blueprint/config"
)

func WithCfg() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("cannot init config: %w", err)
	}
	app.cfg = cfg
	return nil
}
func Cfg() *config.Config {
	mutex.Lock()
	defer mutex.Unlock()
	return app.cfg
}
