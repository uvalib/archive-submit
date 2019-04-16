package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/xid"
)

// DigitalAccession contains data supporting digital file accessions
type DigitalAccession struct {
	ID            int      `json:"-"`
	AccessionID   int      `json:"-" db:"accession_id"`
	UploadID      string   `json:"uploadID" db:"upload_id"`
	Description   string   `json:"description" db:"description"`
	DateRange     string   `json:"dateRange" db:"date_range"`
	RecordTypeIDs []string `json:"selectedTypes" db:"-"`
	Files         []string `json:"uploadedFiles" db:"-"`
	TotalSize     int      `json:"totalSizeBytes" db:"upload_size"`
}

// TableName defines the expected DB table name that holds data for digital accessions
func (da *DigitalAccession) TableName() string {
	return "digital_accessions"
}

// PhysicalAccession contains data supporting digital file accessions
type PhysicalAccession struct {
	ID               int      `json:"-"`
	AccessionID      int      `json:"-" db:"accession_id"`
	DateRange        string   `json:"dateRange" db:"date_range"`
	BoxInfo          string   `json:"boxInfo" db:"box_info"`
	RecordTypeIDs    []string `json:"selectedTypes" db:"-"`
	TransferMethodID int      `json:"transferMethod" db:"transfer_method_id"`
	HasDigital       bool     `json:"hasDigital" db:"has_digital"`
	TechInfo         string   `json:"techInfo" db:"tech_description"`
	MediaCarrierIDs  []string `json:"mediaCarriers" db:"-"`
	MediaCount       string   `json:"mediaCount" db:"media_counts"`
	HasSoftware      bool     `json:"hasSoftware" db:"has_software"`
}

// TableName defines the expected DB table name that holds data for physical accessions
func (pa *PhysicalAccession) TableName() string {
	return "physical_accessions"
}

// Accession wraps the general, physical and digital info for a records transfer
// NOTE: the nested structures cause problems in the insert / update calls and must
// be handled separately. The extra *ID and db:"-" accomplishes this. Lists are also
// a problem and must be blocked and handled separately
type Accession struct {
	ID               int               `json:"id" db:"id"`
	UserID           int               `json:"-" db:"user_id"`
	User             User              `json:"user" db:"-"`
	Summary          string            `json:"summary" binding:"required" db:"description"`
	Activities       string            `json:"activities" db:"activities"`
	Creator          string            `json:"creator" db:"creator"`
	GenreIDs         []string          `json:"selectedGenres" db:"-"`
	Type             string            `json:"accessionType" db:"accession_type"`
	DigitalTransfer  bool              `json:"digitalTransfer" db:"-"`
	Digital          DigitalAccession  `json:"digital" db:"-"`
	PhysicalTransfer bool              `json:"physicalTransfer" db:"-"`
	Physical         PhysicalAccession `json:"physical" db:"-"`
}

// TableName defines the expected DB table name that holds data for users
func (a *Accession) TableName() string {
	return "accessions"
}

// WriteGenres writes genre info for an accession to the DB
func (a *Accession) WriteGenres(tx *dbx.Tx) {
	log.Printf("Commmit genres")
	for _, genreIDStr := range a.GenreIDs {
		genreID, _ := strconv.Atoi(genreIDStr)
		_, err := tx.Insert("accession_genres", dbx.Params{
			"accession_id": a.ID,
			"genre_id":     genreID,
		}).Execute()
		if err != nil {
			log.Printf("WARN: Unable to attach genre %d to accession %d", genreID, a.ID)
		}
	}
}

