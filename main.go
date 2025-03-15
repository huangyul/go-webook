package main

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-webook/internal/repository"
	"github.com/huangyul/go-webook/internal/repository/dao"
	"github.com/huangyul/go-webook/internal/service"
	"github.com/huangyul/go-webook/internal/web"
	"github.com/huangyul/go-webook/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	server := gin.Default()
	db := initDB()
	dao.InitTable(db)

	server.Use(middleware.NewJWTLoginMiddlewareBuild(
		middleware.AddWhiteList("/user/login", "/user/register")).Build())

	userDao := dao.NewUserDAO(db)
	userRepo := repository.NewUserRepository(userDao)
	userService := service.NewUserService(userRepo)
	userHandler := web.NewUserHandler(userService)
	userHandler.RegisterRoutes(server)

	err := server.Run("127.0.0.1:8088")
	if err != nil {
		panic(err)
	}

}

func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	return db
}
