package web

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/service"
)

const (
	Biz = "article"
)

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

func (a *ArticleHandler) RegisterRoutes(g *gin.Engine) {
	ug := g.Group("article")
	{
		ug.POST("save", a.Save)
		ug.POST("publish", a.Publish)
		ug.GET("withdraw", a.Withdraw)
		ug.GET("detail/:id", a.Detail)
		ug.POST("list", a.GetByAuthor)

		pug := ug.Group("pub")
		{
			pug.GET("/detail/:id", a.PubDetail)
			pug.POST("like", a.Like)
			pug.POST("collect", a.Collect)
		}

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
	userId := ctx.MustGet("user_id").(int64)
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
	userId := ctx.MustGet("user_id").(int64)
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
	writeSuccess[any](ctx, nil)
}

func (a *ArticleHandler) Withdraw(ctx *gin.Context) {
	idStr := ctx.Query("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError[any](ctx, errors.New("illegal id"))
		return
	}
	userId := ctx.MustGet("user_id").(int64)
	err = a.svc.WithDraw(ctx, userId, id)
	if err != nil {
		writeError[any](ctx, err)
		return
	}
	writeSuccess[any](ctx, nil)
}

func (a *ArticleHandler) Detail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError[any](ctx, errors.New("illegal id"))
		return
	}
	userId := ctx.MustGet("user_id").(int64)
	art, err := a.svc.GetById(ctx, id, userId)
	if err != nil {
		writeError[any](ctx, err)
		return
	}
	writeSuccess[ArtItemRes](ctx, a.toResItem(art))
}

func (a *ArticleHandler) GetByAuthor(ctx *gin.Context) {
	type Req struct {
		Page     int64 `json:"page"`
		PageSize int64 `json:"page_size"`
	}
	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Page < 1 {
		req.Page = 1
	}
	if req.PageSize < 1 {
		req.PageSize = 10
	}
	userId := ctx.MustGet("user_id").(int64)
	arts, err := a.svc.GetByAuthor(ctx, userId, req.Page, req.PageSize)
	if err != nil {
		writeError[any](ctx, err)
		return
	}

	res := make([]ArtItemRes, 0)
	for _, art := range arts {
		res = append(res, a.toResItem(art))
	}
	writeSuccess[[]ArtItemRes](ctx, res)
}

func (a *ArticleHandler) PubDetail(ctx *gin.Context) {
	idStr := ctx.Param("id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		writeError[any](ctx, errors.New("illegal id"))
		return
	}
	userId := ctx.MustGet("user_id").(int64)
	art, err := a.svc.GetPudById(ctx, id, userId)
	if err != nil {
		writeError[any](ctx, err)
		return
	}
	go func() {
		_ = a.interSvc.IncrReadCnt(ctx, Biz, art.Id)
	}()
	res := a.toResItem(art)
	inter, err := a.interSvc.Get(ctx, Biz, res.Id, userId)
	if err != nil {
		writeError[any](ctx, err)
		return
	}
	res.CollectCnt = inter.CollectCnt
	res.LikeCnt = inter.LikeCnt
	res.ReadCnt = inter.ReadCnt
	res.Liked = inter.Liked
	res.Collected = inter.Collectd

	writeSuccess[ArtItemRes](ctx, res)
}

type ArtItemRes struct {
	Id         int64  `json:"id"`
	Title      string `json:"title"`
	Content    string `json:"content"`
	AuthorId   int64  `json:"author_id"`
	AuthorName string `json:"author_name"`
	Status     uint8  `json:"status"`
	CollectCnt int64  `json:"collect_cnt"`
	LikeCnt    int64  `json:"like_cnt"`
	ReadCnt    int64  `json:"read_cnt"`
	Liked      bool   `json:"liked"`
	Collected  bool   `json:"collected"`
	CreateTime string `json:"create_time"`
	UpdateTime string `json:"update_time"`
}

func (a *ArticleHandler) toResItem(art *domain.Article) ArtItemRes {
	return ArtItemRes{
		Id:         art.Id,
		Title:      art.Title,
		Content:    art.Content,
		AuthorId:   art.Author.Id,
		AuthorName: art.Author.Name,
		Status:     art.Status.ToUint8(),
		CreateTime: art.CreatedAt.Format(time.DateTime),
		UpdateTime: art.UpdatedAt.Format(time.DateTime),
	}
}

func (a *ArticleHandler) Like(ctx *gin.Context) {
	type Req struct {
		Id   int64 `json:"id" binding:"required"`
		Like bool  `json:"like" binding:"required"`
	}
	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := ctx.MustGet("user_id").(int64)
	var err error
	if req.Like {
		err = a.interSvc.Like(ctx, Biz, req.Id, userId)
	} else {
		err = a.interSvc.CancelLike(ctx, Biz, req.Id, userId)
	}
	if err != nil {
		writeError[any](ctx, err)
		return
	}
	writeSuccess[any](ctx, nil)
}

func (a *ArticleHandler) Collect(ctx *gin.Context) {
	type Req struct {
		Id      int64 `json:"id" binding:"required"`
		Collect bool  `json:"collect"`
	}
	var req Req
	if err := ctx.ShouldBind(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	userId := ctx.MustGet("user_id").(int64)
	var err error
	if req.Collect {
		err = a.interSvc.Collect(ctx, Biz, req.Id, userId)
	} else {
		err = a.interSvc.CancelCollect(ctx, Biz, req.Id, userId)
	}
	if err != nil {
		writeError[any](ctx, err)
		return
	}
	writeSuccess[any](ctx, nil)
}
