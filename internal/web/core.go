package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type ApiResponse[T any] struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
	Data T      `json:"data"`
}

func writeSuccess[T any](ctx *gin.Context, data T) {
	ctx.JSON(http.StatusOK, ApiResponse[T]{
		Code: 0,
		Msg:  "success",
		Data: data,
	})
}

func writeError[T any](ctx *gin.Context, err error) {
	ctx.JSON(http.StatusOK, ApiResponse[T]{
		Code: 1,
		Msg:  err.Error(),
	})
}
