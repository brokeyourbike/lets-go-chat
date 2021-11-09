package main

import (
	"fmt"

	"github.com/brokeyourbike/lets-go-chat/api/server"
	"github.com/brokeyourbike/lets-go-chat/configurations"
	"github.com/brokeyourbike/lets-go-chat/models"
	"github.com/caarlos0/env/v6"
	"github.com/go-chi/chi/v5"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	cfg := configurations.Config{}
	if err := env.Parse(&cfg); err != nil {
		fmt.Printf("%+v\n", err)
	}

	orm, err := gorm.Open(postgres.Open(cfg.Database.Url), &gorm.Config{})
	if err != nil {
		fmt.Printf("%+v\n", err)
	}

	orm.AutoMigrate(&models.User{})

	srv := server.NewServer(chi.NewRouter(), orm)
	srv.Handle(&cfg)
}
