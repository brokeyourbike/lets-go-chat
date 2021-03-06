package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/brokeyourbike/lets-go-chat/api/server"
	"github.com/brokeyourbike/lets-go-chat/db"
	"github.com/brokeyourbike/lets-go-chat/models"
	"github.com/brokeyourbike/lets-go-chat/pkg/hasher"
	"github.com/goccy/go-json"
	"github.com/google/uuid"
)

type UsersRepo interface {
	Create(user models.User) (models.User, error)
	GetByUserName(userName string) (models.User, error)
}

type ActiveUsersRepo interface {
	Add(userId uuid.UUID) error
	Delete(userId uuid.UUID) error
	Count() int
}

type TokensRepo interface {
	Create(token models.Token) (models.Token, error)
	Get(id uuid.UUID) (models.Token, error)
	InvalidateByUserId(userId uuid.UUID) error
}

type Users struct {
	usersRepo       UsersRepo
	activeUsersRepo ActiveUsersRepo
	tokensRepo      TokensRepo
}

func NewUsers(u UsersRepo, a ActiveUsersRepo, t TokensRepo) *Users {
	return &Users{usersRepo: u, activeUsersRepo: a, tokensRepo: t}
}

func (u *Users) HandleUserCreate(w http.ResponseWriter, r *http.Request) {
	var data server.CreateUserRequest

	if json.NewDecoder(r.Body).Decode(&data) != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	if len(data.UserName) < 4 {
		http.Error(w, "UserName is too short", http.StatusBadRequest)
		return
	}

	if len(data.Password) < 8 {
		http.Error(w, "Password is too short", http.StatusBadRequest)
		return
	}

	if _, err := u.usersRepo.GetByUserName(data.UserName); err == nil {
		http.Error(w, fmt.Sprintf("User with userName %s already exists", data.UserName), http.StatusBadRequest)
		return
	}

	hashedPassword, err := hasher.HashPassword(data.Password)
	if err != nil {
		http.Error(w, "Password cannot be hashed", http.StatusInternalServerError)
		return
	}

	user, err := u.usersRepo.Create(models.User{ID: uuid.New(), UserName: data.UserName, PasswordHash: hashedPassword})
	if err != nil {
		http.Error(w, "User cannot be created", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(server.CreateUserResponse{Id: user.ID.String(), UserName: user.UserName})
}

func (u *Users) HandleUserLogin(w http.ResponseWriter, r *http.Request) {
	var data server.LoginUserRequest

	if json.NewDecoder(r.Body).Decode(&data) != nil {
		http.Error(w, "Invalid json", http.StatusBadRequest)
		return
	}

	user, err := u.usersRepo.GetByUserName(data.UserName)
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

	token, err := u.tokensRepo.Create(models.Token{ID: uuid.New(), UserID: user.ID, ExpiresAt: time.Now().Add(time.Hour)})
	if err != nil {
		http.Error(w, "Cannot create token for user", http.StatusInternalServerError)
		return
	}

	url := fmt.Sprintf("ws://%s/chat/ws.rtm.start?token=%s", r.Host, token.ID.String())

	w.Header().Set("X-Expires-After", token.ExpiresAt.UTC().String())
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(server.LoginUserResonse{Url: url})
}

func (u *Users) HandleUserActive(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(server.ActiveUsersResponse{Count: u.activeUsersRepo.Count()})
}
