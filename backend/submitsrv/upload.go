package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/xid"
)

// Accession wraps the general, physical and digital info for a records transfer
type Accession struct {
	Summary          string            `json:"summary"`
	Activities       string            `json:"activities"`
	Creator          string            `json:"creator"`
	GenreIDs         []string          `json:"selectedGenres"`
	Type             string            `json:"accessionType"`
	DigitalTransfer  bool              `json:"digitalTransfer"`
	Digital          DigitalAccession  `json:"digital"`
	PhysicalTransfer bool              `json:"physicalTransfer"`
	Physical         PhysicalAccession `json:"physical"`
}

// DigitalAccession contains data supporting digital file accessions
type DigitalAccession struct {
	UploadID      string   `json:"uploadID"`
	Description   string   `json:"description"`
	DateRange     string   `json:"dateRange"`
	RecordTypeIDs []string `json:"selectedTypes"`
	Files         []string `json:"uploadedFiles"`
	TotalSize     int      `json:"totalSizeBytes"`
}

// PhysicalAccession contains data supporting digital file accessions
type PhysicalAccession struct {
	DateRange        string   `json:"dateRange"`
	BoxInfo          string   `json:"boxInfo"`
	RecordTypeIDs    []string `json:"selectedTypes"`
	TransferMethodID int      `json:"transferMethod"`
	HasDigital       bool     `json:"hasDigital"`
	TechInfo         string   `json:"techInfo"`
	MediaCarrierIDs  []string `json:"mediaCarriers"`
	MediaCount       string   `json:"mediaCount"`
	HasSoftware      bool     `json:"hasSoftware"`
}

// Submission wraps all of the data necessary for a records transfer
type Submission struct {
	User      User      `json:"user"`
	Accession Accession `json:"accession"`
}

// Submit accepts a transfer submission, creates a DB record and kicks off submission
// processing. When complete, an email is sent to the submitter
func (svc *ServiceContext) Submit(c *gin.Context) {
	var submission Submission
	err := c.BindJSON(&submission)
	if err != nil {
		log.Printf("ERROR: Unable to parse request: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
	}
	log.Printf("Received: %+v", submission)

	log.Printf("Processing Submitter info...")
	existingUser := User{}
	err = existingUser.FindByEmail(svc.DB, submission.User.Email)
	if err != nil {
		log.Printf("ERROR: Submitter %s not found in DB - %s", submission.User.Email, err.Error())
		c.String(http.StatusBadRequest, "User %s does not exist", submission.User.Email)
		return
	}
	submission.User.ID = existingUser.ID
	submission.User.UpdatedAt = time.Now()
	log.Printf("Found user %+v => updating fields to %+v", existingUser, submission.User)
	err = svc.DB.Model(&submission.User).Exclude("Verified", "VerifyToken", "Admin", "CreatedAt").Update()
	if err != nil {
		log.Printf("WARN: Unable to ubdate %s - %s", submission.User.Email, err.Error())
	}

	log.Printf("Process common accession information...")
	c.String(http.StatusNotImplemented, "WOOF")
}

// GetSubmissionIdentifier will generate an unique token to identify a new submission
// It will be used as a storage subdir for submission files as they are uploaded
func (svc *ServiceContext) GetSubmissionIdentifier(c *gin.Context) {
	id := xid.New()
	c.String(http.StatusOK, id.String())
}

// UploadFile handles raw file uploads from the front end
func (svc *ServiceContext) UploadFile(c *gin.Context) {
	log.Printf("Checking for upload identifier...")
	uploadID := c.PostForm("identifier")
	if uploadID == "" {
		log.Printf("ERROR: No upload identifier submitted")
		c.String(http.StatusBadRequest, "upload identifier missing")
		return
	}
	log.Printf("Identifier received. Create upload directory.")
	uploadDir := fmt.Sprintf("%s/%s", svc.UploadDir, uploadID)
	os.Mkdir(uploadDir, 0777)

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
		dest := fmt.Sprintf("%s/%s", uploadDir, filename)
		if chunkIdx == "0" {
			if _, err := os.Stat(dest); err == nil {
				log.Printf("WARN: File %s already exists; removing", dest)
				os.Remove(dest)
			}
		}
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
		dest := fmt.Sprintf("%s/%s", uploadDir, filename)
		if _, err := os.Stat(dest); err == nil {
			log.Printf("WARN: File %s already exists; removing", dest)
			os.Remove(dest)
		}
		log.Printf("Receiving non-chunked file %s", filename)
		if err := c.SaveUploadedFile(file, dest); err != nil {
			c.String(http.StatusBadRequest, fmt.Sprintf("upload file err: %s", err.Error()))
			return
		}
		log.Printf("Done receiving %s", filename)
		c.String(http.StatusOK, "Submitted")
	}
}

// DeleteUploadedFile will remove a temporary upload file
func (svc *ServiceContext) DeleteUploadedFile(c *gin.Context) {
	tgtFile := c.Param("file")
	uploadID := c.Query("key")
	tgt := fmt.Sprintf("%s/%s/%s", svc.UploadDir, uploadID, tgtFile)
	log.Printf("Request to delete %s", tgt)
	if _, err := os.Stat(tgt); err == nil {
		delErr := os.Remove(tgt)
		if delErr != nil {
			log.Printf("WARN: Unable to delete %s: %s", tgt, delErr.Error())
			c.String(http.StatusInternalServerError, delErr.Error())
			return
		}
	} else {
		log.Printf("WARN: Target file %s does not exist", tgt)
		c.String(http.StatusNotFound, "% not found", tgtFile)
		return
	}
	log.Printf("Deleted %s", tgt)
	c.String(http.StatusOK, "deleted %s", tgtFile)
}
