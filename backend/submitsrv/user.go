package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	dbx "github.com/go-ozzo/ozzo-dbx"
)

// User maps the users table into a structure
type User struct {
	ID          string    `json:"id"`
	FirstName   string    `json:"firstName" db:"first_name"`
	LastName    string    `json:"lastName" db:"last_name"`
	Title       string    `json:"title"`
	Affiliation string    `json:"affiliation"  db:"university_affiliation"`
	Email       string    `json:"email"`
	Phone       string    `json:"phone"`
	Verified    bool      `json:"verified"`
	Admin       bool      `json:"admin"`
	CreatedAt   time.Time `db:"created_at" json:"-"`
	UpdatedAt   time.Time `db:"updated_at" json:"-"`
}

// TableName defines the expected DB table name that holds data for users
func (user *User) TableName() string {
	return "users"
}

// FindByEmail finds a user by email
func (user *User) FindByEmail(db *dbx.DB, email string) error {
	q := db.NewQuery("select id,last_name,first_name,email,title,university_affiliation,phone from users where email={:email} limit 1")
	q.Bind(dbx.Params{"email": email})
	return q.One(user)
}

// Create creates a user record in the DB based in data in the struct
func (user *User) Create(db *dbx.DB) error {
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	return db.Model(user).Insert()
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

// CreateUser creates a new user record from data in the post form
func (svc *ServiceContext) CreateUser(c *gin.Context) {
	c.String(http.StatusNotImplemented, "not implememted")
}
