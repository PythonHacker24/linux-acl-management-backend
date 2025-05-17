package ldap

import (
	"os"
	"crypto/tls"
	"fmt"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
	"github.com/go-ldap/ldap/v3"
	"go.uber.org/zap"
)

/* authenticate a user with ldap */
func AuthenticateUser(username, password, searchbase string) bool {

	/*
		authenticate is a critical functionality
		hence, it's implementation needs to be simplistic
		we only return true or false for authentication
		true is returned only if all the elimination steps are passed
		reducing unauthorized access in edge cases
	*/

	var l *ldap.Conn
	var err error
	ldapAddress := config.BackendConfig.Authentication.LDAPConfig.Address

	if config.BackendConfig.Authentication.LDAPConfig.TLS {
		l, err = ldap.DialURL(ldapAddress, ldap.DialWithTLSConfig(&tls.Config{

			/* true if using self-signed certs (not recommended) */
			InsecureSkipVerify: false,
		}))
	} else {
		l, err = ldap.DialURL(ldapAddress)
	}

	if err != nil {
		zap.L().Error("Failed to connect to LDAP Server",
			zap.Error(err),
		)
		return false
	}
	defer l.Close()

	/* securely fetch LDAP credentials from the environment */
	adminDN := os.Getenv("LDAP_ADMIN_DN")
	adminPassword := os.Getenv("LDAP_ADMIN_PASSWORD")

	/* authenticating with the ldap server with admin */
	err = l.Bind(adminDN, adminPassword)
	if err != nil {
		zap.L().Error("Admin authentication failed",
			zap.Error(err),
		)
		return false
	}

	/* creating a search request for ldap server */
	searchRequest := ldap.NewSearchRequest(
		searchbase,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,

		/* Searching by username */
		fmt.Sprintf("(uid=%s)", username),

		/* We only need the DN */
		[]string{"dn"},
		nil,
	)

	/* searching the ldap server for credentials */
	searchResult, err := l.Search(searchRequest)
	if err != nil {
		zap.L().Error("LDAP search failed",
			zap.Error(err),
		)
		return false
	}

	/* checking if search result is empty */
	if len(searchResult.Entries) == 0 {
		zap.L().Error("User not found in LDAP",
			zap.String("username", username),
			zap.Error(err),
		)
		return false
	}

	userDN := searchResult.Entries[0].DN

	/* checking if the user exists */
	err = l.Bind(userDN, password)
	if err != nil {
		zap.L().Error("User authentication failed",
			zap.String("Username", username),
			zap.Error(err),
		)
		return false
	}

	/* authentication successful */
	zap.L().Info("User authentication successful",
		zap.String("username", username),
	)

	return true
}
