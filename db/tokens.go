package db

import (
	"errors"

	"github.com/brokeyourbike/lets-go-chat/models"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

var ErrTokenNotFound = errors.New("token not found")

type TokensRepo struct {
	db *gorm.DB
}

func NewTokensRepo(db *gorm.DB) *TokensRepo {
	return &TokensRepo{
		db: db,
	}
}

func (t *TokensRepo) Create(token models.Token) (models.Token, error) {
	err := t.db.Create(&token).Error
	return token, err
}

func (t *TokensRepo) Get(id uuid.UUID) (models.Token, error) {
	var token models.Token

	err := t.db.Where("id = ?", id).First(&token).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return token, ErrTokenNotFound
	}

	if err != nil {
		return token, err
	}

	return token, nil
}

func (t *TokensRepo) InvalidateByUserId(userId uuid.UUID) error {
	err := t.db.Where("user_id = ?", userId).Delete(&models.Token{}).Error
	return err
}
