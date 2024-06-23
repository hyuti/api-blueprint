package app

import (
	"fmt"
	pkgGprc "github.com/hyuti/api-blueprint/pkg/grpc"
)

func WithGRPCServer() (*pkgGprc.Server, error) {
	grpcS, err := pkgGprc.New(
		&pkgGprc.SrvConfig{
			Port: app.cfg.GRPC.Port,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("cannot init grpc: %w", err)
	}
	if _, err := pkgGprc.DefaultValidator(); err != nil {
		return nil, fmt.Errorf("cannot init grpc: %w", err)
	}
	return grpcS, nil
}
func GrpcSrv() *pkgGprc.Server {
	mutex.Lock()
	defer mutex.Unlock()
	return app.grpcSrv
}
