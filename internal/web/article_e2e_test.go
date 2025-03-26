//go:build e2e

package web

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestArticleHandler_Save(t *testing.T) {
	db := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook_test?charset=utf8mb4&parseTime=True&loc=Local"))
}
