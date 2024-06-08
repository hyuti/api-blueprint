package app

import (
	"errors"
	"fmt"
	els "github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/hyuti/api-blueprint/config"
	pkgGprc "github.com/hyuti/api-blueprint/pkg/grpc"
	"github.com/hyuti/api-blueprint/pkg/telegram"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"sync"
)

// TODO: serviceKey should be adjusted to be relevant to the project
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
	ginEngine *gin.Engine
	elk       *els.TypedClient
	grpcSrv   *pkgGprc.Server
	tele      *telegram.Tele
	dbDriver  *gorm.DB
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
	// TODO: add more apps here.
	if _, err := pkgGprc.DefaultValidator(); err != nil {
		return fmt.Errorf("cannot init validator: %w", err)
	}
	return nil
}
