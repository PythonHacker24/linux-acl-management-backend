package token

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
)

/* generating jwt token for user identification with specified configs */
func GenerateJWT(username string) (string, error) {
	expiryHours := config.BackendConfig.BackendSecurity.JWTExpiry

	/* generate JWT token with claims */
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * time.Duration(expiryHours)).Unix(),
	})

	return token.SignedString([]byte(config.BackendConfig.BackendSecurity.JWTTokenSecret))
}

/* validate JWT token and return claims */
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	/* parse the token */
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return []byte(config.BackendConfig.BackendSecurity.JWTTokenSecret), nil
	})

	/* check if token is valid */
	if err != nil {
		return nil, fmt.Errorf("JWT parsing error: %w", err)
	}

	/* check if token is valid */
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

/* extracts username from JWT token */
func GetUsernameFromJWT(tokenString string) (string, error) {
	/* get claims from JWT Token */
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return "", fmt.Errorf("JWT validation error: %w", err)
	}

	/* extract username from JWT Token */
	username, ok := claims["username"].(string)
	if !ok {
		return "", fmt.Errorf("username not found in token")
	}

	return username, nil
}

/* extract username from http request (wrapper around GetUsernameFromJWT for http requests) */
func ExtractUsernameFromRequest(r *http.Request) (string, error) {
	/* get the authorization header */
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", fmt.Errorf("authorization header not found")
	}

	/* check if the header is in the correct format */
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "", fmt.Errorf("invalid authorization header format")
	}

	/* extract username from JWT token */
	return GetUsernameFromJWT(parts[1])
}
