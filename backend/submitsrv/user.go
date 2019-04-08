package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/rs/xid"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// User maps the users table into a structure
type User struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"firstName" db:"first_name" form:"fname"`
	LastName    string    `json:"lastName" db:"last_name" form:"lname"`
	Title       string    `json:"title" form:"title"`
	Affiliation string    `json:"affiliation"  db:"university_affiliation" form:"affiliation"`
	Email       string    `json:"email" form:"email"`
	Phone       string    `json:"phone" form:"phone"`
	Verified    bool      `json:"verified"`
	VerifyToken string    `json:"token"  db:"verify_token" `
	Admin       bool      `json:"-"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
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

// TableName defines the expected DB table name that holds data for users
func (user *User) TableName() string {
	return "users"
}

// FindByEmail finds a user by email
func (user *User) FindByEmail(db *dbx.DB, email string) error {
	q := db.NewQuery(`select id,last_name,first_name,email,title,university_affiliation,
		phone,verified,verify_token from users where email={:email} limit 1`)
	q.Bind(dbx.Params{"email": email})
	return q.One(user)
}

// FindByToken finds a user by verfy_token
func (user *User) FindByToken(db *dbx.DB, token string) error {
	q := db.NewQuery(`select id,last_name,first_name,email,title,university_affiliation,
		phone,verified,verify_token from users where verify_token={:token} limit 1`)
	q.Bind(dbx.Params{"token": token})
	return q.One(user)
}

// Create creates a user record in the DB based in data in the struct
func (user *User) Create(db *dbx.DB) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
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
	err = user.Verify(svc.DB)
	if err != nil {
		log.Printf("Unable to verify %s: %s", user.Email, err.Error())
		c.String(http.StatusInternalServerError, err.Error())
		return
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
