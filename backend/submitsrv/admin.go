package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
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

	qs := fmt.Sprintf(`select id, identifier, created_at from accessions 
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
