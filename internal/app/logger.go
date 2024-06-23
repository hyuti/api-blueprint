package app

import (
	"errors"
	"github.com/hyuti/api-blueprint/pkg/logger"
	"golang.org/x/exp/slog"
)

// TODO: serviceKey should be adjusted to be relevant to the project
const serviceKey = "service-name"

var ErrLoggerEmpty = errors.New("logger must not be empty")

func WithLogger() (*slog.Logger, error) {
	if app.cfg == nil {
		return nil, ErrCfgEmpty
	}
	l := logger.FileAndStdLogger(
		app.cfg.PathLogger,
		// TODO: log level should be adjusted in Production deployment
		logger.WithLevelOpt(slog.LevelDebug),
	)
	l = logger.WithServiceName(l, serviceKey, app.cfg.Name)
	l = logger.WithCtxID(l)
	return l, nil
}
func Logger() *slog.Logger {
	mutex.Lock()
	defer mutex.Unlock()
	return app.logger
}
