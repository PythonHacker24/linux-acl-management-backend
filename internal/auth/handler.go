package auth

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/session"
)

/* Handles user login and creates a session */
func LoginHandler(w http.ResponseWriter, r *http.Request) {

	/* POST Request only - specified in routes */

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if user.Username == "" || user.Password == "" {
		http.Error(w, "Username and password are required", http.StatusBadRequest)
		return
	}

	/* authenticate the user */
	authStatus := AuthenticateUser(user.Username,
		user.Password,
		config.BackendConfig.Authentication.LDAPConfig.SearchBase,
	)
	if !authStatus {
		zap.L().Warn("User with invalid credentials attempted to log in")
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	/* after building session manager */
	session.CreateSession(user.Username)

	/* generate JWT for user interaction */
	token, err := GenerateJWT(user.Username)
	if err != nil {
		zap.L().Error("Error generating token",
			zap.Error(err),
		)
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	/* create auth successful response */
	response := map[string]string{"token": token}
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		zap.L().Error("Failed to encode response",
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
