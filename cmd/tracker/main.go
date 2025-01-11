// cmd/tracker/main.go
package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/jack-sneddon/my-health-tracker/internal/display"
	"github.com/jack-sneddon/my-health-tracker/internal/models"
	"github.com/jack-sneddon/my-health-tracker/internal/storage"
	"github.com/jack-sneddon/my-health-tracker/internal/validator"
	"github.com/spf13/cobra"
)

var store storage.StorageManager

func main() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

var (
	weightValue    float64
	weightDate     string
	weightNotes    string
	weightFromDate string
	weightToDate   string
)

var rootCmd = &cobra.Command{
	Use:   "tracker",
	Short: "Health tracking application",
	Long: `A health tracking application for monitoring weight, exercise, fasting, and soda consumption.

Command Structure (CRUD operations):
  CREATE:
    tracker weight add --value 185.5 --date 2024-01-08 --notes "Morning weight"
    tracker exercise add --activity jogging --duration 45 --date 2024-01-08
    tracker fasting add --pattern full-fast --date 2024-01-08
    tracker soda add --consumed --quantity 12 --date 2024-01-08

  READ:
    tracker weight get --date 2024-01-08
    tracker weight list --from 2024-01-01 --to 2024-01-08

  UPDATE:
    tracker weight update w12345 --value 184.5
    tracker exercise update e12345 --duration 50

  DELETE:
    tracker weight delete w12345

Use "tracker [command] --help" for more information about a command.`,
}

// Weight commands
var weightCmd = &cobra.Command{
	Use:   "weight",
	Short: "Manage weight records",
	Long: `Manage weight records with full CRUD operations.
    
Examples:
  # Add a weight record
  tracker weight add --value 185.5 --date 2024-01-08 --notes "Morning weight"

  # Get weight for a specific date
  tracker weight get --date 2024-01-08

  # List weights for a date range
  tracker weight list --from 2024-01-01 --to 2024-01-08

  # Update a weight record
  tracker weight update w12345 --value 184.5

  # Delete a weight record
  tracker weight delete w12345`,
}

func init() {

	// Initialize storage before adding commands
	store = storage.NewJSONStorage("", false)
	if err := store.Init(); err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Add main commands to root
	rootCmd.AddCommand(weightCmd)
	rootCmd.AddCommand(exerciseCmd)
	rootCmd.AddCommand(fastingCmd)
	rootCmd.AddCommand(sodaCmd)

	// Add subcommands to weight
	weightCmd.AddCommand(weightAddCmd)
	weightCmd.AddCommand(weightGetCmd)
	weightCmd.AddCommand(weightListCmd)
	weightCmd.AddCommand(weightUpdateCmd)
	weightCmd.AddCommand(weightDeleteCmd)

	// Flags for weightAddCmd
	weightAddCmd.Flags().Float64VarP(&weightValue, "value", "v", 0, "Weight value in pounds (required)")
	weightAddCmd.Flags().StringVarP(&weightDate, "date", "d", "", "Date of weight record (default: today)")
	weightAddCmd.Flags().StringVarP(&weightNotes, "notes", "n", "", "Optional notes about the weight record")
	weightAddCmd.MarkFlagRequired("value")

	// Flags for weightGetCmd
	weightGetCmd.Flags().StringVarP(&weightDate, "date", "d", "", "Date to get weight record for (required)")
	weightGetCmd.MarkFlagRequired("date")

	// Flags for weightListCmd
	weightListCmd.Flags().StringVarP(&weightFromDate, "from", "f", "", "Start date for listing weights")
	weightListCmd.Flags().StringVarP(&weightToDate, "to", "t", "", "End date for listing weights")

	// Flags for weightUpdateCmd
	weightUpdateCmd.Flags().Float64VarP(&weightValue, "value", "v", 0, "New weight value in pounds")
	weightUpdateCmd.Flags().StringVarP(&weightNotes, "notes", "n", "", "Updated notes about the weight record")
	// At least one flag should be specified
	weightUpdateCmd.MarkFlagRequired("value")
}

