package server

import (
	"log"
	"net/http"
	"time"

	"github.com/brokeyourbike/lets-go-chat/api/handlers"
	"github.com/brokeyourbike/lets-go-chat/api/middlewares"
	"github.com/brokeyourbike/lets-go-chat/configurations"
	"github.com/brokeyourbike/lets-go-chat/db"
	"github.com/eko/gocache/v2/cache"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"gorm.io/gorm"
)

type server struct {
	router *chi.Mux
	cache  *cache.Cache
	db     *gorm.DB
}

func NewServer(router *chi.Mux, cache *cache.Cache, db *gorm.DB) *server {
	s := server{router: router, cache: cache, db: db}
	s.routes()
	return &s
}

func (s *server) Handle(config *configurations.Config) {
	log.Fatal(http.ListenAndServe(config.Host+":"+config.Port, s.router))
}

func (s *server) routes() {
	rl := middlewares.NewRateLimit(middlewares.RateLimitOpts{Limit: 10, Period: time.Hour})

	u := handlers.NewUsers(db.NewUsersRepo(s.db))

	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)

	s.router.Post("/v1/user", u.HandleUserCreate())
	s.router.Post("/v1/user/login", rl.Handle(u.HandleUserLogin()))
	s.router.Post("/v1/user/active", u.HandleUserActive())
}
