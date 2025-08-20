package health

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

/* health handler provides status check on the backend server */
func HealthHandler(w http.ResponseWriter, r *http.Request) {
	var response HealthResponse

	/* set the content type and write the response */
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	/* set the status to ok */
	response.Status = "ok"
	if err := json.NewEncoder(w).Encode(response); err != nil {
		zap.L().Error("Failed to send health response from the handler",
			zap.Error(err),
		)
	}
}
