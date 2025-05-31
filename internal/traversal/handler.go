package traversal

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/auth"
)

/*
	user considers / to be the root of the file path
	the backend transalates / to basepath/ securely
	this translation needs to be done wherever necessary
*/

/* POST handler for listing files in given directory */
func ListFilesInDirectory(w http.ResponseWriter, r *http.Request) {

	/* extracting userID from request */
	userID, err := auth.ExtractUsernameFromRequest(r)
	if err != nil {
		zap.L().Error("Error during getting username in HandleListFiles handler",
			zap.Error(err),
		)
		return
	}

	/* check if the request body is valid */
	var listRequest ListRequest
	err = json.NewDecoder(r.Body).Decode(&listRequest)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	/* list all the files in given filepath */
	entries, err := ListFiles(listRequest.FilePath, userID)
	if err != nil {
		zap.L().Warn("File listing error",
			zap.Error(err),
		)
		http.Error(w, "Failed to list files", http.StatusInternalServerError)
	}

	/* send the response back */
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(entries); err != nil {
		zap.L().Error("Failed to encode response for listing request",
			zap.Error(err),
		)
		http.Error(w, "Failed to encode response for listing request", http.StatusInternalServerError)
		return
	}
}
