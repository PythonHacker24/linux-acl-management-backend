package search

import (
	"crypto/tls"

	"github.com/PythonHacker24/linux-acl-management-backend/config"
	"github.com/go-ldap/ldap/v3"
)

/*
	TODO: Blacklisting
	This needs to be done when admin panel is created.
	Users will be able to add users to blacklist which shouldn't be mentioned to the users.
*/

/* returns all users in LDAP server */
func GetAllUsersFromLDAP() ([]User, error) {
	
	var l *ldap.Conn
	var err error
	ldapAddress := config.BackendConfig.Authentication.LDAPConfig.Address

	/* check if TLS is enabled */
	if config.BackendConfig.Authentication.LDAPConfig.TLS {
		l, err = ldap.DialURL(ldapAddress, ldap.DialWithTLSConfig(&tls.Config{

			/* true if using self-signed certs (not recommended) */
			InsecureSkipVerify: true,
		}))
	} else {
		l, err = ldap.DialURL(ldapAddress)
	}

	if err != nil {
		return nil, err
	}
	defer l.Close()

	/* authenticating with the ldap server with admin */
	err = l.Bind(config.BackendConfig.Authentication.LDAPConfig.AdminDN,
		config.BackendConfig.Authentication.LDAPConfig.AdminPassword,
	)
	if err != nil {
		return nil, err
	}

	/* search for users */
    searchRequest := ldap.NewSearchRequest(
		/* Base DN */
		config.BackendConfig.Authentication.LDAPConfig.AdminDN,
        ldap.ScopeWholeSubtree,
        ldap.NeverDerefAliases,
		/* size limit */
        0,     
		/* time limit */
        0,     
		/* types only */
        false, 
		/* filter */
        "(objectClass=person)",
		/* attributes to retrieve */
        []string{"cn", "mail", "sAMAccountName"}, // 
        nil,
    )

	/* search for request in LDAP Server */
    sr, err := l.Search(searchRequest)
    if err != nil {
        return nil, err
    }

    users := []User{}
    for _, entry := range sr.Entries {
        users = append(users, User{
            CN:       entry.GetAttributeValue("cn"),
            Mail:     entry.GetAttributeValue("mail"),
            Username: entry.GetAttributeValue("sAMAccountName"),
        })
    }

    return users, nil
}
