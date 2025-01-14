// cmd/tracker/commands/exercise/update.go
package exercise

import (
	"fmt"

	"github.com/jack-sneddon/my-health-tracker/cmd/tracker/commands/result"
	"github.com/jack-sneddon/my-health-tracker/internal/display"
	"github.com/jack-sneddon/my-health-tracker/internal/models"
	"github.com/jack-sneddon/my-health-tracker/internal/storage"
	"github.com/jack-sneddon/my-health-tracker/internal/validator"
	"github.com/spf13/cobra"
)

func newUpdateCmd(store storage.StorageManager) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "update",
		Short: "Update an exercise record",
		RunE:  createUpdateCmdRunner(store),
	}

	cmd.Flags().StringVarP(&flags.date, "date", "d", "", "Date of the exercise record to update (required)")
	cmd.Flags().StringVarP(&flags.activity, "activity", "a", "", "Updated activity type")
	cmd.Flags().StringVarP(&flags.otherActivity, "other-activity", "o", "", "Name of activity when type is 'other'")
	cmd.Flags().IntVarP(&flags.duration, "duration", "u", 0, "Updated duration in minutes")
	cmd.Flags().StringVarP(&flags.notes, "notes", "n", "", "Updated notes")
	cmd.Flags().BoolVarP(&flags.completed, "completed", "c", false, "Mark as completed")
	cmd.Flags().BoolVar(&flags.notCompleted, "not-completed", false, "Mark as not completed")

	cmd.MarkFlagRequired("date")

	return cmd
}

func createUpdateCmdRunner(store storage.StorageManager) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// Parse and validate date
		date, err := validator.ParseDate(flags.date)
		if err != nil {
			return result.ValidationFailed(err).Error
		}

		// Get existing record
		record, err := store.GetExercise(date)
		if err != nil {
			return result.StorageError(err).Error
		}
		if record == nil {
			return result.NotFound("Exercise record", flags.date).Error
		}

		// Store original values for comparison
		originalActivity := record.Activity
		originalDuration := record.Duration

		// Update fields if provided
		if cmd.Flags().Changed("activity") {
			activity := models.ActivityType(flags.activity)
			if err := validateActivity(activity, flags.otherActivity); err != nil {
				return result.ValidationFailed(err).Error
			}
			record.Activity = activity
			record.OtherActivity = flags.otherActivity
		}

		if cmd.Flags().Changed("duration") {
			if err := validateDuration(flags.duration); err != nil {
				return result.ValidationFailed(err).Error
			}
			record.Duration = flags.duration

			// Check for significant change
			if flags.duration > originalDuration*2 || flags.duration < originalDuration/2 {
				display.ShowWarning(fmt.Sprintf("Duration change from %d to %d minutes seems unusual",
					originalDuration, flags.duration))
				if !display.ConfirmAction("Do you want to continue?").Confirmed {
					display.ShowInfo("Operation cancelled")
					return result.NewError(fmt.Errorf("operation cancelled")).Error
				}
			}
		}

		if cmd.Flags().Changed("notes") {
			record.Notes = flags.notes
		}

		// Handle completion status
		if cmd.Flags().Changed("completed") && cmd.Flags().Changed("not-completed") {
			return result.ValidationFailed(fmt.Errorf("cannot use both --completed and --not-completed flags")).Error
		}
		if cmd.Flags().Changed("completed") {
			record.Completed = true
		}
		if cmd.Flags().Changed("not-completed") {
			record.Completed = false
		}

		// Display summary of changes
		if err := showUpdateSummary(originalActivity, record.Activity, originalDuration, record.Duration); err != nil {
			return result.ValidationFailed(err).Error
		}
		if !display.ConfirmAction("Do you want to apply these changes?").Confirmed {
			display.ShowInfo("Operation cancelled")
			return result.NewError(fmt.Errorf("operation cancelled")).Error
		}

		// Perform update
		if err := store.UpdateExercise(date, *record); err != nil {
			return result.StorageError(err).Error
		}

		cmdResult := result.NewSuccess(*record, "Exercise record updated successfully")
		display.ShowCommandResult(cmdResult)

		return nil
	}
}

func showUpdateSummary(oldActivity, newActivity models.ActivityType, oldDuration, newDuration int) error {
	display.ShowHeader("\nUpdate Summary:")

	if oldActivity != newActivity {
		fmt.Printf("Activity: %s -> %s\n", oldActivity, newActivity)
	}
	if oldDuration != newDuration {
		fmt.Printf("Duration: %d -> %d minutes\n", oldDuration, newDuration)
	}

	return nil
}
