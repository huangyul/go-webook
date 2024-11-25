package service

import (
	"context"
	"errors"
	"github.com/huangyul/go-blog/internal/domain"
	"github.com/huangyul/go-blog/internal/pkg/errno"
	"github.com/huangyul/go-blog/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Signup(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (domain.User, error)
	UpdateUserInfo(ctx context.Context, user domain.User) error
	GetUserInfo(ctx context.Context, id int64) (domain.User, error)
	GetUserList(ctx context.Context, page, pageSize int) ([]domain.User, int, error)
}

var _ UserService = (*userService)(nil)

type userService struct {
	repo repository.UserRepository
}

func (svc *userService) GetUserList(ctx context.Context, page, pageSize int) ([]domain.User, int, error) {
	return svc.repo.GetUserList(ctx, page, pageSize)
}

func (svc *userService) UpdateUserInfo(ctx context.Context, user domain.User) error {
	return svc.repo.UpdateByID(ctx, user)
}

func (svc *userService) GetUserInfo(ctx context.Context, id int64) (domain.User, error) {
	return svc.repo.FindById(ctx, id)
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// Login
func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if errors.Is(err, errno.ErrNotFoundUser) {
		return domain.User{}, errno.ErrNotFoundUser
	}
	if err != nil {
		return domain.User{}, errno.ErrInternalServer
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, errno.ErrNotFoundUser
	}

	return u, nil
}

// Signup
func (svc *userService) Signup(ctx context.Context, email string, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if errors.Is(err, errno.ErrEmailAlreadyExist) {
		return errno.ErrEmailAlreadyExist
	}
	if err != nil {
		return errno.ErrInternalServer
	}
	return svc.repo.Insert(ctx, domain.User{Email: email, Password: string(hash)})
}
