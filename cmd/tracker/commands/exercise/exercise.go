// cmd/tracker/commands/exercise/exercise.go
package exercise

import (
	"github.com/jack-sneddon/my-health-tracker/internal/storage"
	"github.com/spf13/cobra"
)

// Shared flags across exercise commands
// Shared flags across exercise commands
type exerciseFlags struct {
	// Basic flags for add/update
	activity      string
	otherActivity string
	duration      int
	date          string
	notes         string
	completed     bool

	// List command flags
	fromDate  string
	toDate    string
	lastWeek  bool
	lastMonth bool
}

var flags exerciseFlags

// NewExerciseCmd creates the exercise command and all its subcommands
func NewExerciseCmd(store storage.StorageManager) *cobra.Command {
	exerciseCmd := &cobra.Command{
		Use:   "exercise",
		Short: "Manage exercise records",
		Long: `Manage exercise records with full CRUD operations.

Examples:
  # Add an exercise record
  tracker exercise add --activity jogging --duration 45 --date 2024-01-08 --notes "Morning run"

  # Add other activity type
  tracker exercise add --activity other --other-activity "swimming" --duration 30

  # Add with completion status
  tracker exercise add --activity cycling --duration 60 --completed`,
	}

	// Add subcommands
	exerciseCmd.AddCommand(
		newAddCmd(store),
		newGetCmd(store),
		newListCmd(store),
		// Additional commands will be added here
	)

	return exerciseCmd
}

// Add command implementation
func newAddCmd(store storage.StorageManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add a new exercise record",
		RunE:  createAddCmdRunner(store),
	}

	// Add flags
	cmd.Flags().StringVarP(&flags.activity, "activity", "a", "", "Type of exercise activity (required)")
	cmd.Flags().StringVarP(&flags.otherActivity, "other-activity", "o", "", "Name of activity when type is 'other'")
	cmd.Flags().IntVarP(&flags.duration, "duration", "d", 0, "Duration in minutes (required)")
	cmd.Flags().StringVarP(&flags.date, "date", "t", "", "Date of exercise (default: today)")
	cmd.Flags().StringVarP(&flags.notes, "notes", "n", "", "Optional notes about the exercise")
	cmd.Flags().BoolVarP(&flags.completed, "completed", "c", false, "Mark exercise as completed")

	cmd.MarkFlagRequired("activity")
	cmd.MarkFlagRequired("duration")

	return cmd
}
