package routes

import (
	"net/http"

	"github.com/PythonHacker24/linux-acl-management-backend/api/middleware"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/auth"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/health"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/search"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/session"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/traversal"
)

/* all routes for all features are registered here */
func RegisterRoutes(mux *http.ServeMux, sessionManager *session.Manager) {

	allowedOrigin := []string{"http://localhost:3000"}
	allowedMethods := []string{"GET", "POST", "OPTIONS"}
	allowedHeaders := []string{"Content-Type", "Authorization"}

	/* for logging into the backend and creating a session */
	mux.HandleFunc("POST /login",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				auth.LoginHandler(sessionManager),
			),
			allowedOrigin,
			allowedMethods,
			allowedHeaders,
		),
	)

	/* handle OPTIONS preflight requests for /login */
	mux.HandleFunc("OPTIONS /login",
		middleware.CORSMiddleware(
			func(w http.ResponseWriter, r *http.Request) {
				/*
						This handler will never be called because CORSMiddleware handles OPTIONS
					 	but we need it for the route to be registered
				*/
			},
			allowedOrigin,
			allowedMethods,
			allowedHeaders,
		),
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

	/* for fetching list of all users in the LDAP server */
	mux.Handle("GET /users/ldap/search", http.HandlerFunc(
		middleware.LoggingMiddleware(
			middleware.AuthenticationMiddleware(search.SearchUsersHandler),
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
	mux.Handle("/users/transactions/results", http.HandlerFunc(
		middleware.LoggingMiddleware(
			middleware.AuthenticationQueryMiddleware(sessionManager.StreamUserTransactionsResults),
		),
	))

	/* websocket connection for streaming user transactions data from Redis */
	mux.Handle("/users/transactions/pending", http.HandlerFunc(
		middleware.LoggingMiddleware(
			middleware.AuthenticationQueryMiddleware(sessionManager.StreamUserTransactionsPending),
		),
	))

	/* ARCHIVE WILL BE MADE POST REQUEST */

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
