package models

// These are the models for loading data into the web frontend.

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
	EntraConfig       EntraConfig
	RequestConfig     RequestConfig
	ActivePage        string
	Title             string
}

// EmailData represents the data for the email page.
type EmailData struct {
	ActivePage string
	Title      string
}
