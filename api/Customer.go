package api

import (
	"time"
)

type Customer struct {
	CustomerID     int64
	CompanyName    string
	FirstName      string
	LastName       string
	Address1       string
	Address2       string
	City           string
	State          string
	ZipCode        string
	Country        string
	Province       string
	PostalCode     string
	Email          string
	MainPhone      string
	AlternatePhone string
	DateAdded      time.Time
	DateUpdated    time.Time
}
