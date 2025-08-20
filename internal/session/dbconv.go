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
	/* validate status */
	status := string(session.Status)
	if status != string(StatusActive) && status != string(StatusExpired) && status != string(StatusPending) {
		return nil, fmt.Errorf("invalid session status: %q, must be one of: %q, %q, %q",
			status, StatusActive, StatusExpired, StatusPending)
	}

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

	completedCount := pgtype.Int4{Int32: int32(session.CompletedCount), Valid: true}
	failedCount := pgtype.Int4{Int32: int32(session.FailedCount), Valid: true}

	return &postgresql.StoreSessionPQParams{
		ID:             uuid.MustParse(session.ID.String()),
		Username:       session.Username,
		Ip:             pgtype.Text{String: session.IP, Valid: true},
		UserAgent:      pgtype.Text{String: session.UserAgent, Valid: true},
		Status:         string(session.Status),
		CreatedAt:      createdAt,
		LastActiveAt:   lastActiveAt,
		Expiry:         expiry,
		CompletedCount: completedCount,
		FailedCount:    failedCount,
	}, nil
}

/* converts project transaction struct into PostgreSQL supported format */
func ConvertTransactionPendingtoStoreParams(tx types.Transaction) (postgresql.CreatePendingTransactionPQParams, error) {
	/* marshal ACL entries to JSON bytes */
	entriesJSON, err := json.Marshal(tx.Entries)
	if err != nil {
		return postgresql.CreatePendingTransactionPQParams{}, fmt.Errorf("failed to marshal ACL entries: %w", err)
	}

	/* convert timestamp to pgtype.Timestamptz */
	var timestamp pgtype.Timestamptz
	if err := timestamp.Scan(tx.Timestamp); err != nil {
		return postgresql.CreatePendingTransactionPQParams{}, fmt.Errorf("failed to convert timestamp: %w", err)
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

	return postgresql.CreatePendingTransactionPQParams{
		ID:         tx.ID,
		SessionID:  tx.SessionID,
		Timestamp:  timestamp,
		Operation:  string(tx.Operation),
		TargetPath: tx.TargetPath,
		Entries:    entriesJSON,
		Status:     string(tx.Status),
		Execstatus: tx.ExecStatus,
		ErrorMsg:   errorMsg,
		Output:     output,
		ExecutedBy: tx.ExecutedBy,
		DurationMs: durationMs,
	}, nil
}

/* converts project transaction struct into PostgreSQL supported format */
func ConvertTransactionResulttoStoreParams(tx types.Transaction) (postgresql.CreateResultsTransactionPQParams, error) {
	/* marshal ACL entries to JSON bytes */
	entriesJSON, err := json.Marshal(tx.Entries)
	if err != nil {
		return postgresql.CreateResultsTransactionPQParams{}, fmt.Errorf("failed to marshal ACL entries: %w", err)
	}

	/* convert timestamp to pgtype.Timestamptz */
	var timestamp pgtype.Timestamptz
	if err := timestamp.Scan(tx.Timestamp); err != nil {
		return postgresql.CreateResultsTransactionPQParams{}, fmt.Errorf("failed to convert timestamp: %w", err)
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

	return postgresql.CreateResultsTransactionPQParams{
		ID:         tx.ID,
		SessionID:  tx.SessionID,
		Timestamp:  timestamp,
		Operation:  string(tx.Operation),
		TargetPath: tx.TargetPath,
		Entries:    entriesJSON,
		Status:     string(tx.Status),
		Execstatus: tx.ExecStatus,
		ErrorMsg:   errorMsg,
		Output:     output,
		ExecutedBy: tx.ExecutedBy,
		DurationMs: durationMs,
	}, nil
}
