package error

import (
	"errors"
	"runtime"
	"strings"
)

var (
	LabelErrValidatingRequest   = errors.New("")
	LabelErrInternalServer      = errors.New("")
	LabelErrAuthenticateRequest = errors.New("")
	LabelErrAuthorizeRequest    = errors.New("")
)

type ErrInternalServerOpt func(impl *Error)

var WithNameFuncOpt = func(nameFunc string) ErrInternalServerOpt {
	return func(err *Error) {
		err.nameFunc = nameFunc
	}
}
var WithPayloadOpt = func(payload any) ErrInternalServerOpt {
	return func(impl *Error) {
		impl.payload = payload
	}
}
var WithChainOpt = func(chain ...string) ErrInternalServerOpt {
	return func(impl *Error) {
		impl.chain = chain
	}
}
var MostRecentChainOpt = func(errFunc string) ErrInternalServerOpt {
	return func(err *Error) {
		nameFunc := "TwoMostRecentChainOpt"
		if pc, _, _, ok := runtime.Caller(2); !ok {
			nameFunc = runtime.FuncForPC(pc).Name()
		}
		err.chain = []string{errFunc, nameFunc}
	}
}
var CurrentNameFuncOpt = func(skips ...int) ErrInternalServerOpt {
	return func(err *Error) {
		skip := 0
		if len(skips) > 0 {
			skip = skips[0]
		}
		nameFunc := "CurrentNameFuncOpt"
		if pc, _, _, ok := runtime.Caller(skip + 2); !ok {
			nameFunc = runtime.FuncForPC(pc).Name()
		}
		err.nameFunc = nameFunc
	}
}

type Error struct {
	error
	label    error
	nameFunc string
	payload  any
	chain    []string
	extra    map[string]any
}

func (e *Error) Unwrap() error {
	return e.label
}

func (e *Error) NameFunc() string {
	return e.nameFunc
}

func (e *Error) Payload() any {
	return e.payload
}

func (e *Error) Chain() string {
	return strings.Join(e.chain, "<-")
}

func (e *Error) Get(key string, value ...any) any {
	v, ok := e.extra[key]
	if ok {
		return v
	}
	if len(value) > 0 {
		v = value[0]
	}
	e.extra[key] = v
	return v
}

func (e *Error) Extra() map[string]any {
	return e.extra
}

func ErrValidatingRequest(err error) *Error {
	return &Error{
		error: err,
		label: LabelErrValidatingRequest,
	}
}

func ErrAuthenticateRequest(err error) *Error {
	return &Error{
		error: err,
		label: LabelErrAuthenticateRequest,
	}
}

func ErrAuthorizeRequest(err error) *Error {
	return &Error{
		error: err,
		label: LabelErrAuthorizeRequest,
	}
}

func ErrInternalServer(err error, opts ...ErrInternalServerOpt) *Error {
	errImpl := &Error{
		error: err,
		label: LabelErrInternalServer,
	}
	for _, opt := range opts {
		opt(errImpl)
	}
	return errImpl
}

func DefaultErrInternalServer(err error, errFunc string) *Error {
	return ErrInternalServer(
		err,
		MostRecentChainOpt(errFunc),
		CurrentNameFuncOpt(1),
	)
}
