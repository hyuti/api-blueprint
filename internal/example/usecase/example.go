package usecase

import (
	"context"
	"errors"
	"github.com/hyuti/API-Golang-Template/internal/example/proto"

	"github.com/hyuti/API-Golang-Template/internal/example/entity"
	"github.com/hyuti/API-Golang-Template/internal/example/repo"
)

type ExampleUseCase interface {
	List(ctx context.Context, req *ExampleReq) (*ExampleResp, error)
}

func NewExampleUseCase(r repo.ExampleRepo) ExampleUseCase {
	return &exampleUC{repo: r}
}

type (
	ExampleReq struct {
		PageSize *int32
		Page     *int32
		Search   *string
	}
	ExampleResp struct {
		*entity.PaginatedResponse[*proto.Example]
	}
)

func (e *exampleUC) List(ctx context.Context, req *ExampleReq) (*ExampleResp, error) {
	return nil, errors.New("not implemented")
}

type exampleUC struct {
	repo repo.ExampleRepo
}
