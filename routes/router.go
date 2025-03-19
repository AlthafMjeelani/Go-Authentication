package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
)

// SetupRoutes initializes all routes
func SetupRoutes(db *sql.DB) *mux.Router {
	router := mux.NewRouter()

	// API v1 routes
	apiV1 := router.PathPrefix("/api/v1").Subrouter()

	// Public routes
	public := apiV1.PathPrefix("/public").Subrouter()
	PublicRoutes(public, db)

	// Private routes (Require JWT Auth)
	private := apiV1.PathPrefix("/mobile").Subrouter()
	PrivateRoutes(private, db)

	return router
}
