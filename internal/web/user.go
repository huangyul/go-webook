package web

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/huangyul/go-blog/internal/domain"

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
	ug.GET("/list", h.List)

	// code
	ug.GET("/login-sms", h.SendCode)
	ug.POST("/login-sms", h.LoginSMS)
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
		UserID:    user.ID,
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
	type Req struct {
		Birthday string `json:"birthday" binding:"required"`
		AboutMe  string `json:"about_me" binding:"required"`
		Nickname string `json:"nickname" binding:"required"`
	}
	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage(validator.Translate(err)))
		return
	}

	birthday, err := time.Parse(time.DateOnly, req.Birthday)
	if err != nil {
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage("invalid birthday"))
		return
	}

	id, _ := ctx.Get("user_id")
	uID, ok := id.(int64)
	if !ok {
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage("invalid user id"))
		return
	}

	err = h.svc.UpdateUserInfo(ctx, domain.User{
		ID:       uID,
		Birthday: birthday,
		AboutMe:  req.AboutMe,
		Nickname: req.Nickname,
	})

	if err != nil {
		WriteErrno(ctx, errno.ErrInternalServer.SetMessage(err.Error()))
		return
	}

	WriteSuccess(ctx, nil)

}

func (h *UserHandler) Profile(ctx *gin.Context) {
	id, _ := ctx.Get("user_id")

	uID := id.(int64)

	u, err := h.svc.GetUserInfo(ctx, uID)

	if err != nil {
		WriteErrno(ctx, errno.ErrInternalServer.SetMessage(err.Error()))
		return
	}

	WriteSuccess(ctx, userResp{
		ID:       u.ID,
		Email:    u.Email,
		Nickname: u.Nickname,
		Birthday: u.Birthday.Format(time.DateOnly),
		AboutMe:  u.AboutMe,
		CreateAt: u.CreatedAt.Format(time.DateOnly),
		UpdateAt: u.UpdatedAt.Format(time.DateOnly),
	})

}

func (h *UserHandler) List(ctx *gin.Context) {
	pageStr := ctx.Query("page")
	pageSizeStr := ctx.Query("page_size")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage("page must be number"))
		return
	}
	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil {
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage("page_size must be number"))
		return
	}

	users, count, err := h.svc.GetUserList(ctx, page, pageSize)
	if err != nil {
		WriteErrno(ctx, errno.ErrInternalServer.SetMessage(err.Error()))
		return
	}

	var datas []userResp

	for _, u := range users {
		datas = append(datas, userResp{
			ID:       u.ID,
			Email:    u.Email,
			Nickname: u.Nickname,
			Birthday: u.Birthday.Format(time.DateOnly),
			AboutMe:  u.AboutMe,
			CreateAt: u.CreatedAt.Format(time.DateOnly),
			UpdateAt: u.UpdatedAt.Format(time.DateOnly),
		})
	}

	WriteSuccess(ctx, ListResp[userResp]{
		Datas: datas,
		Total: count,
	})

}

func (h *UserHandler) SendCode(ctx *gin.Context) {
	phone := ctx.Query("phone")
	if phone == "" {
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage("phone must not be empty"))
		return
	}
	err := h.svc.SendCode(ctx, phone)
	if err != nil {
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage(err.Error()))
		return
	}
	WriteErrno(ctx, errno.ErrOK.SetMessage("code send successfully"))
}

func (h *UserHandler) LoginSMS(ctx *gin.Context) {
	type Req struct {
		Phone string `json:"phone_number" binding:"required"`
		Code  string `json:"code" binding:"required"`
	}
	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage(validator.Translate(err)))
		return
	}
	user, err := h.svc.LoginSMS(ctx, req.Phone, req.Code)
	if err != nil {
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage(err.Error()))
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
		UserID:    user.ID,
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

type JWTClaims struct {
	jwt.RegisteredClaims
	UserID    int64
	UserAgent string
}

type userResp struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Nickname string `json:"nickname"`
	Birthday string `json:"birthday"`
	AboutMe  string `json:"about_me"`
	CreateAt string `json:"create_at"`
	UpdateAt string `json:"update_at"`
}
