package controllers

import (
	"fmt"
	"loginApi/database"
	"loginApi/helpers"
	"loginApi/models"
	"net/http"
	"time"
)

func CreateMessage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "http://localhost:3000")
	w.Header().Set("Access-Control-Allow-Methods", "POST")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		return // Handle preflight requests
	}
	var message models.Message

	err := helpers.ParseJSONRequestBody(r, &message)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	fmt.Printf("Decoded message: %+v\n", message)

	if message.Name == "" || message.Email == "" || message.PhoneNumber == "" || message.Subject == "" || message.Message == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		fmt.Printf("Validation failed: Name=%s, Email=%s, PhoneNumber=%s, Subject=%s, Message=%s\n", message.Name, message.Email, message.PhoneNumber, message.Subject, message.Message)
		return
	}

	now := time.Now()
	message.Created_at = now
	message.Updated_at = now

	query := "INSERT INTO messages (name, email, phone_number, subject, message, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?, ?)"
	result, err := database.DB.Exec(query, message.Name, message.Email, message.PhoneNumber, message.Subject, message.Message, message.Created_at, message.Updated_at)
	if err != nil {
		http.Error(w, "Failed to create Message", http.StatusInternalServerError)
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Message created successfully with ID: %d", id)
}
