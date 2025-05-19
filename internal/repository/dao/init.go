package dao

import "gorm.io/gorm"

func InitTable(db *gorm.DB) {
	err := db.AutoMigrate(&User{}, &Article{}, &PubArticle{}, &History{}, &Job{})
	if err != nil {
		panic(err)
	}
}
