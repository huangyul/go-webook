package main

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	_ "github.com/huangyul/go-blog/internal/pkg/ginx/validator"
	"github.com/huangyul/go-blog/internal/repository"
	"github.com/huangyul/go-blog/internal/repository/cache"
	"github.com/huangyul/go-blog/internal/repository/dao"
	"github.com/huangyul/go-blog/internal/service"
	"github.com/huangyul/go-blog/internal/service/sms/localstorage"
	"github.com/huangyul/go-blog/internal/web"
	"github.com/huangyul/go-blog/internal/web/middleware"
	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {

	db := InitDB()

	cmd := InitRedis()

	server := InitServer()

	InitUseWeb(server, db, cmd)

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

func InitRedis() redis.Cmdable {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "127.0.0.1:6379",
		DB:       0,
		Password: "",
	})
	return redisClient
}

func InitServer() *gin.Engine {
	server := gin.Default()

	//cmd := InitRedis()

	//server.Use(ratelimit.NewBuilder(cmd, 10, time.Second*60).Build())

	server.Use(sessions.Sessions("mysession", cookie.NewStore([]byte("secret"))))
	//server.Use((middleware.LoginMiddleBuilder{}).Build())

	server.Use(middleware.NewJWTLoginMiddlewareBuild().AddWhiteList("/user/login", "/user/signup", "/user/login-sms").Build())

	return server
}

func InitUseWeb(server *gin.Engine, db *gorm.DB, cmd redis.Cmdable) {
	cCache := cache.NewRedisCodeCache(cmd)
	cRepo := repository.NewCodeRepository(cCache)
	cSvc := service.NewCodeService(cRepo, localstorage.NewSmsLocalStorageService())

	uDao := dao.NewUserDAOGORM(db)
	uCache := cache.NewRedisUserCache(cmd)
	uRepo := repository.NewUserRepository(uDao, uCache)
	uSvc := service.NewUserService(uRepo, cSvc)
	h := web.NewUserHandler(uSvc)
	h.RegisterRoutes(server)
}
