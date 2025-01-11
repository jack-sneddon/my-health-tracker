// cmd/tracker/commands/soda/soda.go
package soda

import (
	"github.com/jack-sneddon/my-health-tracker/internal/storage"
	"github.com/spf13/cobra"
)

func NewSodaCmd(store storage.StorageManager) *cobra.Command {
	return &cobra.Command{
		Use:   "soda",
		Short: "Manage soda consumption records",
	}
}
