package config

import (
	"errors"

	"github.com/MakeNowJust/heredoc"
)

/* database parameters */
type Database struct {
	TransactionLogRedis TransactionLogRedis `yaml:"transaction_log_redis"`
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

	if r.DB == "" {
		r.DB = "0"
	}

	return nil
}
