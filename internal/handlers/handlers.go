package handlers

import (
	"net/http"
	"encoding/json"

	"go.uber.org/zap"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/models"
)

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

func LoginHandler(w http.ResponseWriter, r *http.Request) {

}
