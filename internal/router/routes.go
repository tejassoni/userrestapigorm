package router

import (
	"userrestapigorm/internal/handlers"

	"github.com/gorilla/mux"
)

func New(userHandler *handlers.UserHandler) *mux.Router {
	router := mux.NewRouter()

	registerHealthRoutes(router)
	registerV1Routes(router, userHandler)
	registerV2Routes(router, userHandler)

	return router
}
