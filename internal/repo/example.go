package repo

import (
	"context"
	"errors"
	els "github.com/elastic/go-elasticsearch/v8"
	"time"
)

type (
	ExampleRepo interface {
		repo[*ExampleWhereReq, *Example]
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

func NewExampleRepo(els *els.TypedClient) ExampleRepo {
	r := &exampleRepo{
		els: els,
	}
	r.repo = &Repo[*ExampleWhereReq, *Example]{
		lister: r,
	}
	return r
}

type exampleRepo struct {
	repo[*ExampleWhereReq, *Example]
	els *els.TypedClient
}

func (e *exampleRepo) List(ctx context.Context, w *Where[*ExampleWhereReq]) ([]*Example, error) {
	return nil, errors.New("not implemented")
}
