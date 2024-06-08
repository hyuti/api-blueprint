package router

import (
	"errors"
	"github.com/go-playground/validator/v10"
	pkgHttp "github.com/hyuti/api-blueprint/pkg/http"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/hyuti/api-blueprint/internal/usecase"
)

const (
	ErrProcessingReq = "err_processing_request"
	ErrValidateReq   = "err_validating_request"
)

var ErrInvalidData = errors.New("some data is invalid")

// @Description Response of API if error occurs
type ErrResponse struct {
	Code  string         `json:"code"`    // Error code
	Msg   string         `json:"message"` // Description error
	Extra map[string]any `json:"extra"`   // Extra info about error
} // @name Error-Response

func ErrResp(err error, code string, extra ...func(e map[string]any)) ErrResponse {
	e := ErrResponse{
		Code:  code,
		Msg:   err.Error(),
		Extra: make(map[string]any),
	}
	for _, f := range extra {
		f(e.Extra)
	}
	return e
}
func ErrProcessing(ctx *gin.Context, err error, extra ...func(e map[string]any)) ErrResponse {
	return ErrResp(err, ErrProcessingReq, extra...)
}
func ErrValidation(ctx *gin.Context, err error, extra ...func(e map[string]any)) ErrResponse {
	if errs, ok := err.(validator.ValidationErrors); ok {
		if trans, _err := pkgHttp.TranslatorCtx(ctx); _err == nil {
			err = ErrInvalidData
			for k, v := range errs.Translate(trans) {
				extra = append(extra, func(e map[string]any) {
					e[k] = v
				})
			}
		}
	}

	return ErrResp(err, ErrValidateReq, extra...)
}
func usecaseRouterErrMapper(ctx *gin.Context, err error) {
	stsCode := http.StatusInternalServerError
	code := ErrProcessingReq
	if e, ok := err.(usecase.Code); ok {
		switch e.Code() {
		case usecase.ErrValidateReq:
			stsCode = http.StatusBadRequest
			code = ErrValidateReq
		}
	}
	extra := make(map[string]any)
	if e, ok := err.(usecase.Extra); ok {
		extra = e.Extra()
	}
	ctx.AbortWithStatusJSON(stsCode, ErrResp(err, code, func(e map[string]any) {
		for k, v := range extra {
			e[k] = v
		}
	}))
}
