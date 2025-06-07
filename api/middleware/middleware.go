package middleware

import (
	"context"
	"net/http"
	"time"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/token"
	"go.uber.org/zap"
)

/* logging middleware for http requests */
func LoggingMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		/* logging recieved request at the instant of receiving */
		zap.L().Info("Recieved request",
			zap.String("Method", r.Method),
			zap.String("Path", r.URL.Path),
		)

		/* return the handler */
		next(w, r)

		/* logging time taken by the request */
		zap.L().Info("Request completed",
			zap.String("Method", r.Method),
			zap.String("Path", r.URL.Path),
			zap.Duration("Duration", time.Since(start)),
		)
	})
}

/* authentication middleware for http requests */
func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/* authenticate the request through JWT */
		username, err := token.ExtractUsernameFromRequest(r)
		if err != nil {
			zap.L().Error("Error during authentication",
				zap.Error(err),
			)
			return
		}

		/* set the header with the username */
		r.Header.Set("X-User", username)

		/* pass username as context */
		ctx := context.WithValue(r.Context(), "username", username)

		/* return the handler */
		next(w, r.WithContext(ctx))
	})
}
