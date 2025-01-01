package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang_projects/model"
	"golang_projects/repository"
	"log"
	"net/http"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// Response defines the structure for both success and error messages
type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type UpdateUserRequest struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
	// Add other fields as needed
}

// HandleRegister handles user registration
func HandleRegister(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSONResponse(w, http.StatusMethodNotAllowed, false, "Method not allowed", nil)
			return
		}

		var user model.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			writeJSONResponse(w, http.StatusBadRequest, false, "Invalid request body", nil)
			log.Printf("Register decode error: %v", err)
			return
		}

		// Validate user input
		if err := validateUser(user, true); err != nil {
			writeJSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
			log.Printf("Validation error: %v", err)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			writeJSONResponse(w, http.StatusInternalServerError, false, "Failed to hash password", nil)
			log.Printf("Hash password error: %v", err)
			return
		}
		user.Password = string(hashedPassword)

		// Save user to the database
		err = repository.CreateUser(db, user)
		if err != nil {
			writeJSONResponse(w, http.StatusInternalServerError, false, "Failed to register user", nil)
			log.Printf("Register insert error: %v", err)
			return
		}

		// Success response
		writeJSONResponse(w, http.StatusCreated, true, "User registered successfully", nil)
	}
}

func writeJSONResponse(w http.ResponseWriter, statusCode int, status bool, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Status:  status,
		Message: message,
		Data:    data,
	})
}

// validateUser validates the user input
func validateUser(user model.User, isRegister bool) error {
	if isRegister {
		if user.Name == "" {
			return fmt.Errorf("name is required")
		}
		if len(user.Name) < 3 {
			return fmt.Errorf("name must be at least 3 characters long")
		}
	}

	if user.Email == "" {
		return fmt.Errorf("email is required")
	}
	if !isValidEmail(user.Email) {
		return fmt.Errorf("email is not valid")
	}

	if user.Password == "" {
		return fmt.Errorf("password is required")
	}
	if len(user.Password) < 6 || !containsSpecialChar(user.Password) {
		return fmt.Errorf("password must be at least 6 characters long and include at least one special character")
	}

	return nil
}

// isValidEmail checks if the email is valid
func isValidEmail(email string) bool {
	re := `^[a-z0-9._%+-]+@[a-z0-9.-]+\.[a-z]{2,}$`
	matched, _ := regexp.MatchString(re, email)
	return matched
}

// containsSpecialChar checks if a string contains a special character
func containsSpecialChar(s string) bool {
	for _, char := range s {
		if unicode.IsSymbol(char) || unicode.IsPunct(char) {
			return true
		}
	}
	return false
}

// HandleLogin handles user login
func HandleLogin(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSONResponse(w, http.StatusMethodNotAllowed, false, "Method not allowed", nil)
			return
		}

		var credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			writeJSONResponse(w, http.StatusMethodNotAllowed, false, "Invalid request body", nil)
			return
		}

		err = validateUser(model.User{Email: credentials.Email, Password: credentials.Password}, false)
		if err != nil {
			writeJSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
			log.Printf("Validation error: %v", err)
			return
		}

		// Retrieve user from the database
		user, err := repository.GetUserByEmail(db, credentials.Email)
		if err != nil {
			writeJSONResponse(w, http.StatusInternalServerError, false, "Invalid password or email", nil)
			log.Printf("Login fetch error: %v", err)
			return
		}

		// Compare passwords
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
		if err != nil {
			writeJSONResponse(w, http.StatusInternalServerError, false, "Provided password is incorrect", nil)
			log.Printf("Login password mismatch: %v", err)
			return
		}

		// Success response
		response := struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
		}{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		}

		// writeJSONResponse(w, http.StatusOK, true, "User logged in successfully")
		writeJSONResponse(w, http.StatusOK, true, "User logged in successfully", response)
		log.Printf("User logged in successfully")

	}

}

// update User by ID
// HandleUpdateUser handles updating user fields
func HandleUpdateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut && r.Method != http.MethodPatch {
			writeJSONResponse(w, http.StatusMethodNotAllowed, false, "Method not allowed", nil)
			return
		}

		// Extract user ID from URL or request context
		// This example assumes the user ID is passed as a query parameter
		// Adjust according to your routing setup (e.g., using URL parameters)
		userIDStr := r.URL.Query().Get("id")
		if userIDStr == "" {
			writeJSONResponse(w, http.StatusBadRequest, false, "User ID is required", nil)
			return
		}

		// Convert userIDStr to integer
		var userID int
		_, err := fmt.Sscanf(userIDStr, "%d", &userID)
		if err != nil {
			writeJSONResponse(w, http.StatusBadRequest, false, "Invalid User ID", nil)
			log.Printf("Invalid User ID: %v", err)
			return
		}

		var updateReq UpdateUserRequest
		err = json.NewDecoder(r.Body).Decode(&updateReq)
		if err != nil {
			writeJSONResponse(w, http.StatusBadRequest, false, "Invalid request body", nil)
			log.Printf("Update decode error: %v", err)
			return
		}

		// validate update request fields
		if updateReq.Name != nil && len(*updateReq.Name) < 3 {
			writeJSONResponse(w, http.StatusBadRequest, false, "Name must be at least 3 characters long", nil)
			return
		}

		// Prepare fields to update
		setClauses := []string{}
		args := []interface{}{}

		if updateReq.Name != nil {
			setClauses = append(setClauses, "name = ?")
			args = append(args, *updateReq.Name)
		}

		if updateReq.Email != nil {
			setClauses = append(setClauses, "email = ?")
			args = append(args, *updateReq.Email)
		}

		if len(setClauses) == 0 {
			writeJSONResponse(w, http.StatusBadRequest, false, "No fields to update", nil)
			return
		}

		// Build the final SQL query
		query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", strings.Join(setClauses, ", "))
		args = append(args, userID)

		// Execute the update
		res, err := db.Exec(query, args...)
		if err != nil {
			writeJSONResponse(w, http.StatusInternalServerError, false, "Failed to update user", nil)
			log.Printf("Update user error: %v", err)
			return
		}

		// Check if any row was affected
		rowsAffected, err := res.RowsAffected()
		if err != nil {
			writeJSONResponse(w, http.StatusInternalServerError, false, "Failed to retrieve update result", nil)
			log.Printf("RowsAffected error: %v", err)
			return
		}

		if rowsAffected == 0 {
			writeJSONResponse(w, http.StatusNotFound, false, "User not found", nil)
			log.Printf("No user found with ID: %d", userID)
			return
		}

		// Success response
		writeJSONResponse(w, http.StatusOK, true, "User updated successfully", nil)
		log.Printf("User with ID %d updated successfully", userID)
	}
}
