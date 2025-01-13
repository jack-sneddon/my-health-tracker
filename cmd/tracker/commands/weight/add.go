// cmd/tracker/commands/weight/add.go
package weight

import (
	"fmt"

	"github.com/jack-sneddon/my-health-tracker/cmd/tracker/commands/result"
	"github.com/jack-sneddon/my-health-tracker/internal/display"
	"github.com/jack-sneddon/my-health-tracker/internal/models"
	"github.com/jack-sneddon/my-health-tracker/internal/storage"
	"github.com/jack-sneddon/my-health-tracker/internal/validator"
	"github.com/spf13/cobra"
)

func newAddCmd(store storage.StorageManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new weight record",
		RunE:  createAddCmdRunner(store),
	}

	// Add flags
	cmd.Flags().Float64VarP(&flags.value, "value", "v", 0, "Weight value in pounds (required)")
	cmd.Flags().StringVarP(&flags.date, "date", "d", "", "Date of weight record (default: today)")
	cmd.Flags().StringVarP(&flags.notes, "notes", "n", "", "Optional notes about the weight record")
	cmd.MarkFlagRequired("value")

	return cmd
}

// cmd/tracker/commands/weight/add.go
func createAddCmdRunner(store storage.StorageManager) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		date, err := validator.ParseDate(flags.date)
		if err != nil {
			return result.ValidationFailed(err).Error
		}

		record := models.WeightRecord{
			Date:   date,
			Weight: flags.value,
			Notes:  flags.notes,
		}

		// Basic validation
		if err := validateWeightRange(record.Weight); err != nil {
			return result.ValidationFailed(err).Error
		}

		// Try to add record
		savedRecord, err := store.AddWeight(record)
		if err != nil {
			if err.Error() == "duplicate_date" {
				display.ShowWarning("Record already exists for %s", date.Format(validator.DateFormat))
				confirmResult := display.ConfirmAction("Do you want to overwrite this record?")
				if !confirmResult.Confirmed {
					display.ShowInfo("Operation cancelled")
					return result.NewError(fmt.Errorf("operation cancelled")).Error
				}
				// If confirmed, use UpdateWeight instead
				existingRecord, _ := store.GetWeight(date)
				if existingRecord != nil {
					record.ID = existingRecord.ID
					if err := store.UpdateWeight(record.ID, record); err != nil {
						return result.StorageError(err).Error
					}
					savedRecord = record
				}
			} else {
				return result.StorageError(err).Error
			}
		}

		// Use CommandResult for success
		cmdResult := result.NewSuccess(savedRecord, "Weight record added successfully")
		display.ShowCommandResult(cmdResult)

		return nil
	}
}
