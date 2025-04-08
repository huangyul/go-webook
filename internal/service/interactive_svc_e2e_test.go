package service

import (
	"context"
	"fmt"
	"github.com/huangyul/go-webook/internal/repository"
	"github.com/huangyul/go-webook/internal/repository/cache"
	"github.com/huangyul/go-webook/internal/repository/dao"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
	"time"
)

type TestSuite struct {
	suite.Suite
	db     *gorm.DB
	rdb    redis.Cmdable
	userId int64
	bizId  int64
	biz    string
}

func (s *TestSuite) SetupSuite() {
	db, err := gorm.Open(mysql.Open("root:root@tcp(localhost:13316)/webook_test?charset=utf8mb4&parseTime=True&loc=Local"))
	if err != nil {
		panic(err)
	}
	s.db = db

	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:16379",
		Password: "",
		DB:       0,
	})
	if er := rdb.Ping(context.Background()).Err(); er != nil {
		panic(err)
	}
	s.rdb = rdb

	dao.InitTable(s.db)

	s.userId = int64(1)
	s.bizId = int64(1)
	s.biz = "article"
}

func (s *TestSuite) TearDownTest() {
	s.db.Exec("TRUNCATE TABLE `interactives`")
	s.db.Exec("TRUNCATE TABLE `user_like_bizs`")
	s.rdb.Del(context.Background(), fmt.Sprintf("interactive:%s:%d", s.biz, s.bizId))
}

type InteractiveSuite struct {
	TestSuite
	svc InteractiveService
}

func (s *InteractiveSuite) SetupSuite() {
	s.TestSuite.SetupSuite()

	interDAO := dao.NewInteractiveDAO(s.db)
	interCache := cache.NewInteractiveCache(s.rdb)
	interRepo := repository.NewInteractiveRepository(interDAO, interCache)
	s.svc = NewInteractiveService(interRepo)
}

func (s *InteractiveSuite) TestInteractiveLike() {
	tests := []struct {
		name    string
		before  func(t *testing.T)
		after   func(t *testing.T)
		wantErr error
	}{
		{
			name:   "success1",
			before: func(t *testing.T) {},
			after: func(t *testing.T) {
				var inter dao.Interactive
				err := s.db.Where("biz_id = ? and biz = ?", s.bizId, s.biz).First(&inter).Error
				assert.NoError(t, err)
				assert.Equal(t, int64(1), inter.LikeCnt)
				var userLikeBiz dao.UserLikeBiz
				err = s.db.Where("biz_id = ? and biz = ? and user_id = ?", s.bizId, s.biz, s.userId).First(&userLikeBiz).Error
				assert.NoError(t, err)
				assert.Equal(t, 1, userLikeBiz.Status)
			},
			wantErr: nil,
		},
		{
			name: "success2",
			before: func(t *testing.T) {
				err := s.db.Model(&dao.Interactive{}).Create(&dao.Interactive{
					BizId:     s.bizId,
					Biz:       s.biz,
					LikeCnt:   10,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}).Error
				assert.NoError(t, err)
				err = s.db.Model(&dao.UserLikeBiz{}).Create(&dao.UserLikeBiz{
					BizId:     s.bizId,
					Biz:       s.biz,
					UserId:    s.userId,
					Status:    0,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var inter dao.Interactive
				err := s.db.Where("biz_id = ? and biz = ?", s.bizId, s.biz).First(&inter).Error
				assert.NoError(t, err)
				assert.Equal(t, int64(11), inter.LikeCnt)
				var userLikeBiz dao.UserLikeBiz
				err = s.db.Where("biz_id = ? and biz = ? and user_id = ?", s.bizId, s.biz, s.userId).First(&userLikeBiz).Error
				assert.NoError(t, err)
				assert.Equal(t, 1, userLikeBiz.Status)
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.before(t)
			err := s.svc.Like(context.Background(), s.biz, s.bizId, s.userId)
			assert.NoError(t, err)
			tt.after(t)
		})
	}
}

func (s *InteractiveSuite) TestInteractiveCancelLike() {
	tests := []struct {
		name    string
		before  func(t *testing.T)
		after   func(t *testing.T)
		wantErr error
	}{
		{
			name: "success",
			before: func(t *testing.T) {
				err := s.db.Model(&dao.Interactive{}).Create(&dao.Interactive{
					BizId:     s.bizId,
					Biz:       s.biz,
					LikeCnt:   10,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}).Error
				assert.NoError(t, err)
				err = s.db.Model(&dao.UserLikeBiz{}).Create(&dao.UserLikeBiz{
					BizId:     s.bizId,
					Biz:       s.biz,
					UserId:    s.userId,
					Status:    1,
					CreatedAt: time.Now(),
					UpdatedAt: time.Now(),
				}).Error
				assert.NoError(t, err)
			},
			after: func(t *testing.T) {
				var inter dao.Interactive
				err := s.db.Where("biz_id = ? and biz = ?", s.bizId, s.biz).First(&inter).Error
				assert.NoError(t, err)
				assert.Equal(t, int64(9), inter.LikeCnt)
				var userLikeBiz dao.UserLikeBiz
				err = s.db.Where("biz_id = ? and biz = ? and user_id = ?", s.bizId, s.biz, s.userId).First(&userLikeBiz).Error
				assert.NoError(t, err)
				assert.Equal(t, 0, userLikeBiz.Status)
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		s.T().Run(tt.name, func(t *testing.T) {
			tt.before(t)
			err := s.svc.CancelLike(context.Background(), s.biz, s.bizId, s.userId)
			assert.NoError(t, err)
			tt.after(t)
		})
	}
}

func TestInteractive(t *testing.T) {
	suite.Run(t, &InteractiveSuite{})
}
