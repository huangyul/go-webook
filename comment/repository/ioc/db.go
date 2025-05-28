package ioc

import (
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	addr := viper.GetString("mysql.addr")
	db, err := gorm.Open(mysql.Open(addr))
	if err != nil {
		panic(err)
	}
	// dao.InitDB(db)
	return db
}
