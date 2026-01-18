package handler

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/puddingtonnn/test/internal/domain"
	"github.com/puddingtonnn/test/internal/service"
	"net/http"
	"strconv"
)

type ChatService interface {
	CreateChat(ctx context.Context, title string) (*domain.Chat, error)
	CreateMessage(ctx context.Context, chatID uint, text string) (*domain.Message, error)
	GetChat(ctx context.Context, chatID uint, limit int) (*domain.Chat, error)
	DeleteChat(ctx context.Context, chatID uint) error
}

type Handler struct {
	svc ChatService
}

func NewHandler(svc ChatService) *Handler {
	return &Handler{svc: svc}
}

func jsonError(w http.ResponseWriter, err error, code int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

// CreateChat POST /chats/
func (h *Handler) CreateChat(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title string `json:"title"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, err, http.StatusBadRequest)
		return
	}

	chat, err := h.svc.CreateChat(r.Context(), req.Title)
	if err != nil {
		if errors.Is(err, service.ErrInvalidTitle) {
			jsonError(w, err, http.StatusBadRequest)
			return
		}
		jsonError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(chat)
}

// CreateMessage POST /chats/{id}/messages/
func (h *Handler) CreateMessage(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	chatID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		jsonError(w, errors.New("invalid chat id"), http.StatusBadRequest)
		return
	}

	var req struct {
		Text string `json:"text"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		jsonError(w, err, http.StatusBadRequest)
		return
	}

	msg, err := h.svc.CreateMessage(r.Context(), uint(chatID), req.Text)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			jsonError(w, err, http.StatusNotFound) // 404
			return
		}
		jsonError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(msg)
}

// GetChat GET /chats/{id}
func (h *Handler) GetChat(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	chatID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		jsonError(w, errors.New("invalid chat id"), http.StatusBadRequest)
		return
	}

	limitStr := r.URL.Query().Get("limit")
	limit := 20
	if limitStr != "" {
		limit, _ = strconv.Atoi(limitStr)
	}

	chat, err := h.svc.GetChat(r.Context(), uint(chatID), limit)
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			jsonError(w, err, http.StatusNotFound)
			return
		}
		jsonError(w, err, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(chat)
}

// DeleteChat DELETE /chats/{id}
func (h *Handler) DeleteChat(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	chatID, err := strconv.ParseUint(idStr, 10, 32)
	if err != nil {
		jsonError(w, errors.New("invalid chat id"), http.StatusBadRequest)
		return
	}

	err = h.svc.DeleteChat(r.Context(), uint(chatID))
	if err != nil {
		if errors.Is(err, service.ErrNotFound) {
			jsonError(w, err, http.StatusNotFound)
			return
		}
		jsonError(w, err, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
