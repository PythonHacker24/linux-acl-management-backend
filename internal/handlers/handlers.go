package handlers

import (
	"net/http"
	"encoding/json"

	"go.uber.org/zap"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/models"
)

/* health handler provides status check on the backend server */
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	var response models.HealthResponse
    
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(http.StatusOK)

    response.Status = "ok"
    if err := json.NewEncoder(w).Encode(response); err != nil {
        zap.L().Error("Failed to send health response from the handler",
			zap.Error(err),
		)
    }
}

/* allows users to create a session */
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user models.User

	/* decode json response*/
	err := json.NewDecoder(r.Body).Decode(&user)
    if err != nil {
		zap.L().Warn("A request with invalid body recieved")
        http.Error(w, "Invalid request body", http.StatusBadRequest)
        return
    }

	/* check if username and password exists */
	if user.Username == "" || user.Password == "" {
		zap.L().Warn("A request with no username or password recieved")
        http.Error(w, "Username and password are required", http.StatusBadRequest)
        return
    }

	/* authenticate the user with ldap */
	authStatus := ldap.AuthenticateUser()
    if !authStatus {
		zap.L().Warn("A request with invalid credentials recieved")
        http.Error(w, "Invalid credentials", http.StatusUnauthorized)
        return
    }

	/* create a session if user exists */
	sessionmanager.CreateSession(user.Username)

	/* create a JWT token for the user */
	token, err := authentication.GenerateJWT(user.Username)
    if err != nil {
        zap.L().Error("Error generating token", 
			zap.Error(err),
		)
        http.Error(w, "Error generating token", http.StatusInternalServerError)
        return
    }

	/* send response with JWT token to the user */
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
