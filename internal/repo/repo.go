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

type repo[W, M any] interface {
	List(ctx context.Context, w *Where[W]) ([]M, error)
	ListBatch(ctx context.Context, w *Where[W], callback func(ctx context.Context, data []M) error) error
	ListBatchNonGreedy(ctx context.Context, w *Where[W], callback func(ctx context.Context, data []M) (bool, error)) error
	ListBatchV2(
		ctx context.Context,
		w *Where[W],
		callback func(ctx context.Context, data []M) error,
		searchAfterCB func(ctx context.Context, lastHit M, searchAfter []any) error,
	) error
	ListBatchNonGreedyV2(
		ctx context.Context,
		w *Where[W],
		callback func(ctx context.Context, data []M) (bool, error),
		searchAfterCB func(ctx context.Context, lastHit M, searchAfter []any) error,
	) error
	Retrieve(ctx context.Context, w *Where[W]) (M, error)
}

var _ repo[any, any] = (*Repo[any, any])(nil)

type Repo[W, M any] struct {
	lister Lister[W, M]
}

func (r *Repo[W, M]) List(ctx context.Context, w *Where[W]) ([]M, error) {
	return r.lister.List(ctx, w)
}

func (r *Repo[W, M]) ListBatch(ctx context.Context, w *Where[W], callback func(ctx context.Context, data []M) error) error {
	return r.ListBatchNonGreedy(ctx, w, func(ctx context.Context, result []M) (bool, error) {
		return false, callback(ctx, result)
	})
}

// ListBatchV2 uses search_after feature of Elastic instead of from param to process big data in batch.
func (r *Repo[W, M]) ListBatchV2(
	ctx context.Context,
	w *Where[W],
	callback func(ctx context.Context, data []M) error,
	searchAfterCB func(ctx context.Context, lastHit M, searchAfter []any) error,
) error {
	return r.ListBatchNonGreedyV2(ctx, w, func(ctx context.Context, result []M) (bool, error) {
		return false, callback(ctx, result)
	}, searchAfterCB)
}

func (r *Repo[W, M]) ListBatchNonGreedyV2(
	ctx context.Context,
	w *Where[W],
	callback func(ctx context.Context, data []M) (bool, error),
	searchAfterCB func(ctx context.Context, lastHit M, searchAfter []any) error,
) error {
	searchAfter := make([]any, 0, 1)
	data, err := r.List(ctx, w)
	for ; err == nil && len(data) > 0; data, err = r.List(ctx, w) {
		stop := false
		stop, err = callback(ctx, data)
		if err != nil {
			return err
		}
		if stop {
			return nil
		}
		//avoiding allocating unnecessary memory
		searchAfter = searchAfter[:0]
		if err := searchAfterCB(ctx, data[len(data)-1], searchAfter); err != nil {
			return err
		}
		w.SearchAfter = searchAfter
	}
	return err
}

func (r *Repo[W, M]) ListBatchNonGreedy(ctx context.Context, w *Where[W], callback func(ctx context.Context, data []M) (bool, error)) error {
	page := int32(0)
	w.Page = &page
	data, err := r.List(ctx, w)
	for ; err == nil && len(data) > 0; data, err = r.List(ctx, w) {
		stop := false
		stop, err = callback(ctx, data)
		if err != nil {
			return err
		}
		if stop {
			return nil
		}
		page += 1
		w.Page = &page
	}
	return err
}

func (r *Repo[W, M]) Retrieve(ctx context.Context, w *Where[W]) (m M, err error) {
	result, err := r.List(ctx, w)
	if err != nil {
		return m, err
	}
	if len(result) == 0 {
		return m, nil
	}
	return result[0], nil
}
