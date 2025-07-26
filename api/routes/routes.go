package routes

import (
	"net/http"

	"github.com/PythonHacker24/linux-acl-management-backend/api/middleware"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/auth"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/health"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/session"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/traversal"
)

/* all routes for all features are registered here */
func RegisterRoutes(mux *http.ServeMux, sessionManager *session.Manager) {

	/* for logging into the backend and creating a session */
	mux.HandleFunc("POST /login",
		middleware.LoggingMiddleware(auth.LoginHandler(sessionManager)),
	)

	/* for monitoring the state of overall server and laclm backend */
	mux.Handle("GET /health", http.HandlerFunc(
		middleware.LoggingMiddleware(health.HealthHandler),
	))

	/* for listing files in a directory */
	mux.Handle("POST /traverse/list-files", http.HandlerFunc(
		middleware.LoggingMiddleware(
			middleware.AuthenticationMiddleware(traversal.ListFilesInDirectory),
		),
	))

	/* for scheduling a transaction */
	mux.Handle("POST /transactions/schedule", http.HandlerFunc(
		middleware.LoggingMiddleware(
			middleware.AuthenticationMiddleware(sessionManager.IssueTransaction),
		),
	))

	/* websocket connection for streaming user session data from Redis */
	mux.Handle("/users/session", http.HandlerFunc(
		middleware.LoggingMiddleware(
			/* you need authentication via query parameter */
			middleware.AuthenticationQueryMiddleware(sessionManager.StreamUserSession),
		),
	))

	/* websocket connection for streaming user transactions data from Redis */
	mux.Handle("/users/transactions", http.HandlerFunc(
		middleware.LoggingMiddleware(
			middleware.AuthenticationMiddleware(sessionManager.StreamUserTransactions),
		),
	))

	/* websocket connection for streaming user session data from PostgreSQL database (archived sessions) */
	mux.Handle("/users/archive/session", http.HandlerFunc(
		middleware.LoggingMiddleware(
			middleware.AuthenticationMiddleware(sessionManager.StreamUserArchiveSessions),
		),
	))

	/* websocket connection for streaming user transactions data from PostgreSQL database (archived sessions) */
	mux.Handle("/users/archive/transactions/pending", http.HandlerFunc(
		middleware.LoggingMiddleware(
			middleware.AuthenticationMiddleware(sessionManager.StreamUserArchivePendingTransactions),
		),
	))
}
