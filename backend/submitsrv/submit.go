package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
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

// InventoryItem contains data to describe a physical inventory item
type InventoryItem struct {
	ID              int    `json:"-"`
	PhysAccessionID int    `json:"-" db:"physical_accession_id"`
	BoxNumber       string `json:"boxNum" db:"box_number"`
	RecordGroup     string `json:"recordGroup" db:"record_group_number"`
	Title           string `json:"title" db:"box_title"`
	Description     string `json:"description" db:"description"`
	Dates           string `json:"dates" db:"dates"`
}

// TableName defines the expected DB table name that holds data for inventory items
func (da *InventoryItem) TableName() string {
	return "inventory_items"
}

// PhysicalAccession contains data supporting physical file accessions
type PhysicalAccession struct {
	ID               int             `json:"-"`
	AccessionID      int             `json:"-" db:"accession_id"`
	DateRange        string          `json:"dateRange" db:"date_range"`
	BoxInfo          string          `json:"boxInfo" db:"box_info"`
	RecordTypeIDs    []string        `json:"selectedTypes" db:"-"`
	TransferMethodID int             `json:"transferMethod" db:"transfer_method_id"`
	HasDigital       bool            `json:"hasDigital" db:"has_digital"`
	TechInfo         string          `json:"techInfo" db:"tech_description"`
	MediaCarrierIDs  []string        `json:"mediaCarriers" db:"-"`
	MediaCount       string          `json:"mediaCount" db:"media_counts"`
	HasSoftware      bool            `json:"hasSoftware" db:"has_software"`
	Inventory        []InventoryItem `json:"inventory" db:"-"`
}

// TableName defines the expected DB table name that holds data for physical accessions
func (pa *PhysicalAccession) TableName() string {
	return "physical_accessions"
}

// Accession wraps the general, physical and digital info for a records transfer
// NOTE: block out lists and associated structures so there is no attempt to write
// the to the DB when the Accession is written. Handle them as separate commits / reads
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
			"accession_type": "digital",
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

	log.Printf("Commmit physical record types")
	for _, IDStr := range a.Physical.RecordTypeIDs {
		ID, _ := strconv.Atoi(IDStr)
		_, err := tx.Insert("accession_record_types", dbx.Params{
			"accession_id":   a.ID,
			"accession_type": "physical",
			"record_type_id": ID,
		}).Execute()
		if err != nil {
			log.Printf("WARN: Unable to attach record type %d to accession %d", ID, a.ID)
		}
	}

	log.Printf("Commmit physical inventory")
	for _, item := range a.Physical.Inventory {
		item.PhysAccessionID = a.Physical.ID
		err = tx.Model(&item).Insert()
		if err != nil {
			log.Printf("WARN: Unable to attach inventory %+v to physical accession %d", item, a.Physical.ID)
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
