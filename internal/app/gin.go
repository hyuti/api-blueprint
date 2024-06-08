package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hyuti/api-blueprint/pkg/http"
)

func WithGinEngine() error {
	if app.cfg == nil {
		return ErrCfgEmpty
	}
	if app.logger == nil {
		return ErrLoggerEmpty
	}
	app.ginEngine = http.New(&http.Cfg{
		Debug:  app.cfg.Debug,
		Logger: app.logger,
		Name:   app.cfg.Name,
	})
	if err := http.DefaultTranslation(); err != nil {
		return fmt.Errorf("cannot init gin: %w", err)
	}
	return nil
}
func Gin() *gin.Engine {
	mutex.Lock()
	defer mutex.Unlock()
	return app.ginEngine
}
