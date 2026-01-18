package service

import (
	"context"
	"errors"
	"github.com/puddingtonnn/test/internal/domain"
	"github.com/puddingtonnn/test/internal/repository"
	"gorm.io/gorm"
	"strings"
)

var (
	ErrInvalidTitle = errors.New("title must be between 1 and 200 chars")
	ErrInvalidText  = errors.New("text must be between 1 and 5000 chars")
	ErrNotFound     = errors.New("not found")
)

type ChatService struct {
	repo *repository.ChatRepository
}

func NewChatService(repo *repository.ChatRepository) *ChatService {
	return &ChatService{repo: repo}
}

func (s *ChatService) CreateChat(ctx context.Context, title string) (*domain.Chat, error) {
	title = strings.TrimSpace(title)
	if len(title) == 0 || len(title) > 200 {
		return nil, ErrInvalidTitle
	}

	chat := &domain.Chat{Title: title}
	if err := s.repo.CreateChat(ctx, chat); err != nil {
		return nil, err
	}
	return chat, nil
}

func (s *ChatService) CreateMessage(ctx context.Context, chatID uint, text string) (*domain.Message, error) {
	if len(text) == 0 || len(text) > 5000 {
		return nil, ErrInvalidText
	}

	msg := &domain.Message{
		ChatID: chatID,
		Text:   text,
	}

	err := s.repo.CreateMessage(ctx, msg)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return msg, err
}

func (s *ChatService) GetChat(ctx context.Context, chatID uint, limit int) (*domain.Chat, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	chat, err := s.repo.GetChatWithMessages(ctx, chatID, limit)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrNotFound
	}
	return chat, err
}

func (s *ChatService) DeleteChat(ctx context.Context, chatID uint) error {
	err := s.repo.DeleteChat(ctx, chatID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}
