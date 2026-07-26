package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"userrestapigorm/internal/models"
	"userrestapigorm/internal/responses"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// UserService defines the application operation required by UserHandler.
type UserService interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
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

// GetUserByID handles GET /users/{id}.
func (h *UserHandler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	h.logger.InfoContext(r.Context(), "HTTP GetUserByID request",
		"method", r.Method,
		"path", r.URL.Path,
		"ip", r.RemoteAddr,
		"user_agent", r.UserAgent(),
	)

	// Extract {id} from the Gorilla Mux route.
	idStr := mux.Vars(r)["id"]
	if idStr == "" {
		h.logger.WarnContext(r.Context(), "missing user ID path parameter")
		writeJSON(w, http.StatusBadRequest, responses.APIResponse{
			Status:  false,
			Message: "User ID is required",
			Error:   "bad request",
		})
		return
	}

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		h.logger.WarnContext(r.Context(), "invalid user ID format", "id", idStr, "error", err)
		writeJSON(w, http.StatusBadRequest, responses.APIResponse{
			Status:  false,
			Message: "Invalid user ID",
			Error:   "invalid ID format",
		})
		return
	}

	user, err := h.userService.GetUserByID(r.Context(), uint(id))
	if err != nil {
		if errors.Is(err, context.Canceled) {
			h.logger.WarnContext(r.Context(), "get user by id cancelled", "error", err)
			writeJSON(w, http.StatusRequestTimeout, responses.APIResponse{
				Status:  false,
				Message: "Request cancelled",
				Error:   "request cancelled",
			})
			return
		}

		if errors.Is(err, context.DeadlineExceeded) {
			h.logger.WarnContext(r.Context(), "get user by id timed out", "error", err)
			writeJSON(w, http.StatusGatewayTimeout, responses.APIResponse{
				Status:  false,
				Message: "Request timed out",
				Error:   "request timeout",
			})
			return
		}

		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.logger.WarnContext(r.Context(), "user not found", "id", id)
			writeJSON(w, http.StatusNotFound, responses.APIResponse{
				Status:  false,
				Message: "User not found",
				Error:   "user not found",
			})
			return
		}

		h.logger.ErrorContext(r.Context(), "get user by id failed", "id", id, "error", err)
		writeJSON(w, http.StatusInternalServerError, responses.APIResponse{
			Status:  false,
			Message: "Unable to get user",
			Error:   "internal server error",
		})
		return
	}

	writeJSON(w, http.StatusOK, responses.APIResponse{
		Status:  true,
		Message: "User retrieved successfully",
		Data:    user,
	})
}

// CreateUser handles POST /users.
func (h *UserHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	// Implementation for creating a new user would go here.
}

// UpdateUser handles PUT /users/{id}.
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	// Implementation for updating an existing user would go here.
}

// DeleteUser handles DELETE /users/{id}.
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Implementation for deleting a user would go here.
}

func writeJSON(w http.ResponseWriter, status int, response responses.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}
