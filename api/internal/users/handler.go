package users

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/google/uuid"
)

type ServiceInterface interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (*User, error)
	ListUsers(ctx context.Context, limit, offset int) ([]*User, error)
}

type Handler struct {
	service ServiceInterface
	logger  *slog.Logger
}

func NewHandler(service ServiceInterface, logger *slog.Logger) *Handler {
	return &Handler{
		service: service,
		logger:  logger,
	}
}

func (h *Handler) HandleGetUser() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Extract user ID from URL path
		idStr := r.PathValue("id")
		if idStr == "" {
			h.respondWithError(w, http.StatusBadRequest, "user ID is required")
			return
		}

		// Parse UUID
		userID, err := uuid.Parse(idStr)
		if err != nil {
			h.respondWithError(w, http.StatusBadRequest, "invalid user ID format")
			return
		}

		// Get user from service
		user, err := h.service.GetUserByID(r.Context(), userID)
		if err != nil {
			if errors.Is(err, ErrUserNotFound) {
				h.respondWithError(w, http.StatusNotFound, "user not found")
				return
			}
			h.logger.Error("failed to get user", "error", err, "user_id", userID)
			h.respondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		// Respond with user
		h.respondWithJSON(w, http.StatusOK, user)
	}
}

func (h *Handler) respondWithJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	if err := json.NewEncoder(w).Encode(payload); err != nil {
		h.logger.Error("failed to encode response", "error", err)
	}
}

func (h *Handler) respondWithError(w http.ResponseWriter, code int, message string) {
	h.respondWithJSON(w, code, map[string]string{"error": message})
}

func (h *Handler) HandleListUsers() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse query parameters
		limitStr := r.URL.Query().Get("limit")
		offsetStr := r.URL.Query().Get("offset")

		limit := 20 // default
		if limitStr != "" {
			parsedLimit, err := strconv.Atoi(limitStr)
			if err != nil || parsedLimit < 0 {
				h.respondWithError(w, http.StatusBadRequest, "invalid limit parameter")
				return
			}
			limit = parsedLimit
		}

		offset := 0 // default
		if offsetStr != "" {
			parsedOffset, err := strconv.Atoi(offsetStr)
			if err != nil || parsedOffset < 0 {
				h.respondWithError(w, http.StatusBadRequest, "invalid offset parameter")
				return
			}
			offset = parsedOffset
		}

		// Get users from service
		users, err := h.service.ListUsers(r.Context(), limit, offset)
		if err != nil {
			h.logger.Error("failed to list users", "error", err)
			h.respondWithError(w, http.StatusInternalServerError, "internal server error")
			return
		}

		// Respond with users
		h.respondWithJSON(w, http.StatusOK, map[string]any{
			"users":  users,
			"limit":  limit,
			"offset": offset,
		})
	}
}
