package handlers

import (
	"net/http"

	"github.com/nromsdahl/squarephish2/internal/dashboard/templates"
	"github.com/nromsdahl/squarephish2/internal/dashboard/utils"
	"github.com/nromsdahl/squarephish2/internal/database"
	"github.com/nromsdahl/squarephish2/internal/models"
)

// ConfigHandler handles the request for the config page
// Parameters:
//   - w: The http.ResponseWriter
//   - r: The http.Request
//
// It returns an error if the template is not found or executed correctly.
func ConfigHandler(w http.ResponseWriter, r *http.Request) {
	data := database.LoadConfig()

	tmpl, err := templates.GetTemplate("config.html")
	if err != nil {
		utils.RespondWithError(w, "Error getting config template", err)
		return
	}

	err = tmpl.ExecuteTemplate(w, "base", data)
	if err != nil {
		utils.RespondWithError(w, "Error executing config template", err)
		return
	}
}

// SaveConfigHandler handles the submission of a configuration
// Parameters:
//   - w: The http.ResponseWriter
//   - r: The http.Request
//
// It returns an error if the form data is not parsed correctly.
func SaveConfigHandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		utils.RespondWithErrorMessage(w, "Failed to save configuration", err)
		return
	}

	smtpHost := r.FormValue("smtpHost")
	smtpPort := r.FormValue("smtpPort")
	smtpUsername := r.FormValue("smtpUsername")
	smtpPassword := r.FormValue("smtpPassword")

	phishHost := r.FormValue("phishHost")
	phishPort := r.FormValue("phishPort")

	emailSender := r.FormValue("emailSender")
	emailSubject := r.FormValue("emailSubject")
	emailBody := r.FormValue("emailBody")

	// Save the configuration data to the database
	err = database.SaveConfig(models.ConfigData{
		SMTPConfig: models.SMTPConfig{
			Host:     smtpHost,
			Port:     smtpPort,
			Username: smtpUsername,
			Password: smtpPassword,
		},
		SquarePhishConfig: models.SquarePhishConfig{
			Host: phishHost,
			Port: phishPort,
		},
		EmailConfig: models.EmailConfig{
			Subject: emailSubject,
			Sender:  emailSender,
			Body:    emailBody,
		},
	})
	if err != nil {
		utils.RespondWithErrorMessage(w, "Failed to save configuration", err)
		return
	}

	utils.RespondWithMessage(w, "Configuration saved successfully")
}
