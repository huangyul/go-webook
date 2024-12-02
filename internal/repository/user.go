package repository

import (
	"context"
	"github.com/huangyul/go-blog/internal/repository/cache"
	"github.com/huangyul/go-blog/pkg/utils"
	"time"

	"github.com/huangyul/go-blog/internal/domain"
	"github.com/huangyul/go-blog/internal/repository/dao"
)

type UserRepository interface {
	FindByEmail(ctx context.Context, email string) (domain.User, error)
	Insert(ctx context.Context, user domain.User) error
	FindById(ctx context.Context, id int64) (domain.User, error)
	UpdateByID(ctx context.Context, user domain.User) error
	GetUserList(ctx context.Context, page, pageSize int) ([]domain.User, int, error)
	FindOrCreateByPhone(ctx context.Context, phone string) (domain.User, error)
}

var _ UserRepository = (*userRepository)(nil)

type userRepository struct {
	dao   dao.UserDAO
	cache cache.UserCache
}

func NewUserRepository(dao dao.UserDAO, cache cache.UserCache) UserRepository {
	return &userRepository{
		dao:   dao,
		cache: cache,
	}
}

func (repo *userRepository) FindOrCreateByPhone(ctx context.Context, phone string) (domain.User, error) {
	dUser, err := repo.dao.FindOrCreateByPhone(ctx, phone)
	if err != nil {
		return domain.User{}, err
	}
	return repo.toDomain(dUser), nil
}

func (repo *userRepository) GetUserList(ctx context.Context, page, pageSize int) ([]domain.User, int, error) {
	users, count, err := repo.dao.GetList(ctx, page, pageSize)
	if err != nil {
		return nil, 0, err
	}
	var userList []domain.User
	for _, user := range users {
		userList = append(userList, repo.toDomain(user))
	}
	return userList, count, nil
}

func (repo *userRepository) UpdateByID(ctx context.Context, user domain.User) error {
	return repo.dao.UpdateByID(ctx, repo.toDao(user))
}

func (repo *userRepository) FindById(ctx context.Context, id int64) (domain.User, error) {
	dUser, err := repo.cache.Get(ctx, id)
	if err == nil {
		return dUser, nil
	}
	u, err := repo.dao.FindByID(ctx, id)
	if err != nil {
		return domain.User{}, err
	}
	dUser = repo.toDomain(u)
	go repo.cache.Set(ctx, dUser)
	return dUser, nil
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
		Email:     utils.PtrToString(u.Email),
		Password:  u.Password,
		AboutMe:   u.AboutMe,
		Phone:     u.Phone,
		Birthday:  time.UnixMilli(u.Birthday),
		Nickname:  u.Nickname,
		CreatedAt: time.UnixMilli(u.CreatedAt),
		UpdatedAt: time.UnixMilli(u.UpdatedAt),
	}
}

func (repo *userRepository) toDao(u domain.User) dao.User {
	return dao.User{
		ID:        u.ID,
		Email:     &u.Email,
		Password:  u.Password,
		AboutMe:   u.AboutMe,
		Phone:     u.Phone,
		Birthday:  u.Birthday.UnixMilli(),
		Nickname:  u.Nickname,
		CreatedAt: u.CreatedAt.UnixMilli(),
		UpdatedAt: u.UpdatedAt.UnixMilli(),
	}
}
