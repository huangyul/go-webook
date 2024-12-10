package intergrate

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-blog/internal/domain"
	"github.com/huangyul/go-blog/internal/intergrate/startup"
	"github.com/huangyul/go-blog/internal/pkg/errno"
	"github.com/huangyul/go-blog/internal/repository/dao"
	"github.com/huangyul/go-blog/internal/web"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

const userID = 1

func TestArticle_Edit(t *testing.T) {
	db := startup.InitDB()
	server := gin.Default()
	handler := startup.InitArticleHandler()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("user_id", int64(userID))
	})
	handler.RegisterRoutes(server)
	tests := []struct {
		name          string
		req           web.ArticleEditReq
		beforeRequest func(t *testing.T)
		afterRequest  func(t *testing.T)
		wantCode      int
		wantBody      web.Response
	}{
		{
			name: "all new success",
			req: web.ArticleEditReq{
				Title:   "test",
				Content: "test",
			},
			beforeRequest: func(t *testing.T) {},
			afterRequest: func(t *testing.T) {
				var art dao.Article
				err := db.Where("id = ? AND author_id = ?", 1, userID).First(&art).Error
				assert.NoError(t, err)
				assert.Equal(t, "test", art.Title)
				assert.Equal(t, "test", art.Content)
				assert.Equal(t, domain.ArticleStatusUnPublished, art.Status)
				assert.Equal(t, int64(userID), art.AuthorID)
			},
			wantCode: http.StatusOK,
			wantBody: web.Response{
				Code:    0,
				Message: "success",
				Data: map[string]interface{}{
					"id": float64(1),
				},
			},
		},
		{
			name: "bad request",
			req: web.ArticleEditReq{
				Content: "test",
			},
			beforeRequest: func(t *testing.T) {},
			afterRequest: func(t *testing.T) {
				var art dao.Article
				err := db.Where("id = ? AND author_id = ?", 1, userID).First(&art).Error
				assert.Equal(t, err, gorm.ErrRecordNotFound)
			},
			wantCode: http.StatusOK,
			wantBody: web.Response{
				Code:    errno.ErrBadRequest.Code,
				Message: "title为必填字段",
				Data:    nil,
			},
		},
		{
			name: "already exist, update success",
			req: web.ArticleEditReq{
				Title:   "new",
				Content: "new",
			},
			beforeRequest: func(t *testing.T) {
				now := time.Now().UnixMilli()
				err := db.Create(&dao.Article{
					Title:     "test",
					Content:   "test",
					AuthorID:  userID,
					Status:    domain.ArticleStatusUnPublished,
					UpdatedAt: now,
					CreatedAt: now,
				}).Error
				assert.NoError(t, err)
			},
			afterRequest: func(t *testing.T) {
				var art dao.Article
				err := db.Where("id = ? AND author_id = ?", 1, userID).First(&art).Error
				assert.NoError(t, err)
				assert.Equal(t, "new", art.Title)
				assert.Equal(t, "new", art.Content)
			},
			wantCode: http.StatusOK,
			wantBody: web.Response{
				Code:    0,
				Message: "success",
				Data: map[string]interface{}{
					"id": float64(1),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer db.Exec("TRUNCATE TABLE articles")
			data, err := json.Marshal(tt.req)
			assert.NoError(t, err)
			request := httptest.NewRequest(http.MethodPost, "/article/edit", bytes.NewReader(data))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, request)
			assert.Equal(t, tt.wantCode, recorder.Code)
			var res web.Response
			err = json.Unmarshal(recorder.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBody, res)
		})
	}
}

