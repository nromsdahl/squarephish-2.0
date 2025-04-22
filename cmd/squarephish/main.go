package main

import (
	"os"

	"github.com/alecthomas/kingpin/v2"
	"github.com/nromsdahl/squarephish2/internal/config"
	"github.com/nromsdahl/squarephish2/internal/dashboard"
	"github.com/nromsdahl/squarephish2/internal/server"
	log "github.com/sirupsen/logrus"
)

var (
	configPath = kingpin.Flag("config", "Path to the config file").Default("config.json").Short('c').String()
	verbose    = kingpin.Flag("verbose", "Enable verbose logging").Short('v').Bool()
)

func main() {
	customFormatter := new(log.TextFormatter)
	customFormatter.TimestampFormat = "2006-01-02 15:04:05"
	customFormatter.FullTimestamp = true
	log.SetFormatter(customFormatter)

	version, err := os.ReadFile("./VERSION")
	if err != nil {
		log.Fatalf("error reading version file: %v", err)
	}
	kingpin.Version(string(version))

	kingpin.CommandLine.HelpFlag.Short('h')
	kingpin.Parse()

	if *verbose {
		log.SetLevel(log.DebugLevel)
	}

	severConfig, err := config.LoadServerConfig(*configPath)
	if err != nil {
		log.Fatalf("error loading config file: %v", err)
	}

	// Run squarephish server in a goroutine
	go server.StartHTTPSServer(&severConfig.PhishConf)

	err = dashboard.ServeDashboard(&severConfig.DashboardConf)
	if err != nil {
		log.Fatalf("error starting dashboard: %v", err)
	}
}
