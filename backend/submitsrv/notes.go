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
