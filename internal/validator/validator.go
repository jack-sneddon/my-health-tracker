// internal/validator/validator.go
package validator

import (
	"fmt"
	"math"
	"regexp"
	"time"

	"github.com/jack-sneddon/my-health-tracker/internal/models"
)

const (
	DateFormat      = "2006-01-02" // Go's reference date format for YYYY-MM-DD
	MinWeight       = 75.0         // Minimum reasonable weight in pounds
	MaxWeight       = 250.0        // Maximum reasonable weight in pounds
	MaxWeightChange = 10.0         // Maximum reasonable weight change in pounds per day
	MaxNoteLength   = 500          // Maximum characters for notes
	WeightIDPattern = `^w\d{5}$`
)

// WeightValidationResult contains validation result and any warnings
type WeightValidationResult struct {
	IsValid     bool
	Error       error
	HasWarning  bool
	Warning     string
	IsDuplicate bool
}

// ValidateWeight checks if weight is within reasonable bounds
func ValidateWeight(weight float64) error {
	if weight < MinWeight || weight > MaxWeight {
		return fmt.Errorf("weight must be between %.1f and %.1f pounds", MinWeight, MaxWeight)
	}
	return nil
}

// ValidateWeightChange checks if weight change is reasonable and returns warnings
func ValidateWeightChange(current float64, lastWeight *models.WeightRecord) WeightValidationResult {
	result := WeightValidationResult{
		IsValid: true,
	}

	if lastWeight == nil {
		return result
	}

	change := math.Abs(current - lastWeight.Weight)
	if change > MaxWeightChange {
		result.HasWarning = true
		result.Warning = fmt.Sprintf("Warning: Weight change of %.1f lbs since last record (%.1f) seems unusual",
			change, lastWeight.Weight)
	}

	return result
}

// ValidateWeightDate checks for duplicate dates and chronological order
func ValidateWeightDate(date time.Time, lastWeight *models.WeightRecord) WeightValidationResult {
	result := WeightValidationResult{
		IsValid: true,
	}

	if lastWeight == nil {
		return result
	}

	// Check for duplicate date
	if date.Equal(lastWeight.Date) {
		result.IsDuplicate = true
		result.HasWarning = true
		result.Warning = fmt.Sprintf("Weight record already exists for %s", date.Format(DateFormat))
	}

	// Check chronological order
	if date.Before(lastWeight.Date) {
		result.HasWarning = true
		result.Warning = "Warning: This record predates the previous record"
	}

	return result
}

// ValidateWeightID validates the format of weight record IDs
func ValidateWeightID(id string) error {
	if !regexp.MustCompile(WeightIDPattern).MatchString(id) {
		return fmt.Errorf("invalid weight record ID format. Must be 'w' followed by 5 digits")
	}
	return nil
}

// ValidateNotes checks the notes field
func ValidateNotes(notes string) error {
	if len(notes) > MaxNoteLength {
		return fmt.Errorf("notes must be %d characters or less (current: %d)",
			MaxNoteLength, len(notes))
	}
	return nil
}

// ValidateWeightUpdate checks if an update would create inconsistent history
func ValidateWeightUpdate(recordID string, newWeight float64, prevRecord, nextRecord *models.WeightRecord) WeightValidationResult {
	result := WeightValidationResult{
		IsValid: true,
	}

	// Check changes against previous record
	if prevRecord != nil {
		change := math.Abs(newWeight - prevRecord.Weight)
		if change > MaxWeightChange {
			result.HasWarning = true
			result.Warning += fmt.Sprintf("Warning: %.1f lb change from previous record seems unusual\n", change)
		}
	}

	// Check changes against next record
	if nextRecord != nil {
		change := math.Abs(nextRecord.Weight - newWeight)
		if change > MaxWeightChange {
			result.HasWarning = true
			result.Warning += fmt.Sprintf("Warning: %.1f lb change to next record seems unusual\n", change)
		}
	}

	return result
}

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

// ValidateDateRange checks if date range is valid
func ValidateDateRange(from, to string) (time.Time, time.Time, error) {
	var fromDate, toDate time.Time
	var err error

	// If no dates provided, default to last 30 days
	if from == "" && to == "" {
		toDate = time.Now()
		fromDate = toDate.AddDate(0, 0, -30)
		return fromDate, toDate, nil
	}

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

func GetDefaultDateRange() (time.Time, time.Time) {
	now := time.Now()
	// Default to last 30 days
	from := now.AddDate(0, 0, -30)
	return from, now
}
