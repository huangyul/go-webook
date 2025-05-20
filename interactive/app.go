package main

import (
	"github.com/huangyul/go-webook/internal/events"
	"github.com/huangyul/go-webook/pkg/grpcx"
)

type App struct {
	server    *grpcx.Server
	consumers []events.Consumer
}
