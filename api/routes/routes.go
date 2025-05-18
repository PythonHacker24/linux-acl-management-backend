package routes

import (
	"net/http"

	"github.com/PythonHacker24/linux-acl-management-backend/api/middleware"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/health"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /health", http.HandlerFunc(
		middleware.LoggingMiddleware(health.HealthHandler),
	))
}
