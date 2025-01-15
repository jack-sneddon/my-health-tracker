// cmd/tracker/commands/exercise/delete.go
package exercise

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
		Use:   "delete",
		Short: "Delete an exercise record",
		RunE:  createDeleteCmdRunner(store),
	}

	cmd.Flags().StringVarP(&flags.date, "date", "d", "", "Date of exercise record to delete (required)")
	cmd.MarkFlagRequired("date")

	return cmd
}

func createDeleteCmdRunner(store storage.StorageManager) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// Parse and validate date
		date, err := validator.ParseDate(flags.date)
		if err != nil {
			return result.ValidationFailed(err).Error
		}

		// Get record to show confirmation
		record, err := store.GetExercise(date)
		if err != nil {
			return result.StorageError(err).Error
		}
		if record == nil {
			return result.NotFound("Exercise record", flags.date).Error
		}

		// Show confirmation with record details
		confirmResult := display.ShowExerciseDeleteConfirmation(
			record.Date.Format(validator.DateFormat),
			string(record.Activity),
			record.OtherActivity,
			record.Duration,
			record.Notes,
			record.Completed,
		)

		if !confirmResult.Confirmed {
			display.ShowInfo("Operation cancelled")
			return result.NewError(fmt.Errorf("operation cancelled")).Error
		}

		if err := store.DeleteExercise(date); err != nil {
			return result.StorageError(err).Error
		}

		cmdResult := result.NewSuccess(nil, "Exercise record deleted successfully")
		display.ShowCommandResult(cmdResult)

		return nil
	}
}
