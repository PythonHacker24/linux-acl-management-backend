package routes

import (
	"net/http"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/handlers"
)

func RegisterRoutes(mux *http.ServeMux) {
	mux.Handle("/health", http.HandlerFunc((handlers.HealthHandler)))
}
