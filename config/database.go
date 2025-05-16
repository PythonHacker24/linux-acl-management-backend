package config

import (
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc"
)

/* database parameters */
type Database struct {
	TransactionLogRedis TransactionLogRedis `yaml:"transaction_logs_redis"`
}

/* transaction log redis parameters */
type TransactionLogRedis struct {
	Address  string `yaml:"address"`
	Password string `yaml:"password"`
	DB       string `yaml:"db"`
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
		fmt.Println("Prefer using password for redis for security purposes")	
	}

	if r.DB == "" {
		r.DB = "0"
	}

	return nil
}
