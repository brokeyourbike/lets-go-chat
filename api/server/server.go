package server

import (
	"log"
	"net/http"
	"time"

	"github.com/brokeyourbike/lets-go-chat/api/handlers"
	"github.com/brokeyourbike/lets-go-chat/configurations"
	"github.com/brokeyourbike/lets-go-chat/db"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/ulule/limiter/v3"
	mhttp "github.com/ulule/limiter/v3/drivers/middleware/stdlib"
	"gorm.io/gorm"
)

type server struct {
	router       *chi.Mux
	limiterStore *limiter.Store
	db           *gorm.DB
}

func NewServer(router *chi.Mux, limiterStore *limiter.Store, db *gorm.DB) *server {
	s := server{router: router, limiterStore: limiterStore, db: db}
	s.routes()
	return &s
}

func (s *server) Handle(config *configurations.Config) {
	log.Fatal(http.ListenAndServe(config.Host+":"+config.Port, s.router))
}

func (s *server) routes() {
	u := handlers.NewUsers(db.NewUsersRepo(s.db))

	s.router.Use(middleware.Logger)
	s.router.Use(middleware.Recoverer)
	s.router.Use(mhttp.NewMiddleware(limiter.New(*s.limiterStore, limiter.Rate{
		Period: time.Hour,
		Limit:  100,
	})).Handler)

	s.router.Post("/v1/user", u.HandleUserCreate())
	s.router.Post("/v1/user/login", u.HandleUserLogin())
}
