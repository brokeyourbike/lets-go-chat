package db

import (
	"fmt"

	"github.com/brokeyourbike/lets-go-chat/models"
)

type Users map[string]models.User

func (u Users) Create(user models.User) error {
	u[user.Id.String()] = user
	return nil
}

func (u Users) GetByUserName(userName string) (models.User, error) {
	for _, user := range u {
		if user.UserName == userName {
			return user, nil
		}
	}

	return models.User{}, fmt.Errorf("user: with UserName %v not found", userName)
}
