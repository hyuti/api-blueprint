package router

import (
	"github.com/hyuti/api-blueprint/pkg/tool"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyuti/api-blueprint/internal/usecase"
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
	uc1 usecase.ExampleUseCase
}

func OnPanic(uc usecase.NotiIfPanicUseCase) gin.RecoveryFunc {
	return func(ctx *gin.Context, err any) {
		var data []byte
		if d, ok := ctx.Get(plKey); ok {
			data = []byte(tool.JSONStringify(d))
		}
		_err := uc.Handle(ctx.Request.Context(), data, err)
		usecaseRouterErrMapper(ctx, _err)
	}
}

func HeathCheck(ctx *gin.Context) {
	ctx.AbortWithStatus(http.StatusOK)
}
