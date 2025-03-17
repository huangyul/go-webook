package service

import (
	"context"
	"errors"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

var (
	ErrUserNotFound          = repository.ErrUserNotFound
	ErrUserAlreadyExists     = repository.ErrUserEmailAlreadyExists
	ErrUserEmailIllegally    = errors.New("user email illegal")
	ErrUserPasswordIllegally = errors.New("user password illegal")
	ErrUserPasswordNotMatch  = errors.New("user password not match")
)

type UserService interface {
	RegisterByEmail(ctx context.Context, email string, password string) error
	LoginByEmail(ctx context.Context, email string, password string) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
	FindById(ctx context.Context, userId int64) (*domain.User, error)
	FindOrCreateByPhone(ctx context.Context, phone string) (*domain.User, error)
}

var _ UserService = (*userService)(nil)

type userService struct {
	repo repository.UserRepository
}

func (u *userService) FindOrCreateByPhone(ctx context.Context, phone string) (*domain.User, error) {
	user, err := u.repo.FindByPhone(ctx, phone)
	if errors.Is(err, ErrUserNotFound) {
		// TODO 调整insert，要返回user
		u.repo.Insert(ctx, &domain.User{
			Phone: phone,
		})
	}

}

func (u *userService) Update(ctx context.Context, user *domain.User) error {
	return u.repo.Update(ctx, user)
}

func (u *userService) FindById(ctx context.Context, userId int64) (*domain.User, error) {
	return u.repo.FindByID(ctx, userId)
}

func (u *userService) RegisterByEmail(ctx context.Context, email string, password string) error {
	passwordStr, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return u.repo.Insert(ctx, &domain.User{
		Email:    email,
		Password: string(passwordStr),
	})
}

func (u *userService) LoginByEmail(ctx context.Context, email string, password string) (*domain.User, error) {
	user, err := u.repo.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return nil, ErrUserPasswordNotMatch
	}
	return user, nil
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}
