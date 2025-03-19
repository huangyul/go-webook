package cache

import (
	"context"
	"errors"
	cachemocks "github.com/huangyul/go-webook/internal/repository/cache/mocks"
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"testing"
)

//go:generate mockgen  -destination=mocks/redis.go -package=cachemocks github.com/redis/go-redis/v9 Cmdable

func TestCodeCache_key(t *testing.T) {
	assert.Equal(t, "phone_code:login:13000000000", (&codeCache{}).key("login", "13000000000"))
}

func TestNewCodeCache(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()
	client := NewCodeCache(cachemocks.NewMockCmdable(ctrl))
	assert.NotNil(t, client)
}

func TestCodeCache_Set(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		biz     string
		phone   string
		code    string
		wantErr error
	}{
		{
			name: "success",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := cachemocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(0))
				res.EXPECT().Eval(gomock.All(), luaSetCode, []string{"phone_code:login:13000000000"}, "123456").Return(cmd)
				return res
			},
			biz:     "login",
			phone:   "13000000000",
			code:    "123456",
			wantErr: nil,
		},
		{
			name: "send too many",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := cachemocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(-1))
				res.EXPECT().Eval(gomock.All(), luaSetCode, []string{"phone_code:login:13000000000"}, "123456").Return(cmd)
				return res
			},
			biz:     "login",
			phone:   "13000000000",
			code:    "123456",
			wantErr: ErrCodeSendTooMany,
		},
		{
			name: "not expiration time",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := cachemocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(-2))
				res.EXPECT().Eval(gomock.All(), luaSetCode, []string{"phone_code:login:13000000000"}, "123456").Return(cmd)
				return res
			},
			biz:     "login",
			phone:   "13000000000",
			code:    "123456",
			wantErr: errors.New("not expiration time"),
		},
		{
			name: "redis error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := cachemocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(errors.New("redis error"))
				cmd.SetVal(int64(-2))
				res.EXPECT().Eval(gomock.All(), luaSetCode, []string{"phone_code:login:13000000000"}, "123456").Return(cmd)
				return res
			},
			biz:     "login",
			phone:   "13000000000",
			code:    "123456",
			wantErr: errors.New("redis error"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()

			cache := NewCodeCache(tt.mock(ctrl))
			err := cache.Set(context.Background(), tt.biz, tt.phone, tt.code)

			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestCodeCache_Verify(t *testing.T) {
	tests := []struct {
		name    string
		mock    func(ctrl *gomock.Controller) redis.Cmdable
		biz     string
		phone   string
		code    string
		wantOk  bool
		wantErr error
	}{
		{
			name: "success",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := cachemocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(0))
				res.EXPECT().Eval(gomock.All(), luaVerifyCode, []string{"phone_code:login:13000000000"}, "123456").Return(cmd)
				return res
			},
			biz:     "login",
			phone:   "13000000000",
			code:    "123456",
			wantOk:  true,
			wantErr: nil,
		},
		{
			name: "redis error",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := cachemocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(errors.New("redis error"))
				cmd.SetVal(int64(0))
				res.EXPECT().Eval(gomock.All(), luaVerifyCode, []string{"phone_code:login:13000000000"}, "123456").Return(cmd)
				return res
			},
			biz:     "login",
			phone:   "13000000000",
			code:    "123456",
			wantOk:  false,
			wantErr: errors.New("redis error"),
		},
		{
			name: "incorrect verification code",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := cachemocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(-2))
				res.EXPECT().Eval(gomock.All(), luaVerifyCode, []string{"phone_code:login:13000000000"}, "123456").Return(cmd)
				return res
			},
			biz:     "login",
			phone:   "13000000000",
			code:    "123456",
			wantOk:  false,
			wantErr: nil,
		},
		{
			name: "code verify too many",
			mock: func(ctrl *gomock.Controller) redis.Cmdable {
				res := cachemocks.NewMockCmdable(ctrl)
				cmd := redis.NewCmd(context.Background())
				cmd.SetErr(nil)
				cmd.SetVal(int64(-1))
				res.EXPECT().Eval(gomock.All(), luaVerifyCode, []string{"phone_code:login:13000000000"}, "123456").Return(cmd)
				return res
			},
			biz:     "login",
			phone:   "13000000000",
			code:    "123456",
			wantOk:  false,
			wantErr: ErrCodeVerifyTooMany,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			cache := NewCodeCache(tt.mock(ctrl))
			ok, err := cache.Verify(context.Background(), tt.biz, tt.phone, tt.code)
			assert.Equal(t, tt.wantOk, ok)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}
