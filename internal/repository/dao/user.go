package dao

import (
	"context"
	"errors"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/huangyul/go-blog/internal/pkg/errno"
	"gorm.io/gorm"
)

type UserDAO interface {
	Insert(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
	FindByID(ctx context.Context, id int64) (User, error)
	UpdateByID(ctx context.Context, user User) error
	GetList(ctx context.Context, page, pageSize int) ([]User, int, error)
}

var _ UserDAO = (*UserDAOGORM)(nil)

type UserDAOGORM struct {
	db *gorm.DB
}

func (dao *UserDAOGORM) GetList(ctx context.Context, page, pageSize int) ([]User, int, error) {
	var users []User
	var count int64
	err := dao.db.Model(&User{}).WithContext(ctx).Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	if err != nil {
		return nil, 0, err
	}
	err = dao.db.Model(&User{}).WithContext(ctx).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	return users, int(count), nil
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &UserDAOGORM{db: db}
}

// UpdateByID update by id,only can update birthday, aboutme, nickname
func (dao *UserDAOGORM) UpdateByID(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	up := dao.db.Model(&User{}).WithContext(ctx).Where("id = ?", user.ID).Updates(map[string]any{
		"about_me":   user.AboutMe,
		"birthday":   user.Birthday,
		"nickname":   user.Nickname,
		"updated_at": now,
	})
	if up.Error != nil {
		return up.Error
	}
	if up.RowsAffected == 0 {
		return errno.ErrNotFoundUser
	}
	return nil
}

// FindByID find user by id
func (dao *UserDAOGORM) FindByID(ctx context.Context, id int64) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).First(&user, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, errno.ErrNotFoundUser
	}
	if err != nil {
		return user, err
	}
	return user, nil
}

// FindByEmail Finds user by email
func (dao *UserDAOGORM) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return User{}, errno.ErrNotFoundUser
	}
	return user, err
}

// Insert implements UserDAO.
func (dao *UserDAOGORM) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.CreatedAt = now
	user.UpdatedAt = now
	err := dao.db.WithContext(ctx).Create(&user).Error
	if e, ok := err.(*mysql.MySQLError); ok {
		if e.Error() == "Error 1062: Duplicate entry" {
			return errno.ErrEmailAlreadyExist
		}
	}
	return err
}

// User user table construct
type User struct {
	ID       int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Password string `gorm:"size:255;not null"`
	Email    string `gorm:"unique;size:255;not null"`

	Nickname string `gorm:"type:varchar(255)"`
	Birthday int64
	AboutMe  string `gorm:"type:varchar(4096)"`

	CreatedAt int64
	UpdatedAt int64
}
