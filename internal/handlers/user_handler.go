package handlers

import (
	"encoding/json"
	"log/slog"
	"net/http"

	"userrestapigorm/internal/config"
	"userrestapigorm/internal/repository"
	"userrestapigorm/internal/responses"
	"userrestapigorm/internal/service"
)

// GetUsers handles GET /users.
func GetUsers(w http.ResponseWriter, r *http.Request) {
	slog.InfoContext(r.Context(), "HTTP GetUsers request",
		"method", r.Method,
		"path", r.URL.Path,
		"ip", r.RemoteAddr,
		"user_agent", r.UserAgent(),
	)

	userRepository := repository.NewUserRepository(config.GormDB)
	userService := service.NewUserService(userRepository)
	users, err := userService.GetUsers(r.Context())
	if err != nil {
		slog.ErrorContext(r.Context(), "get users failed", "error", err)
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
