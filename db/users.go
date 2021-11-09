package db

import (
	"github.com/brokeyourbike/lets-go-chat/models"
	"gorm.io/gorm"
)

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

func (u UsersRepo) GetByUserName(userName string) (models.User, error) {
	var user models.User
	result := u.db.Where("user_name = ?", userName).First(&user)
	return user, result.Error
}
