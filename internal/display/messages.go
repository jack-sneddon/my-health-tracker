// internal/display/messages.go
package display

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"github.com/jack-sneddon/my-health-tracker/cmd/tracker/commands/result"
	"github.com/jack-sneddon/my-health-tracker/internal/models"
	"github.com/jack-sneddon/my-health-tracker/internal/validator"
)

var (
	// Color formatters
	errorColor   = color.New(color.FgRed, color.Bold)
	warningColor = color.New(color.FgYellow)
	successColor = color.New(color.FgGreen)
	infoColor    = color.New(color.FgCyan)
	headerColor  = color.New(color.FgWhite, color.Bold)
)

type ConfirmationResult struct {
	Confirmed bool
	Error     error
}

// ShowError displays formatted error messages
func ShowError(format string, args ...interface{}) {
	errorColor.Printf("✗ Error: "+format+"\n", args...)
}

// ShowWarning displays formatted warning messages
func ShowWarning(format string, args ...interface{}) {
	warningColor.Printf("⚠ Warning: "+format+"\n", args...)
}

// ShowSuccess displays formatted success messages
func ShowSuccess(format string, args ...interface{}) {
	successColor.Printf("✓ "+format+"\n", args...)
}

// ShowInfo displays formatted informational messages
func ShowInfo(format string, args ...interface{}) {
	infoColor.Printf("ℹ "+format+"\n", args...)
}

// ConfirmAction prompts for user confirmation
func ConfirmAction(prompt string) ConfirmationResult {
	fmt.Printf("\n%s (y/N): ", prompt)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return ConfirmationResult{false, err}
	}
	return ConfirmationResult{
		Confirmed: strings.ToLower(response) == "y",
		Error:     nil,
	}
}

// ShowWeightRecord displays a formatted weight record
func ShowWeightRecord(id string, date, weight string, notes string) {
	headerColor.Println("\nWeight Record:")
	fmt.Printf("  ID:     %s\n", id)
	fmt.Printf("  Date:   %s\n", date)
	fmt.Printf("  Weight: %s lbs\n", weight)
	if notes != "" {
		fmt.Printf("  Notes:  %s\n", notes)
	}
}

// handle record details
func ShowDeleteConfirmation(id, date, weight, notes string) ConfirmationResult {
	headerColor.Println("\nDelete Confirmation:")
	fmt.Printf("  ID:     %s\n", id)
	fmt.Printf("  Date:   %s\n", date)
	fmt.Printf("  Weight: %s lbs\n", weight)
	if notes != "" {
		fmt.Printf("  Notes:  %s\n", notes)
	}
	return ConfirmAction("Are you sure you want to delete this record?")
}

func ShowCommandResult(result result.CommandResult) {
	if !result.Success {
		ShowError(result.Error.Error())
		for _, warning := range result.Warnings {
			ShowWarning(warning)
		}
		return
	}

	for _, msg := range result.Messages {
		ShowSuccess(msg)
	}

	if weightRecord, ok := result.Data.(models.WeightRecord); ok {
		ShowWeightRecord(
			weightRecord.ID,
			weightRecord.Date.Format(validator.DateFormat),
			fmt.Sprintf("%.1f", weightRecord.Weight),
			weightRecord.Notes,
		)
	}

	if weightRecord, ok := result.Data.(models.WeightRecord); ok {
		ShowWeightRecord(
			weightRecord.ID,
			weightRecord.Date.Format(validator.DateFormat),
			fmt.Sprintf("%.1f", weightRecord.Weight),
			weightRecord.Notes,
		)
	} else if exerciseRecord, ok := result.Data.(models.ExerciseRecord); ok {
		ShowExerciseRecord(
			exerciseRecord.Date.Format(validator.DateFormat),
			string(exerciseRecord.Activity),
			exerciseRecord.OtherActivity,
			exerciseRecord.Duration,
			exerciseRecord.Notes,
			exerciseRecord.Completed,
		)
	}
}

