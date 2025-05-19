package config

import (
	"errors"

	"github.com/MakeNowJust/heredoc"
)

/* authentication parameters */
type Authentication struct {
	LDAPConfig		LDAPConfig 	`yaml:"ldap,omitempty"`
}

/* ldap authentication parameters */
type LDAPConfig struct {
	TLS             bool		`yaml:"tls,omitempty"`
	Address			string 		`yaml:"address,omitempty"`
	AdminDN 		string		`yaml:"admin_dn,omitempty"`
	AdminPassword	string		`yaml:"admin_password,omitempty"`
	SearchBase		string		`yaml:"search_base,omitempty"`
}

/* normalization function */
func (a *Authentication) Normalize() error {
	return a.LDAPConfig.Normalize()
}

func (l *LDAPConfig) Normalize() error {
	/* TLS will be false by default */

	if l.Address == "" {
		return errors.New(heredoc.Doc(`
			LDAP address is not specified in the configuration file. 

			Please check the docs for more information: 
		`))
	}

	if l.AdminDN == "" {
		return errors.New(heredoc.Doc(`
			LDAP admin DN is not specified in the configuration file. 

			Please check the docs for more information: 
		`))
	}

	if l.AdminPassword == "" {
		return errors.New(heredoc.Doc(`
			LDAP admin password is not specified in the configuration file. 

			Please check the docs for more information: 
		`))
	}

	if l.SearchBase == "" {
		return errors.New(heredoc.Doc(`
			LDAP search base is not specified in the configuration file. 

			Please check the docs for more information: 
		`))
	}
	
	return nil
}
