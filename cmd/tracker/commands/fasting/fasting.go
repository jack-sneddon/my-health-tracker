// cmd/tracker/commands/fasting/fasting.go
package fasting

import (
	"github.com/jack-sneddon/my-health-tracker/internal/storage"
	"github.com/spf13/cobra"
)

func NewFastingCmd(store storage.StorageManager) *cobra.Command {
	return &cobra.Command{
		Use:   "fasting",
		Short: "Manage fasting records",
	}
}
