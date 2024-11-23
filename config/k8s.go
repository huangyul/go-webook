//go:build k8s

package config

var Config = Config{
	DB: DBConfig{
		DSN: "",
	},
	Redis: RedisConfig{
		Addr: "",
	},
}
