// cmd/tracker/commands/weight/get.go
package weight

import (
	"fmt"

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

// cmd/tracker/commands/weight/get.go
func createGetCmdRunner(store storage.StorageManager) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		date, err := validator.ParseDate(flags.date)
		if err != nil {
			return result.ValidationFailed(err).Error
		}

		record, err := store.GetWeight(date)
		if err != nil {
			return result.StorageError(err).Error
		}

		if record == nil {
			return result.NotFound("Weight record", flags.date).Error
		}

		// Directly display the record instead of using CommandResult
		display.ShowWeightRecord(
			record.ID,
			record.Date.Format(validator.DateFormat),
			fmt.Sprintf("%.1f", record.Weight),
			record.Notes,
		)

		// Return success without error
		return nil
	}
}
