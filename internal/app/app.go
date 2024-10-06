package app

import (
	"github.com/gin-gonic/gin"
	"github.com/hyuti/api-blueprint/config"
	pkgGprc "github.com/hyuti/api-blueprint/pkg/grpc"
	"golang.org/x/exp/slog"
	"gorm.io/gorm"
	"sync"
)

var (
	mutex sync.Mutex
	app   *App
)

// TODO: add compiler checker if a specific attribute not initilized but listed, same as the one the stringer pkg did
type App struct {
	cfg       *config.Config
	logger    *slog.Logger
	ginEngine *gin.Engine
	grpcSrv   *pkgGprc.Server
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
