package utils

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-blog/internal/pkg/errno"
)

var Respose struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

func WriteErr(ctx *gin.Context, err *errno.Errno) {
	Respose.Code = err.Code
	Respose.Msg = err.Message
	ctx.JSON(err.HTTP, Respose)
}
