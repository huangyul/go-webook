package main

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-webook/internal/events"
)

type App struct {
	server    *gin.Engine
	consumers []events.Consumer
}
