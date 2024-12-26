package routes

import (
	"database/sql"
	"encoding/json"
	models "golang_projects/model"
	"log"
	"net/http"
)

// HandleUsers handles user-related API requests
func HandleUsers(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getUsers(db, w)
		case http.MethodPost:
			createUser(db, w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	}
}

// Get all users
func getUsers(db *sql.DB, w http.ResponseWriter) {
	users, err := models.GetAllUsers(db)
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		log.Printf("GetUsers error: %v", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// Create a new user
func createUser(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		log.Printf("CreateUser decode error: %v", err)
		return
	}

	err = models.CreateUser(db, user)
	if err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		log.Printf("CreateUser insert error: %v", err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("User created successfully"))
}
