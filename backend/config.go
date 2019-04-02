package main

import (
	"flag"
	"log"
	"os"
)

// ServiceConfig defines all of the archives transfer service configuration paramaters
type ServiceConfig struct {
	DBHost    string
	DBName    string
	DBUser    string
	DBPass    string
	Port      int
	UploadDir string
}

func (cfg *ServiceConfig) load() {
	log.Printf("Loading configuration...")

	flag.StringVar(&cfg.DBHost, "dbhost", "", "DB Host (required)")
	flag.StringVar(&cfg.DBName, "dbname", "", "DB Name (required)")
	flag.StringVar(&cfg.DBUser, "dbuser", "", "DB User (required)")
	flag.StringVar(&cfg.DBPass, "dbpass", "", "DB Password (required)")
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
