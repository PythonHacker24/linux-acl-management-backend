package authentication

import (
	"time"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
	"github.com/golang-jwt/jwt/v5"
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