func TestArticle_Publish(t *testing.T) {
	db := startup.InitDB()
	server := gin.Default()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("user_id", int64(userID))
	})
	handler := startup.InitArticleHandler()
	handler.RegisterRoutes(server)
	tests := []struct {
		name          string
		req           web.ArticleEditReq
		beforeRequest func(t *testing.T)
		afterRequest  func(t *testing.T)
		wantCode      int
		wantBody      web.Response
	}{
		{
			name: "all new success",
			req: web.ArticleEditReq{
				Title:   "test",
				Content: "test",
			},
			beforeRequest: func(t *testing.T) {},
			afterRequest: func(t *testing.T) {
				var art dao.Article
				err := db.Where("id = ? AND author_id = ?", 1, userID).First(&art).Error
				assert.NoError(t, err)
				assert.Equal(t, "test", art.Title)
				assert.Equal(t, "test", art.Content)
				assert.Equal(t, domain.ArticleStatusPublished, art.Status)
				assert.Equal(t, int64(userID), art.AuthorID)
				var pubArt dao.PublishedArticle
				err = db.Where("id = ? AND author_id = ?", 1, userID).First(&pubArt).Error
				assert.NoError(t, err)
				assert.Equal(t, "test", pubArt.Title)
				assert.Equal(t, "test", pubArt.Content)
				assert.Equal(t, domain.ArticleStatusPublished, pubArt.Status)
				assert.Equal(t, int64(userID), pubArt.AuthorID)
			},
			wantCode: http.StatusOK,
			wantBody: web.Response{
				Code:    0,
				Message: "success",
				Data: map[string]interface{}{
					"id": float64(1),
				},
			},
		},
		{
			name: "bad request",
			req: web.ArticleEditReq{
				Content: "test",
			},
			beforeRequest: func(t *testing.T) {},
			afterRequest: func(t *testing.T) {
				var art dao.Article
				err := db.Where("id = ? AND author_id = ?", 1, userID).First(&art).Error
				assert.Equal(t, err, gorm.ErrRecordNotFound)
			},
			wantCode: http.StatusOK,
			wantBody: web.Response{
				Code:    errno.ErrBadRequest.Code,
				Message: "title为必填字段",
				Data:    nil,
			},
		},
		{
			name: "already exist, publish success",
			req: web.ArticleEditReq{
				Title:   "new",
				Content: "new",
			},
			beforeRequest: func(t *testing.T) {
				now := time.Now().UnixMilli()
				err := db.Create(&dao.Article{
					Title:     "test",
					Content:   "test",
					AuthorID:  userID,
					Status:    domain.ArticleStatusUnPublished,
					UpdatedAt: now,
					CreatedAt: now,
				}).Error
				assert.NoError(t, err)
			},
			afterRequest: func(t *testing.T) {
				var art dao.Article
				err := db.Where("id = ? AND author_id = ?", 1, userID).First(&art).Error
				assert.NoError(t, err)
				assert.Equal(t, "new", art.Title)
				assert.Equal(t, "new", art.Content)
				assert.Equal(t, domain.ArticleStatusPublished, art.Status)
				assert.Equal(t, int64(userID), art.AuthorID)
				var pubArt dao.PublishedArticle
				err = db.Where("id = ? AND author_id = ?", 1, userID).First(&pubArt).Error
				assert.NoError(t, err)
				assert.Equal(t, "new", pubArt.Title)
				assert.Equal(t, "new", pubArt.Content)
				assert.Equal(t, domain.ArticleStatusPublished, pubArt.Status)
				assert.Equal(t, int64(userID), pubArt.AuthorID)
			},
			wantCode: http.StatusOK,
			wantBody: web.Response{
				Code:    0,
				Message: "success",
				Data: map[string]interface{}{
					"id": float64(1),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			defer db.Exec("TRUNCATE TABLE articles")
			data, err := json.Marshal(tt.req)
			assert.NoError(t, err)
			request := httptest.NewRequest(http.MethodPost, "/article/publish", bytes.NewReader(data))
			request.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			server.ServeHTTP(recorder, request)
			assert.Equal(t, tt.wantCode, recorder.Code)
			var res web.Response
			err = json.Unmarshal(recorder.Body.Bytes(), &res)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBody, res)
		})
	}
}
