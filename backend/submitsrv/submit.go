package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
	"github.com/rs/xid"
)

// DigitalAccession contains data supporting digital file accessions
type DigitalAccession struct {
	ID          int      `json:"-"`
	AccessionID int      `json:"-" db:"accession_id"`
	Description string   `json:"description" db:"description"`
	DateRange   string   `json:"dateRange" db:"date_range"`
	RecordTypes []string `json:"selectedTypes" db:"-"`
	Files       []string `json:"uploadedFiles" db:"-"`
	TotalSize   int      `json:"totalSizeBytes" db:"upload_size"`
}

// GetFiles retrieves the list of files associated with this accession
func (da *DigitalAccession) GetFiles(db *dbx.DB) {
	q := db.NewQuery("select filename from digital_files where digital_accession_id={:id}")
	q.Bind((dbx.Params{"id": da.ID}))
	rows, _ := q.Rows()
	for rows.Next() {
		var fn string
		rows.Scan(&fn)
		da.Files = append(da.Files, fn)
	}
}

// GetRecordTypes retrieves the list of files associated with this accession
func (da *DigitalAccession) GetRecordTypes(db *dbx.DB) {
	q := db.NewQuery(`select t.name from record_types t 
		inner join accession_record_types a  on a.record_type_id = t.id
		where a.accession_id={:id} and a.accession_type={:type}`)
	q.Bind((dbx.Params{"id": da.ID, "type": "digital"}))
	rows, _ := q.Rows()
	for rows.Next() {
		var fn string
		rows.Scan(&fn)
		da.RecordTypes = append(da.RecordTypes, fn)
	}
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
	Identifier       string            `json:"identifier" db:"identifier"`
	UserID           int               `json:"-" db:"user_id"`
	User             User              `json:"user" db:"-"`
	Summary          string            `json:"summary" binding:"required" db:"description"`
	Activities       string            `json:"activities" db:"activities"`
	Creator          string            `json:"creator" db:"creator"`
	Genres           []string          `json:"genres" db:"-"`
	Type             string            `json:"accessionType" db:"accession_type"`
	CreatedAt        time.Time         `json:"createdAt" db:"created_at"`
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
	for _, genreIDStr := range a.Genres {
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

// GetGenres get genre info for an accession to the DB
func (a *Accession) GetGenres(db *dbx.DB) {
	q := db.NewQuery(`select g.name from genres g inner join accession_genres ag on ag.genre_id = g.id 
		where accession_id={:id}`)
	q.Bind((dbx.Params{"id": a.ID}))
	rows, err := q.Rows()
	if err != nil {
		log.Printf("ERROR: Unable to get generes for %d:%s", a.ID, err.Error())
		return
	}
	for rows.Next() {
		var name string
		rows.Scan(&name)
		a.Genres = append(a.Genres, name)
	}
}

// GetDigitalTransferDetail will get details for a digital transfer
func (a *Accession) GetDigitalTransferDetail(db *dbx.DB) {
	q := db.NewQuery("select * from digital_accessions where accession_id={:id}")
	q.Bind((dbx.Params{"id": a.ID}))
	err := q.One(&a.Digital)
	if err != nil {
		return
	}
	a.DigitalTransfer = true
	a.Digital.GetFiles(db)
	a.Digital.GetRecordTypes(db)
}

// GetPhysicalTransferDetail will get details for a digital transfer
func (a *Accession) GetPhysicalTransferDetail(db *dbx.DB) {
	q := db.NewQuery("select * from physical_accessions where accession_id={:id}")
	q.Bind((dbx.Params{"id": a.ID}))
	err := q.One(&a.Physical)
	if err != nil {
		return
	}
	a.PhysicalTransfer = true
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
	for _, IDStr := range a.Digital.RecordTypes {
		ID, _ := strconv.Atoi(IDStr)
		_, err := tx.Insert("accession_record_types", dbx.Params{
			"accession_id":   a.Digital.ID,
			"accession_type": "digital",
			"record_type_id": ID,
		}).Execute()
		if err != nil {
			log.Printf("WARN: Unable to attach record type %d to accession %d", ID, a.Digital.ID)
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
			"accession_id":   a.Physical.ID,
			"accession_type": "physical",
			"record_type_id": ID,
		}).Execute()
		if err != nil {
			log.Printf("WARN: Unable to attach record type %d to accession %d", ID, a.Physical.ID)
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
	accession.User.FormatPhone()
	err = svc.DB.Model(&accession.User).Exclude("Verified", "VerifyToken", "Admin", "CreatedAt", "email").Update()
	if err != nil {
		log.Printf("WARN: Unable to update %s - %s", accession.User.Email, err.Error())
	}

	log.Printf("Add new accession record")
	tx, _ := svc.DB.Begin()
	accession.UserID = accession.User.ID
	accession.CreatedAt = time.Now()
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

		// Submitted dir gets broken up by YYYY/MM before the submissionID. Make the tree:
		currTime := time.Now()
		xferDir := fmt.Sprintf("%s/%s/%s", svc.UploadDir, "transferred", currTime.Format("2006/01"))
		os.MkdirAll(xferDir, 0777)

		// Move pending into transfer tree
		pendingDir := fmt.Sprintf("%s/%s", svc.UploadDir, "pending")
		uploadDir := fmt.Sprintf("%s/%s", pendingDir, accession.Identifier)
		tgtDir := fmt.Sprintf("%s/%s", xferDir, accession.Identifier)
		log.Printf("Moving pending upload files from %s to %s", uploadDir, tgtDir)
		err = os.Rename(uploadDir, tgtDir)
		if err != nil {
			log.Printf("WARN: Unable to move pending files to submitted: %s", err.Error())
		}
	}
	tx.Commit()

	// Now send recepit to submitter and admins
	accession.User.SendReceiptEmail(svc.DB, svc.SMTP, accession)
	c.String(http.StatusOK, "accepted")
}

// GetAccessionIdentifier will generate an unique token to identify digital content uploads
// It will be used as a storage subdir for files as they are uploaded
func (svc *ServiceContext) GetAccessionIdentifier(c *gin.Context) {
	id := xid.New()
	c.String(http.StatusOK, id.String())
}
