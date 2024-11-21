package dao

import "gorm.io/gorm"

// InitTable init tables
//
// A bad practice
func InitTable(db *gorm.DB) {
	db.AutoMigrate(&User{})
}
