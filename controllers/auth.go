package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"loginApi/database"
	"loginApi/helpers"
	"loginApi/models"
	"loginApi/utils"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

func Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	err := helpers.ParseJSONRequestBody(r, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// Tambahkan log untuk nilai user setelah decoding
	fmt.Printf("Decoded user: %+v\n", user)

	if user.Name == "" || user.Email == "" || user.PhoneNumber == "" || user.Password == "" {
		http.Error(w, "All fields are required", http.StatusBadRequest)
		fmt.Printf("Validation failed: Name=%s, Email=%s, PhoneNumber=%s, Password=%s\n", user.Name, user.Email, user.PhoneNumber, user.Password)
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		http.Error(w, "Failed to hash password", http.StatusInternalServerError)
		fmt.Printf("Error hashing password: %v\n", err) // Tambahkan log error
		return
	}

	// Simpan user ke database
	query := "INSERT INTO users (name, email, phone_number, password) VALUES (?,?, ?, ?, ?)"
	result, err := database.DB.Exec(query, user.Name, user.Email, user.PhoneNumber, hashedPassword)
	if err != nil {
		http.Error(w, "Failed to register user", http.StatusInternalServerError)
		fmt.Printf("Error executing query: %v\n", err) // Tambahkan log error
		return
	}

	// Ambil ID user yang baru saja didaftarkan
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "User registered successfully with ID: %d", id)
}

func Login(w http.ResponseWriter, r *http.Request) {
	var user models.User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	if user.Email == "" || user.Password == "" {
		http.Error(w, "Email and Password are required", http.StatusBadRequest)
		return
	}

	var dbUser models.User
	query := "SELECT id, name, email, phone_number, password FROM users WHERE email = ?"
	row := database.DB.QueryRow(query, user.Email)
	err = row.Scan(&dbUser.ID, &dbUser.Name, &dbUser.Email, &dbUser.PhoneNumber, &dbUser.Password)
	if err == sql.ErrNoRows {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch user", http.StatusInternalServerError)
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(dbUser.Password), []byte(user.Password))
	if err != nil {
		http.Error(w, "Invalid email or password", http.StatusUnauthorized)
		return
	}

	// Generate JWT
	tokenString, err := utils.GenerateJWT(dbUser.ID, dbUser.Name)
	if err != nil {
		http.Error(w, "Failed to generate token", http.StatusInternalServerError)
		return
	}

	// Create the response with user details and JWT
	response := map[string]interface{}{"user": models.LoginResponse{
		ID:          dbUser.ID,
		Name:        dbUser.Name,
		Email:       dbUser.Email,
		PhoneNumber: dbUser.PhoneNumber,
		Token:       tokenString,
	}}

	// Send the response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		fmt.Printf("Error encoding response: %v\n", err)
		return
	}
}
