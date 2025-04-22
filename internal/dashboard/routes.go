package dashboard

import (
	"github.com/gorilla/mux"
	"github.com/nromsdahl/squarephish2/internal/dashboard/handlers"
)

// RegisterRoutes registers the routes for the dashboard
// Parameters:
//   - r: The mux.Router to register the routes with
//
// It returns nothing.
func RegisterRoutes(r *mux.Router) {
	// GET Handlers
	r.HandleFunc("/", handlers.DashboardHandler).Methods("GET")
	r.HandleFunc("/config", handlers.ConfigHandler).Methods("GET")
	r.HandleFunc("/email", handlers.EmailHandler).Methods("GET")
	r.HandleFunc("/token", handlers.TokenHandler).Methods("GET")

	// POST Handlers
	r.HandleFunc("/config", handlers.SaveConfigHandler).Methods("POST")
	r.HandleFunc("/email", handlers.SendEmailHandler).Methods("POST")
}
