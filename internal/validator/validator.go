// internal/validator/validator.go
package validator

import (
	"fmt"
	"time"
)

const (
	DateFormat    = "2006-01-02" // Go's reference date format for YYYY-MM-DD
	MaxNoteLength = 500          // Maximum characters for notes
)

// ParseDate converts string to time.Time and validates format
func ParseDate(date string) (time.Time, error) {
	if date == "" {
		return time.Now(), nil // Default to current date
	}

	parsedDate, err := time.Parse(DateFormat, date)
	if err != nil {
		return time.Time{}, fmt.Errorf("invalid date format. Use YYYY-MM-DD")
	}

	// Don't allow future dates
	if parsedDate.After(time.Now()) {
		return time.Time{}, fmt.Errorf("future dates are not allowed")
	}

	return parsedDate, nil
}

// ValidateNotes checks the notes field
func ValidateNotes(notes string) error {
	if len(notes) > MaxNoteLength {
		return fmt.Errorf("notes must be %d characters or less (current: %d)",
			MaxNoteLength, len(notes))
	}
	return nil
}

// GetDefaultDateRange returns default date range (last 30 days)
func GetDefaultDateRange() (time.Time, time.Time) {
	now := time.Now()
	from := now.AddDate(0, 0, -30)
	return from, now
}

// ValidateDateRange checks if date range is valid
func ValidateDateRange(from, to string) (time.Time, time.Time, error) {
	var fromDate, toDate time.Time
	var err error

	// Parse 'from' date
	if from != "" {
		fromDate, err = ParseDate(from)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid 'from' date: %w", err)
		}
	}

	// Parse 'to' date
	if to != "" {
		toDate, err = ParseDate(to)
		if err != nil {
			return time.Time{}, time.Time{}, fmt.Errorf("invalid 'to' date: %w", err)
		}
	}

	// If only one date is provided, set reasonable defaults
	if from == "" {
		fromDate = toDate.AddDate(0, 0, -30)
	}
	if to == "" {
		toDate = time.Now()
	}

	// Ensure 'from' is before 'to'
	if fromDate.After(toDate) {
		return time.Time{}, time.Time{}, fmt.Errorf("'from' date must be before 'to' date")
	}

	return fromDate, toDate, nil
}
