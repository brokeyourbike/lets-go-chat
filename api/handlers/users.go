package handlers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/brokeyourbike/lets-go-chat/db"
	"github.com/brokeyourbike/lets-go-chat/models"
	"github.com/brokeyourbike/lets-go-chat/pkg/hasher"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type UsersRepo interface {
	Create(user models.User) error
	GetByUserName(userName string) (models.User, error)
}

type Users struct {
	repo UsersRepo
}

func NewUsers(r UsersRepo) *Users {
	u := Users{repo: r}
	return &u
}

func (u Users) HandleUserCreate() http.HandlerFunc {
	type request struct {
		UserName string `json:"userName" validate:"required,min=4"`
		Password string `json:"password" validate:"required,min=8"`
	}
	type response struct {
		Id       string `json:"id"`
		UserName string `json:"userName"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var data request

		if json.NewDecoder(r.Body).Decode(&data) != nil {
			http.Error(w, "Invalid json", http.StatusBadRequest)
			return
		}

		if validator.New().Struct(&data) != nil {
			http.Error(w, "Invalid username or password", http.StatusBadRequest)
			return
		}

		if _, err := u.repo.GetByUserName(data.UserName); err == nil {
			http.Error(w, fmt.Sprintf("User with userName %s already exists", data.UserName), http.StatusBadRequest)
			return
		}

		hashedPassword, err := hasher.HashPassword(data.Password)
		if err != nil {
			http.Error(w, "Password cannot be hashed", http.StatusInternalServerError)
		}

		user := models.User{ID: uuid.New(), UserName: data.UserName, PasswordHash: hashedPassword}

		if u.repo.Create(user) != nil {
			http.Error(w, "User cannot be created", http.StatusInternalServerError)
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response{Id: user.ID.String(), UserName: user.UserName})
	}
}

func (u Users) HandleUserLogin() http.HandlerFunc {
	type request struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}
	type response struct {
		Url string `json:"url"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		var data request

		if json.NewDecoder(r.Body).Decode(&data) != nil {
			http.Error(w, "Invalid json", http.StatusBadRequest)
			return
		}

		user, err := u.repo.GetByUserName(data.UserName)

		if errors.Is(err, db.ErrUserNotFound) {
			http.Error(w, fmt.Sprintf("User with userName %s not found", data.UserName), http.StatusBadRequest)
			return
		}

		if err != nil {
			http.Error(w, fmt.Sprintf("Cannot query user with userName %s", data.UserName), http.StatusInternalServerError)
			return
		}

		if !hasher.CheckPasswordHash(data.Password, user.PasswordHash) {
			http.Error(w, "Invalid password", http.StatusBadRequest)
			return
		}

		url := fmt.Sprintf("ws://%s/ws?token=%s", r.Host, uuid.NewString())

		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Expires-After", time.Now().Add(time.Minute).UTC().String())
		json.NewEncoder(w).Encode(response{Url: url})
	}
}

func (u Users) HandleUserActive() http.HandlerFunc {
	type response struct {
		Count int `json:"count"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response{Count: 0})
	}
}
