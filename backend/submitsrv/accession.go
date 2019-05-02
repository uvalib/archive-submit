package main

import (
	"log"
	"strconv"
	"time"

	dbx "github.com/go-ozzo/ozzo-dbx"
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
	RecordTypes      []string        `json:"selectedTypes" db:"-"`
	TransferMethodID int             `json:"transferMethod" db:"transfer_method_id"`
	TransferMethod   string          `json:"transferMethodName" db:"-"`
	HasDigital       bool            `json:"hasDigital" db:"has_digital"`
	TechInfo         string          `json:"techInfo" db:"tech_description"`
	MediaCarriers    []string        `json:"mediaCarriers" db:"-"`
	MediaCount       string          `json:"mediaCount" db:"media_counts"`
	HasSoftware      bool            `json:"hasSoftware" db:"has_software"`
	Inventory        []InventoryItem `json:"inventory" db:"-"`
}

// TableName defines the expected DB table name that holds data for physical accessions
func (pa *PhysicalAccession) TableName() string {
	return "physical_accessions"
}

// GetRecordTypes retrieves the list of files associated with this accession
func (pa *PhysicalAccession) GetRecordTypes(db *dbx.DB) {
	q := db.NewQuery(`select t.name from record_types t 
		inner join accession_record_types a  on a.record_type_id = t.id
		where a.accession_id={:id} and a.accession_type={:type}`)
	q.Bind(dbx.Params{"id": pa.ID, "type": "physical"})
	rows, _ := q.Rows()
	for rows.Next() {
		var fn string
		rows.Scan(&fn)
		pa.RecordTypes = append(pa.RecordTypes, fn)
	}
}

// GetMediaCarriers retrieves the list of media carriers associated with this accession
func (pa *PhysicalAccession) GetMediaCarriers(db *dbx.DB) {
	q := db.NewQuery(`select c.name from media_carriers c 
		inner join physical_media_carriers p on p.media_carrier_id = c.id
		where p.physical_accession_id={:id}`)
	q.Bind(dbx.Params{"id": pa.ID})
	rows, _ := q.Rows()
	for rows.Next() {
		var name string
		rows.Scan(&name)
		pa.MediaCarriers = append(pa.MediaCarriers, name)
	}
}

// GetInventory retrieves the physical inventory for this accession
func (pa *PhysicalAccession) GetInventory(db *dbx.DB) {
	q := db.NewQuery("select * from inventory_items where physical_accession_id={:id}")
	q.Bind(dbx.Params{"id": pa.ID})
	q.All(&pa.Inventory)
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
	a.Physical.GetRecordTypes(db)
	a.Physical.GetMediaCarriers(db)
	a.Physical.GetInventory(db)
	a.Physical.TransferMethod = GetVocabName(db, "transfer_methods", a.Physical.TransferMethodID)
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
	for _, IDStr := range a.Physical.MediaCarriers {
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
	for _, IDStr := range a.Physical.RecordTypes {
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
