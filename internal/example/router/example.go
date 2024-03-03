package router

import (
	"encoding/json"
	"github.com/hyuti/API-Golang-Template/internal/example/entity"
	"github.com/hyuti/API-Golang-Template/internal/example/proto"
)

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
	*entity.PaginatedResponse[*Example]
}

func (x *Example) MarshalJSON() ([]byte, error) {
	return marshaller.Marshal(x.Example)
}

var _ json.Marshaler = (*Example)(nil)
