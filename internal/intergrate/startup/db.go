package startup

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13316)/go_blog?charset=utf8mb4&parseTime=True&loc=Local"))

	if err != nil {
		panic(err)
	}
	return db
}
