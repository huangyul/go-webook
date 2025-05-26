package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB) {
	err := db.AutoMigrate(&Payment{})
	if err != nil {
		panic(err)
	}
}
