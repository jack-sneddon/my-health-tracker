// cmd/tracker/commands/weight/list.go
package weight

import (
	"fmt"
	"strings"
	"time"

	"github.com/jack-sneddon/my-health-tracker/cmd/tracker/commands/result"
	"github.com/jack-sneddon/my-health-tracker/internal/display"
	"github.com/jack-sneddon/my-health-tracker/internal/models"
	"github.com/jack-sneddon/my-health-tracker/internal/storage"
	"github.com/jack-sneddon/my-health-tracker/internal/validator"
	"github.com/spf13/cobra"
)

type weightStats struct {
	TotalRecords  int
	AverageWeight float64
	MinWeight     float64
	MaxWeight     float64
	TotalChange   float64
}

// cmd/tracker/commands/weight/list.go
func newListCmd(store storage.StorageManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "list",
		Short: "List weight records",
		RunE:  createListCmdRunner(store),
	}

	cmd.Flags().StringVarP(&flags.fromDate, "from", "f", "", "Start date for listing weights")
	cmd.Flags().StringVarP(&flags.toDate, "to", "t", "", "End date for listing weights")
	// Add quick range flags
	cmd.Flags().BoolVarP(&flags.lastWeek, "week", "w", false, "Show last 7 days")
	cmd.Flags().BoolVarP(&flags.lastMonth, "month", "m", false, "Show last month")

	return cmd
}

// cmd/tracker/commands/weight/list.go
func createListCmdRunner(store storage.StorageManager) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		var fromDate, toDate time.Time
		var err error
		var isDefaultRange bool

		// Handle date range selection
		switch {
		case flags.lastWeek:
			toDate = time.Now()
			fromDate = toDate.AddDate(0, 0, -7)
			isDefaultRange = false
		case flags.lastMonth:
			toDate = time.Now()
			fromDate = toDate.AddDate(0, -1, 0)
			isDefaultRange = false
		case flags.fromDate == "" && flags.toDate == "":
			fromDate, toDate = validator.GetDefaultDateRange()
			isDefaultRange = true
		default:
			fromDate, toDate, err = validator.ValidateDateRange(flags.fromDate, flags.toDate)
			if err != nil {
				return result.ValidationFailed(err).Error
			}
			isDefaultRange = false
		}

		// Get records
		records, err := store.GetWeightRange(fromDate, toDate, isDefaultRange)
		if err != nil {
			return result.StorageError(err).Error
		}

		if len(records) == 0 {
			return result.NewError(fmt.Errorf("No weight records found between %s and %s",
				fromDate.Format(validator.DateFormat),
				toDate.Format(validator.DateFormat))).Error
		}

		// Display header
		fmt.Printf("\nWeight Records from %s to %s\n\n",
			fromDate.Format(validator.DateFormat),
			toDate.Format(validator.DateFormat))
		fmt.Printf("%-8s  %-10s  %-7s  %s\n", "ID", "Date", "Weight", "Notes")
		fmt.Println(strings.Repeat("-", 60))

		// Display records and calculate statistics
		var totalWeight float64
		minWeight := records[0].Weight
		maxWeight := records[0].Weight

		for _, record := range records {
			fmt.Printf("%-8s  %-10s  %7.1f  %s\n",
				record.ID,
				record.Date.Format(validator.DateFormat),
				record.Weight,
				record.Notes)

			totalWeight += record.Weight
			if record.Weight < minWeight {
				minWeight = record.Weight
			}
			if record.Weight > maxWeight {
				maxWeight = record.Weight
			}
		}

		// Display summary and statistics
		fmt.Printf("\nSummary:\n")
		fmt.Printf("Total Records : %d\n", len(records))
		fmt.Printf("Average Weight: %.1f lbs\n", totalWeight/float64(len(records)))
		fmt.Printf("Weight Range  : %.1f - %.1f lbs (%.1f lbs)\n",
			minWeight, maxWeight, maxWeight-minWeight)

		// Show trend if more than one record
		if len(records) > 1 {
			firstWeight := records[0].Weight
			lastWeight := records[len(records)-1].Weight
			change := lastWeight - firstWeight
			fmt.Printf("Overall Change: %.1f lbs\n", change)
		}

		return nil
	}
}

func calculateWeightStats(records []models.WeightRecord) weightStats {
	stats := weightStats{
		TotalRecords: len(records),
		MinWeight:    records[0].Weight,
		MaxWeight:    records[0].Weight,
	}

	var totalWeight float64
	for _, record := range records {
		totalWeight += record.Weight
		if record.Weight < stats.MinWeight {
			stats.MinWeight = record.Weight
		}
		if record.Weight > stats.MaxWeight {
			stats.MaxWeight = record.Weight
		}
	}

	stats.AverageWeight = totalWeight / float64(stats.TotalRecords)
	if stats.TotalRecords > 1 {
		stats.TotalChange = records[len(records)-1].Weight - records[0].Weight
	}

	return stats
}

func displayWeightList(records []models.WeightRecord, stats weightStats, fromDate, toDate time.Time) {
	display.ShowHeader(fmt.Sprintf("Weight Records from %s to %s",
		fromDate.Format(validator.DateFormat),
		toDate.Format(validator.DateFormat)))

	display.ShowTable([]string{"ID", "Date", "Weight", "Notes"},
		func() [][]string {
			var rows [][]string
			for _, r := range records {
				rows = append(rows, []string{
					r.ID,
					r.Date.Format(validator.DateFormat),
					fmt.Sprintf("%.1f", r.Weight),
					r.Notes,
				})
			}
			return rows
		}())

	display.ShowStats(map[string]string{
		"Total Records":  fmt.Sprintf("%d", stats.TotalRecords),
		"Average Weight": fmt.Sprintf("%.1f lbs", stats.AverageWeight),
		"Weight Range": fmt.Sprintf("%.1f - %.1f lbs (%.1f lbs)",
			stats.MinWeight, stats.MaxWeight, stats.MaxWeight-stats.MinWeight),
		"Overall Change": fmt.Sprintf("%.1f lbs", stats.TotalChange),
	})
}
