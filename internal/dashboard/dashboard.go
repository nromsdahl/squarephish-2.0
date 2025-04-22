package dashboard

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/nromsdahl/squarephish2/internal/config"
	"github.com/nromsdahl/squarephish2/internal/database"
	log "github.com/sirupsen/logrus"
)

// ServeDashboard starts the HTTP/S server and serves the dashboard.
// Parameters:
//   - config: The server configuration.
//
// It returns an error if the server fails to start (e.g., port conflict, bad certs).
func ServeDashboard(config *config.DashboardServer) error {
	var err error

	// --- Initialize Database ---

	err = database.Initialize()
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}

	// Defer closing the database connection until this server exits.
	defer func() {
		if err := database.Close(); err != nil {
			log.Printf("error closing database: %v", err)
		}
	}()

	// --- Start HTTP/S Server---

	r := mux.NewRouter()

	fs := http.FileServer(http.Dir("static"))
	r.PathPrefix("/static/").Handler(http.StripPrefix("/static/", fs))

	RegisterRoutes(r)

	// Set up the middleware
	r.Use(LoggingMiddleware)

	// Create a new HTTP/S server with timeouts
	srv := &http.Server{
		Addr:               config.ListenURL,
		ReadTimeout:        5 * time.Second,
		ReadHeaderTimeout:  2 * time.Second,
		WriteTimeout:      10 * time.Second,
		IdleTimeout:       30 * time.Second,
		Handler:            r,
	}

	if config.UseTLS {
		err = srv.ListenAndServeTLS(config.CertPath, config.KeyPath)
	} else {
		err = srv.ListenAndServe()
	}

	if err != nil {
		return fmt.Errorf("failed to start dashboard: %v", err)
	}

	return nil
}
