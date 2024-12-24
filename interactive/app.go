package main

import (
	"github.com/huangyul/go-blog/internal/event"
	"github.com/huangyul/go-blog/pkg/grpcx"
)

type App struct {
	consumers []event.Consumer
	server    *grpcx.Server
}
