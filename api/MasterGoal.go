package api

import (
	"time"
)

type MasterGoal struct {
	MasterGoalID int
	Name         string
	Description  string
	ImageURL     string
	Status       int
	DateAdded    time.Time
}
