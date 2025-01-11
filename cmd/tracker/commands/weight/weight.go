// cmd/tracker/commands/weight/weight.go
package weight

import (
	"github.com/jack-sneddon/my-health-tracker/internal/storage"
	"github.com/spf13/cobra"
)

// Shared flags across weight commands
type weightFlags struct {
	value    float64
	date     string
	notes    string
	fromDate string
	toDate   string
}

var flags weightFlags

// NewWeightCmd creates the weight command and all its subcommands
func NewWeightCmd(store storage.StorageManager) *cobra.Command {
	weightCmd := &cobra.Command{
		Use:   "weight",
		Short: "Manage weight records",
		Long: `Manage weight records with full CRUD operations.
    
Examples:
  # Add a weight record
  tracker weight add --value 185.5 --date 2024-01-08 --notes "Morning weight"

  # Get weight for a specific date
  tracker weight get --date 2024-01-08

  # List weights for a date range
  tracker weight list --from 2024-01-01 --to 2024-01-08

  # Update a weight record
  tracker weight update w12345 --value 184.5

  # Delete a weight record
  tracker weight delete w12345`,
	}

	// Add subcommands
	weightCmd.AddCommand(
		newAddCmd(store),
		newGetCmd(store),
		newListCmd(store),
		newUpdateCmd(store),
		newDeleteCmd(store),
	)

	return weightCmd
}
