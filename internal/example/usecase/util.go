package usecase

import (
	"context"

	"github.com/hyuti/API-Golang-Template/internal/example/entity"
)

func paginate[T any](
	_ context.Context,
	req *entity.PaginationListRequest, data []T) (re []T) {
	min := *req.PageSize * (*req.Page)
	l := int32(len(data))
	if min >= l {
		return re
	}
	max := *req.PageSize*(*req.Page+1) + 1
	if max > l {
		max = l
	}
	return data[min:max]
}

func respPaginator[T any](page, size int32, data []T) *entity.PaginatedResponse[T] {
	r := new(entity.PaginatedResponse[T])
	r.Count = int32(len(data))
	r.Data = data
	r.PageSize = size

	if page > 0 {
		r.Prev = page - 1
	}
	if r.Count > size {
		r.Data = data[:len(data)-1]
		r.Count -= 1
		r.Next = page + 1
	}

	return r
}
