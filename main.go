package main

import (
	"fmt"
	"os"

	"github.com/brokeyourbike/lets-go-chat/api/server"
	"github.com/brokeyourbike/lets-go-chat/configurations"
	"github.com/brokeyourbike/lets-go-chat/models"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
		os.Exit(1)
	}
}

func run() error {
	cfg := configurations.Config{}
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("cannot parse config: %v", err)
	}

	orm, err := gorm.Open(postgres.Open(cfg.Database.Url), &gorm.Config{})
	if err != nil {
		return fmt.Errorf("cannot connect to DB: %v", err)
	}

	orm.AutoMigrate(&models.User{})

	srv := server.NewServer(chi.NewRouter(), orm)
	srv.Handle(&cfg)

	return nil
}
