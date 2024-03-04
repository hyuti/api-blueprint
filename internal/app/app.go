package app

import (
	"errors"
	"fmt"
	els "github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/hyuti/API-Golang-Template/config"
	"github.com/hyuti/API-Golang-Template/pkg/elasticsearch"
	pkgGprc "github.com/hyuti/API-Golang-Template/pkg/grpc"
	"github.com/hyuti/API-Golang-Template/pkg/http"
	"github.com/hyuti/API-Golang-Template/pkg/logger"
	"github.com/hyuti/API-Golang-Template/pkg/telegram"
	"golang.org/x/exp/slog"
	"sync"
)

const serviceKey = "service-name"

var (
	mutex sync.Mutex
	app   *App
)

var (
	ErrLoggerEmpty = errors.New("logger expected not to be empty")
	ErrCfgEmpty    = errors.New("config expected not to be empty")
)

type App struct {
	cfg       *config.Config
	logger    *slog.Logger
	elk       *els.TypedClient
	ginEngine *gin.Engine
	grpcSrv   *pkgGprc.Server
	tele      *telegram.Tele
}

func Init() error {
	mutex.Lock()
	defer mutex.Unlock()
	if app != nil {
		return nil
	}
	app = new(App)
	if err := WithCfg(); err != nil {
		return err
	}
	if err := WithLogger(); err != nil {
		return err
	}
	if err := WithGinEngine(); err != nil {
		return err
	}
	if err := WithGRPCServer(); err != nil {
		return err
	}
	if err := WithTele(); err != nil {
		return err
	}
	if err := WithEls(); err != nil {
		return err
	}
	if _, err := pkgGprc.DefaultValidator(); err != nil {
		return fmt.Errorf("unable to init validator: %w", err)
	}
	return nil
}

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
		return fmt.Errorf("unable to init tele: %w", err)
	}
	return nil
}
func WithGRPCServer() error {
	grpcS, err := pkgGprc.New(
		&pkgGprc.SrvConfig{
			Port: app.cfg.GRPC.Port,
		},
	)
	if err != nil {
		return fmt.Errorf("unable to init grpc: %w", err)
	}
	app.grpcSrv = grpcS
	return nil
}

func WithCfg() error {
	cfg, err := config.GetConfig()
	if err != nil {
		return fmt.Errorf("unable to init config: %w", err)
	}
	app.cfg = cfg
	return nil
}

func WithLogger() error {
	if app.cfg == nil {
		return ErrCfgEmpty
	}
	app.logger = logger.FileAndStdLogger(
		app.cfg.PathLogger,
		logger.WithLevelOpt(slog.LevelDebug),
	)
	app.logger = logger.WithServiceName(app.logger, serviceKey, app.cfg.Name)
	app.logger = logger.WithCtxID(app.logger)
	return nil
}

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
		return fmt.Errorf("unable to init gin: %w", err)
	}
	return nil
}
func WithEls() error {
	if app.logger == nil {
		return ErrLoggerEmpty
	}
	if app.cfg == nil {
		return ErrCfgEmpty
	}
	es, err := elasticsearch.New(&elasticsearch.Cfg{
		Addresses: []string{app.cfg.Elastic.URL},
		Username:  app.cfg.Elastic.Username,
		Password:  app.cfg.Elastic.Password,
	})
	if err != nil {
		return fmt.Errorf("unable to init els: %w", err)
	}
	app.elk = es
	return nil
}
func Cfg() *config.Config {
	mutex.Lock()
	defer mutex.Unlock()
	return app.cfg
}
func Logger() *slog.Logger {
	mutex.Lock()
	defer mutex.Unlock()
	return app.logger
}
func Gin() *gin.Engine {
	mutex.Lock()
	defer mutex.Unlock()
	return app.ginEngine
}
func GrpcSrv() *pkgGprc.Server {
	mutex.Lock()
	defer mutex.Unlock()
	return app.grpcSrv
}
func Tele() *telegram.Tele {
	mutex.Lock()
	defer mutex.Unlock()
	return app.tele
}
func Els() *els.TypedClient {
	mutex.Lock()
	defer mutex.Unlock()
	return app.elk
}
