package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hyuti/api-blueprint/pkg/http"
)

func WithGinEngine() (*gin.Engine, error) {
	if app.cfg == nil {
		return nil, ErrCfgEmpty
	}
	if app.logger == nil {
		return nil, ErrLoggerEmpty
	}
	restful := http.New(&http.Cfg{
		Debug:  app.cfg.Debug,
		Logger: app.logger,
		Name:   app.cfg.Name,
	})
	if err := http.DefaultTranslation(); err != nil {
		return nil, fmt.Errorf("cannot init gin: %w", err)
	}
	return restful, nil
}
func Gin() *gin.Engine {
	mutex.Lock()
	defer mutex.Unlock()
	return app.ginEngine
}