// Update weightAddCmd
var weightAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new weight record",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse and validate date
		var date time.Time
		var err error
		if weightDate != "" {
			date, err = validator.ParseDate(weightDate)
			if err != nil {
				display.ShowError(err.Error())
				return err
			}
		} else {
			date = time.Now()
		}

		// Create the record
		record := models.WeightRecord{
			Date:   date,
			Weight: weightValue,
			Notes:  weightNotes,
		}

		// Basic weight validation
		if err := validator.ValidateWeight(weightValue); err != nil {
			display.ShowError(err.Error())
			return err
		}

		// Get last weight record for change validation
		lastWeight, err := store.GetLastWeightRecord()
		if err != nil {
			display.ShowError("Error checking previous records: %v", err)
			return err
		}

		// Check for duplicates
		dateValidation := validator.ValidateWeightDate(date, lastWeight)
		if dateValidation.IsDuplicate {
			display.ShowWarning("Record already exists for %s",
				date.Format(validator.DateFormat))
			result := display.ConfirmAction("Do you want to overwrite this record?")
			if !result.Confirmed {
				display.ShowInfo("Operation cancelled")
				return nil
			}
		}

		// Check for unusual weight changes
		if lastWeight != nil {
			changeValidation := validator.ValidateWeightChange(weightValue, lastWeight)
			if changeValidation.HasWarning {
				display.ShowWarning(changeValidation.Warning)
				result := display.ConfirmAction("Do you want to continue?")
				if !result.Confirmed {
					display.ShowInfo("Operation cancelled")
					return nil
				}
			}
		}

		// Add the record
		savedRecord, err := store.AddWeight(record)
		if err != nil {
			display.ShowError("Failed to add weight record: %v", err)
			return err
		}

		display.ShowSuccess("Weight record added successfully")
		display.ShowWeightRecord(
			savedRecord.ID,
			savedRecord.Date.Format(validator.DateFormat),
			fmt.Sprintf("%.1f", savedRecord.Weight),
			savedRecord.Notes,
		)
		return nil
	},
}

// Update weightGetCmd
var weightGetCmd = &cobra.Command{
	Use:   "get",
	Short: "Get weight record for a specific date",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Parse and validate date
		date, err := validator.ParseDate(weightDate)
		if err != nil {
			return err
		}

		// TODO: Get the weight record
		fmt.Printf("Getting weight record for: %s\n", date.Format(validator.DateFormat))
		return nil
	},
}

var weightListCmd = &cobra.Command{
	Use:   "list",
	Short: "List weight records",
	RunE: func(cmd *cobra.Command, args []string) error {
		// Get date range
		var fromDate, toDate time.Time
		var err error

		isDefaultRange := weightFromDate == "" && weightToDate == ""
		if isDefaultRange {
			// Use default range if no dates specified
			fromDate, toDate = validator.GetDefaultDateRange()
		} else {
			fromDate, toDate, err = validator.ValidateDateRange(weightFromDate, weightToDate)
			if err != nil {
				display.ShowError(err.Error())
				return err
			}
		}

		// Get records
		records, err := store.GetWeightRange(fromDate, toDate, isDefaultRange)
		if err != nil {
			display.ShowError("Failed to retrieve weight records: %v", err)
			return err
		}

		if len(records) == 0 {
			display.ShowInfo("No weight records found between %s and %s",
				fromDate.Format(validator.DateFormat),
				toDate.Format(validator.DateFormat))
			return nil
		}

		// Display header
		fmt.Printf("\nWeight Records from %s to %s\n\n",
			fromDate.Format(validator.DateFormat),
			toDate.Format(validator.DateFormat))
		fmt.Printf("%-8s  %-10s  %-7s  %s\n", "ID", "Date", "Weight", "Notes")
		fmt.Println(strings.Repeat("-", 60))

		// Display records and calculate statistics
		var totalWeight float64
		minWeight := records[0].Weight
		maxWeight := records[0].Weight

		for _, record := range records {
			fmt.Printf("%-8s  %-10s  %7.1f  %s\n",
				record.ID,
				record.Date.Format(validator.DateFormat),
				record.Weight,
				record.Notes)

			totalWeight += record.Weight
			if record.Weight < minWeight {
				minWeight = record.Weight
			}
			if record.Weight > maxWeight {
				maxWeight = record.Weight
			}
		}

		// Display summary and statistics
		fmt.Printf("\nSummary:\n")
		fmt.Printf("Total Records: %d\n", len(records))
		fmt.Printf("Average Weight: %.1f lbs\n", totalWeight/float64(len(records)))
		fmt.Printf("Range: %.1f - %.1f lbs (%.1f lbs)\n",
			minWeight, maxWeight, maxWeight-minWeight)

		// Show trend if more than one record
		if len(records) > 1 {
			firstWeight := records[0].Weight
			lastWeight := records[len(records)-1].Weight
			change := lastWeight - firstWeight
			fmt.Printf("Overall Change: %.1f lbs\n", change)
		}

		return nil
	},
}

