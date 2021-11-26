package db

import (
	"errors"

	"github.com/brokeyourbike/lets-go-chat/models"
	"gorm.io/gorm"
)

var ErrUserNotFound = errors.New("user not found")

type UsersRepo struct {
	db *gorm.DB
}

func NewUsersRepo(db *gorm.DB) *UsersRepo {
	return &UsersRepo{
		db: db,
	}
}

func (u *UsersRepo) Create(user models.User) error {
	result := u.db.Create(&user)
	return result.Error
}

func (u *UsersRepo) GetByUserName(userName string) (models.User, error) {
	var user models.User

	err := u.db.Where("user_name = ?", userName).First(&user).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return user, ErrUserNotFound
	}

	if err != nil {
		return user, err
	}

	return user, nil
}
