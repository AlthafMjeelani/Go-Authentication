package routes

import (
	"database/sql"

	"github.com/gorilla/mux"
)

// PublicRoutes registers routes accessible without authentication
func PublicRoutes(r *mux.Router, db *sql.DB) {
	r.HandleFunc("/register", HandleRegister(db)).Methods("POST")
	r.HandleFunc("/login", HandleLogin(db)).Methods("POST")
	r.HandleFunc("/get_all_users", HandleUsers(db)).Methods("GET")
}
