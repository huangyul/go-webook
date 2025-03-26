package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {

	dbAddr := viper.GetString("mysql.addr")

	if dbAddr == "" {
		dbAddr = "root:root@tcp(localhost:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local"
	}

	db, err := gorm.Open(mysql.Open(dbAddr))
	if err != nil {
		panic(err)
	}
	return db
}
