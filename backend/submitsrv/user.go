package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// User maps the users table into a structure
type User struct {
	ID          int       `json:"id"`
	FirstName   string    `json:"firstName" binding:"required" db:"first_name"`
	LastName    string    `json:"lastName" binding:"required" db:"last_name"`
	Title       string    `json:"title" binding:"required" form:"title"`
	Affiliation string    `json:"affiliation"  binding:"required" db:"university_affiliation"`
	Email       string    `json:"email" binding:"required" form:"email"`
	Phone       string    `json:"phone" binding:"required" form:"phone"`
	Verified    bool      `json:"verified"`
	VerifyToken string    `json:"token"  db:"verify_token" `
	Admin       bool      `json:"-"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}

// TableName defines the expected DB table name that holds data for users
func (user *User) TableName() string {
	return "users"
}

// FormatPhone will ensure phone numbers are in standard format (no spaces)
func (user *User) FormatPhone() {
	user.Phone = strings.Replace(user.Phone, " ", "", -1)
}

// IsValid makes sure all fields are set and look right
func (user *User) IsValid() bool {
	if user.FirstName == "" || user.LastName == "" || user.Title == "" ||
		user.Affiliation == "" || user.Email == "" || user.Phone == "" {
		return false
	}
	return true
}

// FullName return the full name of the user
func (user *User) FullName() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}

// SendReceiptEmail will send the user (and admins) a transfer receipt email
func (user *User) SendReceiptEmail(smtpCfg SMTPConfig, accession Accession, bcc []string) {
	log.Printf("Rendering receipt email body")
	var renderedEmail bytes.Buffer
	tpl := template.Must(template.ParseFiles("templates/receipt_email.html"))
	err := tpl.Execute(&renderedEmail, accession)
	if err != nil {
		log.Printf("ERROR: Unable to render receipt email: %s", err.Error())
		return
	}

	log.Printf("Generate SMTP message")
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: UVA Archives Transfer Receipt\n"
	toHdr := fmt.Sprintf("To: %s\n", accession.User.Email)
	BCC := fmt.Sprintf("Bcc: %s\n", strings.Join(bcc, ","))
	msg := []byte(subject + toHdr + BCC + mime + renderedEmail.String())

	if smtpCfg.DevMode {
		log.Printf("Email is in dev mode. Logging message instead of sending")
		log.Printf("==================================================")
		log.Printf("%s", msg)
		log.Printf("==================================================")
	} else {
		log.Printf("Send verify email to %s", accession.User.Email)
		to := []string{user.Email}
		err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpCfg.Host, smtpCfg.Port), nil, "no-reply@virginia.edu", to, msg)
		if err != nil {
			log.Printf("ERROR: Unable to send receipt email: %s", err.Error())
			return
		}
	}
}

// SendVerifyEmail will send a verify email to a new user
func (user *User) SendVerifyEmail(baseURL string, smtpCfg SMTPConfig) {
	log.Printf("Rendering verification email body")
	var renderedEmail bytes.Buffer
	var data struct {
		Name string
		URL  string
	}
	data.Name = user.FullName()
	data.URL = fmt.Sprintf("https://%s/verify/%s", baseURL, user.VerifyToken)
	tpl := template.Must(template.ParseFiles("templates/verify_email.html"))
	err := tpl.Execute(&renderedEmail, data)
	if err != nil {
		log.Printf("ERROR: Unable to render verify email: %s", err.Error())
		return
	}

	log.Printf("Generate SMTP message")
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"
	subject := "Subject: UVA Archives Transfer Verification\n"
	to := fmt.Sprintf("To: %s\n", user.Email)
	msg := []byte(subject + to + mime + renderedEmail.String())

	if smtpCfg.DevMode {
		log.Printf("Email is in dev mode. Logging message instead of sending")
		log.Printf("==================================================")
		log.Printf("%s", msg)
		log.Printf("==================================================")
	} else {
		log.Printf("Send verify email to %s", user.Email)
		to := []string{user.Email}
		err := smtp.SendMail(fmt.Sprintf("%s:%d", smtpCfg.Host, smtpCfg.Port), nil, "no-reply@virginia.edu", to, msg)
		if err != nil {
			log.Printf("ERROR: Unable to send verify email: %s", err.Error())
			return
		}
	}
}

// FindByEmail finds a user by email
func (user *User) FindByEmail(db *dbx.DB, email string) error {
	q := db.NewQuery(`select * from users where email={:email} limit 1`)
	q.Bind(dbx.Params{"email": email})
	return q.One(user)
}

// FindByToken finds a user by verfy_token
func (user *User) FindByToken(db *dbx.DB, token string) error {
	q := db.NewQuery(`select * from users where verify_token={:token} limit 1`)
	q.Bind(dbx.Params{"token": token})
	return q.One(user)
}

// Create creates a user record in the DB based in data in the struct
func (user *User) Create(db *dbx.DB) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.FormatPhone()
	return db.Model(user).Insert()
}

// Verify will mark this user account as verified
func (user *User) Verify(db *dbx.DB) error {
	user.Verified = true
	user.UpdatedAt = time.Now()
	return db.Model(user).Update()
}

// UserSearch will find users by a variety of search tearms.
func (svc *ServiceContext) UserSearch(c *gin.Context) {
	// Note: currently support just email lookup
	email := c.Query("email")
	if email == "" {
		c.String(http.StatusBadRequest, "missing required email query param")
		return
	}
	log.Printf("Checking for presence of user with email %s", email)
	user := User{}
	err := user.FindByEmail(svc.DB, email)
	if err != nil {
		log.Printf("%s not found: %s", email, err.Error())
		c.String(http.StatusNotFound, "%s not found", email)
		return
	}
	log.Printf("%s found", email)
	c.JSON(http.StatusOK, user)
}

// CreateUser creates a new user record from data in the post form. Mark
func (svc *ServiceContext) CreateUser(c *gin.Context) {
	var user User
	c.Bind(&user)
	if user.IsValid() == false {
		log.Printf("ERROR: Request to create user with missing fields")
		c.String(http.StatusBadRequest, "All fields are required")
		return
	}
	user.VerifyToken = xid.New().String()
	err := user.Create(svc.DB)
	if err != nil {
		log.Printf("ERROR: User create failed: %s", err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	user.SendVerifyEmail(svc.Hostname, svc.SMTP)

	c.JSON(http.StatusOK, user)
}

// VerifyUser accepts a token, finds the assiciated user and marks them as validated
func (svc *ServiceContext) VerifyUser(c *gin.Context) {
	token := c.Param("token")
	log.Printf("Verify user with token [%s]", token)
	user := User{}
	err := user.FindByToken(svc.DB, token)
	if err != nil {
		log.Printf("ERROR: No user found for %s: %s", token, err.Error())
		c.String(http.StatusNotFound, err.Error())
		return
	}
	if user.Verified == false {
		log.Printf("Marking %s as verified", user.Email)
		err = user.Verify(svc.DB)
		if err != nil {
			log.Printf("Unable to verify %s: %s", user.Email, err.Error())
			c.String(http.StatusInternalServerError, err.Error())
			return
		}
	} else {
		log.Printf("User %s already verified; nothing to do.", user.Email)
	}
	c.JSON(http.StatusOK, user)
}

// ResendVerification accepts a token, finds the associated user, and resends the validation email
func (svc *ServiceContext) ResendVerification(c *gin.Context) {
	var data struct{ Token string }
	c.Bind(&data)
	log.Printf("Resend verification for [%s]", data.Token)
	user := User{}
	err := user.FindByToken(svc.DB, data.Token)
	if err != nil {
		log.Printf("ERROR: No user found for %s: %s", data.Token, err.Error())
		c.String(http.StatusNotFound, err.Error())
		return
	}
	log.Printf("[%s] found; resending verification email", data.Token)
	user.SendVerifyEmail(svc.Hostname, svc.SMTP)
	c.String(http.StatusOK, "email resent")
}
