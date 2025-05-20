package main

import "github.com/spf13/viper"

func main() {
	initViper()

	app := InitApp()

	for _, c := range app.consumers {
		c.Start()
	}

	err := app.server.Serve()
	if err != nil {
		panic(err)
	}
}

func initViper() {
	viper.AddConfigPath("./config")
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		panic(err)
	}
}
