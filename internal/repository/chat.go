package repository

import (
	"context"

	"gorm.io/gorm"
)

type ChatRepository interface {
	Create(ctx context.Context, chat *models.Chat) error
}

type chatRepo struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *chatRepo {
	return &chatRepo{db: db}
}

func (r *chatRepo) Create(ctx context.Context, title string) (int, error) {
	chat := &Chat{
	}