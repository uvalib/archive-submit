package main

import (
	"fmt"
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
		log.Printf("ERROR: Expected auth header not present in request. Not authorized.")
		c.String(http.StatusForbidden, "You are not authorized to access this site")
		return
	}

	email := fmt.Sprintf("%s@virginia.edu", computingID)
	user := User{}
	err := user.FindByEmail(svc.DB, email)
	if err != nil {
		log.Printf("No user record found for authorized computing ID; creating one")
		user.Email = email
		user.Verified = true
		user.Affiliation = "UVA Library"
		createErr := user.Create(svc.DB)
		if createErr != nil {
			log.Printf("ERROR: Unable to create NetBadge authorized user record: %s", createErr.Error())
		}
	}

	log.Printf("Authentication successful for %s", computingID)
	c.JSON(http.StatusOK, user)
}
