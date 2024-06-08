package app

import (
	"github.com/hyuti/api-blueprint/pkg/logger"
	"golang.org/x/exp/slog"
)

func WithLogger() error {
	if app.cfg == nil {
		return ErrCfgEmpty
	}
	app.logger = logger.FileAndStdLogger(
		app.cfg.PathLogger,
		// TODO: log level should be adjusted in Production deployment
		logger.WithLevelOpt(slog.LevelDebug),
	)
	app.logger = logger.WithServiceName(app.logger, serviceKey, app.cfg.Name)
	app.logger = logger.WithCtxID(app.logger)
	return nil
}
func Logger() *slog.Logger {
	mutex.Lock()
	defer mutex.Unlock()
	return app.logger
}
