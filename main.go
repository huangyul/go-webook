package main

import (
	"fmt"
	"github.com/spf13/viper"
)

func main() {
	//log.Init()
	//defer log.Sync()
	initViper()

	app := InitApp()

	consumers := app.consumers

	for _, c := range consumers {
		err := c.Start()
		if err != nil {
			panic(err)
		}
	}

	jobs := app.jobs
	jobs.Start()
	defer func() {
		<-jobs.Stop().Done()
	}()

	addr := viper.GetString("server.addr")
	if addr == "" {
		addr = "8088"
	}
	server := app.server
	err := server.Run(fmt.Sprintf("127.0.0.1:%s", addr))
	if err != nil {
		panic(err)
	}
}

func initViper() {
	viper.SetConfigName("config")
	viper.AddConfigPath("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
