package web

import (
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/service"
	"net/http"
)

type ArticleHandler struct {
	svc service.ArticleService
}

func NewArticleHandler(svc service.ArticleService) *ArticleHandler {
	return &ArticleHandler{
		svc: svc,
	}
}

func (a *ArticleHandler) Register(g *gin.Engine) {
	ug := g.Group("article")
	{
		ug.POST("save", a.Save)
		ug.POST("publish", a.Publish)
	}
}

func (a *ArticleHandler) Save(ctx *gin.Context) {
	type Req struct {
		Id      int64  `json:"id"`
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := ctx.MustGet("userId").(int64)
	_, err := a.svc.Save(ctx, &domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: userId,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	writeSuccess[any](ctx, nil)
}

func (a *ArticleHandler) Publish(ctx *gin.Context) {
	type Req struct {
		Id      int64  `json:"id"`
		Title   string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := ctx.MustGet("userId").(int64)
	err := a.svc.Publish(ctx, &domain.Article{
		Id:      req.Id,
		Title:   req.Title,
		Content: req.Content,
		Author: domain.Author{
			Id: userId,
		},
	})
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"msg": "publish success"})
}
