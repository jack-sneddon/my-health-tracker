// internal/validation/weight/validator.go
package weight

import (
	"fmt"
	"math"
	"regexp"
	"time"

	"github.com/jack-sneddon/my-health-tracker/internal/models"
	"github.com/jack-sneddon/my-health-tracker/internal/validator"
)

// cmd/tracker/commands/weight/validator.go
// Add these constants and function:

const (
	WeightIDPattern = `^w\d{5}$`
)

// ValidationResult encapsulates the result of a validation operation
type ValidationResult struct {
	IsValid     bool
	Error       error
	Warnings    []string
	IsDuplicate bool
}

// Weight specific validation constants
const (
	MinWeight       = 75.0  // Minimum reasonable weight in pounds
	MaxWeight       = 250.0 // Maximum reasonable weight in pounds
	MaxWeightChange = 10.0  // Maximum reasonable weight change in pounds per day
)

// ValidateWeightID validates the format of weight record IDs
func ValidateWeightID(id string) error {
	if !regexp.MustCompile(WeightIDPattern).MatchString(id) {
		return fmt.Errorf("invalid weight record ID format. Must be 'w' followed by 5 digits")
	}
	return nil
}

// WeightRecord validates a complete weight record
func WeightRecord(record models.WeightRecord, lastRecord *models.WeightRecord) ValidationResult {
	result := ValidationResult{
		IsValid: true,
	}

	// Basic weight validation
	if err := validateWeightRange(record.Weight); err != nil {
		result.IsValid = false
		result.Error = err
		return result
	}

	// Date validation
	if lastRecord != nil {
		dateResult := validateDate(record.Date, lastRecord.Date)
		if !dateResult.IsValid {
			return dateResult
		}
		result.IsDuplicate = dateResult.IsDuplicate

		// Weight change validation
		changeResult := validateWeightChange(record.Weight, lastRecord.Weight)
		if !changeResult.IsValid {
			return changeResult
		}
		result.Warnings = append(result.Warnings, changeResult.Warnings...)
	}

	return result
}

// ValidationRequest encapsulates all data needed for validation
type ValidationRequest struct {
	Record     models.WeightRecord
	LastRecord *models.WeightRecord
	NextRecord *models.WeightRecord
}

// ValidationContext provides additional context for validation decisions
type ValidationContext struct {
	IsUpdate    bool
	AllowFuture bool
	StrictMode  bool
}

// ValidateWithContext performs validation with additional context
func ValidateWithContext(req ValidationRequest, ctx ValidationContext) ValidationResult {
	result := ValidationResult{
		IsValid: true,
	}

	// Base validations
	baseResult := WeightRecord(req.Record, req.LastRecord)
	if !baseResult.IsValid {
		return baseResult
	}
	result.Warnings = append(result.Warnings, baseResult.Warnings...)

	// Context-specific validations
	if ctx.IsUpdate {
		updateResult := validateUpdate(req)
		if !updateResult.IsValid {
			return updateResult
		}
		result.Warnings = append(result.Warnings, updateResult.Warnings...)
	}

	if !ctx.AllowFuture && isFutureDate(req.Record.Date) {
		result.IsValid = false
		result.Error = fmt.Errorf("future dates are not allowed")
		return result
	}

	return result
}

// Helper functions
func validateWeightRange(weight float64) error {
	if weight < MinWeight || weight > MaxWeight {
		return fmt.Errorf("weight must be between %.1f and %.1f pounds", MinWeight, MaxWeight)
	}
	return nil
}

func validateDate(current, last time.Time) ValidationResult {
	result := ValidationResult{
		IsValid: true,
	}

	if current.Equal(last) {
		result.IsDuplicate = true
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("record already exists for %s", current.Format(validator.DateFormat)))
	}

	if current.Before(last) {
		result.Warnings = append(result.Warnings,
			"this record predates the previous record")
	}

	return result
}

func validateWeightChange(current, last float64) ValidationResult {
	result := ValidationResult{
		IsValid: true,
	}

	change := math.Abs(current - last)
	if change > MaxWeightChange {
		result.Warnings = append(result.Warnings,
			fmt.Sprintf("weight change of %.1f lbs seems unusual", change))
	}

	return result
}

func validateUpdate(req ValidationRequest) ValidationResult {
	result := ValidationResult{
		IsValid: true,
	}

	if req.LastRecord != nil {
		change := math.Abs(req.Record.Weight - req.LastRecord.Weight)
		if change > MaxWeightChange {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("%.1f lb change from previous record seems unusual", change))
		}
	}

	if req.NextRecord != nil {
		change := math.Abs(req.NextRecord.Weight - req.Record.Weight)
		if change > MaxWeightChange {
			result.Warnings = append(result.Warnings,
				fmt.Sprintf("%.1f lb change to next record seems unusual", change))
		}
	}

	return result
}

func isFutureDate(date time.Time) bool {
	return date.After(time.Now())
}
