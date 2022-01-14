package db

import (
	"time"

	"github.com/brokeyourbike/lets-go-chat/models"
	"github.com/google/uuid"
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

func (t *MessagesRepo) GetAfterDateExcludingUserId(after time.Time, userId uuid.UUID) ([]models.Message, error) {
	var messages []models.Message
	err := t.db.Where("user_id != ?", userId).Where("created_at < ?", after).Find(&messages).Error
	return messages, err
}
