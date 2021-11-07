package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/brokeyourbike/lets-go-chat/internal/user"
	"github.com/brokeyourbike/lets-go-chat/pkg/hasher"
	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type server struct {
	router *chi.Mux
	users  user.Users
}

func (s *server) routes() {
	s.router.Post("/v1/user", s.handleUserCreate())
	s.router.Post("/v1/user/login", s.handleUserLogin())
}

func (s *server) decode(w http.ResponseWriter, r *http.Request, v interface{}) error {
	return json.NewDecoder(r.Body).Decode(v)
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, data interface{}, status int) {
	if data != nil {
		w.Header().Set("Content-type", "application/json")
	}

	w.WriteHeader(status)

	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func (s *server) handleUserCreate() http.HandlerFunc {
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
		err := s.decode(w, r, &data)

		if err != nil {
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		err = validator.New().Struct(&data)

		if err != nil {
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		_, err = s.users.GetUserByUserName(data.UserName)

		if err == nil {
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		hashedPassword, _ := hasher.HashPassword(data.Password)
		userId := s.users.AddUser(user.User{UserName: data.UserName, PasswordHash: hashedPassword})

		s.respond(w, r, response{Id: userId, UserName: data.UserName}, http.StatusOK)
	}
}

func (s *server) handleUserLogin() http.HandlerFunc {
	type request struct {
		UserName string `json:"userName"`
		Password string `json:"password"`
	}
	type response struct {
		Url string `json:"url"`
	}
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("X-Rate-Limit", "60")

		var data request
		err := s.decode(w, r, &data)

		if err != nil {
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		usr, err := s.users.GetUserByUserName(data.UserName)

		if err != nil {
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		if !hasher.CheckPasswordHash(data.Password, usr.PasswordHash) {
			s.respond(w, r, nil, http.StatusBadRequest)
			return
		}

		url := fmt.Sprintf("ws://%s/ws?token=%s", r.Host, uuid.NewString())

		w.Header().Set("X-Expires-After", time.Now().Add(time.Minute).UTC().String())

		s.respond(w, r, response{Url: url}, http.StatusOK)
	}
}

func main() {
	host := os.Getenv("HOST")
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	srv := server{
		router: chi.NewRouter(),
		users:  make(user.Users),
	}

	srv.routes()

	log.Fatal(http.ListenAndServe(host+":"+port, srv.router))
}
