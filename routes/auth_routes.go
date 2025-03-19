package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"golang_projects/model"
	"golang_projects/repository"
	utils "golang_projects/utility"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// HandleRegister handles user registration
func HandleRegister(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, false, "Method not allowed", nil)
			return
		}

		var user model.User

		// Determine Content-Type and parse accordingly
		contentType := r.Header.Get("Content-Type")
		if strings.Contains(contentType, "application/json") {
			// Parse JSON raw data
			err := json.NewDecoder(r.Body).Decode(&user)
			if err != nil {
				utils.WriteJSONResponse(w, http.StatusBadRequest, false, "Invalid JSON request body", nil)
				log.Printf("Register JSON decode error: %v", err)
				return
			}
		} else if strings.Contains(contentType, "application/x-www-form-urlencoded") || strings.Contains(contentType, "multipart/form-data") {
			// Parse Form data
			err := r.ParseMultipartForm(10 << 20) // Max memory 10 MB
			if err != nil {
				utils.WriteJSONResponse(w, http.StatusBadRequest, false, "Failed to parse form data", nil)
				log.Printf("Form data parse error: %v", err)
				return
			}

			// Assign form values to user
			user.Name = r.Form.Get("name")
			user.Email = r.Form.Get("email")
			user.Password = r.Form.Get("password")
			user.Phone = r.Form.Get("phone")
			user.Address = r.Form.Get("address")
		} else {
			utils.WriteJSONResponse(w, http.StatusUnsupportedMediaType, false, "Unsupported content type", nil)
			return
		}

		// Validate user input
		if err := validateUser(user, true); err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, false, err.Error(), nil)
			log.Printf("Validation error: %v", err)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, false, "Failed to hash password", nil)
			log.Printf("Hash password error: %v", err)
			return
		}
		user.Password = string(hashedPassword)

		// Save user to the database
		err = repository.CreateUser(db, user)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, false, err.Error(), nil)
			log.Printf("Register insert error: %v", err)
			return
		}

		// Success response
		utils.WriteJSONResponse(w, http.StatusCreated, true, "User registered successfully", nil)
	}
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
			utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, false, "Method not allowed", nil)
			return
		}

		var credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, false, "Invalid request body", nil)
			return
		}

		// Retrieve user from DB
		user, err := repository.GetUserLogin(db, credentials.Email)
		if err != nil {
			log.Printf("User not found: %v", err)
			utils.WriteJSONResponse(w, http.StatusUnauthorized, false, "Invalid email or password", nil)
			return
		}

		// Log the retrieved user details for debugging
		log.Printf("Retrieved User: %+v", user)

		// Check if the password matches
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
		if err != nil {
			log.Printf("Password mismatch: %v", err)
			utils.WriteJSONResponse(w, http.StatusUnauthorized, false, "Invalid email or password", nil)
			return
		}

		// Generate JWT token
		token, err := utils.GenerateJWT(user.ID)
		if err != nil {
			log.Printf("JWT generation error: %v", err)
			utils.WriteJSONResponse(w, http.StatusInternalServerError, false, "Failed to generate token", nil)
			return
		}

		// Success response
		response := struct {
			ID    int    `json:"id"`
			Name  string `json:"name"`
			Email string `json:"email"`
			Token string `json:"access_token"`
		}{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Token: token,
		}

		utils.WriteJSONResponse(w, http.StatusOK, true, "Login successful", response)
	}
}

// update User by ID
// HandleUpdateUser handles updating user fields
// HandleUpdateUser handles updating user fields using the repository pattern
func HandleUpdateUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut && r.Method != http.MethodPatch {
			utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, false, "Method not allowed", nil)
			return
		}

		// Get user ID from query parameters
		userIDStr := r.URL.Query().Get("id")
		if userIDStr == "" {
			utils.WriteJSONResponse(w, http.StatusBadRequest, false, "User ID is required", nil)
			return
		}

		// Convert userIDStr to integer
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, false, "Invalid User ID", nil)
			log.Printf("Invalid User ID: %v", err)
			return
		}

		var updateReq model.User
		err = json.NewDecoder(r.Body).Decode(&updateReq)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, false, "Invalid request body", nil)
			log.Printf("Update decode error: %v", err)
			return
		}

		// Create a map for fields to update
		updateFields := make(map[string]interface{})

		// Validate and add fields to update
		if updateReq.Name != "" && len(updateReq.Name) < 3 {
			utils.WriteJSONResponse(w, http.StatusBadRequest, false, "Name must be at least 3 characters long", nil)
			return
		}
		if updateReq.Name != "" {
			updateFields["name"] = updateReq.Name
		}

		if updateReq.Email != "" {
			updateFields["email"] = updateReq.Email
		}

		if updateReq.Phone != "" {
			updateFields["phone"] = updateReq.Phone
		}

		if updateReq.Address != "" {
			updateFields["address"] = updateReq.Address
		}

		// Check if password is being updated
		if updateReq.Password != "" {
			// Hash the new password
			hashedPassword, err := bcrypt.GenerateFromPassword([]byte(updateReq.Password), bcrypt.DefaultCost)
			if err != nil {
				utils.WriteJSONResponse(w, http.StatusInternalServerError, false, "Failed to hash password", nil)
				log.Printf("Hash password error: %v", err)
				return
			}
			updateFields["password"] = string(hashedPassword)
		}

		// If no fields to update, return an error
		if len(updateFields) == 0 {
			utils.WriteJSONResponse(w, http.StatusBadRequest, false, "No fields to update", nil)
			return
		}

		// Use repository to update the user
		rowsAffected, err := repository.UpdateUserByID(db, userID, updateFields)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, false, "Failed to update user", nil)
			log.Printf("Update user error: %v", err)
			return
		}

		if rowsAffected == 0 {
			utils.WriteJSONResponse(w, http.StatusNotFound, false, "User not found", nil)
			log.Printf("No user found with ID: %d", userID)
			return
		}

		// Success response
		utils.WriteJSONResponse(w, http.StatusOK, true, "User updated successfully", nil)
		log.Printf("User with ID %d updated successfully", userID)
	}
}

func HandleDeleteUser(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			utils.WriteJSONResponse(w, http.StatusMethodNotAllowed, false, "Method not allowed", nil)
			return
		}

		// Get user ID from query parameters
		userIDStr := r.URL.Query().Get("id")
		if userIDStr == "" {
			utils.WriteJSONResponse(w, http.StatusBadRequest, false, "User ID is required", nil)
			return
		}

		// Convert userIDStr to integer
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusBadRequest, false, "Invalid User ID", nil)
			log.Printf("Invalid User ID: %v", err)
			return
		}

		// Use repository to delete the user
		rowsAffected, err := repository.DeleteUserByID(db, userID)
		if err != nil {
			utils.WriteJSONResponse(w, http.StatusInternalServerError, false, "Failed to delete user", nil)
			log.Printf("Delete user error: %v", err)
			return
		}

		if rowsAffected == 0 {
			utils.WriteJSONResponse(w, http.StatusNotFound, false, "User not found", nil)
			log.Printf("No user found with ID: %d", userID)
			return
		}

		// Success response
		utils.WriteJSONResponse(w, http.StatusOK, true, "User deleted successfully", nil)
		log.Printf("User with ID %d deleted successfully", userID)
	}
}
