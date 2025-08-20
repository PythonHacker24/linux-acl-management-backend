package search

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"
)

/* handler to return list of users that match the query in LDAP server */
func SearchUsersHandler(w http.ResponseWriter, r *http.Request) {
	/* fetch all users from LDAP server */
	query := r.URL.Query().Get("q")
    users, err := GetAllUsersFromLDAP(query)
    if err != nil {
		zap.L().Error("LDAP error",
			zap.Error(err),
		)
        http.Error(w, "LDAP error", http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}
