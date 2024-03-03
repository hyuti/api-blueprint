package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	protovalidatemiddleware "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/ratelimit"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/hyuti/API-Golang-Template/docs"
	"github.com/hyuti/API-Golang-Template/internal/app"
	srvGrpc "github.com/hyuti/API-Golang-Template/internal/example/grpc"
	"github.com/hyuti/API-Golang-Template/internal/example/proto"
	"github.com/hyuti/API-Golang-Template/internal/example/repo"
	"github.com/hyuti/API-Golang-Template/internal/example/router"
	"github.com/hyuti/API-Golang-Template/internal/example/usecase"
	pkgGprc "github.com/hyuti/API-Golang-Template/pkg/grpc"
	"github.com/hyuti/API-Golang-Template/pkg/http/middleware"
	"github.com/hyuti/API-Golang-Template/pkg/tool"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"google.golang.org/grpc"
	"log"
)

// @title Example API
// @version 1.0

// @BasePath /api
// @schemes http https

// @securityDefinitions.apikey AuthToken
// @in header
// @name Authorization

// @description Example API.
func main() {
	if err := app.Init(); err != nil {
		log.Fatalln(err)
	}

	re := repo.NewExampleRepo(app.Els())
	_uc := usecase.NewExampleUseCase(re)
	_uc1 := usecase.NewNotiIfPanicUseCase(app.Tele(), app.Cfg().Name)

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

	grpcS := app.GrpcSrv()
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

	wo := tool.New[error]()
	wo.AddJob(func() error {
		return app.Gin().Run(fmt.Sprintf("0.0.0.0:%v", app.Cfg().Gin.Port))
	}).AddJob(func() error {
		return grpcS.Run()
	}).Start()

	if err := wo.Error(); err != nil {
		log.Fatalln(err)
	}
}
