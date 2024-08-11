package routes

import (
	// Adjust the import path as necessary
	"loginApi/controllers"
	"loginApi/middleware"
	"net/http"
)

func RegisterRoutes() {
	// Auth
	http.HandleFunc("/register", controllers.Register)
	http.HandleFunc("/login", controllers.Login)

	// Products
	http.HandleFunc("/products", controllers.GetProduct)
	http.Handle("/create/product", middleware.JWTAuth(http.HandlerFunc(controllers.CreateProduct)))
	http.Handle("/update/products/", middleware.JWTAuth(http.HandlerFunc(controllers.UpdateProduct)))

	// Categories
	http.HandleFunc("/categories", controllers.GetCategory)
	http.Handle("/create/categories", middleware.JWTAuth(http.HandlerFunc(controllers.CreateCategory)))
	http.Handle("/update/categories/", middleware.JWTAuth(http.HandlerFunc(controllers.UpdateCategory)))

	http.HandleFunc("/create/message", controllers.CreateMessage)

}
