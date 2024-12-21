package dao

import "gorm.io/gorm"

func InitTable(db *gorm.DB) {
	err := db.AutoMigrate(&Interactive{}, &UserLikeBiz{}, &UserCollectionBiz{})
	if err != nil {
		panic(err)
	}
}
