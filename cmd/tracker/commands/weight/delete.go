// cmd/tracker/commands/weight/delete.go
package weight

import (
	"fmt"

	"github.com/jack-sneddon/my-health-tracker/cmd/tracker/commands/result"
	"github.com/jack-sneddon/my-health-tracker/internal/display"
	"github.com/jack-sneddon/my-health-tracker/internal/storage"
	"github.com/jack-sneddon/my-health-tracker/internal/validator"
	"github.com/spf13/cobra"
)

func newDeleteCmd(store storage.StorageManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete [record-id]",
		Short: "Delete a weight record",
		Args:  cobra.ExactArgs(1),
		RunE:  createDeleteCmdRunner(store),
	}

	return cmd
}

// cmd/tracker/commands/weight/delete.go
func createDeleteCmdRunner(store storage.StorageManager) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		recordID := args[0]

		// Validate record ID format
		if err := validator.ValidateWeightID(recordID); err != nil {
			return result.ValidationFailed(err).Error
		}

		// Get record to show confirmation
		record, err := store.GetWeightByID(recordID)
		if err != nil {
			return result.StorageError(err).Error
		}
		if record == nil {
			return result.NotFound("Weight record", recordID).Error
		}

		// Show confirmation with record details
		confirmResult := display.ShowDeleteConfirmation(
			record.ID,
			record.Date.Format(validator.DateFormat),
			fmt.Sprintf("%.1f", record.Weight),
			record.Notes,
		)

		if !confirmResult.Confirmed {
			display.ShowInfo("Operation cancelled")
			return result.NewError(fmt.Errorf("operation cancelled")).Error
		}

		if err := store.DeleteWeight(recordID); err != nil {
			return result.StorageError(err).Error
		}

		display.ShowSuccess("Weight record deleted successfully")
		return nil
	}
}
