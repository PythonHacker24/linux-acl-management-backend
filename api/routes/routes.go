package routes

import (
	"net/http"

	"github.com/PythonHacker24/linux-acl-management-backend/api/middleware"
	"github.com/PythonHacker24/linux-acl-management-backend/internal/handlers"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("GET /health", http.HandlerFunc(
		middleware.LoggingMiddleware(handlers.HealthHandler),
	))
}
