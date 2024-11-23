package web

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"net/http"
	"time"

	"github.com/huangyul/go-blog/internal/pkg/errno"
	"github.com/huangyul/go-blog/internal/pkg/ginx/validator"

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
		WriteError(ctx, errno.ErrBadRequest.Code, validator.Translate(err))
		return
	}

	if req.Password != req.ConfirmPassword {
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage("password not match"))
		return
	}

	ok, err := h.emailRexExp.MatchString(req.Email)
	if err != nil {
		WriteErrno(ctx, errno.ErrInternalServer)
		return
	}
	if !ok {
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage("illegal email"))
		return
	}

	ok, err = h.passwordRexExp.MatchString(req.Password)
	if err != nil {
		WriteErrno(ctx, errno.ErrInternalServer.SetMessage(err.Error()))
		return
	}
	if !ok {
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage("Password are only contain letters and numbers"))
		return
	}

	err = h.svc.Signup(ctx, req.Email, req.Password)
	if err != nil {
		WriteErrno(ctx, errno.ErrInternalServer.SetMessage(err.Error()))
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
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage(validator.Translate(err)))
		return
	}
	user, err := h.svc.Login(ctx, req.Email, req.Password)
	if errors.Is(err, errno.ErrNotFoundUser) {
		WriteErrno(ctx, errno.ErrEmailOrPasswordError)
		return
	}
	if err != nil {
		WriteErrno(ctx, errno.ErrInternalServer.SetMessage(err.Error()))
		return
	}

	// set login token
	// type 1 session
	sess := sessions.Default(ctx)
	sess.Set("user_id", user.ID)
	sess.Options(sessions.Options{
		MaxAge: 86400,
	})
	if err := sess.Save(); err != nil {
		WriteErrno(ctx, errno.ErrInternalServer.SetMessage(err.Error()))
		return
	}
	// type 2 jwt
	c := JWTClaims{
		UserID:    int(user.ID),
		UserAgent: ctx.Request.UserAgent(),
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 2)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, c)
	tokenStr, err := token.SignedString([]byte("JWT_TOKEN_KEY"))
	if err != nil {
		WriteErrno(ctx, errno.ErrInternalServer.SetMessage(err.Error()))
	}

	type LoginResp struct {
		Token string `json:"token"`
	}

	WriteSuccess(ctx, LoginResp{
		Token: tokenStr,
	})
}

func (h *UserHandler) Edit(ctx *gin.Context) {

}

func (h *UserHandler) Profile(ctx *gin.Context) {
	ctx.JSON(http.StatusOK, gin.H{
		"msg": "profile",
	})

}

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID    int
	UserAgent string
}
