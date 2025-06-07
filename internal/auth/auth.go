package auth

import (
	"net/http"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/token"
)

/* extract username from http request */
func ExtractUsernameFromRequest(r *http.Request) (string, error) {
	return token.ExtractUsernameFromRequest(r)
}
