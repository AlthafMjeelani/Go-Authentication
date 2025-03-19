package repository

import (
	"database/sql"
	"fmt"
	model "golang_projects/model"
	"strings"

	"github.com/mattn/go-sqlite3"
)

// CreateUser adds a new user to the database
func CreateUser(db *sql.DB, user model.User) error {
	query := `INSERT INTO users (name, email, password, phone, address) 
	          VALUES (?, ?, ?, ?, ?)`
	_, err := db.Exec(query, user.Name, user.Email, user.Password, user.Phone, user.Address)

	if err != nil {
		if sqliteErr, ok := err.(sqlite3.Error); ok && sqliteErr.Code == sqlite3.ErrConstraint {
			return fmt.Errorf("email already exists")
		}
		return err
	}
	return nil
}

// GetUserByEmail retrieves a user by email
func GetUserByEmail(db *sql.DB, email string) (model.User, error) {
	var user model.User
	err := db.QueryRow("SELECT id, name, email, phone, address FROM users WHERE email = ?", email).Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Address)
	return user, err
}

func GetUserLogin(db *sql.DB, email string) (model.User, error) {
	var user model.User
	err := db.QueryRow("SELECT id, name, email, password, phone, address FROM users WHERE email = ?", email).Scan(&user.ID, &user.Name, &user.Email, &user.Password, &user.Phone, &user.Address)
	return user, err
}

// GetAllUsers retrieves all users from the database
func GetAllUsers(db *sql.DB) ([]model.User, error) {
	rows, err := db.Query("SELECT id, name, email, phone, address FROM users")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []model.User
	for rows.Next() {
		var user model.User
		err := rows.Scan(&user.ID, &user.Name, &user.Email, &user.Phone, &user.Address)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// UpdateUserByID updates the user fields in the database based on the provided map
func UpdateUserByID(db *sql.DB, userID int, updateFields map[string]interface{}) (int64, error) {
	setClauses := []string{}
	args := []interface{}{}

	for field, value := range updateFields {
		setClauses = append(setClauses, fmt.Sprintf("%s = ?", field))
		args = append(args, value)
	}
	args = append(args, userID)

	query := fmt.Sprintf("UPDATE users SET %s WHERE id = ?", strings.Join(setClauses, ", "))

	res, err := db.Exec(query, args...)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}

func DeleteUserByID(db *sql.DB, userID int) (int64, error) {
	res, err := db.Exec("DELETE FROM users WHERE id = ?", userID)
	if err != nil {
		return 0, err
	}

	return res.RowsAffected()
}
