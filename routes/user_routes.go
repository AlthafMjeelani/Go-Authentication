package routes

import (
	"database/sql"
	"fmt"
	"golang_projects/repository"
	utils "golang_projects/utility"
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
func HandleGetUserByEmail(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		email := r.URL.Query().Get("email")
		if email == "" {
			utils.WriteJSONResponse(w, http.StatusBadRequest, false, "Email query parameter is required", nil)
			return
		}

		user, err := repository.GetUserByEmail(db, email)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, false, "Failed to fetch user", nil)
			log.Printf("GetUserByEmail error: %v", err)
			return
		}

		utils.WriteJSONResponse(w, http.StatusOK, true, "Success", user)
	}
}

// Get all users
func getUsers(db *sql.DB, w http.ResponseWriter) {
	users, err := repository.GetAllUsers(db)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		log.Printf("GetUsers error: %v", err)
		return
	}

	utils.WriteJSONResponse(w, http.StatusOK, true, "Success", users)
}
