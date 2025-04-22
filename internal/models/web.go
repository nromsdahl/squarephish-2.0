package models

// DashboardData represents the data for the dashboard page.
type DashboardData struct {
	EmailsSent    int
	EmailsScanned int
	Credentials   []Credential
	ActivePage    string
	Title         string
}

// ConfigData represents the configuration data for SquarePhish.
type ConfigData struct {
	SMTPConfig        SMTPConfig
	SquarePhishConfig SquarePhishConfig
	EmailConfig       EmailConfig
	ActivePage        string
	Title             string
}

// EmailData represents the data for the email page.
type EmailData struct {
	ActivePage string
	Title      string
}
