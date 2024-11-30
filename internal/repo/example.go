package repo

import (
	"context"
	"errors"
	pkgerr "github.com/hyuti/api-blueprint/pkg/error"
	"time"
)

type (
	ExampleRepo interface {
		Repo[*ExampleWhereReq, *Example]
	}
	ExampleWhereReq struct {
		Search *string
		Page   *int32
	}
	Example struct {
		Id        int32
		Name      string
		CreatedAt *time.Time
		UpdatedAt *time.Time
		DeletedAt *time.Time
	}
)

func NewExampleRepo() *ExampleRepoImpl {
	r := &ExampleRepoImpl{}
	r.Repo = &RepoImpl[*ExampleWhereReq, *Example]{
		lister: r,
	}
	return r
}

type ExampleRepoImpl struct {
	Repo[*ExampleWhereReq, *Example]
}

var _ ExampleRepo = (*ExampleRepoImpl)(nil)

func (e *ExampleRepoImpl) List(ctx context.Context, w *Where[*ExampleWhereReq]) ([]*Example, error) {
	return nil, pkgerr.DefaultErrInternalServer(
		errors.New("not implemented"),
		"e.ErrFunc",
	)
}
