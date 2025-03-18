package dao

import (
	"context"
	"errors"
	"github.com/go-sql-driver/mysql"
	"gorm.io/gorm"
	"time"
)

var (
	ErrUserEmailAlreadyExists = errors.New("user email already exists")
	ErrUserNotFound           = errors.New("user not found")
)

type UserDAO interface {
	Insert(ctx context.Context, user *User) error
	FindByEmail(ctx context.Context, email string) (*User, error)
	FindById(ctx context.Context, id int64) (*User, error)
	Update(ctx context.Context, user *User) error
	FindByPhone(ctx context.Context, phone string) (*User, error)
}

var _ UserDAO = (*GORMUserDAO)(nil)

type GORMUserDAO struct {
	db *gorm.DB
}

func (dao *GORMUserDAO) FindByPhone(ctx context.Context, phone string) (*User, error) {
	var user User
	err := dao.db.WithContext(ctx).Where("phone = ?", phone).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func (dao *GORMUserDAO) FindById(ctx context.Context, id int64) (*User, error) {
	var user User
	err := dao.db.WithContext(ctx).First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrUserNotFound
	}
	return &user, err
}

func (dao *GORMUserDAO) Update(ctx context.Context, user *User) error {
	updates := map[string]interface{}{
		"nickname":   user.Nickname,
		"about_me":   user.AboutMe,
		"created_at": time.Now(),
	}

	if user.Birthday != nil {
		if user.Birthday.IsZero() {
			updates["birthday"] = nil
		} else {
			updates["birthday"] = user.Birthday
		}
	}

	return dao.db.WithContext(ctx).Model(&User{}).Where("id = ?", user.ID).Updates(updates).Error
}

func (dao *GORMUserDAO) Insert(ctx context.Context, user *User) error {
	now := time.Now()
	user.CreatedAt = now
	user.UpdatedAt = now
	if user.Birthday != nil && user.Birthday.IsZero() {
		user.Birthday = nil
	}
	err := dao.db.WithContext(ctx).Create(user).Error
	if err != nil {
		var mysqlErr *mysql.MySQLError
		if errors.As(err, &mysqlErr) && mysqlErr.Number == 1062 {
			return ErrUserEmailAlreadyExists
		}
		return err
	}
	return nil
}

func (dao *GORMUserDAO) FindByEmail(ctx context.Context, email string) (*User, error) {
	var user User
	err := dao.db.WithContext(ctx).First(&user, "email = ?", email).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}
	return &user, nil
}

func NewUserDAO(db *gorm.DB) UserDAO {
	return &GORMUserDAO{db: db}
}

type User struct {
	ID        int64      `gorm:"primary_key;auto_increment"`
	Email     string     `gorm:"type:varchar(255);uniqueIndex"`
	Password  string     `gorm:"type:varchar(255);"`
	Nickname  string     `gorm:"type:varchar(255);"`
	Birthday  *time.Time `gorm:"type:datetime;"`
	AboutMe   string     `gorm:"type:varchar(255);"`
	Phone     string
	CreatedAt time.Time
	UpdatedAt time.Time
}