func ShowTable(headers []string, rows [][]string) {
	// Calculate column widths
	widths := make([]int, len(headers))
	for i, h := range headers {
		widths[i] = len(h)
	}
	for _, row := range rows {
		for i, cell := range row {
			if len(cell) > widths[i] {
				widths[i] = len(cell)
			}
		}
	}

	// Print headers
	headerColor.Printf("\n")
	for i, header := range headers {
		headerColor.Printf("%-*s", widths[i]+2, header)
	}
	fmt.Println()

	// Print separator
	totalWidth := 0
	for _, w := range widths {
		totalWidth += w + 2
	}
	fmt.Println(strings.Repeat("-", totalWidth))

	// Print rows
	for _, row := range rows {
		for i, cell := range row {
			fmt.Printf("%-*s", widths[i]+2, cell)
		}
		fmt.Println()
	}
	fmt.Println()
}

func ShowStats(stats map[string]string) {
	headerColor.Printf("\nSummary:\n")
	maxKeyLength := 0
	for k := range stats {
		if len(k) > maxKeyLength {
			maxKeyLength = len(k)
		}
	}

	for k, v := range stats {
		fmt.Printf("%-*s: %s\n", maxKeyLength, k, v)
	}
}

func ShowHeader(text string) {
	headerColor.Printf("\n%s\n", text)
}

// internal/display/messages.go
func ShowWeightList(records []models.WeightRecord) {
	fmt.Printf("%-8s  %-10s  %-7s  %s\n", "ID", "Date", "Weight", "Notes")
	fmt.Println(strings.Repeat("-", 60))

	for _, record := range records {
		fmt.Printf("%-8s  %-10s  %7.1f  %s\n",
			record.ID,
			record.Date.Format(validator.DateFormat),
			record.Weight,
			record.Notes)
	}
	fmt.Println()
}

// ShowExerciseRecord displays a formatted exercise record
func ShowExerciseRecord(date string, activity string, otherActivity string, duration int, notes string, completed bool) {
	headerColor.Println("\nExercise Record:")
	fmt.Printf("  Date:       %s\n", date)
	if activity == "other" {
		fmt.Printf("  Activity:   %s (%s)\n", activity, otherActivity)
	} else {
		fmt.Printf("  Activity:   %s\n", activity)
	}
	fmt.Printf("  Duration:   %d minutes\n", duration)
	if notes != "" {
		fmt.Printf("  Notes:      %s\n", notes)
	}
	fmt.Printf("  Completed:  %v\n", completed)
}

func ShowExerciseList(records []models.ExerciseRecord) {
	headerColor.Printf("\n%-12s %-15s %-8s %-30s %s\n",
		"Date",
		"Activity",
		"Duration",
		"Notes",
		"Completed")
	fmt.Println(strings.Repeat("-", 80))

	for _, record := range records {
		// Format activity string
		activityStr := string(record.Activity)
		if record.Activity == models.Other {
			activityStr = fmt.Sprintf("%s (%s)", record.Activity, record.OtherActivity)
		}

		fmt.Printf("%-12s %-15s %-8d %-30s %v\n",
			record.Date.Format(validator.DateFormat),
			activityStr,
			record.Duration,
			truncateString(record.Notes, 30),
			record.Completed)
	}
	fmt.Println()
}

// Helper function to truncate strings for display
func truncateString(str string, length int) string {
	if len(str) <= length {
		return str
	}
	return str[:length-3] + "..."
}

func ShowExerciseDeleteConfirmation(date string, activity string, otherActivity string, duration int, notes string, completed bool) ConfirmationResult {
	headerColor.Println("\nDelete Confirmation:")
	fmt.Printf("  Date:       %s\n", date)
	if activity == "other" {
		fmt.Printf("  Activity:   %s (%s)\n", activity, otherActivity)
	} else {
		fmt.Printf("  Activity:   %s\n", activity)
	}
	fmt.Printf("  Duration:   %d minutes\n", duration)
	if notes != "" {
		fmt.Printf("  Notes:      %s\n", notes)
	}
	fmt.Printf("  Completed:  %v\n", completed)
	return ConfirmAction("Are you sure you want to delete this record?")
}
