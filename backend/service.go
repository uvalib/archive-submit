package main

import (
	"fmt"
	"log"
	"os"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// ServiceContext contains the data
type ServiceContext struct {
	UploadDir   string
	DevAuthUser string
	DB          *dbx.DB
}

// Init will initialize the service context based on the config parameters
func (svc *ServiceContext) Init(cfg *ServiceConfig) {
	log.Printf("Initializing Service...")
	svc.UploadDir = cfg.UploadDir
	svc.DevAuthUser = cfg.DevAuthUser

	log.Printf("Init DB connection to %s...", cfg.DBHost)
	connectStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBName)
	db, err := dbx.Open("mysql", connectStr)
	if err != nil {
		log.Printf("FATAL: Unable to make connection: %s", err.Error())
		os.Exit(1)
	}
	svc.DB = db
}
