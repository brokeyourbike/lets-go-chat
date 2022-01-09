package db

import (
	"github.com/brokeyourbike/lets-go-chat/models"
	"gorm.io/gorm"
)

type MessagesRepo struct {
	db *gorm.DB
}

func NewMessagesRepo(db *gorm.DB) *MessagesRepo {
	return &MessagesRepo{
		db: db,
	}
}

func (t *MessagesRepo) Create(msg models.Message) (models.Message, error) {
	err := t.db.Create(&msg).Error
	return msg, err
}
