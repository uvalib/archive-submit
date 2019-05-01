package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// GetAccessions is an asmin API call that returns a paged list of accessions
func (svc *ServiceContext) GetAccessions(c *gin.Context) {
	pageStr := c.Query("page")
	const pageSize = 50
	page := 1
	if pageStr != "" {
		pageInt, err := strconv.Atoi(pageStr)
		if err == nil {
			page = pageInt
		}
	}
	start := (page - 1) * pageSize

	type AccessionRow struct {
		ID          int       `json:"id" db:"id"`
		AccessionID string    `json:"accessionID" db:"identifier"`
		Submitter   string    `json:"submitter" db:"submitter"`
		Description string    `json:"description" db:"description"`
		Type        string    `json:"type" db:"accession_type"`
		Genres      string    `json:"genres" db:"genres"`
		Digital     bool      `json:"digital" db:"digital"`
		Physical    bool      `json:"physical" db:"physical"`
		SubmittedAt time.Time `json:"submittedAt" db:"created_at"`
	}
	type SubmissionsPage struct {
		Total      int            `json:"total"`
		Page       int            `json:"page"`
		PageSize   int            `json:"pageSize"`
		Accessions []AccessionRow `json:"accessions"`
	}
	out := SubmissionsPage{Total: 0, Page: page, PageSize: pageSize}

	log.Printf("Get total accessions")
	tq := svc.DB.NewQuery("select count(*) as total from accessions")
	tq.One(&out)

	qs := fmt.Sprintf(`select a.id as id, identifier, u.email as submitter, description, 
		accession_type, group_concat(g.name) genres, 
		(select count(*) from digital_accessions da where da.accession_id=a.id) as digital,
		(select count(*) from physical_accessions pa where pa.accession_id=a.id) as physical,
		a.created_at from accessions a 
			inner join users u on u.id = user_id
			inner join accession_genres ag on ag.accession_id = a.id
			inner join genres g on g.id = ag.genre_id 
		group by a.id
		order by created_at desc limit %d,%d`, start, pageSize)
	q := svc.DB.NewQuery(qs)
	err := q.All(&out.Accessions)
	if err != nil {
		log.Printf("ERROR: Unable to get accessions: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, out)
}

// GetAccessionDetail is an admin API call that returns the full detail of an accession
func (svc *ServiceContext) GetAccessionDetail(c *gin.Context) {
	ID := c.Param("id")
	q := svc.DB.NewQuery("select * from accessions where id={:id}")
	q.Bind((dbx.Params{"id": ID}))
	var accession Accession
	err := q.One(&accession)
	if err != nil {
		log.Printf("ERROR: Unable to get accession %s: %s", ID, err.Error())
		return
	}
	c.JSON(http.StatusOK, accession)
}
