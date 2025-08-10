package auth

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/session"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/token"
)

/* handles user login and creates a session */
func LoginHandler(sessionManager *session.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		/* POST Request only - specified in routes */

		/* decode the request body */
		var user User
		err := json.NewDecoder(r.Body).Decode(&user)
		if err != nil {
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		/* check if username and password are specified */
		if user.Username == "" || user.Password == "" {
			http.Error(w, "Username and password are required", http.StatusBadRequest)
			return
		}

		/* authenticate the user */
		authStatus := AuthenticateUser(user.Username,
			user.Password,
			config.BackendConfig.Authentication.LDAPConfig.SearchBase,
		)

		/* check if authentication is successful */
		if !authStatus {
			zap.L().Warn("User with invalid credentials attempted to log in")
			http.Error(w, "Invalid credentials", http.StatusUnauthorized)
			return
		}

		/* 
			check if the session already exists in the manager.
			if it exists, refresh it's timer and return a jwt token
		*/

		/* create session for the user */
		sessionID, err := sessionManager.CreateSession(user.Username, r.RemoteAddr, r.UserAgent())
		if err != nil {
			zap.L().Error("Error creating session",
				zap.Error(err),
			)
			http.Error(w, "Error creating session", http.StatusInternalServerError)
			return
		}

		/* generate JWT for user interaction */
		token, err := token.GenerateJWT(user.Username, sessionID)
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
			zap.L().Error("Failed to encode response for login request",
				zap.Error(err),
			)
			http.Error(w, "Failed to encode response for login request", http.StatusInternalServerError)
			return
		}
	}
}

/* handles user logout and expire session */
func LogoutHandler(sessionManager *session.Manager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		
		/* authenticate the request through JWT */
		username, _, err := token.ExtractDataFromRequest(r)
		if err != nil {
			zap.L().Info("Error during token extraction in logout",
				zap.Error(err),
			)
			http.Error(w, "Error during token extraction in logout", http.StatusInternalServerError)
			return
		}

		err = sessionManager.ExpireSession(username)
		if err != nil {
			zap.L().Error("Failed to expire session during logout",
				zap.Error(err),
			)
			http.Error(w, "Failed to expire session during logout", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

/* validate a token */
func ValidateToken(w http.ResponseWriter, r *http.Request) {

	/* authenticate the request through JWT */
	_, _, err := token.ExtractDataFromRequest(r)
	if err != nil {
		zap.L().Info("Error during authentication",
			zap.Error(err),
		)
		http.Error(w, "Authentication Failed", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
