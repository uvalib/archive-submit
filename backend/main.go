package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/contrib/static"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

// Version of the service
const version = "1.0.0"

// favHandler is a dummy handler to silence browser API requests that look for /favicon.ico
func favHandler(c *gin.Context) {
}

// versionHandler reports the version of the serivce
func versionHandler(c *gin.Context) {
	c.String(http.StatusOK, "Archive Submission System version %s", version)
}

// healthCheckHandler reports the health of the serivce
func healthCheckHandler(c *gin.Context) {
	hcMap := make(map[string]string)
	hcMap["alive"] = "true"
	c.JSON(http.StatusOK, hcMap)
}

// Upload handles raw file uploads from the front end
func (svc *ServiceContext) Upload(c *gin.Context) {

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

	/**
		WITH CHNKING, LOOKS LIKE THIS (other form params):
		dzuuid: 57030696-46bf-4aa6-8575-9693ef9d7896
	dzchunkindex: 318
	dztotalfilesize: 637439656
	dzchunksize:     2000000
	dztotalchunkcount: 319
	dzchunkbyteoffset: 636000000
	*/
	// file, err := c.FormFile("file")
	// if err != nil {
	// 	c.String(http.StatusBadRequest, fmt.Sprintf("get form err: %s", err.Error()))
	// 	return
	// }

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

/**
 * MAIN
 */
func main() {
	log.Printf("===> Archive Submission System staring up <===")

	// Get config params; service port, directories, DB
	cfg := ServiceConfig{}
	cfg.load()
	svc := ServiceContext{}
	svc.Init(&cfg)

	log.Printf("Setup routes...")
	gin.SetMode(gin.ReleaseMode)
	gin.DisableConsoleColor()
	router := gin.Default()
	router.GET("/favicon.ico", favHandler)
	router.GET("/version", versionHandler)
	router.GET("/healthcheck", healthCheckHandler)
	api := router.Group("/api")
	{
		api.POST("/upload", svc.Upload)
	}
	router.Use(static.Serve("/", static.LocalFile("./public", true)))

	// add a catchall route that renders the index page.
	// based on no-history config setup info here:
	//    https://router.vuejs.org/guide/essentials/history-mode.html#example-server-configurations
	router.NoRoute(func(c *gin.Context) {
		c.File("./public/index.html")
	})

	portStr := fmt.Sprintf(":%d", cfg.Port)
	log.Printf("Start service v%s on port %s", version, portStr)
	log.Fatal(router.Run(portStr))
}
