package response

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

type response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func HandleSuccess(ctx *gin.Context, data interface{}) {
	if data == nil {
		data = map[string]string{}
	}
	resp := response{Code: ErrSuccess.Code, Message: ErrSuccess.Message, Data: data}
	ctx.JSON(http.StatusOK, resp)
}

func HandleError(ctx *gin.Context, httpCode int, err error, data interface{}) {
	if data == nil {
		data = map[string]string{}
	}
	resp := response{Code: errorCodeMap[err], Message: err.Error(), Data: data}
	ctx.JSON(httpCode, resp)
}

type Error struct {
	Code    int
	Message string
}

var errorCodeMap = map[error]int{}

func newError(code int, msg string) *Error {
	err := errors.New(msg)
	errorCodeMap[err] = code
	return &Error{
		Code:    code,
		Message: msg,
	}
}
func (e Error) Error() string {
	return e.Message
}
