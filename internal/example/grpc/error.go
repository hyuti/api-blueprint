package grpc

import (
	"github.com/hyuti/API-Golang-Template/internal/example/usecase"
	"github.com/hyuti/API-Golang-Template/pkg/tool"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const (
	ErrProcessingReq = "err_processing_request"
	ErrValidateReq   = "err_validating_request"
)

func usecaseGrpcMapper(err error) error {
	code := codes.Internal
	codeInternal := ErrProcessingReq
	if e, ok := err.(usecase.Code); ok {
		switch e.Code() {
		case usecase.ErrValidateReq:
			code = codes.InvalidArgument
			codeInternal = ErrValidateReq
		}
	}
	extra := make(map[string]any)
	if e, ok := err.(usecase.Extra); ok {
		extra = e.Extra()
	}
	return status.Error(code, ErrResponse{
		Code:  codeInternal,
		Msg:   err.Error(),
		Extra: extra,
	}.Error())
}

type ErrResponse struct {
	Code  string         `json:"code"`    // Error code
	Msg   string         `json:"message"` // Description error
	Extra map[string]any `json:"extra"`   // Extra info about error
}

func (e ErrResponse) Error() string {
	return tool.JSONStringify(&e)
}
