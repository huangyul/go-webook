package dao

import "gorm.io/gorm"

func InitTables(db *gorm.DB) {
	err := db.AutoMigrate(&Interactive{}, &UserLikeBiz{}, &UserCollectBiz{})
	if err != nil {
		panic(err)
	}
}
