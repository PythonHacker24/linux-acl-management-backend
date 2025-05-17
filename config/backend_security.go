package config

import (
	"errors"

	"github.com/MakeNowJust/heredoc"
)

/* backend security configs */
type BackendSecurity struct {
	JWTTokenSecret 	string 	`yaml:"jwt_secret_token"`
	JWTExpiry 		int 	`yaml:"jwt_expiry"`
}

/* normalization function */
func (b *BackendSecurity) Normalize() error {
	if b.JWTTokenSecret == "" {
		return errors.New(heredoc.Doc(`
			JWT Token Security is not specified in the configuration file. 

			Please check the docs for more information: 
		`)) 
	}

	if b.JWTExpiry == 0 {
		b.JWTExpiry = 1
	}

	return nil
}
