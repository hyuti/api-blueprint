package router

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/hyuti/api-blueprint/internal/usecase"
	pkgerr "github.com/hyuti/api-blueprint/pkg/error"
	"github.com/hyuti/api-blueprint/pkg/tool"
	"golang.org/x/exp/slog"
	"net/http"
	"runtime"
)

const plKey = "payloadKey"

func New(router gin.IRouter,
	u1 usecase.ExampleUseCase,
) {
	_router := new(route)
	_router.uc1 = u1

	router.POST("/list", _router.list)
}

type route struct {
	lgr *slog.Logger
	uc1 usecase.ExampleUseCase
}

func OnPanic(lgr *slog.Logger) gin.RecoveryFunc {
	return func(ctx *gin.Context, errObj any) {
		var data []byte
		if d, ok := ctx.Get(plKey); ok {
			// force to be json format no matter underlying format is
			data = []byte(tool.JSONStringify(d))
		}
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
			pkgerr.WithPayloadOpt(data),
		)
		handleError(ctx, lgr, myErr)
	}
}

func HeathCheck(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusOK)
}
