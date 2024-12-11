package web

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-blog/internal/domain"
	"github.com/huangyul/go-blog/internal/pkg/errno"
	"github.com/huangyul/go-blog/internal/pkg/ginx/validator"
	"github.com/huangyul/go-blog/internal/service"
	"strconv"
)

const biz = "article"

type ArticleHandler struct {
	svc      service.ArticleService
	interSvc service.InteractiveService
}

func NewArticleHandler(svc service.ArticleService, interSvc service.InteractiveService) *ArticleHandler {
	return &ArticleHandler{
		svc:      svc,
		interSvc: interSvc,
	}
}

func (h *ArticleHandler) RegisterRoutes(g *gin.Engine) {
	ug := g.Group("/article")
	ug.POST("/edit", h.Edit)
	ug.POST("/publish", h.Publish)
	ug.GET("/withdraw/:id", h.Withdraw)
	ug.POST("/list", h.List)
	ug.GET("/detail/:id", h.Detail)
	pg := ug.Group("/pub")
	pg.GET("/detail/:id", h.PubDetail)
}

type ArticleEditReq struct {
	ID      int64  `json:"id"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}

func (h *ArticleHandler) Edit(ctx *gin.Context) {

	var req ArticleEditReq
	if err := ctx.ShouldBind(&req); err != nil {
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage(validator.Translate(err)))
		return
	}
	userId := ctx.MustGet("user_id").(int64)
	id, err := h.svc.Save(ctx, domain.Article{
		ID:      req.ID,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			ID: userId,
		},
	})
	if err != nil {
		WriteErrno(ctx, errno.ErrInternalServer.SetMessage(err.Error()))
		return
	}
	WriteSuccess(ctx, gin.H{"id": id})

}

func (h *ArticleHandler) Publish(ctx *gin.Context) {
	type Req struct {
		ID      int64  `json:"id"`
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage(validator.Translate(err)))
		return
	}
	userId := ctx.MustGet("user_id").(int64)
	id, err := h.svc.Publish(ctx, domain.Article{
		ID:      req.ID,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			ID: userId,
		},
	})
	if err != nil {
		WriteErrno(ctx, errno.ErrInternalServer.SetMessage(err.Error()))
		return
	}
	WriteSuccess(ctx, gin.H{"id": id})
}

func (h *ArticleHandler) Withdraw(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage("id illegal"))
		return
	}
	userId := ctx.MustGet("user_id").(int64)
	err = h.svc.Withdraw(ctx, userId, id)
	if err != nil {
		WriteErrno(ctx, errno.ErrInternalServer.SetMessage(err.Error()))
		return
	}
	WriteSuccess(ctx, nil)
}

func (h *ArticleHandler) List(ctx *gin.Context) {
	var req Page
	if err := ctx.ShouldBind(&req); err != nil {
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage(validator.Translate(err)))
		return
	}
	req.SetDefault()

	userId := ctx.MustGet("user_id").(int64)
	data, err := h.svc.List(ctx, userId, req.Page, req.PageSize)
	if err != nil {
		WriteErrno(ctx, errno.ErrInternalServer.SetMessage(err.Error()))
		return
	}
	WriteSuccess(ctx, gin.H{"data": data})
}

func (h *ArticleHandler) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage("id illegal"))
		return
	}
	userId := ctx.MustGet("user_id").(int64)
	art, err := h.svc.Detail(ctx, userId, id)
	if err != nil {
		WriteErrno(ctx, errno.ErrInternalServer.SetMessage(err.Error()))
		return
	}
	WriteSuccess(ctx, gin.H{"data": art})
}

func (h *ArticleHandler) PubDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		WriteErrno(ctx, errno.ErrBadRequest.SetMessage("id illegal"))
		return
	}
	userId := ctx.MustGet("user_id").(int64)
	art, err := h.svc.PubDetail(ctx, userId, id)
	if err != nil {
		WriteErrno(ctx, errno.ErrInternalServer.SetMessage(err.Error()))
		return
	}
	go func() {
		h.interSvc.IncrReadCnt(ctx, art.ID, biz)
	}()
	WriteSuccess(ctx, gin.H{"data": art})
}
