// internal/models/weight.go
package models

import (
	"fmt"
	"time"
)

// internal/models/weight.go
type WeightRecord struct {
	ID     string    `json:"id"`
	Date   time.Time `json:"date"`
	Weight float64   `json:"weight"` // in pounds
	Notes  string    `json:"notes,omitempty"`
}

func (w WeightRecord) GetDate() time.Time {
	return w.Date
}

func (w WeightRecord) Validate() error {
	if w.Weight <= 0 {
		return fmt.Errorf("weight must be greater than 0")
	}
	if w.Weight > 1000 { // reasonable upper limit
		return fmt.Errorf("weight seems unreasonably high")
	}
	return nil
}

func (w WeightRecord) IsCompliant() bool {
	return true // Weight records are always compliant - they're measurements
}
