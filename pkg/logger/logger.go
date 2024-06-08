package logger

import (
	"context"
	ctx2 "github.com/hyuti/api-blueprint/pkg/ctx"
	"golang.org/x/exp/slog"
	"io"
	"os"
	"time"
)

func FileAndStdLogger(path string, opts ...func(options *slog.HandlerOptions)) *slog.Logger {
	return New(NewJsonHandler(FileAndStdWriter(path), opts...))
}
func Default(opts ...func(options *slog.HandlerOptions)) *slog.Logger {
	return New(NewJsonHandler(DefaultWriter(), opts...))
}

func NetCallWrapper[T any](
	ctx context.Context,
	l *slog.Logger,
	name,
	before, after string,
	invoker func() T, args ...any) T {
	if before == "" {
		before = "started network calling"
	}
	if after == "" {
		after = "done network calling"
	}
	l.InfoContext(ctx, before, args...)
	r := invoker()
	args = append(args, "resp", r)
	l.InfoContext(ctx, after, args...)
	return r
}

func WithCtxAttrs(h slog.Handler, attrs ...func(ctx context.Context) slog.Attr) slog.Handler {
	return &ctxHandler{
		Handler: h,
		attrs:   attrs,
	}
}
func WithAttrs(h slog.Handler, attrs ...slog.Attr) slog.Handler {
	h.WithAttrs(attrs)
	return h
}
func WithServiceName(l *slog.Logger, srvKey, srvName string) *slog.Logger {
	WithAttrs(l.Handler(), slog.String(srvKey, srvName))
	return l
}
func WithCtxID(l *slog.Logger) *slog.Logger {
	h := WithCtxAttrs(l.Handler(), func(ctx context.Context) slog.Attr {
		return slog.String(ctx2.CtxIDKey, ctx2.GetCtxID(ctx))
	})
	return slog.New(h)
}
func WithLevelOpt(level slog.Level) func(opt *slog.HandlerOptions) {
	return func(opt *slog.HandlerOptions) {
		opt.Level = level
	}
}
func WithTimeFormatOpt(format string) func(opt *slog.HandlerOptions) {
	return func(opt *slog.HandlerOptions) {
		opt.ReplaceAttr = func(groups []string, a slog.Attr) slog.Attr {
			if a.Key != slog.TimeKey {
				return a
			}
			v, ok := a.Value.Any().(time.Time)
			if !ok {
				return a
			}
			f := time.RFC3339
			if format != "" {
				f = format
			}
			a.Value = slog.StringValue(v.Format(f))
			return a
		}
	}
}
func WithFileWriter(w io.Writer, path string) io.Writer {
	if w == nil {
		w = os.Stdout
	}
	logFile, err := os.OpenFile(
		path,
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0664,
	)
	if err != nil {
		panic(err)
	}
	return io.MultiWriter(w, logFile)
}
func DefaultWriter() io.Writer {
	return os.Stdout
}
func FileAndStdWriter(path string) io.Writer {
	return WithFileWriter(DefaultWriter(), path)
}
func New(handler slog.Handler) *slog.Logger {
	return slog.New(handler)
}
func NewJsonHandler(w io.Writer, opts ...func(options *slog.HandlerOptions)) slog.Handler {
	opt := slog.HandlerOptions{}
	for _, o := range opts {
		o(&opt)
	}
	return slog.NewJSONHandler(w, &opt)
}
