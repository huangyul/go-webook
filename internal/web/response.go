package web

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-blog/internal/pkg/errno"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type ListResp[T any] struct {
	Datas []T `json:"datas"`
	Total int `json:"total"`
}

func WriteResponse(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func WriteError(ctx *gin.Context, code int, message string) {
	ctx.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}

func WriteErrno(ctx *gin.Context, err *errno.Errno) {
	WriteError(ctx, err.Code, err.Message)
}

func WriteSuccess(ctx *gin.Context, data interface{}) {
	WriteResponse(ctx, 0, "success", data)
}