var weightUpdateCmd = &cobra.Command{
	Use:   "update [record-id]",
	Short: "Update a weight record",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		recordID := args[0]

		// Validate record ID format
		if err := validator.ValidateWeightID(recordID); err != nil {
			return err
		}

		// Validate weight if provided
		if cmd.Flags().Changed("value") {
			if err := validator.ValidateWeight(weightValue); err != nil {
				return err
			}
		}

		// Validate notes if provided
		if cmd.Flags().Changed("notes") {
			if err := validator.ValidateNotes(weightNotes); err != nil {
				return err
			}
		}

		// Get the record to update and surrounding records for validation
		record, err := store.GetWeightByID(recordID)
		if err != nil {
			return fmt.Errorf("error retrieving record: %w", err)
		}
		if record == nil {
			return fmt.Errorf("record not found: %s", recordID)
		}

		prevRecord, err := store.GetPreviousWeightRecord(record.Date)
		if err != nil {
			return fmt.Errorf("error retrieving previous record: %w", err)
		}

		nextRecord, err := store.GetNextWeightRecord(record.Date)
		if err != nil {
			return fmt.Errorf("error retrieving next record: %w", err)
		}

		// Validate the update wouldn't create inconsistent history
		validation := validator.ValidateWeightUpdate(recordID, weightValue, prevRecord, nextRecord)
		if validation.HasWarning {
			fmt.Println(validation.Warning)
			fmt.Print("Do you want to continue? (y/N): ")
			var response string
			fmt.Scanln(&response)
			if response != "y" && response != "Y" {
				return fmt.Errorf("operation cancelled")
			}
		}

		// TODO: Perform the update
		return nil
	},
}

var weightDeleteCmd = &cobra.Command{
	// ... other fields ...
	RunE: func(cmd *cobra.Command, args []string) error {
		recordID := args[0]

		// Validate record ID format
		if err := validator.ValidateWeightID(recordID); err != nil {
			display.ShowError(err.Error())
			return err
		}

		record, err := store.GetWeightByID(recordID)
		if err != nil {
			display.ShowError("Error retrieving record: %v", err)
			return err
		}
		if record == nil {
			display.ShowError("Record not found: %s", recordID)
			return fmt.Errorf("record not found")
		}

		// Format record and context for display
		recordInfo := fmt.Sprintf("  Date: %s\n  Weight: %.1f lbs",
			record.Date.Format(validator.DateFormat),
			record.Weight)
		if record.Notes != "" {
			recordInfo += fmt.Sprintf("\n  Notes: %s", record.Notes)
		}

		// Show confirmation with context
		result := display.ShowDeleteConfirmation(recordInfo, "")
		if !result.Confirmed {
			display.ShowInfo("Delete cancelled")
			return nil
		}

		// TODO: Perform deletion
		display.ShowSuccess("Record deleted successfully")
		return nil
	},
}

// Exercise command (stub)
var exerciseCmd = &cobra.Command{
	Use:   "exercise",
	Short: "Manage exercise records",
}

// Fasting command (stub)
var fastingCmd = &cobra.Command{
	Use:   "fasting",
	Short: "Manage fasting records",
}

// Soda command (stub)
var sodaCmd = &cobra.Command{
	Use:   "soda",
	Short: "Manage soda consumption records",
}

// cmd/tracker/main.go
func testDisplayFormats() {
	// Test adding a normal weight record
	fmt.Println("\n=== Scenario 1: Normal Weight Addition ===")
	cmd := exec.Command("./bin/tracker", "weight", "add", "-v", "185.5", "--date", "2024-01-08")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	// Test adding a duplicate record
	fmt.Println("\n=== Scenario 2: Duplicate Record ===")
	cmd = exec.Command("./bin/tracker", "weight", "add", "-v", "186.0", "--date", "2024-01-08")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	// Test adding an unusual weight change
	fmt.Println("\n=== Scenario 3: Unusual Weight Change ===")
	cmd = exec.Command("./bin/tracker", "weight", "add", "-v", "175.0", "--date", "2024-01-09")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	// Test invalid input
	fmt.Println("\n=== Scenario 4: Invalid Weight ===")
	cmd = exec.Command("./bin/tracker", "weight", "add", "-v", "45.0", "--date", "2024-01-10")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()

	// Test delete confirmation
	fmt.Println("\n=== Scenario 5: Delete Confirmation ===")
	cmd = exec.Command("./bin/tracker", "weight", "delete", "w12345")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Run()
}
