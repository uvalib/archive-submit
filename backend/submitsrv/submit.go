package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"
)

// Submit accepts a transfer submission, creates a DB record and kicks off submission
// processing. When complete, an email is sent to the submitter
func (svc *ServiceContext) Submit(c *gin.Context) {
	var accession Accession
	err := c.ShouldBindJSON(&accession)
	if err != nil {
		log.Printf("ERROR: Unable to parse request: %s", err.Error())
		c.String(http.StatusBadRequest, err.Error())
		return
	}
	log.Printf("Received: %+v", accession)

	log.Printf("Update existing user %d:%s", accession.User.ID, accession.User.Email)
	accession.User.UpdatedAt = time.Now()
	accession.User.FormatPhone()
	err = svc.DB.Model(&accession.User).Exclude("Verified", "VerifyToken", "Admin", "CreatedAt", "email").Update()
	if err != nil {
		log.Printf("WARN: Unable to update %s - %s", accession.User.Email, err.Error())
	}

	log.Printf("Add new accession record")
	tx, _ := svc.DB.Begin()
	accession.UserID = accession.User.ID
	accession.CreatedAt = time.Now()
	err = tx.Model(&accession).Insert()
	if err != nil {
		log.Printf("ERROR: Unable to add accession %s", err.Error())
		tx.Rollback()
		c.String(http.StatusInternalServerError, "Unable to create accession record")
		return
	}

	accession.WriteGenres(tx)
	if accession.PhysicalTransfer {
		perr := accession.WritePhysicalTransfer(tx)
		if perr != nil {
			log.Printf("ERROR: Unable to write physical xfer: %s", perr.Error())
			tx.Rollback()
			c.String(http.StatusInternalServerError, "Unable to create physical transfer record")
			return
		}
	}
	if accession.DigitalTransfer {
		derr := accession.WriteDigitalTransfer(tx)
		if derr != nil {
			log.Printf("ERROR: Unable to write digital xfer: %s", derr.Error())
			tx.Rollback()
			c.String(http.StatusInternalServerError, "Unable to create digital transfer record")
			return
		}

		// Submitted dir gets broken up by YYYY/MM before the submissionID. Make the tree:
		currTime := time.Now()
		xferDir := fmt.Sprintf("%s/%s/%s", svc.UploadDir, "transferred", currTime.Format("2006/01"))
		os.MkdirAll(xferDir, 0777)

		// Move pending into transfer tree
		pendingDir := fmt.Sprintf("%s/%s", svc.UploadDir, "pending")
		uploadDir := fmt.Sprintf("%s/%s", pendingDir, accession.Identifier)
		tgtDir := fmt.Sprintf("%s/%s", xferDir, accession.Identifier)
		log.Printf("Moving pending upload files from %s to %s", uploadDir, tgtDir)
		err = os.Rename(uploadDir, tgtDir)
		if err != nil {
			log.Printf("WARN: Unable to move pending files to submitted: %s", err.Error())
		}
	}
	tx.Commit()

	// Now send recepit to submitter and admins
	accession.User.SendReceiptEmail(svc.DB, svc.SMTP, accession)
	c.String(http.StatusOK, "accepted")
}

// GetAccessionIdentifier will generate an unique token to identify digital content uploads
// It will be used as a storage subdir for files as they are uploaded
func (svc *ServiceContext) GetAccessionIdentifier(c *gin.Context) {
	id := xid.New()
	c.String(http.StatusOK, id.String())
}
