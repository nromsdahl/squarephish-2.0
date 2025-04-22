// filepath: /internal/dashboard/handlers/dashboard.go
package handlers

import (
	"net/http"

	"github.com/nromsdahl/squarephish2/internal/dashboard/templates"
	"github.com/nromsdahl/squarephish2/internal/dashboard/utils"
	"github.com/nromsdahl/squarephish2/internal/database"
)

// DashboardHandler handles the request for the dashboard
// Parameters:
//   - w: The http.ResponseWriter
//   - r: The http.Request
//
// It returns an error if the template is not found or executed correctly.
func DashboardHandler(w http.ResponseWriter, r *http.Request) {
	data := database.LoadDashboardData()

	tmpl, err := templates.GetTemplate("dashboard.html")
	if err != nil {
		utils.RespondWithError(w, "Error getting dashboard template", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		utils.RespondWithError(w, "Error executing dashboard template", err)
		return
	}
}
