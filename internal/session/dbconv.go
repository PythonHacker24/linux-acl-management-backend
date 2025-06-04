package session

import (
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/postgresql"
)

/* converts project session struct into PostgreSQL supported format */
func ConvertSessionToStoreParams(session *Session) (*postgresql.StoreSessionPQParams, error) {
	createdAt := pgtype.Timestamp{}
	if err := createdAt.Scan(session.CreatedAt); err != nil {
		return nil, err
	}

	lastActiveAt := pgtype.Timestamp{}
	if err := lastActiveAt.Scan(session.LastActiveAt); err != nil {
		return nil, err
	}

	expiry := pgtype.Timestamp{}
	if err := expiry.Scan(session.Expiry); err != nil {
		return nil, err
	}

	return &postgresql.StoreSessionPQParams{
		ID:           uuid.MustParse(session.ID),
		Username:     session.Username,
		Ip:           pgtype.Text{String: session.IP, Valid: true},
		UserAgent:    pgtype.Text{String: session.UserAgent, Valid: true},
		Status:       string(session.Status),
		CreatedAt:    createdAt,
		LastActiveAt: lastActiveAt,
		Expiry:       expiry,
	}, nil
}
