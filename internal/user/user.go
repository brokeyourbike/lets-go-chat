package user

import (
	"fmt"
)

type User struct {
	UserName     string
	PasswordHash string
}

type Users map[string]User

func (u Users) AddUser(id string, user User) string {
	u[id] = user
	return id
}

func (u Users) GetUserByUserName(userName string) (User, error) {
	for _, user := range u {
		if user.UserName == userName {
			return user, nil
		}
	}

	return User{}, fmt.Errorf("user: with UserName %v not found", userName)
}
