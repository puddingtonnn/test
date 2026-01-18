package handler

import (
	"bytes"
	"context"
	"github.com/puddingtonnn/test/internal/domain"
	"github.com/puddingtonnn/test/internal/service"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type MockChatService struct {
	mock.Mock
}

func (m *MockChatService) CreateChat(ctx context.Context, title string) (*domain.Chat, error) {
	args := m.Called(ctx, title)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Chat), args.Error(1)
}

func (m *MockChatService) CreateMessage(ctx context.Context, chatID uint, text string) (*domain.Message, error) {
	args := m.Called(ctx, chatID, text)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Message), args.Error(1)
}

func (m *MockChatService) GetChat(ctx context.Context, chatID uint, limit int) (*domain.Chat, error) {
	args := m.Called(ctx, chatID, limit)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Chat), args.Error(1)
}

func (m *MockChatService) DeleteChat(ctx context.Context, chatID uint) error {
	args := m.Called(ctx, chatID)
	return args.Error(0)
}

func TestCreateChat(t *testing.T) {
	tests := []struct {
		name           string
		inputBody      string
		mockBehavior   func(m *MockChatService)
		expectedStatus int
		expectedBody   string
	}{
		{
			name:      "Success",
			inputBody: `{"title": "Go Developers"}`,
			mockBehavior: func(m *MockChatService) {
				m.On("CreateChat", mock.Anything, "Go Developers").
					Return(&domain.Chat{ID: 1, Title: "Go Developers"}, nil)
			},
			expectedStatus: http.StatusCreated,
			expectedBody:   `"title":"Go Developers"`,
		},
		{
			name:      "Empty Title Error",
			inputBody: `{"title": ""}`,
			mockBehavior: func(m *MockChatService) {
				m.On("CreateChat", mock.Anything, "").
					Return(nil, service.ErrInvalidTitle)
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   `error`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockChatService)
			tt.mockBehavior(mockSvc)

			h := NewHandler(mockSvc)

			req := httptest.NewRequest("POST", "/chats/", bytes.NewBufferString(tt.inputBody))
			w := httptest.NewRecorder()

			h.CreateChat(w, req)

			resp := w.Result()
			assert.Equal(t, tt.expectedStatus, resp.StatusCode)

			mockSvc.AssertExpectations(t)
		})
	}
}

func TestCreateMessage(t *testing.T) {
	tests := []struct {
		name           string
		chatID         string
		inputBody      string
		mockBehavior   func(m *MockChatService)
		expectedStatus int
	}{
		{
			name:      "Success",
			chatID:    "1",
			inputBody: `{"text": "Hello World"}`,
			mockBehavior: func(m *MockChatService) {
				m.On("CreateMessage", mock.Anything, uint(1), "Hello World").
					Return(&domain.Message{ID: 10, ChatID: 1, Text: "Hello World"}, nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:      "Chat Not Found",
			chatID:    "999",
			inputBody: `{"text": "Where am I?"}`,
			mockBehavior: func(m *MockChatService) {
				m.On("CreateMessage", mock.Anything, uint(999), "Where am I?").
					Return(nil, service.ErrNotFound)
			},
			expectedStatus: http.StatusNotFound,
		},
		{
			name:      "Invalid Chat ID",
			chatID:    "abc",
			inputBody: `{"text": "Text"}`,
			mockBehavior: func(m *MockChatService) {
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockSvc := new(MockChatService)
			tt.mockBehavior(mockSvc)

			h := NewHandler(mockSvc)

			url := "/chats/" + tt.chatID + "/messages/"
			req := httptest.NewRequest("POST", url, bytes.NewBufferString(tt.inputBody))
			w := httptest.NewRecorder()

			mux := http.NewServeMux()
			mux.HandleFunc("POST /chats/{id}/messages/", h.CreateMessage)
			mux.ServeHTTP(w, req)

			assert.Equal(t, tt.expectedStatus, w.Result().StatusCode)
			mockSvc.AssertExpectations(t)
		})
	}
}
