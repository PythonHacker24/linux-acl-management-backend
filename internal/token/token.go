package token

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
)

/* generating jwt token for user identification with specified configs */
func GenerateJWT(username string, sessionID uuid.UUID) (string, error) {
	expiryHours := config.BackendConfig.BackendSecurity.JWTExpiry

	/* generate JWT token with claims */
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"sessionID": sessionID, 
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

/* extracts username and sessionID from JWT token */
func GetDataFromJWT(tokenString string) (string, string, error) {
	/* get claims from JWT Token */
	claims, err := ValidateJWT(tokenString)
	if err != nil {
		return "", "", fmt.Errorf("JWT validation error: %w", err)
	}

	/* extract username from JWT Token */
	username, ok := claims["username"].(string)
	if !ok {
		return "", "", fmt.Errorf("username not found in token")
	}

	/* extract sessionID from JWT Token */
	sessionID, ok := claims["sessionID"].(string)
	if !ok {
		return "", "", fmt.Errorf("sessionID not found in token")
	}
	return username, sessionID, nil
}

/* extract username and sessionID from http request (wrapper around GetUsernameFromJWT for http requests) */
func ExtractDataFromRequest(r *http.Request) (string, string, error) {
	/* get the authorization header */
	authHeader := r.Header.Get("Authorization")
	if authHeader == "" {
		return "", "", fmt.Errorf("authorization header not found")
	}

	/* check if the header is in the correct format */
	parts := strings.Split(authHeader, " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return "","",  fmt.Errorf("invalid authorization header format")
	}

	/* extract username and sessionID from JWT token */
	return GetDataFromJWT(parts[1])
}
