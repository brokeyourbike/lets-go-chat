package server

import (
	"log"
	"net/http"
	"time"

	"github.com/brokeyourbike/lets-go-chat/api/handlers"
	"github.com/brokeyourbike/lets-go-chat/api/middlewares"
	"github.com/brokeyourbike/lets-go-chat/cache"
	"github.com/brokeyourbike/lets-go-chat/configurations"
	"github.com/brokeyourbike/lets-go-chat/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type server struct {
	logger *logrus.Logger
	router *chi.Mux
	db     *gorm.DB
}

func NewServer(logger *logrus.Logger, router *chi.Mux, db *gorm.DB) *server {
	s := server{router: router, db: db}
	s.routes()
	return &s
}

func (s *server) Handle(config *configurations.Config) {
	log.Fatal(http.ListenAndServe(config.Host+":"+config.Port, s.router))
}

func (s *server) routes() {
	rl := middlewares.NewRateLimit(middlewares.RateLimitOpts{Limit: 10, Period: time.Hour})

	u := handlers.NewUsers(db.NewUsersRepo(s.db), cache.NewActiveUsersRepo(), db.NewTokensRepo(s.db))

	s.router.Use(middleware.Logger)
	s.router.Use(middlewares.Recoverer(s.logger))

	s.router.Post("/v1/user", u.HandleUserCreate())
	s.router.Post("/v1/user/login", rl.Handle(u.HandleUserLogin()))
	s.router.Get("/v1/user/active", u.HandleUserActive())
	s.router.Get("/v1/chat/ws.rtm.start", u.HandleChat())
}
