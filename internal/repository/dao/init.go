package dao

import "gorm.io/gorm"

func InitTable(db *gorm.DB) {
	err := db.AutoMigrate(&User{}, &Article{}, &PubArticle{}, &Interactive{}, &UserLikeBiz{}, &UserCollectBiz{}, &History{})
	if err != nil {
		panic(err)
	}
}
