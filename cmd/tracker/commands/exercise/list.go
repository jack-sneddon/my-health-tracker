// cmd/tracker/commands/exercise/list.go
package exercise

import (
	"fmt"
	"time"

	"github.com/jack-sneddon/my-health-tracker/cmd/tracker/commands/result"
	"github.com/jack-sneddon/my-health-tracker/internal/display"
	"github.com/jack-sneddon/my-health-tracker/internal/models"
	"github.com/jack-sneddon/my-health-tracker/internal/storage"
	"github.com/jack-sneddon/my-health-tracker/internal/validator"
	"github.com/spf13/cobra"
)

type exerciseStats struct {
	TotalRecords     int
	TotalDuration    int
	AverageDuration  float64
	CompletedRecords int
}

func newListCmd(store storage.StorageManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List exercise records",
		RunE:  createListCmdRunner(store),
	}

	cmd.Flags().StringVarP(&flags.fromDate, "from", "f", "", "Start date for listing exercises")
	cmd.Flags().StringVarP(&flags.toDate, "to", "t", "", "End date for listing exercises")
	cmd.Flags().BoolVarP(&flags.lastWeek, "week", "w", false, "Show last 7 days")
	cmd.Flags().BoolVarP(&flags.lastMonth, "month", "m", false, "Show last month")

	return cmd
}

// cmd/tracker/commands/exercise/list.go

func createListCmdRunner(store storage.StorageManager) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var fromDate, toDate time.Time
		var err error
		var isDefaultRange bool

		// Handle date range selection
		switch {
		case flags.lastWeek:
			if store.IsTestMode() {
				// In test mode, use fixed date range
				toDate, _ = time.Parse(validator.DateFormat, "2024-01-14")
				fromDate = toDate.AddDate(0, 0, -7)
			} else {
				toDate = time.Now()
				fromDate = toDate.AddDate(0, 0, -7)
			}
			isDefaultRange = false
		case flags.lastMonth:
			if store.IsTestMode() {
				// In test mode, use fixed date range
				toDate, _ = time.Parse(validator.DateFormat, "2024-01-31")
				fromDate, _ = time.Parse(validator.DateFormat, "2024-01-01")
			} else {
				toDate = time.Now()
				fromDate = toDate.AddDate(0, -1, 0)
			}
			isDefaultRange = false
		case flags.fromDate == "" && flags.toDate == "":
			if store.IsTestMode() {
				// In test mode, use fixed date range
				toDate, _ = time.Parse(validator.DateFormat, "2024-01-31")
				fromDate = toDate.AddDate(0, 0, -30)
			} else {
				fromDate, toDate = validator.GetDefaultDateRange()
			}
			isDefaultRange = true
		default:
			fromDate, toDate, err = validator.ValidateDateRange(flags.fromDate, flags.toDate)
			if err != nil {
				return result.ValidationFailed(err).Error
			}
			isDefaultRange = false
		}

		// Get records
		records, err := store.GetExerciseRange(fromDate, toDate, isDefaultRange)
		if err != nil {
			return result.StorageError(err).Error
		}

		if len(records) == 0 {
			return result.NewError(fmt.Errorf("No exercise records found between %s and %s",
				fromDate.Format(validator.DateFormat),
				toDate.Format(validator.DateFormat))).Error
		}

		// Calculate statistics
		stats := calculateExerciseStats(records)

		// Display results
		display.ShowHeader(fmt.Sprintf("Exercise Records from %s to %s",
			fromDate.Format(validator.DateFormat),
			toDate.Format(validator.DateFormat)))

		display.ShowExerciseList(records)

		display.ShowStats(map[string]string{
			"Average Duration":  fmt.Sprintf("%.1f minutes", stats.AverageDuration),
			"Completed Records": fmt.Sprintf("%d", stats.CompletedRecords),
			"Completion Rate":   fmt.Sprintf("%.1f%%", float64(stats.CompletedRecords)/float64(stats.TotalRecords)*100),
			"Total Records":     fmt.Sprintf("%d", stats.TotalRecords),
			"Total Duration":    fmt.Sprintf("%d minutes", stats.TotalDuration),
		})

		return nil
	}
}

func calculateExerciseStats(records []models.ExerciseRecord) exerciseStats {
	stats := exerciseStats{
		TotalRecords: len(records),
	}

	for _, record := range records {
		stats.TotalDuration += record.Duration
		if record.Completed {
			stats.CompletedRecords++
		}
	}

	if stats.TotalRecords > 0 {
		stats.AverageDuration = float64(stats.TotalDuration) / float64(stats.TotalRecords)
	}

	return stats
}
