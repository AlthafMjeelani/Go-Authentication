package routes

import (
	"database/sql"
	"golang_projects/repository"
	"log"
	"net/http"
)

// HandleUsers handles user-related API requests
func HandleUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getUsers(db, w)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// Get all users
func getUsers(db *sql.DB, w http.ResponseWriter) {
	users, err := repository.GetAllUsers(db)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		log.Printf("GetUsers error: %v", err)
		return
	}

	writeJSONResponse(w, http.StatusOK, true, "Success", users)
}
