package middleware

import (
	"context"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/PythonHacker24/linux-acl-management-backend/internal/token"
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

/* 
	authentication middleware for http requests 
	return username and sessionID with context
*/
func AuthenticationMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/* authenticate the request through JWT */
		username, sessionID, err := token.ExtractDataFromRequest(r)
		if err != nil {
			zap.L().Info("Error during authentication",
				zap.Error(err),
			)
			http.Error(w, "Authentication Failed", http.StatusInternalServerError)
			return
		}

		/* set the header with the username */
		r.Header.Set("X-User", username)

		/* pass username and sessionID as context */
		ctx := context.WithValue(r.Context(), ContextKeyUsername, username)
		ctx = context.WithValue(ctx, ContextKeySessionID, sessionID)

		/* return the handler */
		next(w, r.WithContext(ctx))
	})
}

/* 
	authentication middleware for http requests with query 
	return username and sessionID with context
*/
func AuthenticationQueryMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		/* get the HTTP query */
		query := r.URL.Query()

		/* get the token query */
		tokenQ := query.Get("token")
        if tokenQ == "" {
           	zap.L().Info("Query authentication without token value")
			http.Error(w, "Missing 'token' query parameter value", http.StatusBadRequest)
			return
        } 

		/* extract username and sessionID from the token */
		username, sessionID, err := token.GetDataFromJWT(tokenQ)
		if err != nil {
			zap.L().Info("Error during authentication",
				zap.Error(err),
			)
			http.Error(w, "Authentication Failed", http.StatusInternalServerError)
			return
		}

		/* set the header with the username */
		r.Header.Set("X-User", username)

		/* pass username and sessionID as context */
		ctx := context.WithValue(r.Context(), ContextKeyUsername, username)
		ctx = context.WithValue(ctx, ContextKeySessionID, sessionID)

		/* return the handler */
		next(w, r.WithContext(ctx))
	})
}
