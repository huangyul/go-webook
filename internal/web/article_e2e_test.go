//go:build e2e

package web

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/repository/cache"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/huangyul/go-webook/internal/repository"
	"github.com/huangyul/go-webook/internal/repository/dao"
	"github.com/huangyul/go-webook/internal/service"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// E2ETestSuite 是端到端测试的基础套件，可以被其他测试复用
type E2ETestSuite struct {
	suite.Suite
	db     *gorm.DB
	rdb    redis.Cmdable
	server *gin.Engine
	userId int64 // 测试用户ID
}

func (s *E2ETestSuite) SetupSuite() {
	// 初始化数据库连接
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook_test?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		s.T().Fatal(err)
	}
	s.db = db

	// 初始化Redis连接
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:16379",
		Password: "",
		DB:       0,
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		s.T().Fatal(err)
	}
	s.rdb = rdb

	// 初始化数据库表
	dao.InitTable(s.db)

	// 默认用户ID
	s.userId = int64(1)

	// 初始化Gin服务器
	server := gin.New()
	server.Use(func(ctx *gin.Context) {
		ctx.Set("userId", s.userId)
	})
	s.server = server
}

func (s *E2ETestSuite) TearDownSuite() {
	// 清理资源
}

func (s *E2ETestSuite) TearDownTest() {
	s.db.Exec("TRUNCATE TABLE `articles`")
	s.db.Exec("TRUNCATE TABLE `pub_articles`")
}

// SetUserId 设置测试用户ID
func (s *E2ETestSuite) SetUserId(id int64) {
	s.userId = id
}

// SendRequest 发送HTTP请求并返回响应
func (s *E2ETestSuite) SendRequest(method, path string, body interface{}) *httptest.ResponseRecorder {
	data, err := json.Marshal(body)
	assert.NoError(s.T(), err)

	req := httptest.NewRequest(method, path, bytes.NewBuffer(data))
	req.Header.Set("Content-Type", "application/json")

	resp := httptest.NewRecorder()
	s.server.ServeHTTP(resp, req)
	return resp
}

// CleanTable 清理表数据
func (s *E2ETestSuite) CleanTable(tableName string) {
	err := s.db.Exec(fmt.Sprintf("TRUNCATE TABLE `%s`", tableName)).Error
	assert.NoError(s.T(), err)
}

// ArticleTestSuite 文章测试套件，继承基础测试套件
type ArticleTestSuite struct {
	E2ETestSuite
	articleHandler *ArticleHandler
}

func (s *ArticleTestSuite) SetupSuite() {
	s.E2ETestSuite.SetupSuite()

	// 初始化文章相关依赖
	articleDao := dao.NewArticleDAO(s.db)
	articleCache := cache.NewArticleCache(s.rdb)
	articleRepo := repository.NewArticleRepository(articleDao, articleCache)
	userDao := dao.NewUserDAO(s.db)
	userCache := cache.NewRedisUserCache(s.rdb)
	userRepo := repository.NewUserRepository(userDao, userCache)
	articleService := service.NewArticleService(articleRepo, userRepo)
	s.articleHandler = NewArticleHandler(articleService, nil)

	// 注册路由
	s.articleHandler.RegisterRoutes(s.server)
}

// TestArticleSave
func (s *ArticleTestSuite) TestArticleSave() {
	type Req struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}

	tests := []struct {
		name           string
		before         func(t *testing.T)
		after          func(t *testing.T)
		req            Req
		wantStatusCode int
		wantBody       ApiResponse[any]
	}{
		{
			name: "success",
			before: func(t *testing.T) {
				s.CleanTable("articles")
			},
			req: Req{
				Title:   "title",
				Content: "content",
			},
			after: func(t *testing.T) {
				defer s.CleanTable("articles")
				var art dao.Article
				err := s.db.Where("id = ? AND author_id = ?", 1, s.userId).First(&art).Error
				assert.NoError(t, err)
				assert.Equal(t, int64(1), art.Id)
				assert.Equal(t, "title", art.Title)
				assert.Equal(t, "content", art.Content)
			},
			wantStatusCode: http.StatusOK,
			wantBody: ApiResponse[any]{
				Code: 0,
				Msg:  "success",
				Data: nil,
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}

			resp := s.SendRequest(http.MethodPost, "/article/save", tt.req)

			assert.Equal(t, tt.wantStatusCode, resp.Code)
			var rep ApiResponse[any]
			err := json.Unmarshal(resp.Body.Bytes(), &rep)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBody.Code, rep.Code)
			assert.Equal(t, tt.wantBody.Msg, rep.Msg)
			assert.Equal(t, tt.wantBody.Data, rep.Data)

			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}

