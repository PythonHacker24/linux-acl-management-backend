package session

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/postgresql"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/types"
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

/* converts project transaction struct into PostgreSQL supported format */
func ConvertTransactiontoStoreParams(tx types.Transaction) (postgresql.CreateTransactionPQParams, error) {
	/* marshal ACL entries to JSON bytes */
	entriesJSON, err := json.Marshal(tx.Entries)
	if err != nil {
		return postgresql.CreateTransactionPQParams{}, fmt.Errorf("failed to marshal ACL entries: %w", err)
	}

	/* convert timestamp to pgtype.Timestamptz */
	var timestamp pgtype.Timestamptz
	if err := timestamp.Scan(tx.Timestamp); err != nil {
		return postgresql.CreateTransactionPQParams{}, fmt.Errorf("failed to convert timestamp: %w", err)
	}

	/* handle optional error message */
	var errorMsg pgtype.Text
	if tx.ErrorMsg != "" {
		errorMsg = pgtype.Text{String: tx.ErrorMsg, Valid: true}
	}

	/* handle optional output */
	var output pgtype.Text
	if tx.Output != "" {
		output = pgtype.Text{String: tx.Output, Valid: true}
	}

	/* Handle optional duration */
	var durationMs pgtype.Int8
	if tx.DurationMs > 0 {
		durationMs = pgtype.Int8{Int64: tx.DurationMs, Valid: true}
	}

	return postgresql.CreateTransactionPQParams{
		ID:         tx.ID,
		SessionID:  tx.SessionID,
		Timestamp:  timestamp,
		Operation:  string(tx.Operation),
		TargetPath: tx.TargetPath,
		Entries:    entriesJSON,
		Status:     string(tx.Status),
		ErrorMsg:   errorMsg,
		Output:     output,
		ExecutedBy: tx.ExecutedBy,
		DurationMs: durationMs,
	}, nil
}
