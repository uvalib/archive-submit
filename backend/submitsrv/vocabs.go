package main

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// ControlledVocab contains the data common to controlled vocabulary tables
type ControlledVocab struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Digital     bool   `json:"digitalOnly,omitempty" db:"digital"`
}

// GetVocabNamesCSV will convert a list og string identifiers into a comma separated list of
// controlled value names
func GetVocabNamesCSV(db *dbx.DB, table string, ids []string) string {
	IDs := strings.Join(ids, ",")
	if IDs == "" {
		// not IDs present, return a blank string
		return ""
	}
	qs := fmt.Sprintf("select name from %s where id in (%s)", table, IDs)
	q := db.NewQuery(qs)
	var out []string
	rows, _ := q.Rows()
	for rows.Next() {
		var val string
		rows.Scan(&val)
		out = append(out, val)
	}
	return strings.Join(out, ", ")
}

// GetVocabName will get the text name for the controlled vocabulary ID
func GetVocabName(db *dbx.DB, table string, id int) string {
	qs := fmt.Sprintf("select name from %s where id = %d", table, id)
	q := db.NewQuery(qs)
	var out struct{ Name string }
	q.One(&out)
	return out.Name
}

// GetGenres returns a list of genres as JSON
func (svc *ServiceContext) GetGenres(c *gin.Context) {
	q := svc.DB.NewQuery("SELECT id, name FROM genres")
	var genres []ControlledVocab
	err := q.All(&genres)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to retrive genres: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, genres)
}

// GetTransferMethods returns a list of genres as JSON
func (svc *ServiceContext) GetTransferMethods(c *gin.Context) {
	q := svc.DB.NewQuery("SELECT id, name FROM transfer_methods")
	var methods []ControlledVocab
	err := q.All(&methods)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to retrive transfer methods: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, methods)
}

// GetMediaCarriers returns a list of genres as JSON
func (svc *ServiceContext) GetMediaCarriers(c *gin.Context) {
	q := svc.DB.NewQuery("SELECT id, name FROM media_carriers")
	var carriers []ControlledVocab
	err := q.All(&carriers)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to retrive carriers: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, carriers)
}

// GetTypes returns a list of object types as JSON
func (svc *ServiceContext) GetTypes(c *gin.Context) {
	q := svc.DB.NewQuery("SELECT id, name, description, digital FROM record_types")
	var objects []ControlledVocab
	err := q.All(&objects)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to retrive record types: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, objects)
}
