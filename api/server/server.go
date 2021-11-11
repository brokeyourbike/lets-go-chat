package server

import (
	"log"
	"net/http"

	"github.com/brokeyourbike/lets-go-chat/api/handlers"
	"github.com/brokeyourbike/lets-go-chat/configurations"
	"github.com/brokeyourbike/lets-go-chat/db"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type server struct {
	router *chi.Mux
	db     *gorm.DB
}

func NewServer(router *chi.Mux, db *gorm.DB) *server {
	s := server{router: router, db: db}
	s.routes()
	return &s
}

func (s *server) Handle(config *configurations.Config) {
	log.Fatal(http.ListenAndServe(config.Host+":"+config.Port, s.router))
}

func (s *server) routes() {
	u := handlers.NewUsers(db.NewUsersRepo(s.db))

	s.router.Post("/v1/user", u.HandleUserCreate())
	s.router.Post("/v1/user/login", u.HandleUserLogin())
}
