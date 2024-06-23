package app

import (
	"fmt"
	"github.com/hyuti/api-blueprint/pkg/telegram"
)

func WithTele() (*telegram.Tele, error) {
	if app.cfg == nil {
		return nil, ErrCfgEmpty
	}
	t, err := telegram.New(
		&telegram.TeleCfg{
			Token:        app.cfg.Telegram.Token,
			ChatID:       app.cfg.Telegram.ChatID,
			Debug:        app.cfg.Debug,
			FailSilently: app.cfg.Telegram.FailSilently,
		})
	if err != nil {
		return nil, fmt.Errorf("cannot to init telegram: %w", err)
	}
	return t, nil
}
func Tele() *telegram.Tele {
	mutex.Lock()
	defer mutex.Unlock()
	return app.tele
}
