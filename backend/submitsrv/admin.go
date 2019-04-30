package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetAccessions is an asmin API call that returns a paged list of accessions
func (svc *ServiceContext) GetAccessions(c *gin.Context) {
	c.String(http.StatusOK, "accessions")
}
