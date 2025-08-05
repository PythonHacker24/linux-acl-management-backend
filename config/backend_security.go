package config

import (
	"errors"

	"github.com/MakeNowJust/heredoc"
)

/* backend security configs */
type BackendSecurity struct {
	JWTTokenSecret string `yaml:"jwt_secret_token,omitempty"`

	/* make this obselete */
	JWTExpiry      int    `yaml:"jwt_expiry,omitempty"`
}

/* normalization function */
func (b *BackendSecurity) Normalize() error {

	/* check if JWT token secret is specified */
	if b.JWTTokenSecret == "" {
		return errors.New(heredoc.Doc(`
			JWT Token Security is not specified in the configuration file. 

			Please check the docs for more information: 
		`))
	}

	/* set default JWT expiry to 24 hours */
	if b.JWTExpiry == 0 {
		b.JWTExpiry = 24
	}

	return nil
}
