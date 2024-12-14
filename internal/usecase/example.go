package usecase

import (
	"context"
	"github.com/hyuti/api-blueprint/internal/proto"

	"github.com/hyuti/api-blueprint/internal/repo"
)

type ExampleUseCase interface {
	List(ctx context.Context, req *ExampleReq) (*ExampleResp, error)
}

func NewExampleUseCase(r repo.ExampleRepo) *ExampleUseCaseImpl {
	return &ExampleUseCaseImpl{repo: r}
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

var _ ExampleUseCase = (*ExampleUseCaseImpl)(nil)

func (e *ExampleUseCaseImpl) List(ctx context.Context, req *ExampleReq) (*ExampleResp, error) {
	if _, err := e.repo.List(ctx, &repo.Where[*repo.ExampleWhereReq]{
		W: nil,
	}); err != nil {
		return nil, err
	}
	return &ExampleResp{}, nil
}

type ExampleUseCaseImpl struct {
	repo repo.ExampleRepo
}
