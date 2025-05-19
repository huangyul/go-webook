package ioc

import (
	dao2 "github.com/huangyul/go-webook/interactive/repository/dao"
	"github.com/huangyul/go-webook/internal/repository/dao"
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
	dao.InitTable(db)
	dao2.InitTables(db)
	return db
}
