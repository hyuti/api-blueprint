package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/hyuti/api-blueprint/config"
	"github.com/hyuti/api-blueprint/pkg/http"
	"golang.org/x/exp/slog"
)

func WithGinEngine(cfg *config.Config, logger *slog.Logger) (*gin.Engine, error) {
	if cfg == nil {
		return nil, ErrCfgEmpty
	}
	if logger == nil {
		return nil, ErrLoggerEmpty
	}
	restful := http.New(&http.Cfg{
		Debug:  cfg.Debug,
		Logger: logger,
		Name:   cfg.Name,
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
