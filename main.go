package main

import (
	"fmt"

	"github.com/brokeyourbike/lets-go-chat/api/server"
	"github.com/brokeyourbike/lets-go-chat/configurations"
	"github.com/brokeyourbike/lets-go-chat/db"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
)

func main() {
	cfg := configurations.Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	users := make(db.Users)
	srv := server.NewServer(chi.NewRouter(), &users)
	srv.Handle(&cfg)
}
