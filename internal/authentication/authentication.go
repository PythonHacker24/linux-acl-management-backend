package authentication

import (
	"fmt"
	"time"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/PythonHacker24/linux-acl-management-backend/config"
)

/* generating jwt token for user identification with specified configs */
func GenerateJWT(username string) (string, error) {
	expiryHours := config.BackendConfig.BackendSecurity.JWTExpiry
	if expiryHours == 0 {
		expiryHours = 24
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * time.Duration(expiryHours)).Unix(),
	})

	return token.SignedString([]byte(config.EnvConfig.JWTSecret))
}

/* validate JWT token and return claims */
func ValidateJWT(tokenString string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return config.EnvConfig.JWTSecret, nil
	})

	if err != nil {
		return nil, err
	}

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
        return "", fmt.Errorf("invalid token: %v", err)
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

	/* extract authentication hearder from http request */
    authHeader := r.Header.Get("Authorization")
    if authHeader == "" {
        return "", fmt.Errorf("missing Authorization header")
    }
	
	/* parse the token from the header */
    tokenParts := strings.Split(authHeader, " ")
    if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
        return "", fmt.Errorf("invalid Authorization header format")
    }

    return GetUsernameFromJWT(tokenParts[1])
}
