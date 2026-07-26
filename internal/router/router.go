package router

import (
	"encoding/json"
	"net/http"
	"userrestapigorm/internal/config"

	"github.com/gorilla/mux"
)

func registerHealthRoutes(router *mux.Router, cfg *config.Config) {
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		response := map[string]interface{}{
			"status":      true,
			"message":     "User REST API is running",
			"data":        nil,
			"app_name":    cfg.App.Name,
			"app_version": cfg.App.Version,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}).Methods("GET")
}
