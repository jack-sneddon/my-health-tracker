// cmd/tracker/commands/weight/get.go
package weight

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
		Short: "Get weight record for a specific date",
		RunE:  createGetCmdRunner(store),
	}

	cmd.Flags().StringVarP(&flags.date, "date", "d", "", "Date to get weight record for (required)")
	cmd.MarkFlagRequired("date")

	return cmd
}

func createGetCmdRunner(store storage.StorageManager) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// 1. Parse and validate date
		date, err := validator.ParseDate(flags.date)
		if err != nil {
			return result.ValidationFailed(err).Error
		}

		// 2. Get record from storage
		record, err := store.GetWeight(date)
		if err != nil {
			return result.StorageError(err).Error
		}

		// 3. Handle not found
		if record == nil {
			return result.NotFound("Weight record", flags.date).Error
		}

		// 4. Create success result and display
		cmdResult := result.NewSuccess(*record, "Found weight record")
		display.ShowCommandResult(cmdResult)

		return nil
	}
}
