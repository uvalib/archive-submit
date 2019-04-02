package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/go-sql-driver/mysql"
)

// ServiceContext contains the data
type ServiceContext struct {
	UploadDir string
	DB        *dbx.DB
}

// Init will initialize the service context based on the config parameters
func (svc *ServiceContext) Init(cfg *ServiceConfig) {
	log.Printf("Initializing Service...")
	svc.UploadDir = cfg.UploadDir

	log.Printf("Init DB connection to %s...", cfg.DBHost)
	connectStr := fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", cfg.DBUser, cfg.DBPass, cfg.DBHost, cfg.DBName)
	db, err := dbx.Open("mysql", connectStr)
	if err != nil {
		log.Printf("FATAL: Unable to make connection: %s", err.Error())
		os.Exit(1)
	}
	svc.DB = db
}

// GetGenres returns a list of genres as JSON
func (svc *ServiceContext) GetGenres(c *gin.Context) {
	q := svc.DB.NewQuery("SELECT id, name FROM genres")
	var genres []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	err := q.All(&genres)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to retrive genres: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, genres)
}

// UploadFile handles raw file uploads from the front end
func (svc *ServiceContext) UploadFile(c *gin.Context) {

	// fmt.Println(header.Filename)
	// log.Printf("Received request to upload %s; create temp file", filename)
	// out, err := os.Create("./uploads/" + filename)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer out.Close()
	// log.Printf("Receiving file")
	// _, err = io.Copy(out, file)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf("DONEReceiving file")
	// c.String(http.StatusOK, "Submitted")

	// when chunking is being used, there will be additional form params:
	// dzchunkindex,  dztotalfilesize, dzchunksize,  dztotalchunkcount
	// All sizes are in bytes
	chunkIdx := c.PostForm("dzchunkindex")
	if chunkIdx != "" {
		// this is a chunked file; open in append mode and write to it
		file, header, err := c.Request.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("Unable to get form file: %s", err.Error()))
			return
		}
		filename := header.Filename
		log.Printf("Received CHUNKED request to upload %s, chunk %s size %s", filename, chunkIdx, c.PostForm("dzchunksize"))
		dest := fmt.Sprintf("./uploads/%s", filename)
		outFile, err := os.OpenFile(dest, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0777)
		if err != nil {
			c.String(http.StatusInternalServerError, fmt.Sprintf("unable to receive file %s", err.Error()))
		}
		defer outFile.Close()
		_, err = io.Copy(outFile, file)
	} else {
		// not chunked; just save the file in the temp dir
		file, err := c.FormFile("file")
		if err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
			return
		}
		filename := filepath.Base(file.Filename)
		dest := fmt.Sprintf("%s/%s", svc.UploadDir, filename)
		log.Printf("Receiving non-chunked file %s", filename)
		if err := c.SaveUploadedFile(file, dest); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		log.Printf("Done receiving %s", filename)
		c.String(http.StatusOK, "Submitted")
	}

}
