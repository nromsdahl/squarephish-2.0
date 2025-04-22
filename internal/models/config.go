package models

// SMTPConfig holds the configuration details for connecting and
// authenticating to an SMTP server.
type SMTPConfig struct {
	Host     string `json:"host"`
	Port     string `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// SquarePhishConfig represents the SquarePhish server configuration.
type SquarePhishConfig struct {
	Host string `json:"host"`
	Port string `json:"port"`
}

// EmailConfig represents an email configuration.
type EmailConfig struct {
	Sender  string `json:"sender"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// Credential represents a user's email and authentication token.
type Credential struct {
	Email string        `json:"email"`
	Token TokenResponse `json:"token"`
}
