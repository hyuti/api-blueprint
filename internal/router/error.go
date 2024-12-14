package router

import (
	"errors"
	pkgerr "github.com/hyuti/api-blueprint/pkg/error"
	"golang.org/x/exp/slog"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

const (
	internalErrorMsg    = "Something went wrong, please check server logs for detail"
	emptyBodyRequestMsg = "Body request must not be empty"
)

// @Description Response of API if error occurs
type ErrResponse struct {
	Code  any            `json:"code"`    // Error code
	Msg   string         `json:"message"` // Description error
	Extra map[string]any `json:"extra"`   // Extra info about error
} // @name Error-Response

var _ error = (*ErrResponse)(nil)

func (r *ErrResponse) Error() string {
	return r.Msg
}

func handleError(
	ctx *gin.Context,
	lgr *slog.Logger,
	err error,
) {
	code := http.StatusInternalServerError

	var myErr *pkgerr.Error
	if !errors.As(err, &myErr) {
		myErr = pkgerr.ErrInternalServer(err)
		if errors.Is(err, io.EOF) {
			myErr = pkgerr.ErrValidatingRequest(errors.New(emptyBodyRequestMsg))
		}
	}

	switch {
	case errors.Is(myErr, pkgerr.LabelErrValidatingRequest):
		code = http.StatusBadRequest
	case errors.Is(myErr, pkgerr.LabelErrAuthenticateRequest):
		code = http.StatusForbidden
	case errors.Is(myErr, pkgerr.LabelErrAuthorizeRequest):
		code = http.StatusUnauthorized
	}

	if code != http.StatusInternalServerError {
		ctx.AbortWithStatusJSON(code, &ErrResponse{
			Code:  code,
			Msg:   myErr.Error(),
			Extra: myErr.Extra(),
		})
		return
	}

	// TODO: trigger github issue creation flow
	lgr.ErrorContext(
		ctx.Request.Context(),
		"error internal server",
		"error", myErr.Error(),
		"func", myErr.NameFunc(),
		"payload", myErr.Payload(),
		"chain", myErr.Chain(),
		"path", ctx.FullPath(),
		"controller", ctx.HandlerName(),
		"params", ctx.Params,
		"query", ctx.Request.URL.String(),
	)

	ctx.AbortWithStatusJSON(http.StatusInternalServerError, ErrResponse{
		Code: http.StatusInternalServerError,
		Msg:  internalErrorMsg,
	})
	return
}
