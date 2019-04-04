package main

import (
	"flag"
	"log"
	"os"
)

// ServiceConfig defines all of the archives transfer service configuration paramaters
type ServiceConfig struct {
	DBHost      string
	DBName      string
	DBUser      string
	DBPass      string
	Port        int
	UploadDir   string
	DevAuthUser string
}

// Load will load the service configuration from env/cmdline
func (cfg *ServiceConfig) Load() {
	log.Printf("Loading configuration...")

	flag.StringVar(&cfg.DBHost, "dbhost", "", "DB Host (required)")
	flag.StringVar(&cfg.DBName, "dbname", "", "DB Name (required)")
	flag.StringVar(&cfg.DBUser, "dbuser", "", "DB User (required)")
	flag.StringVar(&cfg.DBPass, "dbpass", "", "DB Password (required)")
	flag.StringVar(&cfg.DevAuthUser, "devuser", "", "Authorized computing id for dev")
	flag.IntVar(&cfg.Port, "port", 8080, "Service port (default 8080)")
	flag.StringVar(&cfg.UploadDir, "upload", "./uploads", "Upload directory")

	flag.Parse()
	log.Printf("%#v", cfg)

	// if anything is still not set, die
	if cfg.DBHost == "" || cfg.DBUser == "" ||
		cfg.DBPass == "" || cfg.DBName == "" {
		flag.Usage()
		log.Printf("FATAL: Missing DB configuration")
		os.Exit(1)
	}
}
