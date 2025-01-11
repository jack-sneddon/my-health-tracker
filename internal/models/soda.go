// internal/models/soda.go
package models

import (
	"fmt"
	"time"
)

type SodaRecord struct {
	Date     time.Time `json:"date"`
	Consumed bool      `json:"consumed"`
	Quantity float64   `json:"quantity,omitempty"` // in oz
	Notes    string    `json:"notes,omitempty"`
}

func (s SodaRecord) GetDate() time.Time {
	return s.Date
}

func (s SodaRecord) Validate() error {
	if s.Consumed && s.Quantity <= 0 {
		return fmt.Errorf("if soda was consumed, quantity must be greater than 0")
	}
	if !s.Consumed && s.Quantity > 0 {
		return fmt.Errorf("if no soda was consumed, quantity should be 0")
	}
	if s.Quantity > 64 { // reasonable upper limit for daily consumption
		return fmt.Errorf("soda quantity seems unreasonably high")
	}
	return nil
}

func (s SodaRecord) IsCompliant() bool {
	weekday := s.Date.Weekday()
	switch {
	case weekday >= time.Monday && weekday <= time.Thursday:
		return !s.Consumed
	case weekday >= time.Friday && weekday <= time.Sunday:
		return !s.Consumed || s.Quantity <= 12.0
	}
	return false
}
