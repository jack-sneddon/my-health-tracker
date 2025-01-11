// internal/display/messages.go
package display

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
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

// ShowDeleteConfirmation displays a formatted delete confirmation
func ShowDeleteConfirmation(record string, context string) ConfirmationResult {
	headerColor.Println("\nDelete Confirmation:")
	fmt.Printf("%s\n", record)
	if context != "" {
		fmt.Printf("\nContext:\n%s\n", context)
	}
	return ConfirmAction("Are you sure you want to delete this record?")
}
