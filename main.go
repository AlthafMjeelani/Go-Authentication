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

	// Setup router
	router := routes.SetupRoutes(db)

	// Start the server
	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
