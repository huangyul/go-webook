package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	addr := viper.GetString("mysql.addr")
	db, er := gorm.Open(mysql.Open(addr))
	if er != nil {
		panic(er)
	}
	return db
}
