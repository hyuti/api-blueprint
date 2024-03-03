package usecase

import (
	"errors"
	"strings"
)

var ErrProcessingRequest = errors.New("something wrong occurred")

const (
	ErrProcessingReq = "err_processing_request"
	ErrValidateReq   = "err_validating_request"
)

type (
	Code interface {
		Code() string
	}
	Extra interface {
		Extra() map[string]any
	}
)

type errValidation struct {
	key   string
	value string
	msg   string
}

func (e errValidation) Code() string {
	return ErrValidateReq
}
func (e errValidation) Error() string {
	var s strings.Builder
	s.WriteString(e.msg)
	s.WriteByte(' ')
	s.WriteString("field=")
	s.WriteString(e.key)
	s.WriteByte(' ')
	s.WriteString("value=")
	s.WriteString(e.value)
	return s.String()
}

type errProcessingReq struct {
	extra map[string]any
	err   error
}

func (e errProcessingReq) Code() string {
	return ErrProcessingReq
}

func (e errProcessingReq) Extra() map[string]any {
	return e.extra
}
func (e errProcessingReq) Error() string {
	if e.err == nil {
		e.err = ErrProcessingRequest
	}
	return e.err.Error()
}
