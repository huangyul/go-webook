package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-webook/internal/repository"
	"github.com/huangyul/go-webook/internal/repository/dao"
	"github.com/huangyul/go-webook/internal/service"
	"github.com/huangyul/go-webook/internal/web"
	"github.com/huangyul/go-webook/internal/web/middleware"
	"github.com/huangyul/go-webook/pkg/ginx/middleware/ratelimit"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	server := gin.Default()
	db := initDB()
	dao.InitTable(db)
	redis := initRedis()

	server.Use(
		ratelimit.NewBuilder(redis).Build(),
		middleware.NewJWTLoginMiddlewareBuild(
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

func initRedis() redis.Cmdable {
	cmd := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:16379",
		Password: "",
		DB:       0,
	})
	if err := cmd.Ping(context.Background()).Err(); err != nil {
		panic(err)
	}
	return cmd
}
