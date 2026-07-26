package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"

	"userrestapigorm/internal/models"
	"userrestapigorm/internal/responses"
)

// UserService defines the application operation required by UserHandler.
type UserService interface {
	GetUsers(ctx context.Context) ([]models.User, error)
}

// UserHandler contains HTTP handlers for user resources.
type UserHandler struct {
	userService UserService
	logger      *slog.Logger
}

func NewUserHandler(userService UserService, logger *slog.Logger) *UserHandler {
	return &UserHandler{
		userService: userService,
		logger:      logger,
	}
}

// GetUsers handles GET /users.
func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
	h.logger.InfoContext(r.Context(), "HTTP GetUsers request",
		"method", r.Method,
		"path", r.URL.Path,
		"ip", r.RemoteAddr,
		"user_agent", r.UserAgent(),
	)

	users, err := h.userService.GetUsers(r.Context())
	if err != nil {
		if errors.Is(err, context.Canceled) {
			h.logger.WarnContext(r.Context(), "get users cancelled", "error", err)
			writeJSON(w, http.StatusRequestTimeout, responses.APIResponse{
				Status:  false,
				Message: "Request cancelled",
				Error:   "request cancelled",
			})
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			h.logger.WarnContext(r.Context(), "get users timed out", "error", err)
			writeJSON(w, http.StatusGatewayTimeout, responses.APIResponse{
				Status:  false,
				Message: "Request timed out",
				Error:   "request timeout",
			})
			return
		}

		h.logger.ErrorContext(r.Context(), "get users failed", "error", err)
		writeJSON(w, http.StatusInternalServerError, responses.APIResponse{
			Status:  false,
			Message: "Unable to get users",
			Error:   "internal server error",
		})
		return
	}

	writeJSON(w, http.StatusOK, responses.APIResponse{
		Status:  true,
		Message: "Users retrieved successfully",
		Data:    users,
	})
}

func writeJSON(w http.ResponseWriter, status int, response responses.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}
