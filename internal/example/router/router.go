package router

import (
	"github.com/hyuti/API-Golang-Template/internal/example/entity"
	"github.com/hyuti/API-Golang-Template/internal/example/proto"
	"github.com/hyuti/API-Golang-Template/pkg/collection"
	"github.com/hyuti/API-Golang-Template/pkg/tool"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyuti/API-Golang-Template/internal/example/usecase"
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

// @Summary List
// @Description Get list
// @Tags list
// @Accept  json
// @Produce json
// @Param payload body ListExampleReq true "Request body"
// @Success 200 {object} ListExampleResp
// @Failure 500 {object} ErrResponse
// @Failure 400 {object} ErrResponse
// @Router /list [post] .
func (r *route) list(ctx *gin.Context) {
	var req ListExampleReq
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.AbortWithStatusJSON(http.StatusBadRequest, ErrValidation(ctx, err))
		return
	}
	ctx.Set(plKey, req)
	resp, err := r.uc1.List(ctx.Request.Context(), &usecase.ExampleReq{
		PageSize: req.PageSize,
		Page:     req.Page,
		Search:   req.Search,
	})
	if err != nil {
		usecaseRouterErrMapper(ctx, err)
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, &ListExampleResp{
		PaginatedResponse: &entity.PaginatedResponse[*Example]{
			Next:     resp.Next,
			Prev:     resp.Prev,
			PageSize: resp.PageSize,
			Count:    resp.Count,
			Data: collection.Map(resp.Data, func(item *proto.Example, index int) *Example {
				return &Example{
					Example: item,
				}
			}),
		},
	})
}
