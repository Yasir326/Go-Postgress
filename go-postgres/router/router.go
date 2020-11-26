package router

import (
	"go-postgres/go-postgres/middleware"
	"github.com/gorilla/mux"
)

// Router is exported and used in main.go
func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("api/user/{id}", middleware.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("api/user", middleware.GetAllUsers).Methods("GET", "OPTIONS")
	router.HandleFunc("api/newuser", middleware.CreateUser).Methods("POSTS", "OPTIONS")
	router.HandleFunc("api/user/{id}", middleware.UpdateUser).Methods("PUT", "OPTIONS")
	router.HandleFunc("api/deleteuser/{id}", middleware.DeleteUser).Methods("DELETE", "OPTIONS")

	return router
}