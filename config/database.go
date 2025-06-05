package config

import (
	"errors"
	"fmt"

	"github.com/MakeNowJust/heredoc"
)

/* database parameters */
type Database struct {
	TransactionLogRedis TransactionLogRedis `yaml:"transaction_log_redis,omitempty"`
	ArchivalPQ			ArchivalPQ			`yaml:"archival_postgres,omitempty"`
}

/* transaction log redis parameters */
type TransactionLogRedis struct {
	Address  string `yaml:"address,omitempty"`
	Password string `yaml:"password,omitempty"`
	DB       int    `yaml:"db,omitempty"`
}

/* archival PostgreSQL parameters */
type ArchivalPQ struct {
	Host		string 	`yaml:"host,omitempty"` 
	Port		int		`yaml:"port,omitempty"`
	User		string 	`yaml:"user,omitempty"`
	Password	string 	`yaml:"password,omitempty"`
	DBName		string 	`yaml:"dbname,omitempty"`
	SSLMode		string 	`yaml:"sslmode,omitempty"`
}

/* normalization function for database */
func (d *Database) Normalize() error {
	/* check if Redis parameters are valid */
	err := d.TransactionLogRedis.Normalize()
	if err != nil {
		return err
	}

	/* check if PostgreSQL parameters are valid */
	err = d.ArchivalPQ.Normalize()
	if err != nil {
		return err 
	}

	return nil
}

/* transaction log redis normalization function */
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

/* archival PostgreSQL parameters */
func (a *ArchivalPQ) Normalize() error {

	/* return localhost if empty */
	if a.Host == "" {
        a.Host = "localhost"
    }

	/* return default port if empty */
    if a.Port == 0 {
        a.Port = 5432
    }

	/* username is mandatory */
    if a.User == "" {
        return errors.New("Database username is not set in the configuration.")
    }

	/* dbname is mandatory */
    if a.DBName == "" {
        return errors.New("Database name (dbname) is not set in the configuration.")
    }

	/* sslmode is disabled by default */
    if a.SSLMode == "" {
        a.SSLMode = "disable"
    }

	/* empty password but give a warning */
    if a.Password == "" {
        fmt.Printf("Warning: Connecting to PostgreSQL without a password. Consider using one for security.\n\n")
    }

    return nil
}
