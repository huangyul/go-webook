package repository

import (
	"context"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/repository/dao"
)

var (
	ErrUserNotFound           = dao.ErrUserNotFound
	ErrUserEmailAlreadyExists = dao.ErrUserEmailAlreadyExists
)

type UserRepository interface {
	Insert(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
}

var _ UserRepository = (*userRepository)(nil)

type userRepository struct {
	dao dao.UserDAO
}

func NewUserRepository(dao dao.UserDAO) UserRepository {
	return &userRepository{dao: dao}
}

func (repo *userRepository) Insert(ctx context.Context, user *domain.User) error {
	return repo.dao.Insert(ctx, repo.toEntity(user))
}

func (repo *userRepository) FindByEmail(ctx context.Context, email string) (*domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return nil, err
	}
	return repo.toDomain(u), nil
}

func (repo *userRepository) toDomain(user *dao.User) *domain.User {
	return &domain.User{
		ID:        user.ID,
		Email:     user.Email,
		Password:  user.Password,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (repo *userRepository) toEntity(user *domain.User) *dao.User {
	return &dao.User{
		ID:       user.ID,
		Email:    user.Email,
		Password: user.Password,
	}
}
