package postgresql

import (
	"database/sql"
	"fmt"
)

/*
	PostgreSQL database stores logs of all expired sessions and processed transactions results
	this is an archival database and is used to serve archived information to users and admin panel
*/

/* 
	TODO: when application shutsdown, all sessions in Redis are marked expired
	make sure to push those expired sessions to PostreSQL database for archival
*/

/* PostgreSQL client */
type PGClient struct {
	db *sql.DB
}

/* create new PostgreSQL client */
func NewPGClient(connStr string) (*PGClient, error) {

	/* initiate PostgreSQL connection */
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open PostgreSQL connection: %w", err)
	}

	/* ping PostgreSQL connection */
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping PostgreSQL: %w", err)
	}

	/* return the PostgreSQL client */
	return &PGClient{db: db}, nil
}
