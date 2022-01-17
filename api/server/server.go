package server

import (
	"net/http"

	"github.com/brokeyourbike/lets-go-chat/api/middlewares"
	"github.com/brokeyourbike/lets-go-chat/configurations"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

type UsersHandler interface {
	HandleUserCreate(w http.ResponseWriter, r *http.Request)
	HandleUserLogin(w http.ResponseWriter, r *http.Request)
	HandleUserActive(w http.ResponseWriter, r *http.Request)
}

type ChatHandler interface {
	HandleChat(w http.ResponseWriter, r *http.Request, params WsRTMStartParams)
}

type server struct {
	users UsersHandler
	chat  ChatHandler
}

func NewServer(users UsersHandler, chat ChatHandler) *server {
	s := server{users: users, chat: chat}
	return &s
}

func (s *server) WsRTMStart(w http.ResponseWriter, r *http.Request, params WsRTMStartParams) {
	s.chat.HandleChat(w, r, params)
}

func (s *server) CreateUser(w http.ResponseWriter, r *http.Request) {
	s.users.HandleUserCreate(w, r)
}

func (s *server) GetActiveUsers(w http.ResponseWriter, r *http.Request) {
	s.users.HandleUserActive(w, r)
}

func (s *server) LoginUser(w http.ResponseWriter, r *http.Request) {
	s.users.HandleUserLogin(w, r)
}

func (s *server) Handle(config *configurations.Config, router *chi.Mux) {
	router.Use(middlewares.Logger)
	router.Use(middlewares.ErrorLogger)
	router.Use(middlewares.Recoverer)

	log.Fatal(http.ListenAndServe(config.Host+":"+config.Port, HandlerFromMux(s, router)))
}
