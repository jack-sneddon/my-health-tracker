// cmd/tracker/commands/exercise/add.go
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

func createAddCmdRunner(store storage.StorageManager) func(*cobra.Command, []string) error {
	return func(cmd *cobra.Command, args []string) error {
		// Parse and validate date
		date, err := validator.ParseDate(flags.date)
		if err != nil {
			return result.ValidationFailed(err).Error
		}

		// Validate activity type
		activity := models.ActivityType(flags.activity)
		if err := validateActivity(activity, flags.otherActivity); err != nil {
			return result.ValidationFailed(err).Error
		}

		// Validate duration
		if err := validateDuration(flags.duration); err != nil {
			return result.ValidationFailed(err).Error
		}

		// Create record
		record := models.ExerciseRecord{
			Date:          date,
			Activity:      activity,
			OtherActivity: flags.otherActivity,
			Duration:      flags.duration,
			Notes:         flags.notes,
			Completed:     flags.completed,
		}

		// Try to add record
		if err := store.AddExercise(record); err != nil {
			if err.Error() == "duplicate_date" {
				display.ShowWarning("Record already exists for %s", date.Format(validator.DateFormat))
				confirmResult := display.ConfirmAction("Do you want to overwrite this record?")
				if !confirmResult.Confirmed {
					display.ShowInfo("Operation cancelled")
					return result.NewError(fmt.Errorf("operation cancelled")).Error
				}
				// If confirmed, try to update instead
				existingRecord, _ := store.GetExercise(date)
				if existingRecord != nil {
					// Update the existing record
					if err := store.UpdateExercise(existingRecord.Date, record); err != nil {
						return result.StorageError(err).Error
					}
				}
			} else {
				return result.StorageError(err).Error
			}
		}

		// Use CommandResult for success
		cmdResult := result.NewSuccess(record, "Exercise record added successfully")
		display.ShowCommandResult(cmdResult)

		return nil
	}
}

// Validation functions
func validateActivity(activity models.ActivityType, otherActivity string) error {
	validActivities := map[models.ActivityType]bool{
		models.Jogging:      true,
		models.Skiing:       true,
		models.Walking:      true,
		models.Cycling:      true,
		models.MountainBike: true,
		models.Pickleball:   true,
		models.Other:        true,
	}

	if !validActivities[activity] {
		return fmt.Errorf("invalid activity type: %s", activity)
	}

	if activity == models.Other && otherActivity == "" {
		return fmt.Errorf("other-activity flag is required when activity type is 'other'")
	}

	return nil
}

func validateDuration(duration int) error {
	if duration <= 0 {
		return fmt.Errorf("duration must be greater than 0 minutes")
	}
	if duration > 480 { // 8 hours seems like a reasonable maximum
		return fmt.Errorf("duration seems unreasonably high (max 480 minutes / 8 hours)")
	}
	return nil
}
