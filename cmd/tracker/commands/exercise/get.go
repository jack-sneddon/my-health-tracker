// cmd/tracker/commands/exercise/get.go
package exercise

import (
	"github.com/jack-sneddon/my-health-tracker/cmd/tracker/commands/result"
	"github.com/jack-sneddon/my-health-tracker/internal/display"
	"github.com/jack-sneddon/my-health-tracker/internal/storage"
	"github.com/jack-sneddon/my-health-tracker/internal/validator"
	"github.com/spf13/cobra"
)

func newGetCmd(store storage.StorageManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Get exercise record for a specific date",
		RunE:  createGetCmdRunner(store),
	}

	cmd.Flags().StringVarP(&flags.date, "date", "d", "", "Date to get exercise record for (required)")
	cmd.MarkFlagRequired("date")

	return cmd
}

func createGetCmdRunner(store storage.StorageManager) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// Parse and validate date
		date, err := validator.ParseDate(flags.date)
		if err != nil {
			return result.ValidationFailed(err).Error
		}

		// Get record from storage
		record, err := store.GetExercise(date)
		if err != nil {
			return result.StorageError(err).Error
		}

		// Handle not found
		if record == nil {
			return result.NotFound("Exercise record", flags.date).Error
		}

		// Create success result and display
		cmdResult := result.NewSuccess(*record, "Found exercise record")
		display.ShowCommandResult(cmdResult)

		return nil
	}
}
