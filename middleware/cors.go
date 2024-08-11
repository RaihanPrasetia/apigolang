package middleware

import (
	"fmt"
	"log"
	"loginApi/controllers" // Adjust the import path according to your project structure
	"net/http"

	"github.com/rs/cors"
)

func main() {
	// Create a new router or use your existing router setup
	mux := http.NewServeMux()

	// Define your routes
	mux.HandleFunc("/create/message", controllers.CreateMessage)

	// Setup CORS options
	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:3000"}, // Adjust as needed
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type"},
		AllowCredentials: true,
	})

	// Wrap your router with the CORS handler
	handler := corsHandler.Handler(mux)

	// Start the server
	fmt.Println("Server running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
