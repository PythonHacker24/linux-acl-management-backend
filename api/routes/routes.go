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
}
