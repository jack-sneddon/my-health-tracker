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

		// Get validation context
		lastWeight, err := store.GetLastWeightRecord()
		if err != nil {
			return result.StorageError(err).Error
		}

		// Validate with context
		validationResult := ValidateWithContext(
			ValidationRequest{
				Record:     record,
				LastRecord: lastWeight,
			},
			ValidationContext{
				IsUpdate:    false,
				AllowFuture: false,
			},
		)

		if !validationResult.IsValid {
			return result.ValidationFailed(validationResult.Error, validationResult.Warnings...).Error
		}

		// Handle warnings
		if len(validationResult.Warnings) > 0 {
			for _, warning := range validationResult.Warnings {
				display.ShowWarning(warning)
			}
			if !display.ConfirmAction("Do you want to continue?").Confirmed {
				return result.NewError(fmt.Errorf("operation cancelled")).Error
			}
		}

		// Save record
		savedRecord, err := store.AddWeight(record)
		if err != nil {
			return result.StorageError(err).Error
		}

		cmdResult := result.NewSuccess(savedRecord, "Weight record added successfully")
		display.ShowCommandResult(cmdResult)
		return nil
	}
}
