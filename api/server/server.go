package server

import (
	"net/http"
	"time"

	"github.com/brokeyourbike/lets-go-chat/api/middlewares"
	"github.com/brokeyourbike/lets-go-chat/configurations"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

type UsersHandler interface {
	HandleUserCreate() http.HandlerFunc
	HandleUserLogin() http.HandlerFunc
	HandleUserActive() http.HandlerFunc
	HandleChat() http.HandlerFunc
}

type server struct {
	router *chi.Mux
}

func NewServer(router *chi.Mux) *server {
	s := server{router: router}
	return &s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) Handle(config *configurations.Config) {
	log.Fatal(http.ListenAndServe(config.Host+":"+config.Port, s.router))
}

func (s *server) Routes(u UsersHandler) {
	rl := middlewares.NewRateLimit(middlewares.RateLimitOpts{Limit: 10, Period: time.Hour})

	s.router.Use(middlewares.Logger)
	s.router.Use(middlewares.ErrorLogger)
	s.router.Use(middlewares.Recoverer)

	s.router.Post("/v1/user", u.HandleUserCreate())
	s.router.Post("/v1/user/login", rl.Handle(u.HandleUserLogin()))
	s.router.Get("/v1/user/active", u.HandleUserActive())
	s.router.Get("/v1/chat/ws.rtm.start", u.HandleChat())
}
