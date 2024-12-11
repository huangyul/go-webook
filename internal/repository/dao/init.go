package dao

import "gorm.io/gorm"

// InitTable init tables
// A bad practice
func InitTable(db *gorm.DB) {
	err := db.AutoMigrate(&User{}, &Article{}, &PublishedArticle{}, &Interactive{})
	if err != nil {
		panic(err)
	}
}
