package routes

import (
	"database/sql"
	"golang_projects/middleware"

	"github.com/gorilla/mux"
)

// PrivateRoutes registers routes that require authentication
func PrivateRoutes(r *mux.Router, db *sql.DB) {

	r.HandleFunc("/users_details", middleware.JWTAuthMiddleware(HandleGetUserByEmail(db))).Methods("GET")
	r.HandleFunc("/update_user", middleware.JWTAuthMiddleware(HandleUpdateUser(db))).Methods("PUT", "PATCH")
	r.HandleFunc("/delete_user", middleware.JWTAuthMiddleware(HandleDeleteUser(db))).Methods("DELETE")
}
