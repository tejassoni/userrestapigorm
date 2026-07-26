package router

import (
	"userrestapigorm/internal/handlers"

	"github.com/gorilla/mux"
)

/*
Routes for the API are registered version 1
*/
func registerUserV1Routes(route *mux.Router, userHandler *handlers.UserHandler) {
	// Register user-related routes
	route.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	route.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	// route.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	// route.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	// route.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")
}

/*
Routes for the API are registered version 2
*/
func registerUserV2Routes(route *mux.Router, userHandler *handlers.UserHandler) {
	// Register user-related routes
	route.HandleFunc("/users", userHandler.GetUsers).Methods("GET")
	// route.HandleFunc("/users/{id}", userHandler.GetUserByID).Methods("GET")
	// route.HandleFunc("/users", userHandler.CreateUser).Methods("POST")
	// route.HandleFunc("/users/{id}", userHandler.UpdateUser).Methods("PUT")
	// route.HandleFunc("/users/{id}", userHandler.DeleteUser).Methods("DELETE")
}
