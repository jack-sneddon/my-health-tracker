// cmd/tracker/commands/root.go
package commands

import (
	"log"

	"github.com/jack-sneddon/my-health-tracker/cmd/tracker/commands/exercise"
	"github.com/jack-sneddon/my-health-tracker/cmd/tracker/commands/fasting"
	"github.com/jack-sneddon/my-health-tracker/cmd/tracker/commands/soda"
	"github.com/jack-sneddon/my-health-tracker/cmd/tracker/commands/weight"
	"github.com/jack-sneddon/my-health-tracker/internal/storage"
	"github.com/spf13/cobra"
)

var (
	rootCmd = &cobra.Command{
		Use:   "tracker",
		Short: "Health tracking application",
		Long: `A health tracking application for monitoring weight, exercise, fasting, and soda consumption.

Command Structure (CRUD operations):
  CREATE:
    tracker weight add --value 185.5 --date 2024-01-08 --notes "Morning weight"
    tracker exercise add --activity jogging --duration 45 --date 2024-01-08
    tracker fasting add --pattern full-fast --date 2024-01-08
    tracker soda add --consumed --quantity 12 --date 2024-01-08

  READ:
    tracker weight get --date 2024-01-08
    tracker weight list --from 2024-01-01 --to 2024-01-08

  UPDATE:
    tracker weight update w12345 --value 184.5
    tracker exercise update e12345 --duration 50

  DELETE:
    tracker weight delete w12345

Use "tracker [command] --help" for more information about a command.`,
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute(testMode bool) error {
	// Initialize storage
	store := storage.NewJSONStorage("", testMode)
	if err := store.Init(); err != nil {
		log.Fatalf("Failed to initialize storage: %v", err)
	}

	// Add main command groups
	rootCmd.AddCommand(weight.NewWeightCmd(store))
	rootCmd.AddCommand(exercise.NewExerciseCmd(store))
	rootCmd.AddCommand(fasting.NewFastingCmd(store))
	rootCmd.AddCommand(soda.NewSodaCmd(store))

	return rootCmd.Execute()
}
