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

	/* move it to config file */
	allowedOrigin := []string{"http://localhost:3000"}
	allowedMethods := []string{"GET", "POST", "OPTIONS"}
	allowedHeaders := []string{"*"}

	/* for monitoring the state of overall server and laclm backend */
	mux.Handle("GET /health", http.HandlerFunc(
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(health.HealthHandler),
			allowedOrigin,
			allowedMethods,
			allowedHeaders,
		),
	))

	/* handle OPTIONS preflight requests for /health */
	mux.HandleFunc("OPTIONS /health",
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

	/* for logging into the backend and creating a session */
	mux.HandleFunc("POST /auth/login",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				auth.LoginHandler(sessionManager),
			),
			allowedOrigin,
			allowedMethods,
			allowedHeaders,
		),
	)

	/* handle OPTIONS preflight requests for /auth/login */
	mux.HandleFunc("OPTIONS /auth/login",
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

	/* for logging out of the backend and expiring the session */
	mux.HandleFunc("GET /auth/logout",
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				auth.LogoutHandler(sessionManager),
			),
			allowedOrigin,
			allowedMethods,
			allowedHeaders,
		),
	)

	/* handle OPTIONS preflight requests for /auth/logout */
	mux.HandleFunc("OPTIONS /auth/logout",
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

	/* for verifying if a token is valid or not */
	mux.Handle("GET /auth/token/validate", http.HandlerFunc(
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				auth.ValidateToken,
			),
			allowedOrigin,
			allowedMethods,
			allowedHeaders,
		),
	))

	/* handle OPTIONS preflight requests for /auth/token/validate */
	mux.HandleFunc("OPTIONS /auth/token/validate",
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

	/* for listing files in a directory */
	mux.Handle("POST /traverse/list-files", http.HandlerFunc(
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthenticationMiddleware(traversal.ListFilesInDirectory),
			),
			allowedOrigin,
			allowedMethods,
			allowedHeaders,
		),
	))

	/* handle OPTIONS preflight requests for /traverse/list-files */
	mux.HandleFunc("OPTIONS /traverse/list-files",
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

	/* for scheduling a transaction */
	mux.Handle("POST /transactions/schedule", http.HandlerFunc(
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthenticationMiddleware(sessionManager.IssueTransaction),
			),
			allowedOrigin,
			allowedMethods,
			allowedHeaders,
		),
	))

	/* handle OPTIONS preflight requests for /transactions/schedule */
	mux.HandleFunc("OPTIONS /transactions/schedule",
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

	/*
		for fetching list of users matching the query in the LDAP server
		supports URL params: q (Query)
	*/
	mux.Handle("GET /users/ldap/search", http.HandlerFunc(
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthenticationMiddleware(search.SearchUsersHandler),
			),
			allowedOrigin,
			allowedMethods,
			allowedHeaders,
		),
	))

	/* handle OPTIONS preflight requests for /users/ldap/search */
	mux.HandleFunc("OPTIONS /users/ldap/search",
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

	/*
		websocket connection for streaming user session data from Redis
		supports URL pamars: token (JWT authentication)
	*/
	mux.Handle("/users/session", http.HandlerFunc(
		middleware.LoggingMiddleware(
			/* you need authentication via query parameter */
			middleware.AuthenticationQueryMiddleware(sessionManager.StreamUserSession),
		),
	))

	/*
		websocket connection for streaming user transactions data from Redis
		supports URL pamars: token (JWT authentication)
	*/
	mux.Handle("/users/transactions/results", http.HandlerFunc(
		middleware.LoggingMiddleware(
			middleware.AuthenticationQueryMiddleware(sessionManager.StreamUserTransactionsResults),
		),
	))

	/*
		websocket connection for streaming user transactions data from Redis
		supports URL pamars: token (JWT authentication)
	*/
	mux.Handle("/users/transactions/pending", http.HandlerFunc(
		middleware.LoggingMiddleware(
			middleware.AuthenticationQueryMiddleware(sessionManager.StreamUserTransactionsPending),
		),
	))

	/* ARCHIVE WILL BE MADE POST REQUEST -> Header based Authentication */

	/* websocket connection for streaming user session data from PostgreSQL database (archived sessions) */
	mux.Handle("POST /users/archive/session", http.HandlerFunc(
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthenticationMiddleware(sessionManager.StreamUserArchiveSessions),
			),
			allowedOrigin,
			allowedMethods,
			allowedHeaders,
		),
	))

	/* handle OPTIONS preflight requests for /users/archive/session */
	mux.HandleFunc("OPTIONS /users/archive/session",
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

	/* websocket connection for streaming user transactions data from PostgreSQL database (archived sessions) */
	mux.Handle("POST /users/archive/transactions/results", http.HandlerFunc(
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthenticationMiddleware(sessionManager.StreamUserArchiveResultsTransactions),
			),
			allowedOrigin,
			allowedMethods,
			allowedHeaders,
		),
	))

	/* handle OPTIONS preflight requests for /users/archive/transactions/results */
	mux.HandleFunc("OPTIONS /users/archive/transactions/results",
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

	/* websocket connection for streaming user transactions data from PostgreSQL database (archived sessions) */
	mux.Handle("POST /users/archive/transactions/pending", http.HandlerFunc(
		middleware.CORSMiddleware(
			middleware.LoggingMiddleware(
				middleware.AuthenticationMiddleware(sessionManager.StreamUserArchivePendingTransactions),
			),
			allowedOrigin,
			allowedMethods,
			allowedHeaders,
		),
	))

	/* handle OPTIONS preflight requests for /users/archive/transactions/pending */
	mux.HandleFunc("OPTIONS /users/archive/transactions/pending",
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
}
