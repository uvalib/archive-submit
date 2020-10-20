package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// ServiceContext contains the data
type ServiceContext struct {
	UploadDir   string
	DevAuthUser string
	Hostname    string
	DB          *dbx.DB
	SMTP        SMTPConfig
}

// Init will initialize the service context based on the config parameters
func (svc *ServiceContext) Init(cfg *ServiceConfig) {
	log.Printf("Initializing Service...")
	svc.UploadDir = cfg.UploadDir
	svc.DevAuthUser = cfg.DevAuthUser
	svc.Hostname = cfg.Hostname
	svc.SMTP = cfg.SMTP

	log.Printf("Init DB connection to %s...", cfg.DBHost)
	connectStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBName)
	db, err := dbx.Open("mysql", connectStr)
	if err != nil {
		log.Printf("FATAL: Unable to make connection: %s", err.Error())
		os.Exit(1)
	}
	db.LogFunc = log.Printf
	svc.DB = db
}

// GetVersion reports the version of the serivce
func (svc *ServiceContext) GetVersion(c *gin.Context) {

	build := "unknown"

	// cos our CWD is the bin directory
	files, _ := filepath.Glob("../buildtag.*")
	if len(files) == 1 {
		build = strings.Replace(files[0], "../buildtag.", "", 1)
	}

	vMap := make(map[string]string)
	vMap["version"] = version
	vMap["build"] = build
	c.JSON(http.StatusOK, vMap)
}

// HealthCheck reports the health of the serivce
func (svc *ServiceContext) HealthCheck(c *gin.Context) {
	q := svc.DB.NewQuery("select version from versions order by created_at desc limit 1")
	var version string
	err := q.One(&version)
	if err != nil {
		// gin.H is a shortcut for map[string]interface{}
		c.JSON(http.StatusInternalServerError, gin.H{"alive": "true", "mysql": "false"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"alive": "true", "mysql": "true"})
}
