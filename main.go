package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	protovalidatemiddleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/ratelimit"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/hyuti/api-blueprint/docs"
	"github.com/hyuti/api-blueprint/internal/app"
	srvGrpc "github.com/hyuti/api-blueprint/internal/grpc"
	"github.com/hyuti/api-blueprint/internal/proto"
	"github.com/hyuti/api-blueprint/internal/repo"
	"github.com/hyuti/api-blueprint/internal/router"
	"github.com/hyuti/api-blueprint/internal/usecase"
	pkgGprc "github.com/hyuti/api-blueprint/pkg/grpc"
	"github.com/hyuti/api-blueprint/pkg/http/middleware"
	"github.com/hyuti/minion"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"log"
	"sync"
)

var (
	mutex sync.Mutex
	_uc   usecase.ExampleUseCase
	_uc1  usecase.NotiIfPanicUseCase
	_jobs []func() error
)

func init() {
	if err := app.Init(); err != nil {
		log.Fatalln(err)
	}
	// TODO: add more usecases and repositories as well
	// usecases and repos should be cross servers level because plenty of servers could depend on them
	re := repo.NewExampleRepo(app.Els())
	_uc = usecase.NewExampleUseCase(re)
	_uc1 = usecase.NewNotiIfPanicUseCase(app.Tele(), app.Cfg().Name)

	// TODO: add more jobs here
	initRestfulServer()
	// initGrpcServer()
}

// @title Example API
// @version 1.0

// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey AuthToken
// @in header
// @name Authorization

// @description Example API.
func main() {
	gru := minion.New[error]()
	defer gru.Clean()
	gru.Start(_jobs...)

	if err := gru.Error(); err != nil {
		log.Fatalln(err)
	}
}

func sharedInit(initiator func(), job func() error) {
	mutex.Lock()
	defer mutex.Unlock()
	initiator()
	_jobs = append(_jobs, job)
}
func initRestfulServer() {
	sharedInit(func() {
		app.Gin().Use(gin.CustomRecovery(router.OnPanic(_uc1)))
		app.Gin().Use(middleware.LimiterMiddleware(pkgGprc.RateLimit().Limiter()))
		basePath := app.Cfg().Gin.BasePath + "/api/v1"
		docs.SwaggerInfo.BasePath = basePath
		group := app.Gin().Group(basePath)
		{
			group.GET("/docs/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))
			group.GET("/healthcheck", router.HeathCheck)
			router.New(group, _uc)
		}
	}, func() error {
		return app.Gin().Run(fmt.Sprintf("0.0.0.0:%v", app.Cfg().Gin.Port))
	})
}
func initGrpcServer() {
	grpcS := app.GrpcSrv()
	sharedInit(func() {
		grpcS.WithOpt(
			grpc.ChainUnaryInterceptor(
				recovery.UnaryServerInterceptor(recovery.WithRecoveryHandlerContext(_uc1.GrpcHandle)),
				logging.UnaryServerInterceptor(pkgGprc.InterceptorLogger(app.Logger())),
				protovalidatemiddleware.UnaryServerInterceptor(pkgGprc.Validator()),
				ratelimit.UnaryServerInterceptor(pkgGprc.RateLimit()),
			),
			grpc.ChainStreamInterceptor(
				recovery.StreamServerInterceptor(recovery.WithRecoveryHandlerContext(_uc1.GrpcHandle)),
				logging.StreamServerInterceptor(pkgGprc.InterceptorLogger(app.Logger())),
			),
		)
		proto.RegisterApiGolangTemplateServer(grpcS.Server(), srvGrpc.New(
			_uc,
		))
	}, func() error {
		return grpcS.Run()
	})
}
