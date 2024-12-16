package main

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-blog/internal/event"
)

type App struct {
	server    *gin.Engine
	consumers []event.Consumer
}
