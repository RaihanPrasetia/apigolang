package middleware

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"loginApi/utils" // Adjust the import path as necessary
)

func JWTAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the token from the Authorization header
		tokenString := extractToken(r)
		if tokenString == "" {
			http.Error(w, "Missing or invalid token", http.StatusUnauthorized)
			return
		}

		// Parse and validate the token
		token, claims, err := utils.ParseJWT(tokenString)
		if err != nil || !token.Valid {
			http.Error(w, "Invalid token", http.StatusUnauthorized)
			return
		}

		// Convert the Subject claim (userID) to an integer
		userIDStr := claims.Subject
		userID, err := strconv.Atoi(userIDStr)
		if err != nil {
			fmt.Printf("Error converting subject to integer: %v\n", err)
			http.Error(w, "Invalid token claims", http.StatusUnauthorized)
			return
		}

		// Add the userID to the request context
		ctx := context.WithValue(r.Context(), "userID", userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// Helper function to extract the token from the Authorization header
func extractToken(r *http.Request) string {
	bearerToken := r.Header.Get("Authorization")
	if len(bearerToken) > 7 && strings.ToUpper(bearerToken[0:7]) == "BEARER " {
		return bearerToken[7:]
	}
	return ""
}
