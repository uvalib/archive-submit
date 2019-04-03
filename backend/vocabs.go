package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
