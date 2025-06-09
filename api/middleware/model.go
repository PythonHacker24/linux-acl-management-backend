package middleware

/* contextKey type for middleware context value passing */
type contextKey string

/* defining contextKey types */
const (
	ContextKeyUsername  contextKey = "username"
	ContextKeySessionID contextKey = "session_id"
)
