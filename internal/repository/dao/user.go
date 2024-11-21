package dao

import (
	"context"
	"time"

	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
)

var (
	ErrDuplicate = gorm.ErrDuplicatedKey
)

type UserDAO interface {
	Insert(ctx context.Context, user User) error
	FindByEmail(ctx context.Context, email string) (User, error)
}

var _ UserDAO = (*UserDAOGORM)(nil)

type UserDAOGORM struct {
	db *gorm.DB
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &UserDAOGORM{db: db}
}

// FindByEmail Finds user by email
func (dao *UserDAOGORM) FindByEmail(ctx context.Context, email string) (User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	return user, err
}

// Insert implements UserDAO.
func (dao *UserDAOGORM) Insert(ctx context.Context, user User) error {
	now := time.Now().UnixMilli()
	user.CreatedAt = now
	user.UpdatedAt = now
	err := dao.db.WithContext(ctx).Create(&user).Error
	if err != nil {
		if e, ok := err.(*mysql.MySQLError); ok {
			if e.Error() == "Error 1062: Duplicate entry" {
				return ErrDuplicate
			}
		}
	}
	return err
}

// User user table construct
type User struct {
	ID       int64  `gorm:"primary_key;AUTO_INCREMENT"`
	Password string `gorm:"size:255;not null"`
	Email    string `gorm:"unique;size:255;not null"`

	CreatedAt int64
	UpdatedAt int64
}
