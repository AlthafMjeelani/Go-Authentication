package routes

import (
	"database/sql"
	"encoding/json"
	"fmt"
	models "golang_projects/model"
	"log"
	"net/http"
	"regexp"
	"unicode"

	"golang.org/x/crypto/bcrypt"
)

// Response defines the structure for both success and error messages
type Response struct {
	Status  bool   `json:"status"`
	Message string `json:"message"`
}

// HandleRegister handles user registration
func HandleRegister(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeJSONResponse(w, http.StatusMethodNotAllowed, false, "Method not allowed")
			return
		}

		var user models.User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			writeJSONResponse(w, http.StatusBadRequest, false, "Invalid request body")
			log.Printf("Register decode error: %v", err)
			return
		}

		// Validate user input
		if err := validateUser(user); err != nil {
			writeJSONResponse(w, http.StatusBadRequest, false, err.Error())
			log.Printf("Validation error: %v", err)
			return
		}

		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
		if err != nil {
			writeJSONResponse(w, http.StatusInternalServerError, false, "Failed to hash password")
			log.Printf("Hash password error: %v", err)
			return
		}
		user.Password = string(hashedPassword)

		// Save user to the database
		err = models.CreateUser(db, user)
		if err != nil {
			writeJSONResponse(w, http.StatusInternalServerError, false, "Failed to register user")
			log.Printf("Register insert error: %v", err)
			return
		}

		// Success response
		writeJSONResponse(w, http.StatusCreated, true, "User registered successfully")
	}
}

// writeJSONResponse writes a structured JSON response
func writeJSONResponse(w http.ResponseWriter, statusCode int, status bool, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(Response{
		Status:  status,
		Message: message,
	})
}

// validateUser validates the user input
func validateUser(user models.User) error {
	if user.Name == "" {
		return fmt.Errorf("name is required")
	}
	if len(user.Name) < 3 {
		return fmt.Errorf("name must be at least 3 characters long")
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
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var credentials struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		err := json.NewDecoder(r.Body).Decode(&credentials)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			log.Printf("Login decode error: %v", err)
			return
		}

		// Retrieve user from the database
		user, err := models.GetUserByEmail(db, credentials.Email)
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			log.Printf("Login fetch error: %v", err)
			return
		}

		// Compare passwords
		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(credentials.Password))
		if err != nil {
			http.Error(w, "Invalid email or password", http.StatusUnauthorized)
			log.Printf("Login password mismatch: %v", err)
			return
		}

		w.Write([]byte("Login successful"))
	}
}
