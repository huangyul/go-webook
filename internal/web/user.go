package web

import (
	"errors"
	"fmt"
	regexp "github.com/dlclark/regexp2"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/service"
	"github.com/huangyul/go-webook/internal/web/middleware"
	"net/http"
	"time"
)

var (
	passwordPattern = `^(?=.*[A-Za-z])(?=.*\d)[A-Za-z\d]{6,18}$`
	emailPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
)

const biz = "login"

type UserHandler struct {
	svc     service.UserService
	codeSvc service.CodeService
}

func NewUserHandler(svc service.UserService, code service.CodeService) *UserHandler {
	return &UserHandler{svc: svc, codeSvc: code}
}

func (hdl *UserHandler) RegisterRoutes(g *gin.Engine) {
	ug := g.Group("/user")
	{
		ug.POST("/register", hdl.Register)
		//ug.POST("/login", hdl.Login)
		ug.POST("/login", hdl.JWTLogin)
		ug.GET("/profile", hdl.Profile)
		ug.POST("/edit", hdl.Edit)
		ug.POST("/sms/send", hdl.SendSMS)
		ug.POST("/sms/login", hdl.SMSLogin)
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

	tokenStr, err := hdl.setToken(user)

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
	userId := ctx.GetInt64("user_id")
	user, err := hdl.svc.FindById(ctx, userId)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"user": struct {
		ID        int64  `json:"id"`
		Nickname  string `json:"nickname"`
		Birthday  string `json:"birthday"`
		AboutMe   string `json:"about_me"`
		CreatedAt string `json:"created_at"`
	}{
		ID:        user.ID,
		Nickname:  user.Nickname,
		Birthday:  formatTime(user.Birthday),
		AboutMe:   user.AboutMe,
		CreatedAt: formatTime(user.CreatedAt),
	}})
}

func (hdl *UserHandler) Edit(ctx *gin.Context) {
	type Req struct {
		Nickname string `json:"nickname"`
		Birthday string `json:"birthday"`
		AboutMe  string `json:"about_me"`
	}
	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userId := ctx.GetInt64("user_id")
	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "birthday format error"})
		return
	}
	err = hdl.svc.Update(ctx, &domain.User{
		ID:       userId,
		Nickname: req.Nickname,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "user successfully profile"})
}

func (hdl *UserHandler) SendSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone" binging:"required" binding:"required"`
	}
	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := hdl.codeSvc.Send(ctx, biz, req.Phone)
	if errors.Is(err, service.ErrCodeSendTooMany) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "code send too many"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "send sms successfully"})
}

func (hdl *UserHandler) SMSLogin(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	ok, err := hdl.codeSvc.Verify(ctx, biz, req.Phone, req.Code)
	if errors.Is(err, service.ErrCodeVerifyTooMany) {
		ctx.JSON(http.StatusOK, gin.H{"msg": "code verify too many"})
		return
	}
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		ctx.JSON(http.StatusOK, gin.H{"msg": "code verify failed"})
		return
	}
	user, err := hdl.svc.FindOrCreateByPhone(ctx, req.Phone)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	tokenStr, err := hdl.setToken(user)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"token": tokenStr})
}

func (hdl *UserHandler) setToken(user *domain.User) (string, error) {
	loginClaim := middleware.LoginClaims{
		UserId: user.ID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
		},
	}
	return jwt.NewWithClaims(jwt.SigningMethodHS256, loginClaim).SignedString([]byte("secret"))
}
