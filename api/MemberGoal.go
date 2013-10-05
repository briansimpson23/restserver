package api

import (
	"time"
)

type MemberGoal struct {
	MemberID      int64
	GoalID        int
	MasterGoalID  int
	GoalName      string
	GoalStatus    int
	DateAdded     time.Time
	DateCompleted time.Time
}
