package user

import (
	"fmt"

	"github.com/google/uuid"
)

type User struct {
	UserName     string
	PasswordHash string
}

type Users map[string]User

func (u Users) AddUser(user User) string {
	key := uuid.NewString()
	u[key] = user
	return key
}

func (u Users) GetUserByUserName(userName string) (User, error) {
	for _, user := range u {
		if user.UserName == userName {
			return user, nil
		}
	}

	return User{}, fmt.Errorf("user: with UserName %v not found", userName)
}
