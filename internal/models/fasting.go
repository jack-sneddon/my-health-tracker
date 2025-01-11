// internal/models/fasting.go
package models

import (
	"fmt"
	"time"
)

type FastingRecord struct {
	Date            time.Time   `json:"date"`
	ExpectedPattern MealPattern `json:"expected_pattern"`
	ActualPattern   MealPattern `json:"actual_pattern"`
	Notes           string      `json:"notes,omitempty"`
}

func (f FastingRecord) GetDate() time.Time {
	return f.Date
}

func (f FastingRecord) Validate() error {
	switch f.ActualPattern {
	case FullFast, OneMeal, Regular:
		// valid pattern
	default:
		return fmt.Errorf("invalid meal pattern: %s", f.ActualPattern)
	}
	return nil
}

func (f FastingRecord) IsCompliant() bool {
	weekday := f.Date.Weekday()
	switch weekday {
	case time.Monday, time.Tuesday:
		return f.ActualPattern == FullFast
	case time.Wednesday, time.Thursday:
		return f.ActualPattern == OneMeal
	case time.Friday, time.Saturday, time.Sunday:
		return f.ActualPattern == Regular
	}
	return false
}
