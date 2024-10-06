package repo

import (
	"context"
)

type Where[W any] struct {
	W           W
	Page        *int32
	SearchAfter []any
}

type Lister[W, M any] interface {
	List(ctx context.Context, w *Where[W]) ([]M, error)
}

type Repo[W, M any] interface {
	List(ctx context.Context, w *Where[W]) ([]M, error)
	Retrieve(ctx context.Context, w *Where[W]) (M, error)
}

var _ Repo[any, any] = (*RepoImpl[any, any])(nil)

type RepoImpl[W, M any] struct {
	lister Lister[W, M]
}

func (r *RepoImpl[W, M]) List(ctx context.Context, w *Where[W]) ([]M, error) {
	return r.lister.List(ctx, w)
}

func (r *RepoImpl[W, M]) Retrieve(ctx context.Context, w *Where[W]) (m M, err error) {
	result, err := r.List(ctx, w)
	if err != nil {
		return m, err
	}
	if len(result) == 0 {
		return m, nil
	}
	return result[0], nil
}
