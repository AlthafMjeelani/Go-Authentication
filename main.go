package main

import (
	"log"
	"net/http"
	"golang_projects/database"
	"golang_projects/routes"
)

func main() {
	// Initialize the database connection
	db := database.InitDB()
	defer db.Close()

	// Register routes
	http.HandleFunc("/users", routes.HandleUsers(db))
	http.HandleFunc("/register", routes.HandleRegister(db))
	http.HandleFunc("/login", routes.HandleLogin(db))

	// Start the server
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
