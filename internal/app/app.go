package app

import (
	els "github.com/elastic/go-elasticsearch/v8"
	"github.com/gin-gonic/gin"
	"github.com/hyuti/api-blueprint/config"
	pkgGprc "github.com/hyuti/api-blueprint/pkg/grpc"
	"github.com/hyuti/api-blueprint/pkg/telegram"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"sync"
)

var (
	mutex sync.Mutex
	app   *App
)

// TODO: add compiler checker if a specific attribute not initilized but listed, same as one as the stringer pkg did
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
	a, err := initializeApp()
	if err != nil {
		return err
	}
	app = a
	return nil
}
