package dashboard

import (
	"net/http"

	log "github.com/sirupsen/logrus"
)

// LoggingMiddleware is a middleware that logs the incoming requests.
// It logs the request method and URL path for debugging.
// Parameters:
//   - next: The next http.Handler in the chain.
//
// It returns an http.Handler that wraps the next handler with logging functionality.
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.WithFields(log.Fields{
			"method": r.Method,
			"path":   r.URL.Path,
		}).Debug("Dashboard request")

		next.ServeHTTP(w, r)
	})
}
