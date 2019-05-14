package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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
		Total         int            `json:"total"`
		FilteredTotal int            `json:"filteredTotal"`
		Page          int            `json:"page"`
		PageSize      int            `json:"pageSize"`
		Accessions    []AccessionRow `json:"accessions"`
	}
	out := SubmissionsPage{Total: 0, Page: page, PageSize: pageSize}

	log.Printf("Get total accessions")
	tq := svc.DB.NewQuery("select count(*) as total from accessions")
	tq.One(&out)

	selQS := fmt.Sprintf(`select a.id as id, identifier, concat(u.first_name, ' ', u.last_name) as submitter, 
		a.description as description, accession_type, group_concat(g.name) genres, 
		(select count(*) from digital_accessions da where da.accession_id=a.id) as digital,
		(select count(*) from physical_accessions pa where pa.accession_id=a.id) as physical,
		a.created_at`)
	fromQS := ` from accessions a 
			inner join users u on u.id = user_id
			inner join accession_genres ag on ag.accession_id = a.id
			inner join genres g on g.id = ag.genre_id
			left outer join digital_accessions da on da.accession_id = a.id
			 left outer join physical_accessions pa on pa.accession_id = a.id`
	qs := selQS + fromQS
	groupQS := " group by a.id"
	pageQS := fmt.Sprintf(" order by created_at desc limit %d,%d", start, pageSize)

	// Check for and apply and filter / query params
	qParam := strings.TrimSpace(c.Query("q"))
	qQuery := ""
	if qParam != "" {
		log.Printf("Filter accessions by query string [%s]", qParam)
		qParam = "%" + qParam + "%"
		qQuery = ` where (a.description like {:q} or da.description like {:q} or tech_description like {:q}
			or first_name like {:q} or last_name like {:q}
			or pa.date_range like {:q} or da.date_range like {:q} or a.created_at like {:q})`
		qs += qQuery
	}

	gParam := strings.TrimSpace(c.Query("g"))
	if gParam != "" {
		log.Printf("Filter submission by genre [%s]", gParam)
		// To ensure all tags are included in result, can't use where clause.
		// If used, only the single matching tag is returned. Instead add a having
		// clause to the group by. The is leaves all tags in the results and matches
		// on the CSV tag list instead.
		groupQS += " having Find_In_Set({:g}, genres)"
	}

	if qParam != "" || gParam != "" {
		countQS := "select count(distinct a.id) as filtered_cnt " + fromQS
		if qQuery != "" {
			countQS += qQuery
		}

		// Since all of the tags are not required for a simple match count,
		// the weird group by and having find_in_set is not needed.
		// Just a simple where will work.
		if gParam != "" {
			if strings.Contains(countQS, " where ") {
				countQS += " and "
			} else {
				countQS += " where "
			}
			countQS += " g.name={:t}"
		}

		log.Printf("Get filtered total")
		cq := svc.DB.NewQuery(countQS)
		cq.Bind(dbx.Params{"q": qParam, "g": gParam})
		cq.Row(&out.FilteredTotal)
	}

	log.Printf("Get one page of submission data")
	qs = qs + groupQS + pageQS
	q := svc.DB.NewQuery(qs)
	q.Bind(dbx.Params{"q": qParam, "g": gParam})
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
	accession.GetGenres(svc.DB)
	accession.GetDigitalTransferDetail(svc.DB)
	accession.GetPhysicalTransferDetail(svc.DB)

	c.JSON(http.StatusOK, accession)
}
