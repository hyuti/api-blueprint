package usecase

import (
	"context"
	"errors"
	"github.com/hyuti/api-blueprint/internal/proto"

	"github.com/hyuti/api-blueprint/internal/repo"
)

type ExampleUseCase interface {
	List(ctx context.Context, req *ExampleReq) (*ExampleResp, error)
}

func NewExampleUseCase(r repo.ExampleRepo) ExampleUseCase {
	return &exampleUC{repo: r}
}

type (
	Example struct {
		*proto.Example
	}
	ExampleReq struct {
		PaginatedRequest
	}
	ExampleResp struct {
		PaginatedResponse[Example]
	}
)

func (e *exampleUC) List(ctx context.Context, req *ExampleReq) (*ExampleResp, error) {
	return nil, errors.New("not implemented")
}

type exampleUC struct {
	repo repo.ExampleRepo
}
