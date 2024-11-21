package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/huangyul/go-blog/internal/pkg/ginx/validator"
	"github.com/huangyul/go-blog/internal/repository"
	"github.com/huangyul/go-blog/internal/repository/dao"
	"github.com/huangyul/go-blog/internal/service"
	"github.com/huangyul/go-blog/internal/web"
	"github.com/huangyul/go-blog/internal/web/middleware"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	db := InitDB()

	server := InitServer()

	InitUseWeb(server, db)

	server.Run("127.0.0.1:8088")
}

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13306)/go_blog?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	dao.InitTable(db)
	return db
}

func InitServer() *gin.Engine {
	server := gin.Default()

	server.Use(sessions.Sessions("mysession", cookie.NewStore([]byte("secret"))))
	server.Use((middleware.LoginMiddleBuilder{}).Build())

	return server
}

func InitUseWeb(server *gin.Engine, db *gorm.DB) {
	uDao := dao.NewUserDAO(db)
	uRepo := repository.NewUserRepository(uDao)
	uSvc := service.NewUserService(uRepo)
	h := web.NewUserHandler(uSvc)
	h.RegisterRoutes(server)
}
