package service

import (
	"context"
	"github.com/huangyul/go-blog/internal/domain"
	"github.com/huangyul/go-blog/internal/repository"
	repomocks "github.com/huangyul/go-blog/internal/repository/mocks"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
	"golang.org/x/crypto/bcrypt"
	"testing"
)

func TestUserService_Login(t *testing.T) {
	tests := []struct {
		name     string
		email    string
		password string
		before   func(ctrl *gomock.Controller) repository.UserRepository
	}{
		{
			name:     "success",
			email:    "123123@qq.com",
			password: "123123",
			before: func(ctrl *gomock.Controller) repository.UserRepository {
				encryptPass, _ := bcrypt.GenerateFromPassword([]byte("123123"), bcrypt.DefaultCost)
				repo := repomocks.NewMockUserRepository(ctrl)
				repo.EXPECT().FindByEmail(gomock.Any(), "123123@qq.com").Return(domain.User{
					Password: string(encryptPass),
					ID:       int64(1),
				}, nil)
				return repo
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			defer ctrl.Finish()
			repo := tt.before(ctrl)
			svc := NewUserService(repo)
			user, err := svc.Login(context.Background(), tt.email, tt.password)
			assert.NoError(t, err)
			assert.Equal(t, user.ID, int64(1))
		})
	}
}
