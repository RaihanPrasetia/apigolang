package main

import (
	"loginApi/database"
	"loginApi/routes"
	"net/http"
)

func main() {
	database.Connect()
	routes.RegisterRoutes()
	http.ListenAndServe(":8080", nil)
}
