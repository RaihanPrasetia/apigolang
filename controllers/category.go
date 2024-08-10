package controllers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"loginApi/database"
	"loginApi/helpers"
	"loginApi/models"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func GetCategory(w http.ResponseWriter, r *http.Request) {
	query := "SELECT * FROM categories"
	rows, err := database.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch categories", http.StatusInternalServerError)
		fmt.Printf("Error executing query: %v\n", err)
		return
	}
	defer rows.Close()

	var categories []models.Category

	for rows.Next() {
		var category models.Category
		var createdAt []byte
		var updatedAt []byte

		err := rows.Scan(&category.ID, &category.Name, &createdAt, &updatedAt)
		if err != nil {
			http.Error(w, "Failed to scan category", http.StatusInternalServerError)
			fmt.Printf("Error scanning row: %v\n", err)
			return
		}

		category.Created_at, err = helpers.ParseDatetime(createdAt)
		if err != nil {
			http.Error(w, "Failed to parse created_at", http.StatusInternalServerError)
			fmt.Printf("Error parsing created_at: %v\n", err)
			return
		}

		category.Updated_at, err = helpers.ParseNullableDatetime(updatedAt)
		if err != nil {
			http.Error(w, "Failed to parse updated_at", http.StatusInternalServerError)
			fmt.Printf("Error parsing updated_at: %v\n", err)
			return
		}

		categories = append(categories, category)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error during row iteration", http.StatusInternalServerError)
		fmt.Printf("Error during row iteration: %v\n", err)
		return
	}

	response := map[string]interface{}{
		"categories": categories,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode categories to JSON", http.StatusInternalServerError)
		fmt.Printf("Error encoding JSON: %v\n", err)
		return
	}
}

func CreateCategory(w http.ResponseWriter, r *http.Request) {
	var category models.Category

	err := helpers.ParseJSONRequestBody(r, &category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	fmt.Printf("Decoded category: %+v\n", category)

	if category.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		fmt.Printf("Validation failed: Name=%s\n", category.Name)
		return
	}

	now := time.Now()
	category.Created_at = now
	category.Updated_at = sql.NullTime{Valid: false} // Set Updated_at to NULL

	query := "INSERT INTO categories (name, created_at, updated_at) VALUES (?, ?, ?)"
	result, err := database.DB.Exec(query, category.Name, category.Created_at, category.Updated_at)
	if err != nil {
		http.Error(w, "Failed to create Category", http.StatusInternalServerError)
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Category created successfully with ID: %d", id)
}

func UpdateCategory(w http.ResponseWriter, r *http.Request) {
	// Extract the ID from the URL
	idStr := strings.TrimPrefix(r.URL.Path, "/update/categories/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid category ID", http.StatusBadRequest)
		fmt.Printf("Invalid ID: %v\n", idStr)
		return
	}

	var category models.Category

	// Parse JSON body
	err = helpers.ParseJSONRequestBody(r, &category)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// Validate input
	if category.Name == "" {
		http.Error(w, "Name is required", http.StatusBadRequest)
		fmt.Printf("Validation failed: Name=%s\n", category.Name)
		return
	}

	// Set updated_at to current time
	now := time.Now()
	category.Updated_at = sql.NullTime{Time: now, Valid: true}

	// Query to update category
	query := "UPDATE categories SET name = ?, updated_at = ? WHERE id = ?"
	result, err := database.DB.Exec(query, category.Name, category.Updated_at, id)
	if err != nil {
		http.Error(w, "Failed to update category", http.StatusInternalServerError)
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		http.Error(w, "Failed to check affected rows", http.StatusInternalServerError)
		fmt.Printf("Error checking affected rows: %v\n", err)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Category not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Category updated successfully")
}