// WriteDigitalTransfer writes digital xfer info for an accession to the DB
func (a *Accession) WriteDigitalTransfer(tx *dbx.Tx) error {
	log.Printf("Commmit digital transfer details for accession %d", a.ID)
	a.Digital.AccessionID = a.ID
	err := tx.Model(&a.Digital).Insert()
	if err != nil {
		return err
	}

	log.Printf("Commmit digital files")
	for _, file := range a.Digital.Files {
		_, err := tx.Insert("digital_files", dbx.Params{
			"digital_accession_id": a.Digital.ID,
			"filename":             file,
		}).Execute()
		if err != nil {
			log.Printf("ERROR: Unable to attach file %s to digital accession %d", file, a.Digital.ID)
			return err
		}
	}

	log.Printf("Commmit digital record types")
	for _, IDStr := range a.Digital.RecordTypeIDs {
		ID, _ := strconv.Atoi(IDStr)
		_, err := tx.Insert("accession_record_types", dbx.Params{
			"accession_id":   a.ID,
			"record_type_id": ID,
		}).Execute()
		if err != nil {
			log.Printf("WARN: Unable to attach record type %d to accession %d", ID, a.ID)
		}
	}
	return nil
}

// WritePhysicalTransfer writes physical xfer info for an accession to the DB
func (a *Accession) WritePhysicalTransfer(tx *dbx.Tx) error {
	log.Printf("Commmit physical transfer details")
	a.Physical.AccessionID = a.ID
	err := tx.Model(&a.Physical).Insert()
	if err != nil {
		return err
	}

	log.Printf("Commmit physical transfer media carriers")
	for _, IDStr := range a.Physical.MediaCarrierIDs {
		ID, _ := strconv.Atoi(IDStr)
		_, err := tx.Insert("physical_media_carriers", dbx.Params{
			"physical_accession_id": a.Physical.ID,
			"media_carrier_id":      ID,
		}).Execute()
		if err != nil {
			log.Printf("WARN: Unable to attach media carrier %d to physical accession %d", ID, a.Physical.ID)
		}
	}

	log.Printf("Commmit physical resord types")
	for _, IDStr := range a.Physical.RecordTypeIDs {
		ID, _ := strconv.Atoi(IDStr)
		_, err := tx.Insert("accession_record_types", dbx.Params{
			"accession_id":   a.ID,
			"record_type_id": ID,
		}).Execute()
		if err != nil {
			log.Printf("WARN: Unable to attach record type %d to accession %d", ID, a.ID)
		}
	}
	return nil
}

// Submit accepts a transfer submission, creates a DB record and kicks off submission
// processing. When complete, an email is sent to the submitter
func (svc *ServiceContext) Submit(c *gin.Context) {
	var accession Accession
	err := c.ShouldBindJSON(&accession)
	if err != nil {
		log.Printf("ERROR: Unable to parse request: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("Received: %+v", accession)

	log.Printf("Update existing user %d:%s", accession.User.ID, accession.User.Email)
	accession.User.UpdatedAt = time.Now()
	err = svc.DB.Model(&accession.User).Exclude("Verified", "VerifyToken", "Admin", "CreatedAt", "email").Update()
	if err != nil {
		log.Printf("WARN: Unable to update %s - %s", accession.User.Email, err.Error())
	}

	log.Printf("Add new accession record")
	tx, _ := svc.DB.Begin()
	accession.UserID = accession.User.ID
	err = tx.Model(&accession).Insert()
	if err != nil {
		log.Printf("ERROR: Unable to add accession %s", err.Error())
		tx.Rollback()
		c.String(http.StatusInternalServerError, "Unable to create accession record")
		return
	}

	accession.WriteGenres(tx)
	if accession.PhysicalTransfer {
		perr := accession.WritePhysicalTransfer(tx)
		if perr != nil {
			log.Printf("ERROR: Unable to write physical xfer: %s", perr.Error())
			tx.Rollback()
			c.String(http.StatusInternalServerError, "Unable to create physical transfer record")
			return
		}
	}
	if accession.DigitalTransfer {
		derr := accession.WriteDigitalTransfer(tx)
		if derr != nil {
			log.Printf("ERROR: Unable to write digital xfer: %s", derr.Error())
			tx.Rollback()
			c.String(http.StatusInternalServerError, "Unable to create digital transfer record")
			return
		}
	}
	tx.Commit()
	c.String(http.StatusOK, "accepted")
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
