package usecase

import (
	"context"
)

type PaginatedRequest struct {
	PageSize *int32
	Page     *int32
}

type PaginatedResponse[T any] struct {
	Next     int32
	Prev     int32
	PageSize int32
	Count    int32
	Data     []T
}

func paginate[T any](
	_ context.Context,
	req *PaginatedRequest, data []T) (re []T) {
	mi := *req.PageSize * (*req.Page)
	l := int32(len(data))
	if mi >= l {
		return re
	}
	ma := *req.PageSize*(*req.Page+1) + 1
	if ma > l {
		ma = l
	}
	return data[mi:ma]
}

func respPaginator[T any](page, size int32, data []T) *PaginatedResponse[T] {
	r := new(PaginatedResponse[T])
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
