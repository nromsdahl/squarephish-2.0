package config

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// DashboardServer represents the Dashboard server configuration details.
type DashboardServer struct {
	ListenURL string `json:"listen_url"`
	CertPath  string `json:"cert_path"`
	KeyPath   string `json:"key_path"`
	UseTLS    bool   `json:"use_tls"`
}

// PhishServer represents the SquarePhish server configuration details.
type PhishServer struct {
	ListenURL string `json:"listen_url"`
	CertPath  string `json:"cert_path"`
	KeyPath   string `json:"key_path"`
	UseTLS    bool   `json:"use_tls"`
}

// Config represents the configuration information.
type Config struct {
	DashboardConf  DashboardServer `json:"dashboard_server"`
	PhishConf      PhishServer     `json:"phish_server"`
}

// LoadConfigFile loads the configuration from the specified filepath
// and returns a Config struct.
// Parameters:
//   - filepath: The path to the configuration file.
//
// It returns a pointer to the Config struct or an error if the file does
// not exist or if the file is not valid JSON.
func LoadServerConfig(fp string) (*Config, error) {
	// Get the config file
	fp = filepath.Clean(fp)
	configFile, err := os.ReadFile(fp)
	if err != nil {
		return nil, err
	}

	config := &Config{}
	err = json.Unmarshal(configFile, config)
	if err != nil {
		return nil, err
	}

	return config, nil
}
