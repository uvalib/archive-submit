package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
	dbx "github.com/go-ozzo/ozzo-dbx"
)

// GetGenres returns a list of genres as JSON
func (svc *ServiceContext) GetGenres(c *gin.Context) {
	q := svc.DB.NewQuery("SELECT id, name FROM genres")
	var genres []struct {
		ID   string `json:"id"`
		Name string `json:"name"`
	}
	err := q.All(&genres)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to retrive genres: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, genres)
}

// GetTypes returns a list of object types as JSON
func (svc *ServiceContext) GetTypes(c *gin.Context) {
	isDigital := c.Query("digital")
	if isDigital == "" {
		isDigital = "0"
	}
	q := svc.DB.NewQuery("SELECT id, name, description FROM record_types where digital={:digital}")
	q.Bind(dbx.Params{"digital": isDigital})
	var objects []struct {
		ID          string `json:"id"`
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	err := q.All(&objects)
	if err != nil {
		c.String(http.StatusInternalServerError, "Unable to retrive record types: %s", err.Error())
		return
	}
	c.JSON(http.StatusOK, objects)
}
