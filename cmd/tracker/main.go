// cmd/tracker/main.go
package main

import (
	"log"
	"os"

	"github.com/jack-sneddon/my-health-tracker/cmd/tracker/commands"
)

func main() {
	// Check for test mode
	testMode := os.Getenv("TEST_MODE") == "true"

	// Execute with test mode setting
	if err := commands.Execute(testMode); err != nil {
		log.Printf("Error: %v\n", err)
		os.Exit(1)
	}
}
