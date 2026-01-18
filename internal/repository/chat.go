package repository

import (
	"context"
	"github.com/puddingtonnn/test/internal/domain"
	"gorm.io/gorm"
)

type ChatRepository struct {
	db *gorm.DB
}

func NewChatRepository(db *gorm.DB) *ChatRepository {
	return &ChatRepository{db: db}
}

func (r *ChatRepository) CreateChat(ctx context.Context, chat *domain.Chat) error {
	return r.db.WithContext(ctx).Create(chat).Error
}

func (r *ChatRepository) CreateMessage(ctx context.Context, msg *domain.Message) error {
	var count int
	if err := r.db.WithContext(ctx).Model(&domain.Chat{}).Where("id = ?", msg.ChatID).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		return gorm.ErrRecordNotFound
	}
	return r.db.WithContext(ctx).Create(msg).Error
}

func (r *ChatRepository) GetChatWithMessages(ctx context.Context, chatID uint, limit int) (*domain.Chat, error) {
	var chat domain.Chat

	err := r.db.WithContext(ctx).
		Preload("Messages", func(db *gorm.DB) *gorm.DB {
			return db.Order("created_at DESC").Limit(limit)
		}).
		First(&chat, chatID).Error

	return &chat, err
}

func (r *ChatRepository) DeleteChat(ctx context.Context, chatID uint) error {
	result := r.db.WithContext(ctx).Delete(&domain.Chat{}, chatID)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
