// cmd/tracker/commands/weight/update.go
package weight

import (
	"fmt"
	"math"

	"github.com/jack-sneddon/my-health-tracker/cmd/tracker/commands/result"
	"github.com/jack-sneddon/my-health-tracker/internal/display"
	"github.com/jack-sneddon/my-health-tracker/internal/storage"
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

		// Get existing record
		record, err := store.GetWeightByID(recordID)
		if err != nil {
			return result.StorageError(err).Error
		}
		if record == nil {
			return result.NotFound("Weight record", recordID).Error
		}

		// Store original weight for comparison
		originalWeight := record.Weight

		// Update fields if provided
		if cmd.Flags().Changed("value") {
			// Validate weight range first
			if err := validateWeightRange(flags.value); err != nil {
				return result.ValidationFailed(err).Error
			}
			record.Weight = flags.value

			// Then check for significant change
			change := math.Abs(record.Weight - originalWeight)
			if change > MaxWeightChange {
				display.ShowWarning(fmt.Sprintf("Weight change of %.1f lbs seems unusual", change))
				if !display.ConfirmAction("Do you want to continue?").Confirmed {
					display.ShowInfo("Operation cancelled")
					return result.NewError(fmt.Errorf("operation cancelled")).Error
				}
			}
		}
		if cmd.Flags().Changed("notes") {
			record.Notes = flags.notes
		}

		// Perform update
		if err := store.UpdateWeight(recordID, *record); err != nil {
			return result.StorageError(err).Error
		}

		cmdResult := result.NewSuccess(*record, "Weight record updated successfully")
		display.ShowCommandResult(cmdResult)

		return nil
	}
}
