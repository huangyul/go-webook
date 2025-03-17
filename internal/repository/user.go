package repository

import (
	"context"
	"github.com/huangyul/go-webook/internal/domain"
	"github.com/huangyul/go-webook/internal/repository/cache"
	"github.com/huangyul/go-webook/internal/repository/dao"
)

var (
	ErrUserNotFound           = dao.ErrUserNotFound
	ErrUserEmailAlreadyExists = dao.ErrUserEmailAlreadyExists
)

type UserRepository interface {
	Insert(ctx context.Context, user *domain.User) error
	FindByEmail(ctx context.Context, email string) (*domain.User, error)
	FindByID(ctx context.Context, id int64) (*domain.User, error)
	Update(ctx context.Context, user *domain.User) error
}

var _ UserRepository = (*userRepository)(nil)

type userRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func (repo *userRepository) FindByID(ctx context.Context, id int64) (*domain.User, error) {

	du, err := repo.cache.Get(ctx, id)
	if err == nil {
		return du, nil
	}

	u, err := repo.dao.FindById(ctx, id)
	if err != nil {
		return nil, err
	}
	go func() {
		_ = repo.cache.Set(ctx, repo.toDomain(u))
	}()
	return repo.toDomain(u), nil
}

func (repo *userRepository) Update(ctx context.Context, user *domain.User) error {
	return repo.dao.Update(ctx, repo.toEntity(user))
}

func NewUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &userRepository{
		dao:   dao,
		cache: cache,
	}
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
		Nickname:  user.Nickname,
		Birthday:  user.Birthday,
		AboutMe:   user.AboutMe,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}

func (repo *userRepository) toEntity(user *domain.User) *dao.User {
	return &dao.User{
		ID:       user.ID,
		Email:    user.Email,
		Nickname: user.Nickname,
		Birthday: user.Birthday,
		AboutMe:  user.AboutMe,
		Password: user.Password,
	}
}
