package app

import (
	"fmt"
	pkgGprc "github.com/hyuti/api-blueprint/pkg/grpc"
)

func WithGRPCServer() error {
	grpcS, err := pkgGprc.New(
		&pkgGprc.SrvConfig{
			Port: app.cfg.GRPC.Port,
		},
	)
	if err != nil {
		return fmt.Errorf("cannot init grpc: %w", err)
	}
	app.grpcSrv = grpcS
	return nil
}
func GrpcSrv() *pkgGprc.Server {
	mutex.Lock()
	defer mutex.Unlock()
	return app.grpcSrv
}
