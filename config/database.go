package config

import (
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc"
)

/* database parameters */
type Database struct {
	TransactionLogRedis TransactionLogRedis `yaml:"transaction_log_redis,omitempty"`
}

/* transaction log redis parameters */
type TransactionLogRedis struct {
	Address  string `yaml:"address,omitempty"`
	Password string `yaml:"password,omitempty"`
	DB       int    `yaml:"db,omitempty"`
}

/* normalization function */
func (d *Database) Normalize() error {
	return d.TransactionLogRedis.Normalize()
}

func (r *TransactionLogRedis) Normalize() error {
	if r.Address == "" {
		return errors.New(heredoc.Doc(`
			Transaction Log Redis Address is not specified in the configuration file. 

			Please check the docs for more information: 
		`))
	}
	
	/* password can be empty */
	if r.Password == "" {
		/* just warn users to use password protected redis */
		fmt.Printf("Prefer using password for redis for security purposes\n\n")	
	}

	/* r.DB default value can be 0 */

	return nil
}
