package ioc

import (
	"log"
	"os"
	"time"

	"github.com/huangyul/go-blog/internal/repository/dao"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func InitDB() *gorm.DB {

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second, // Slow SQL threshold
			LogLevel:                  logger.Info, // Log level
			IgnoreRecordNotFoundError: true,        // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      false,       // Don't include params in the SQL log
			Colorful:                  false,       // Disable color
		},
	)

	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local"), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		panic(err)
	}
	dao.InitTable(db)
	return db
}
