package grpc

import (
	"context"
	"errors"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	pkgerr "github.com/hyuti/api-blueprint/pkg/error"
	"golang.org/x/exp/slog"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"runtime"
)

func handleError(
	ctx context.Context,
	lgr *slog.Logger,
	err error,
) error {
	code := codes.Internal

	var myErr *pkgerr.Error
	if !errors.As(err, &myErr) {
		myErr = pkgerr.ErrInternalServer(err)
	}
	switch {
	case errors.Is(myErr, pkgerr.LabelErrValidatingRequest):
		code = codes.InvalidArgument
	case errors.Is(myErr, pkgerr.LabelErrAuthenticateRequest):
		code = codes.Unauthenticated
	case errors.Is(myErr, pkgerr.LabelErrAuthorizeRequest):
		code = codes.PermissionDenied
	}
	if code != codes.Internal {
		return status.Error(code, myErr.Error())
	}

	// TODO: trigger github issue creation flow
	lgr.ErrorContext(
		ctx,
		"error internal server",
		"error", myErr.Error(),
		"func", myErr.NameFunc(),
		"payload", myErr.Payload(),
		"chain", myErr.Chain(),
		"path", "",
		"controller", "",
		"params", "",
		"query", "",
	)

	return status.Error(code, "Something went wrong, please check server logs for detail")
}

func OnPanic(lgr *slog.Logger) recovery.RecoveryHandlerFuncContext {
	return func(ctx context.Context, errObj any) error {
		var chain []string
		for i := 5; ; i++ {
			pc, _, _, ok := runtime.Caller(i)
			if !ok {
				break
			}
			chain = append(chain, runtime.FuncForPC(pc).Name())
		}
		err, ok := errObj.(error)
		if !ok {
			err = errors.New(chain[0])
		}

		myErr := pkgerr.ErrInternalServer(
			err,
			pkgerr.WithChainOpt(chain...),
			pkgerr.WithNameFuncOpt(chain[len(chain)-1]),
		)
		return handleError(ctx, lgr, myErr)
	}
}
