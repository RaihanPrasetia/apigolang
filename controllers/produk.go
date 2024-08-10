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

func CreateProduct(w http.ResponseWriter, r *http.Request) {
	// Get the userID from the context
	userID := r.Context().Value("userID").(int)

	var product models.Product

	// Decode the request body
	err := helpers.ParseJSONRequestBody(r, &product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	// Validate input
	if product.Name == "" || product.Price <= 0 || product.Category_id <= 0 {
		http.Error(w, "All fields are required and must be valid", http.StatusBadRequest)
		fmt.Printf("Validation failed: Name=%s, Price=%d, Category_id=%d\n", product.Name, product.Price, product.Category_id)
		return
	}

	// Check if category_id exists in categories table
	var exists bool
	query := "SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?)"
	err = database.DB.QueryRow(query, product.Category_id).Scan(&exists)
	if err != nil {
		http.Error(w, "Failed to check category existence", http.StatusInternalServerError)
		fmt.Printf("Error checking category existence: %v\n", err)
		return
	}

	if !exists {
		http.Error(w, "Category not found", http.StatusBadRequest)
		fmt.Printf("Validation failed: Category_id=%d not found\n", product.Category_id)
		return
	}

	// Set creation and update times
	now := time.Now()
	product.Created_at = now
	product.Updated_at = sql.NullTime{Valid: false} // Set Updated_at to NULL
	product.User_id = userID

	// SQL query to insert product
	query = "INSERT INTO products (name, price, user_id, category_id, created_at, updated_at) VALUES (?, ?, ?, ?, ?, ?)"
	result, err := database.DB.Exec(query, product.Name, product.Price, product.User_id, product.Category_id, product.Created_at, product.Updated_at)
	if err != nil {
		http.Error(w, "Failed to create product", http.StatusInternalServerError)
		fmt.Printf("Error executing query: %v\n", err)
		return
	}

	// Retrieve the last insert ID
	id, err := result.LastInsertId()
	if err != nil {
		http.Error(w, "Failed to retrieve last insert ID", http.StatusInternalServerError)
		return
	}

	// Respond with success
	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Product created successfully with ID: %d", id)
}

func GetProduct(w http.ResponseWriter, r *http.Request) {
	query := "SELECT id, name, price, category_id, created_at, updated_at FROM products"
	rows, err := database.DB.Query(query)
	if err != nil {
		http.Error(w, "Failed to fetch products", http.StatusInternalServerError)
		fmt.Printf("Error executing query: %v\n", err)
		return
	}
	defer rows.Close()

	var products []models.Product

	for rows.Next() {
		var product models.Product
		var createdAt []byte
		var updatedAt []byte

		err := rows.Scan(&product.ID, &product.Name, &product.Price, &product.Category_id, &createdAt, &updatedAt)
		if err != nil {
			http.Error(w, "Failed to scan product", http.StatusInternalServerError)
			fmt.Printf("Error scanning row: %v\n", err)
			return
		}

		// Parse created_at
		product.Created_at, err = helpers.ParseDatetime(createdAt)
		if err != nil {
			http.Error(w, "Failed to parse created_at", http.StatusInternalServerError)
			fmt.Printf("Error parsing created_at: %v\n", err)
			return
		}

		// Parse updated_at
		product.Updated_at, err = helpers.ParseNullableDatetime(updatedAt)
		if err != nil {
			http.Error(w, "Failed to parse updated_at", http.StatusInternalServerError)
			fmt.Printf("Error parsing updated_at: %v\n", err)
			return
		}

		products = append(products, product)
	}

	if err := rows.Err(); err != nil {
		http.Error(w, "Error during row iteration", http.StatusInternalServerError)
		fmt.Printf("Error during row iteration: %v\n", err)
		return
	}

	// Wrap products in a map
	response := map[string]interface{}{
		"products": products,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		http.Error(w, "Failed to encode products to JSON", http.StatusInternalServerError)
		fmt.Printf("Error encoding JSON: %v\n", err)
		return
	}
}

func UpdateProduct(w http.ResponseWriter, r *http.Request) {
	// Extract userID from the request context (set by middleware)
	userID := r.Context().Value("userID").(int)

	// Extract the product ID from the URL
	idStr := strings.TrimPrefix(r.URL.Path, "/update/products/")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "Invalid product ID", http.StatusBadRequest)
		fmt.Printf("Invalid ID: %v\n", idStr)
		return
	}

	// Check if the product belongs to the user
	var existingUserID int
	err = database.DB.QueryRow("SELECT user_id FROM products WHERE id = ?", id).Scan(&existingUserID)
	if err == sql.ErrNoRows {
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Failed to fetch product", http.StatusInternalServerError)
		fmt.Printf("Error fetching product: %v\n", err)
		return
	}

	if existingUserID != userID {
		http.Error(w, "Unauthorized: You do not own this product", http.StatusUnauthorized)
		return
	}

	var product models.Product
	err = helpers.ParseJSONRequestBody(r, &product)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		fmt.Printf("Error parsing JSON: %v\n", err)
		return
	}

	if product.Category_id > 0 {
		var categoryExists bool
		err := database.DB.QueryRow("SELECT EXISTS(SELECT 1 FROM categories WHERE id = ?)", product.Category_id).Scan(&categoryExists)
		if err != nil {
			http.Error(w, "Failed to validate category ID", http.StatusInternalServerError)
			fmt.Printf("Error checking category ID: %v\n", err)
			return
		}
		if !categoryExists {
			http.Error(w, "Category not found", http.StatusBadRequest)
			return
		}
	}

	// Prepare to build the query dynamically
	var setClauses []string
	var args []interface{}

	if product.Name != "" {
		setClauses = append(setClauses, "name = ?")
		args = append(args, product.Name)
	}

	if product.Price > 0 {
		setClauses = append(setClauses, "price = ?")
		args = append(args, product.Price)
	}

	if product.Category_id > 0 {
		setClauses = append(setClauses, "category_id = ?")
		args = append(args, product.Category_id)
	}

	// Set Updated_at to the current time if any fields are updated
	if len(setClauses) > 0 {
		now := time.Now()
		setClauses = append(setClauses, "updated_at = ?")
		args = append(args, sql.NullTime{Time: now, Valid: true})
	}

	if len(setClauses) == 0 {
		http.Error(w, "No fields to update", http.StatusBadRequest)
		return
	}

	// Build the final query
	query := fmt.Sprintf("UPDATE products SET %s WHERE id = ?", strings.Join(setClauses, ", "))
	args = append(args, id)

	// Execute the query
	result, err := database.DB.Exec(query, args...)
	if err != nil {
		http.Error(w, "Failed to update product", http.StatusInternalServerError)
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
		http.Error(w, "Product not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Product updated successfully")
}
