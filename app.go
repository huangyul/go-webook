package main

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-webook/internal/events"
	"github.com/robfig/cron/v3"
)

// App
type App struct {
	server    *gin.Engine
	consumers []events.Consumer
	cron      *cron.Cron
}
