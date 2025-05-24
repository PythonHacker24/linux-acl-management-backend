package config

import (
	"errors"

	"github.com/MakeNowJust/heredoc"
)

/* authentication parameters */
type Authentication struct {
	LDAPConfig LDAPConfig `yaml:"ldap,omitempty"`
}

/* ldap authentication parameters */
type LDAPConfig struct {
	TLS           bool   `yaml:"tls,omitempty"`
	Address       string `yaml:"address,omitempty"`
	AdminDN       string `yaml:"admin_dn,omitempty"`
	AdminPassword string `yaml:"admin_password,omitempty"`
	SearchBase    string `yaml:"search_base,omitempty"`
}

/* normalization function */
func (a *Authentication) Normalize() error {
	return a.LDAPConfig.Normalize()
}

/* ldap authentication normalization function */
func (l *LDAPConfig) Normalize() error {
	/* TLS will be false by default */

	/* check if address is specified */
	if l.Address == "" {
		return errors.New(heredoc.Doc(`
			LDAP address is not specified in the configuration file. 

			Please check the docs for more information: 
		`))
	}

	/* check if admin DN is specified */
	if l.AdminDN == "" {
		return errors.New(heredoc.Doc(`
			LDAP admin DN is not specified in the configuration file. 

			Please check the docs for more information: 
		`))
	}

	/* check if admin password is specified */
	if l.AdminPassword == "" {
		return errors.New(heredoc.Doc(`
			LDAP admin password is not specified in the configuration file. 

			Please check the docs for more information: 
		`))
	}

	/* check if search base is specified */
	if l.SearchBase == "" {
		return errors.New(heredoc.Doc(`
			LDAP search base is not specified in the configuration file. 

			Please check the docs for more information: 
		`))
	}

	return nil
}
