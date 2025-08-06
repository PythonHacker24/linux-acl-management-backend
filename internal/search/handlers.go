package search

import (
	"fmt"
	"net/http"
	"encoding/json"
)

/* handler to return list of all users in LDAP server */
func SearchUsersHandler(w http.ResponseWriter, r *http.Request) {
	/* fetch all users from LDAP server */
    users, err := GetAllUsersFromLDAP()
    if err != nil {
        http.Error(w, fmt.Sprintf("LDAP error: %v", err), http.StatusInternalServerError)
        return
    }

    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(users)
}
