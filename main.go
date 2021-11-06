package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
)

type server struct {
	router *chi.Mux
}

func (s *server) routes() {
	s.router.Post("/v1/user", s.handleUserCreate())
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

		s.respond(w, r, response{Id: uuid.New().String(), UserName: data.UserName}, http.StatusOK)
	}
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	srv := server{
		router: chi.NewRouter(),
	}
	srv.routes()
	log.Fatal(http.ListenAndServe("127.0.0.1:"+port, srv.router))
}
