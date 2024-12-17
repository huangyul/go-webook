package main

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-blog/internal/event"
	"github.com/robfig/cron/v3"
)

type App struct {
	server    *gin.Engine
	consumers []event.Consumer
	jobs      *cron.Cron
}
