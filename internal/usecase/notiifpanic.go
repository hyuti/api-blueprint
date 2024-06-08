package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	pkgCtx "github.com/hyuti/api-blueprint/pkg/ctx"
	"github.com/hyuti/api-blueprint/pkg/telegram"
	"github.com/hyuti/api-blueprint/pkg/tool"
)

const (
	errKey  = "err"
	dataKey = "data"
)

type NotiIfPanicUseCase interface {
	Handle(context.Context, []byte, any) error
	GrpcHandle(context.Context, any) error
}

func NewNotiIfPanicUseCase(
	tele *telegram.Tele,
	name string,
) NotiIfPanicUseCase {
	return &notiIfPanicUC{
		tele: tele,
		name: name,
	}
}

type notiIfPanicUC struct {
	tele *telegram.Tele
	name string
}

func (u *notiIfPanicUC) handle(ctx context.Context) {
	msg := ctx.Value(errKey)
	if _msg, ok := msg.(error); ok {
		msg = _msg.Error()
	}
	pl := ctx.Value(dataKey)

	go func() {
		m, ok := msg.(string)
		if !ok {
			m = ""
		}
		e := telegram.ErrorMsg(m, u.name, "", tool.JSONStringify(map[string]any{
			"payload": pl,
			"id":      pkgCtx.GetCtxID(ctx),
		}))
		_ = u.tele.SendWithTeleMsg(e)
	}()
}

func (u *notiIfPanicUC) GrpcHandle(ctx context.Context, p any) error {
	_ctx := context.WithValue(ctx, errKey, p)
	u.handle(_ctx)
	return fmt.Errorf("%w: %s", ErrProcessingRequest, p)
}

func (u *notiIfPanicUC) Handle(ctx context.Context, data []byte, err any) error {
	var pl any
	if err := json.Unmarshal(data, &pl); err != nil {
		pl = string(data)
	}
	_ctx := context.WithValue(ctx, errKey, err)
	_ctx = context.WithValue(_ctx, dataKey, pl)
	u.handle(_ctx)
	return errProcessingReq{
		err: ErrProcessingRequest,
		extra: map[string]any{
			"id":      pkgCtx.GetCtxID(ctx),
			"payload": pl,
		},
	}
}
