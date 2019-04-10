package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// ControlledVocab contains the data common to controlled vocabulary tables
type ControlledVocab struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Digital     bool   `json:"digitalOnly,omitempty" db:"digital"`
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
