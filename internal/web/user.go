package web

import (
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/huangyul/go-webook/internal/service"
	"github.com/huangyul/go-webook/internal/web/middleware"
	"net/http"
	"time"
)

var (
	passwordPattern = `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{6,18}$`
	emailPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

type UserHandler struct {
	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{svc: svc}
}

func (hdl *UserHandler) RegisterRoutes(g *gin.Engine) {
	ug := g.Group("/user")
	{
		ug.POST("/register", hdl.Register)
		//ug.POST("/login", hdl.Login)
		ug.POST("/login", hdl.JWTLogin)
		ug.GET("/profile", hdl.Profile)
	}
}

func (hdl *UserHandler) Register(ctx *gin.Context) {
	type RegisterReq struct {
		Email           string `json:"email" binding:"required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}
	var req RegisterReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok, _ := regexp.MustCompile(emailPattern, regexp.None).MatchString(req.Email)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "email illegal"})
		return
	}
	ok, _ = regexp.MustCompile(passwordPattern, regexp.None).MatchString(req.Password)
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "password illegal"})
		return
	}
	if req.Password != req.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "password not match"})
		return
	}
	err := hdl.svc.RegisterByEmail(ctx, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "user successfully registered"})
}

func (hdl *UserHandler) JWTLogin(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var req LoginReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := hdl.svc.LoginByEmail(ctx, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	loginClaim := middleware.LoginClaims{
		UserId: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	tokenStr, err := jwt.NewWithClaims(jwt.SigningMethodHS256, loginClaim).SignedString([]byte("secret"))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"token": tokenStr})

}

func (hdl *UserHandler) Login(ctx *gin.Context) {
	type LoginReq struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	var req LoginReq
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := hdl.svc.LoginByEmail(ctx, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sess := sessions.Default(ctx)
	sess.Set("userId", user.ID)
	sess.Options(sessions.Options{
		MaxAge: 86400 * 30,
	})
	err = sess.Save()
	if err != nil {
		fmt.Sprintf("session save error: %v", err)
	}

	ctx.JSON(http.StatusOK, gin.H{"msg": "user successfully login"})
}

func (hdl *UserHandler) Profile(ctx *gin.Context) {

	id := ctx.GetInt64("user_id")
	if id == 0 {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "user id illegal"})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": fmt.Sprintf("user successfully get user %d", id)})
}
