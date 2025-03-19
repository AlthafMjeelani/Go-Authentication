package model

// User represents a user in the database
type User struct {
	ID       int    `json:"id" db:"id"`
	Name     string `json:"name" db:"name" validate:"required,min=3"`
	Email    string `json:"email" db:"email" validate:"required,email"`
	Password string `json:"password,omitempty" db:"password" validate:"required,min=6"`
	Phone    string `json:"phone" db:"phone" validate:"min=10"`
	Address  string `json:"address" db:"address" validate:"min=5"`
}
