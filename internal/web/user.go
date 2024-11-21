package web

import (
	"net/http"

	regexp "github.com/dlclark/regexp2"
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-blog/internal/service"

	"github.com/gin-contrib/sessions"
)

var (
	emailPattern    = `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	passwordPattern = `^[a-zA-Z0-9]{6,18}$`
)

type UserHandler struct {
	emailRexExp    *regexp.Regexp
	passwordRexExp *regexp.Regexp

	svc service.UserService
}

func NewUserHandler(svc service.UserService) *UserHandler {
	return &UserHandler{
		emailRexExp:    regexp.MustCompile(emailPattern, regexp.None),
		passwordRexExp: regexp.MustCompile(passwordPattern, regexp.None),
		svc:            svc,
	}
}

func (h *UserHandler) RegisterRoutes(server *gin.Engine) {
	ug := server.Group("/user")
	ug.POST("/signup", h.Signup)
	ug.POST("/login", h.Login)
	ug.POST("/edit", h.Edit)
	ug.GET("/profile", h.Profile)
}

func (h *UserHandler) Signup(ctx *gin.Context) {

	type Req struct {
		Email           string `json:"email" binding:"required"`
		Password        string `json:"password" binding:"required"`
		ConfirmPassword string `json:"confirm_password" binding:"required"`
	}

	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if req.Password != req.ConfirmPassword {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "password does not match"})
		return
	}

	ok, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "illegal email"})
		return
	}

	ok, err = h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if !ok {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Password are only contain letters and numbers"})
		return
	}

	err = h.svc.Signup(ctx, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{})
}

func (h *UserHandler) Login(ctx *gin.Context) {
	type Req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := h.svc.Login(ctx, req.Email, req.Password)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	sess := sessions.Default(ctx)
	sess.Set("user_id", user.ID)
	sess.Options(sessions.Options{
		MaxAge: 86400,
	})
	if err := sess.Save(); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "login success",
	})
}

func (h *UserHandler) Edit(ctx *gin.Context) {

}

func (h *UserHandler) Profile(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "profile",
	})

}
