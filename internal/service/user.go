package service

import (
	"context"

	"github.com/huangyul/go-blog/internal/domain"
	"github.com/huangyul/go-blog/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	Signup(ctx context.Context, email, password string) error
	Login(ctx context.Context, email, password string) (domain.User, error)
}

var _ UserService = (*userService)(nil)

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{repo: repo}
}

// Login
func (svc *userService) Login(ctx context.Context, email string, password string) (domain.User, error) {
	u, err := svc.repo.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	if err != nil {
		return domain.User{}, err
	}

	return u, nil
}

// Signup
func (svc *userService) Signup(ctx context.Context, email string, password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	return svc.repo.Insert(ctx, domain.User{Eamil: email, Password: string(hash)})
}
