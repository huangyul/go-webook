package repository

import (
	"context"
	"time"

	"github.com/huangyul/go-blog/internal/domain"
	"github.com/huangyul/go-blog/internal/repository/dao"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	Insert(ctx context.Context, user domain.User) error
}

var _ UserRepository = (*userRepository)(nil)

type userRepository struct {
	dao dao.UserDAO
}

func NewUserRepository(dao dao.UserDAO) UserRepository {
	return &userRepository{
		dao: dao,
	}
}

// FindByEmail
func (repo *userRepository) FindByEmail(ctx context.Context, email string) (domain.User, error) {
	u, err := repo.dao.FindByEmail(ctx, email)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(u), nil
}

// Insert
func (repo *userRepository) Insert(ctx context.Context, user domain.User) error {
	return repo.dao.Insert(ctx, repo.toDao(user))
}

func (repo *userRepository) toDomain(u dao.User) domain.User {
	return domain.User{
		ID:        u.ID,
		Eamil:     u.Email,
		Password:  u.Password,
		CreatedAt: time.UnixMilli(u.CreatedAt),
		UpdatedAt: time.UnixMilli(u.UpdatedAt),
	}
}

func (repo *userRepository) toDao(u domain.User) dao.User {
	return dao.User{
		Email:     u.Eamil,
		Password:  u.Password,
		CreatedAt: u.CreatedAt.UnixMilli(),
		UpdatedAt: u.UpdatedAt.UnixMilli(),
	}
}
