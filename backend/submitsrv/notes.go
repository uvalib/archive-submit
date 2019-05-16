package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// Note contains DB and JSON mappingss for a note
type Note struct {
	ID        int       `json:"id" db:"id"`
	Title     string    `json:"title" db:"title"`
	Note      string    `json:"note" db:"note"`
	UserID    int       `json:"userID" db:"user_id"`
	UserName  string    `json:"userName" db:"user_name"`
	CreatedAt time.Time ` json:"createdAt" db:"created_at"`
}

// TableName defines the expected DB table name that holds data for users
func (note *Note) TableName() string {
	return "notes"
}

// AddAccessionNote adds a note to the specified accession
func (svc *ServiceContext) AddAccessionNote(c *gin.Context) {
	accessionID := c.Param("id")
	log.Printf("Add new note to accession %s", accessionID)
	var note Note
	err := c.ShouldBindJSON(&note)
	if err != nil {
		log.Printf("ERROR: Unable to parse POST data: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	note.CreatedAt = time.Now()
	tx, _ := svc.DB.Begin()
	err = tx.Model(&note).Exclude("UserName").Insert()
	if err != nil {
		log.Printf("ERROR: Add note failed: %s", err.Error())
		tx.Rollback()
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	log.Printf("Note %d added, linking to accession %s", note.ID, accessionID)
	nq := tx.NewQuery("insert into accession_notes (accession_id,note_id) values ({:aid},{:nid})")
	nq.Bind(dbx.Params{"aid": accessionID, "nid": note.ID})
	_, err = nq.Execute()
	if err != nil {
		log.Printf("ERROR: Add accession_note failed: %s", err.Error())
		tx.Rollback()
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	tx.Commit()

	q := svc.DB.NewQuery("select concat(first_name,' ',last_name) as user_name from users where id = {:id}")
	q.Bind(dbx.Params{"id": note.UserID})
	q.One(&note)
	c.JSON(http.StatusOK, note)
}

// GetAccessionNotes returns the details about notes for a particular accession
func (svc *ServiceContext) GetAccessionNotes(c *gin.Context) {
	accessionID := c.Param("id")
	q := svc.DB.NewQuery(`select n.*,concat(u.first_name,' ',u.last_name) as user_name from notes n
			inner join accession_notes an on an.note_id = n.id
			inner join users u on u.id = n.user_id
		where an.accession_id={:id}`)
	q.Bind(dbx.Params{"id": accessionID})
	var notes []Note
	err := q.All(&notes)
	if err != nil {
		log.Printf("ERROR: Unable to get notes for accessions %s:%s", accessionID, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}
	c.JSON(http.StatusOK, notes)
}
