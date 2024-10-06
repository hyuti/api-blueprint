package router

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/hyuti/api-blueprint/internal/proto"
	"github.com/hyuti/api-blueprint/internal/usecase"
	"github.com/hyuti/api-blueprint/pkg/collection"
	"net/http"
)

type PaginatedResponse[T any] struct {
	Next     int32 `json:"next"`
	Prev     int32 `json:"previous"`
	PageSize int32 `json:"page_size"`
	Count    int32 `json:"count"`
	Data     []T   `json:"data"`
} // @name PaginatedResponse

type PaginationListRequest struct {
	PageSize *int32 `json:"page_size" example:"100"` // default is 100
	Page     *int32 `json:"page" example:"0"`        // default is 0
}
type Example struct {
	*proto.Example
}
type ListExampleReq struct {
	// default is 100
	PageSize *int32 `json:"page_size"`
	// default is 0
	Page   *int32  `json:"page"`
	Search *string `json:"search" binding:"required"`
}
type ListExampleResp struct {
	*PaginatedResponse[*Example]
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
		handleError(ctx, r.lgr, err)
		return
	}
	ctx.Set(plKey, req)
	resp, err := r.uc1.List(ctx.Request.Context(), &usecase.ExampleReq{
		PaginatedRequest: usecase.PaginatedRequest{
			PageSize: req.PageSize,
			Page:     req.Page,
		},
	})
	if err != nil {
		handleError(ctx, r.lgr, err)
		return
	}

	ctx.AbortWithStatusJSON(http.StatusOK, &ListExampleResp{
		PaginatedResponse: &PaginatedResponse[*Example]{
			Next:     resp.Next,
			Prev:     resp.Prev,
			PageSize: resp.PageSize,
			Count:    resp.Count,
			Data: collection.Map(resp.Data, func(item usecase.Example, index int) *Example {
				return &Example{
					Example: item.Example,
				}
			}),
		},
	})
}

func (x *Example) MarshalJSON() ([]byte, error) {
	return marshaller.Marshal(x.Example)
}

var _ json.Marshaler = (*Example)(nil)
