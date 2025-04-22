package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/nromsdahl/squarephish2/internal/dashboard/utils"
	"github.com/nromsdahl/squarephish2/internal/database"
)

// TokenHandler handles the request for the view token page
// Parameters:
//   - w: The http.ResponseWriter
//   - r: The http.Request
//
// It returns an error if the template is not found or executed correctly.
func TokenHandler(w http.ResponseWriter, r *http.Request) {
	email := r.URL.Query().Get("email")
	if email == "" {
		utils.RespondWithError(w, "No email provided", nil)
		return
	}

	token, err := database.LoadToken(email)
	if err != nil {
		utils.RespondWithError(w, "Failed to retrieve token", err)
		return
	}

	// Format the token as pretty JSON
	var prettyJSON bytes.Buffer
	err = json.Indent(&prettyJSON, []byte(token), "", "  ")
	if err != nil {
		utils.RespondWithError(w, "Failed to format token as JSON", err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	utils.RespondWithMessage(w, prettyJSON.String())
}
