//go:build wireinject
// +build wireinject

package main

import (
	"github.com/brokeyourbike/lets-go-chat/api/handlers"
	"github.com/brokeyourbike/lets-go-chat/cache"
	"github.com/brokeyourbike/lets-go-chat/db"
	"github.com/google/wire"
	"gorm.io/gorm"
)

func NewUsers(orm *gorm.DB) *handlers.Users {
	wire.Build(
		db.NewUsersRepo,
		db.NewTokensRepo,
		cache.NewActiveUsersRepo,
		wire.Bind(new(handlers.UsersRepo), new(*db.UsersRepo)),
		wire.Bind(new(handlers.TokensRepo), new(*db.TokensRepo)),
		wire.Bind(new(handlers.ActiveUsersRepo), new(*cache.ActiveUsersRepo)),
		handlers.NewUsers,
	)
	return &handlers.Users{}
}

func NewChat(orm *gorm.DB) *handlers.Chat {
	wire.Build(
		handlers.NewHub,
		cache.NewActiveUsersRepo,
		db.NewTokensRepo,
		db.NewMessagesRepo,
		wire.Bind(new(handlers.ActiveUsersRepo), new(*cache.ActiveUsersRepo)),
		wire.Bind(new(handlers.TokensRepo), new(*db.TokensRepo)),
		wire.Bind(new(handlers.MessagesRepo), new(*db.MessagesRepo)),
		handlers.NewChat,
	)
	return &handlers.Chat{}
}
