package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Authenticate will authenticate a user based on Shibboleth headers
func (svc *ServiceContext) Authenticate(c *gin.Context) {
	log.Printf("Checking authentication headers...")
	computingID := c.GetHeader("remote_user")
	if svc.DevAuthUser != "" {
		computingID = svc.DevAuthUser
	}
	if computingID == "" {
		c.String(http.StatusForbidden, "You are not authorized to access this site")
		return
	}

	log.Printf("Authentication OK for %s", computingID)
	c.JSON(http.StatusOK, "authorized")
}
