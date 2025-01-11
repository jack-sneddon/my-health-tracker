// cmd/tracker/commands/weight/update.go
package weight

import (
	"fmt"

	"github.com/jack-sneddon/my-health-tracker/cmd/tracker/commands/result"
	"github.com/jack-sneddon/my-health-tracker/internal/display"
	"github.com/jack-sneddon/my-health-tracker/internal/storage"
	"github.com/jack-sneddon/my-health-tracker/internal/validator"
	"github.com/spf13/cobra"
)

func newUpdateCmd(store storage.StorageManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update [record-id]",
		Short: "Update a weight record",
		Args:  cobra.ExactArgs(1),
		RunE:  createUpdateCmdRunner(store),
	}

	cmd.Flags().Float64VarP(&flags.value, "value", "v", 0, "New weight value in pounds")
	cmd.Flags().StringVarP(&flags.notes, "notes", "n", "", "Updated notes about the weight record")

	return cmd
}

func createUpdateCmdRunner(store storage.StorageManager) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		recordID := args[0]

		// Validate record ID format
		if err := ValidateWeightID(recordID); err != nil {
			return result.ValidationFailed(err).Error
		}

		// Get existing record
		record, err := store.GetWeightByID(recordID)
		if err != nil {
			return result.StorageError(err).Error
		}
		if record == nil {
			return result.NotFound("Weight record", recordID).Error
		}

		// Update fields if provided
		if cmd.Flags().Changed("value") {
			record.Weight = flags.value
		}
		if cmd.Flags().Changed("notes") {
			record.Notes = flags.notes
		}

		// Validate the updated weight
		validationResult := ValidateWithContext(
			ValidationRequest{
				Record:     *record,
				LastRecord: nil, // Will be fetched by storage layer
			},
			ValidationContext{
				IsUpdate:    true,
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

		// Perform update
		if err := store.UpdateWeight(recordID, *record); err != nil {
			return result.StorageError(err).Error
		}

		display.ShowSuccess("Weight record updated successfully")
		display.ShowWeightRecord(
			record.ID,
			record.Date.Format(validator.DateFormat),
			fmt.Sprintf("%.1f", record.Weight),
			record.Notes,
		)

		return nil
	}
}
