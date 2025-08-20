package auth

import (
	"net/http"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/token"
)

/* extract username from http request */
func ExtractDataFromRequest(r *http.Request) (string, string, error) {
	return token.ExtractDataFromRequest(r)
}
