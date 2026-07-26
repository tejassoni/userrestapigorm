package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"log/slog"
	"net/mail"
	"net/http"
	"strconv"
	"strings"
	"time"

	apperrors "userrestapigorm/internal/errors"
	"userrestapigorm/internal/models"
	"userrestapigorm/internal/requests"
	"userrestapigorm/internal/responses"

	"github.com/gorilla/mux"
	"gorm.io/gorm"
)

// UserService defines the application operation required by UserHandler.
type UserService interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUserByID(ctx context.Context, id uint) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, id uint, user *models.User) (*models.User, error)
	DeleteUser(ctx context.Context, id uint) error
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
	var request requests.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, responses.APIResponse{Status: false, Message: "Invalid request body", Error: "bad request"})
		return
	}

	user, err := createUserFromRequest(request)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, responses.APIResponse{Status: false, Message: err.Error(), Error: "validation failed"})
		return
	}

	createdUser, err := h.userService.CreateUser(r.Context(), user)
	if err != nil {
		h.writeUserMutationError(w, r, "create user", err)
		return
	}

	writeJSON(w, http.StatusCreated, responses.APIResponse{Status: true, Message: "User created successfully", Data: createdUser})
}

// UpdateUser handles PUT /users/{id}.
func (h *UserHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	id, ok := h.userIDFromPath(w, r)
	if !ok {
		return
	}

	var request requests.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		writeJSON(w, http.StatusBadRequest, responses.APIResponse{Status: false, Message: "Invalid request body", Error: "bad request"})
		return
	}

	user, err := updateUserFromRequest(request)
	if err != nil {
		writeJSON(w, http.StatusBadRequest, responses.APIResponse{Status: false, Message: err.Error(), Error: "validation failed"})
		return
	}

	updatedUser, err := h.userService.UpdateUser(r.Context(), id, user)
	if err != nil {
		h.writeUserMutationError(w, r, "update user", err)
		return
	}

	writeJSON(w, http.StatusOK, responses.APIResponse{Status: true, Message: "User updated successfully", Data: updatedUser})
}

// DeleteUser handles DELETE /users/{id}.
func (h *UserHandler) DeleteUser(w http.ResponseWriter, r *http.Request) {
	id, ok := h.userIDFromPath(w, r)
	if !ok {
		return
	}

	if err := h.userService.DeleteUser(r.Context(), id); err != nil {
		h.writeUserMutationError(w, r, "delete user", err)
		return
	}

	writeJSON(w, http.StatusOK, responses.APIResponse{Status: true, Message: "User deleted successfully"})
}

func (h *UserHandler) userIDFromPath(w http.ResponseWriter, r *http.Request) (uint, bool) {
	id, err := strconv.ParseUint(mux.Vars(r)["id"], 10, 64)
	if err != nil || id == 0 {
		writeJSON(w, http.StatusBadRequest, responses.APIResponse{Status: false, Message: "Invalid user ID", Error: "bad request"})
		return 0, false
	}

	return uint(id), true
}

func (h *UserHandler) writeUserMutationError(w http.ResponseWriter, r *http.Request, operation string, err error) {
	if apperrors.IsDuplicateEmailError(err) {
		writeJSON(w, http.StatusConflict, responses.APIResponse{Status: false, Message: "Email already exists", Error: "conflict"})
		return
	}
	if errors.Is(err, gorm.ErrRecordNotFound) {
		writeJSON(w, http.StatusNotFound, responses.APIResponse{Status: false, Message: "User not found", Error: "user not found"})
		return
	}
	if errors.Is(err, context.DeadlineExceeded) {
		writeJSON(w, http.StatusGatewayTimeout, responses.APIResponse{Status: false, Message: "Request timed out", Error: "request timeout"})
		return
	}

	h.logger.ErrorContext(r.Context(), operation+" failed", "error", err)
	writeJSON(w, http.StatusInternalServerError, responses.APIResponse{Status: false, Message: "Unable to process user", Error: "internal server error"})
}

func createUserFromRequest(request requests.CreateUserRequest) (*models.User, error) {
	if request.Password != request.ConfirmPassword {
		return nil, errors.New("passwords do not match")
	}
	if len(request.Password) < 8 {
		return nil, errors.New("password must be at least 8 characters")
	}

	user, err := userFromFields(request.Name, request.Email, request.Gender, request.Birthdate, request.IsActive)
	if err != nil {
		return nil, err
	}
	user.Password = request.Password
	return user, nil
}

func updateUserFromRequest(request requests.UpdateUserRequest) (*models.User, error) {
	return userFromFields(request.Name, request.Email, request.Gender, request.Birthdate, request.IsActive)
}

func userFromFields(name, email, gender, birthdate string, isActive bool) (*models.User, error) {
	if len(strings.TrimSpace(name)) < 3 || len(name) > 100 {
		return nil, errors.New("name must be between 3 and 100 characters")
	}
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, errors.New("invalid email")
	}
	if gender != "male" && gender != "female" && gender != "other" {
		return nil, errors.New("gender must be male, female, or other")
	}

	parsedBirthdate, err := time.Parse("2006-01-02", birthdate)
	if err != nil || !parsedBirthdate.Before(time.Now()) {
		return nil, errors.New("birthdate must be a past date in YYYY-MM-DD format")
	}

	return &models.User{Name: strings.TrimSpace(name), Email: email, Gender: gender, Birthdate: parsedBirthdate, IsActive: isActive}, nil
}

func writeJSON(w http.ResponseWriter, status int, response responses.APIResponse) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(response)
}
