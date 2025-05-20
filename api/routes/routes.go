package routes

import (
	"net/http"

	"github.com/PythonHacker24/linux-acl-management-backend/api/middleware"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/auth"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/health"
)

/* all routes for all features are registered here */
func RegisterRoutes(mux *http.ServeMux) {

	/* for logging into the backend and creating a session */
	mux.Handle("POST /login", http.HandlerFunc(
		middleware.LoggingMiddleware(auth.LoginHandler),
	))

	/* for monitoring the state of overall server and laclm backend */
	mux.Handle("GET /health", http.HandlerFunc(
		middleware.LoggingMiddleware(health.HealthHandler),
	))
}
