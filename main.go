package main

import (
	"golang_projects/database"
	"golang_projects/routes"
	"log"
	"net/http"
)

func main() {
	// Initialize the database connection
	db := database.InitDB()
	defer db.Close()

	//Get all users
	http.HandleFunc("/users", routes.HandleUsers(db))
	// Register routes
	http.HandleFunc("/register", routes.HandleRegister(db))
	// Login route
	http.HandleFunc("/login", routes.HandleLogin(db))

	//update user
	http.HandleFunc("/update_user", routes.HandleUpdateUser(db))

	// Start the server
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
