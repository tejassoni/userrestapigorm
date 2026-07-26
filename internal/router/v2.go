package router

import (
	"userrestapigorm/internal/config"

	"github.com/gorilla/mux"
)

func registerV2Routes(route *mux.Router) {
	v2 := route.PathPrefix(config.API_PREFIX + "/" + config.API_VERSION).Subrouter()

	registerUserV2Routes(v2)
	// registerAuthV2Routes(v2)
}
