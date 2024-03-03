package entity

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
