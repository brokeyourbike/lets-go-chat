package main

import (
	"fmt"
	"os"

	"github.com/brokeyourbike/lets-go-chat/api/handlers"
	"github.com/brokeyourbike/lets-go-chat/api/server"
	"github.com/brokeyourbike/lets-go-chat/cache"
	"github.com/brokeyourbike/lets-go-chat/configurations"
	"github.com/brokeyourbike/lets-go-chat/db"
	"github.com/brokeyourbike/lets-go-chat/models"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	log.SetFormatter(&log.JSONFormatter{})

	if err := run(); err != nil {
		log.Fatalf("%s\n", err)
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
	orm.AutoMigrate(&models.Token{})

	users := handlers.NewUsers(db.NewUsersRepo(orm), cache.NewActiveUsersRepo(), db.NewTokensRepo(orm))

	srv := server.NewServer(chi.NewRouter())
	srv.Routes(users)
	srv.Handle(&cfg)

	return nil
}
