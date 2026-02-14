package middleware

import (
	"net/http"
	"runtime/debug"

	"github.com/rs/zerolog/log"
)

// Recovery recovers from panics in HTTP handlers, logs the stack trace,
// and returns a 500 Internal Server Error.
func Recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if rec := recover(); rec != nil {
				log.Error().
					Interface("panic", rec).
					Str("stack", string(debug.Stack())).
					Str("method", r.Method).
					Str("path", r.URL.Path).
					Msg("recovered from panic")

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte(`{"error":"internal_server_error","message":"An unexpected error occurred"}`))
			}
		}()
		next.ServeHTTP(w, r)
	})
}
