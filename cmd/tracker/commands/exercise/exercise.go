// cmd/tracker/commands/exercise/exercise.go
package exercise

import (
	"github.com/jack-sneddon/my-health-tracker/internal/storage"
	"github.com/spf13/cobra"
)

func NewExerciseCmd(store storage.StorageManager) *cobra.Command {
	return &cobra.Command{
		Use:   "exercise",
		Short: "Manage exercise records",
	}
}
