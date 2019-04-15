package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

// Authenticate will authenticate a user based on Shibboleth headers, then redirect to other pages for content
// NOTE: this is called directly from the front-end as a transient page with  window.location.href = "/authenticate"
// to force it through NetBadge authentication. This route does nothing more than ensure users have been thru
// authentication and create a user if one does not exist. End result is a redirect.
func (svc *ServiceContext) Authenticate(c *gin.Context) {
	log.Printf("Checking authentication headers...")
	computingID := c.GetHeader("remote_user")
	if svc.DevAuthUser != "" {
		computingID = svc.DevAuthUser
	}
	if computingID == "" {
		log.Printf("ERROR: Expected auth header not present in request. Not authorized.")
		c.Redirect(http.StatusFound, "/forbidden")
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
	json, _ := json.Marshal(user)
	log.Printf("Authenticated user: %s", string(json))
	c.SetCookie("archives_xfer_user", string(json), 3600, "/", "", false, false)
	c.Redirect(http.StatusFound, "/submit")
}
