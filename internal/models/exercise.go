// internal/models/exercise.go
package models

import (
	"fmt"
	"time"
)

type ActivityType string

const (
	Jogging      ActivityType = "jogging"
	Skiing       ActivityType = "skiing"
	Walking      ActivityType = "walking"
	Cycling      ActivityType = "cycling"
	MountainBike ActivityType = "mountain_biking"
	Pickleball   ActivityType = "pickleball"
	Other        ActivityType = "other"
)

type ExerciseRecord struct {
	Date          time.Time    `json:"date"`
	Activity      ActivityType `json:"activity"`
	OtherActivity string       `json:"other_activity,omitempty"` // name of activity if Activity is Other
	Duration      int          `json:"duration"`                 // in minutes
	Notes         string       `json:"notes,omitempty"`
	Completed     bool         `json:"completed"`
}

func (e ExerciseRecord) Validate() error {
	// Validate activity type
	validActivity := false
	for _, a := range []ActivityType{Jogging, Skiing, Walking, Cycling, MountainBike, Pickleball, Other} {
		if e.Activity == a {
			validActivity = true
			break
		}
	}
	if !validActivity {
		return fmt.Errorf("invalid activity type: %s", e.Activity)
	}

	// If activity is Other, ensure OtherActivity is provided
	if e.Activity == Other && e.OtherActivity == "" {
		return fmt.Errorf("other_activity must be specified when activity type is 'other'")
	}

	if e.Duration <= 0 {
		return fmt.Errorf("duration must be greater than 0")
	}
	if e.Duration > 480 { // 8 hours - reasonable upper limit
		return fmt.Errorf("duration seems unreasonably high")
	}

	return nil
}

func (e ExerciseRecord) IsCompliant() bool {
	return e.Completed && e.Duration >= 45 // 45 minutes daily goal
}

func (e ExerciseRecord) GetDate() time.Time {
	return e.Date
}
