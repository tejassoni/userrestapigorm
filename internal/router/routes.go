package router

import (
	"userrestapigorm/internal/config"
	"userrestapigorm/internal/handlers"

	"github.com/gorilla/mux"
)

func New(userHandler *handlers.UserHandler, cfg *config.Config) *mux.Router {
	router := mux.NewRouter()

	registerHealthRoutes(router, cfg)
	registerV1Routes(router, userHandler, cfg)
	registerV2Routes(router, userHandler, cfg)

	return router
}
