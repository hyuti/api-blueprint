package app

import (
	"fmt"
	"github.com/hyuti/api-blueprint/pkg/telegram"
)

func WithTele() error {
	if app.cfg == nil {
		return ErrCfgEmpty
	}
	var err error
	app.tele, err = telegram.New(
		&telegram.TeleCfg{
			Token:        app.cfg.Telegram.Token,
			ChatID:       app.cfg.Telegram.ChatID,
			Debug:        app.cfg.Debug,
			FailSilently: app.cfg.Telegram.FailSilently,
		})
	if err != nil {
		return fmt.Errorf("cannot to init telegram: %w", err)
	}
	return nil
}
func Tele() *telegram.Tele {
	mutex.Lock()
	defer mutex.Unlock()
	return app.tele
}
