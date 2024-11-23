package grpc

import (
	"context"
	"google.golang.org/grpc"
)

const (
	keyControllerContext = "keyControllerContext"
	keyPathContext       = "keyPathContext"
)

func UnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		ctx = context.WithValue(ctx, keyControllerContext, info.FullMethod)
		ctx = context.WithValue(ctx, keyPathContext, info.FullMethod)
		return handler(ctx, req)
	}
}

func StreamServerInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		ctx := ss.Context()
		ctx = context.WithValue(ctx, keyControllerContext, info.FullMethod)
		ctx = context.WithValue(ctx, keyPathContext, info.FullMethod)
		return handler(srv, ss)
	}
}

func PathContext(ctx context.Context) string {
	return ctx.Value(keyPathContext).(string)
}

func ControllerContext(ctx context.Context) string {
	return ctx.Value(keyControllerContext).(string)
}
