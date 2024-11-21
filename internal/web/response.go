package web

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func WriteResponse(ctx *gin.Context, code int, message string, data interface{}) {
	ctx.JSON(http.StatusOK, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}

func WriteError(ctx *gin.Context, httpCode int, code int, message string, data interface{}) {
	ctx.JSON(httpCode, Response{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