// TestArticlePublish
func (s *ArticleTestSuite) TestArticlePublish() {
	type Req struct {
		Id      int64  `json:"id"`
		Title   string `json:"title"`
		Content string `json:"content"`
	}
	tests := []struct {
		name           string
		before         func(t *testing.T)
		after          func(t *testing.T)
		req            Req
		wantStatusCode int
		wantBody       ApiResponse[any]
	}{
		{
			name: "success: none of them exists",
			before: func(t *testing.T) {
				s.CleanTable("articles")
				s.CleanTable("pub_articles")
			},
			after: func(t *testing.T) {
				var art dao.Article
				err := s.db.Where("id = ? AND author_id = ?", 1, s.userId).First(&art).Error
				assert.NoError(t, err)
				assert.Equal(t, int64(1), art.Id)
				assert.Equal(t, "title", art.Title)
				assert.Equal(t, "content", art.Content)
				var pArt dao.PubArticle
				err = s.db.Where("id = ? AND author_id = ?", 1, s.userId).First(&pArt).Error
				assert.NoError(t, err)
				assert.Equal(t, int64(1), pArt.Id)
				assert.Equal(t, "title", pArt.Title)
				assert.Equal(t, "content", pArt.Content)
			},
			req: Req{
				Title:   "title",
				Content: "content",
			},
			wantStatusCode: http.StatusOK,
			wantBody: ApiResponse[any]{
				Code: 0,
				Msg:  "success",
				Data: nil,
			},
		},
		{
			name: "success: article exists, pub_articles not exists",
			before: func(t *testing.T) {
				s.CleanTable("articles")
				s.CleanTable("pub_articles")
				now := time.Now()
				err := s.db.Create(&dao.Article{
					Title:     "title",
					Content:   "content",
					AuthorId:  s.userId,
					CreatedAt: now,
					UpdatedAt: now,
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art dao.Article
				err := s.db.Where("id = ? AND author_id = ?", 1, s.userId).First(&art).Error
				assert.NoError(t, err)
				assert.Equal(t, int64(1), art.Id)
				assert.Equal(t, "new_title", art.Title)
				assert.Equal(t, "new_content", art.Content)
				var pArt dao.PubArticle
				err = s.db.Where("id = ? AND author_id = ?", 1, s.userId).First(&pArt).Error
				assert.NoError(t, err)
				assert.Equal(t, int64(1), pArt.Id)
				assert.Equal(t, "new_title", pArt.Title)
				assert.Equal(t, "new_content", pArt.Content)
			},
			req: Req{
				Id:      int64(1),
				Title:   "new_title",
				Content: "new_content",
			},
			wantStatusCode: http.StatusOK,
			wantBody: ApiResponse[any]{
				Code: 0,
				Msg:  "success",
				Data: nil,
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}
			resp := s.SendRequest(http.MethodPost, "/article/publish", tt.req)
			assert.Equal(t, tt.wantStatusCode, resp.Code)
			var rep ApiResponse[any]
			err := json.Unmarshal(resp.Body.Bytes(), &rep)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBody.Code, rep.Code)
			assert.Equal(t, tt.wantBody.Msg, rep.Msg)
			assert.Equal(t, tt.wantBody.Data, rep.Data)
			if tt.after != nil {
				tt.after(t)
			}
		})
	}
}

func (s *ArticleTestSuite) TestArticleWithdraw() {
	tests := []struct {
		name           string
		before         func(t *testing.T)
		after          func(t *testing.T)
		id             int64
		wantStatusCode int
		wantBody       ApiResponse[any]
	}{
		{
			name: "success",
			before: func(t *testing.T) {
				s.CleanTable("articles")
				s.CleanTable("pub_articles")
				now := time.Now()
				art := dao.Article{
					Title:     "title",
					Content:   "content",
					AuthorId:  s.userId,
					CreatedAt: now,
					UpdatedAt: now,
					Status:    domain.ArticleStatusPublished.ToUint8(),
				}
				err := s.db.Create(&art).Error
				assert.NoError(t, err)
				pubArticle := dao.PubArticle{
					Title:     "title",
					Content:   "content",
					AuthorId:  s.userId,
					CreatedAt: now,
					UpdatedAt: now,
					Status:    domain.ArticleStatusPublished.ToUint8(),
				}
				err = s.db.Create(&pubArticle).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var art dao.Article
				err := s.db.Where("id = ? AND author_id = ?", 1, s.userId).First(&art).Error
				assert.NoError(t, err)
				assert.Equal(t, domain.ArticleStatusPrivate.ToUint8(), art.Status)
				var pubArticle dao.PubArticle
				err = s.db.Where("id = ? AND author_id = ?", 1, s.userId).First(&pubArticle).Error
				assert.NoError(t, err)
				assert.Equal(t, domain.ArticleStatusPrivate.ToUint8(), pubArticle.Status)
			},
			id:             1,
			wantStatusCode: http.StatusOK,
			wantBody: ApiResponse[any]{
				Code: 0,
				Msg:  "success",
				Data: nil,
			},
		},
		{
			name: "article not exists, pub_articles not exists",
			before: func(t *testing.T) {
				s.CleanTable("articles")
				s.CleanTable("pub_articles")
			},
			after: func(t *testing.T) {
			},
			id:             1,
			wantStatusCode: http.StatusOK,
			wantBody: ApiResponse[any]{
				Code: 1,
				Msg:  "article not found",
				Data: nil,
			},
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			if tt.before != nil {
				tt.before(t)
			}
			if tt.after != nil {
				defer tt.after(t)
			}
			resp := s.SendRequest(http.MethodGet, fmt.Sprintf("/article/withdraw?id=%d", tt.id), nil)
			assert.Equal(t, tt.wantStatusCode, resp.Code)
			var rep ApiResponse[any]
			err := json.Unmarshal(resp.Body.Bytes(), &rep)
			assert.NoError(t, err)
			assert.Equal(t, tt.wantBody.Code, rep.Code)
			assert.Equal(t, tt.wantBody.Msg, rep.Msg)
			assert.Equal(t, tt.wantBody.Data, rep.Data)
		})
	}
}

// TestArticle 启动文章测试套件
func TestArticle(t *testing.T) {
	suite.Run(t, &ArticleTestSuite{})
}
