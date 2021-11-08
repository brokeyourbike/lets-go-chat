package server

import (
	"log"
	"net/http"

	"github.com/brokeyourbike/lets-go-chat/api/handlers"
	"github.com/brokeyourbike/lets-go-chat/configurations"
	"github.com/go-chi/chi/v5"
)

type server struct {
	router    *chi.Mux
	usersRepo handlers.UsersRepo
}

func NewServer(router *chi.Mux, usersRepo handlers.UsersRepo) *server {
	s := server{router: router, usersRepo: usersRepo}
	s.routes()
	return &s
}

func (s *server) Handle(config *configurations.Config) {
	log.Fatal(http.ListenAndServe(config.Server.Host+":"+config.Server.Port, s.router))
}

func (s *server) routes() {
	u := handlers.NewUsers(s.usersRepo)

	s.router.Post("/v1/user", u.HandleUserCreate())
	s.router.Post("/v1/user/login", u.HandleUserLogin())
}
