package web

import (
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	intrv1 "github.com/huangyul/go-webook/api/proto/gen/intr/v1"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/service"
)

const (
	Biz = "article"
)

type ArticleHandler struct {
	svc        service.ArticleService
	interSvc   intrv1.InteractiveServiceClient
	historySvc service.HistoryService
	rankingSvc service.RankingService
}

func NewArticleHandler(svc service.ArticleService, interSvc intrv1.InteractiveServiceClient, historySvc service.HistoryService, rankingSvc service.RankingService) *ArticleHandler {
	return &ArticleHandler{
		svc:        svc,
		interSvc:   interSvc,
		historySvc: historySvc,
		rankingSvc: rankingSvc,
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
			pug.GET("history", a.History)
			pug.GET("top100", a.TopN)
		}

	}
}

func (a *ArticleHandler) TopN(ctx *gin.Context) {
	arts, err := a.rankingSvc.GetTopN(ctx)
	if err != nil {
		writeError[any](ctx, err)
		return
	}
	type Resp struct {
		Id        int64  `json:"id"`
		Title     string `json:"title"`
		Abstract  string `json:"abstract"`
		CreatedAt string `json:"created_at"`
	}
	var resp []Resp
	for _, item := range arts {
		resp = append(resp, Resp{
			Id:        item.Id,
			Title:     item.Title,
			Abstract:  item.Abstract(),
			CreatedAt: item.CreatedAt.Format(time.DateTime),
		})
	}
	writeSuccess(ctx, resp)

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
	art, err := a.svc.GetPudById(ctx, id, userId, Biz)
	if err != nil {
		writeError[any](ctx, err)
		return
	}
	res := a.toResItem(art)
	grpcRes, err := a.interSvc.Get(ctx, &intrv1.GetRequest{
		Biz:    Biz,
		BizId:  res.Id,
		UserId: userId,
	})
	if err != nil {
		writeError[any](ctx, err)
		return
	}
	inter := grpcRes.Intr
	res.CollectCnt = inter.CollectCnt
	res.LikeCnt = inter.LikeCnt
	res.ReadCnt = inter.ReadCnt
	res.Liked = inter.Liked
	res.Collected = inter.Collected

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
		_, err = a.interSvc.Like(ctx, &intrv1.LikeRequest{
			Biz:    Biz,
			BizId:  req.Id,
			UserId: userId,
		})
	} else {
		_, err = a.interSvc.CancelLike(ctx, &intrv1.CancelLikeRequest{
			Biz:    Biz,
			BizId:  req.Id,
			UserId: userId,
		})
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
		_, err = a.interSvc.Collect(ctx, &intrv1.CollectRequest{
			Biz:    biz,
			BizId:  req.Id,
			UserId: userId,
		})
	} else {
		_, err = a.interSvc.CancelCollect(ctx, &intrv1.CancelCollectRequest{
			Biz:    biz,
			BizId:  req.Id,
			UserId: userId,
		})
	}
	if err != nil {
		writeError[any](ctx, err)
		return
	}
	writeSuccess[any](ctx, nil)
}

func (a *ArticleHandler) History(ctx *gin.Context) {
	userId := ctx.MustGet("user_id").(int64)
	res, err := a.historySvc.GetListByUserId(ctx, userId)
	if err != nil {
		writeError[any](ctx, err)
		return
	}
	type Resp struct {
		Id         int64  `json:"id"`
		Title      string `json:"title"`
		ArticleId  int64  `json:"article_id"`
		AuthorName string `json:"author_name"`
		AuthorId   int64  `json:"author_id"`
		CreatedAt  string `json:"created_at"`
	}
	var resp []Resp
	for _, item := range res {
		resp = append(resp, Resp{
			Id:         item.Id,
			Title:      item.ArticleTitle,
			ArticleId:  item.ArticleId,
			AuthorName: item.AuthorName,
			AuthorId:   item.AuthorId,
			CreatedAt:  item.CreatedAt.Format(time.DateTime),
		})
	}
	writeSuccess[[]Resp](ctx, resp)
}
