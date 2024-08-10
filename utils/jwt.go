package utils

import (
	"errors"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("your_secret_key") // Change this to a secure key

// GenerateJWT generates a JWT token
func GenerateJWT(userID int, userName string) (string, error) {
	claims := &jwt.StandardClaims{
		Subject:   strconv.Itoa(userID),
		Issuer:    userName,
		ExpiresAt: time.Now().Add(24 * time.Hour * 30).Unix(), // Token expires in 24 hours
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// ParseJWT parses and validates a JWT token
func ParseJWT(tokenString string) (*jwt.Token, *jwt.StandardClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwt.StandardClaims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		return nil, nil, err
	}

	if claims, ok := token.Claims.(*jwt.StandardClaims); ok && token.Valid {
		return token, claims, nil
	}

	return nil, nil, errors.New("invalid token")
}
