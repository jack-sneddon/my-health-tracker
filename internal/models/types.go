// internal/models/types.go
package models

import "time"

type WeekDay int

const (
	Monday WeekDay = iota
	Tuesday
	Wednesday
	Thursday
	Friday
	Saturday
	Sunday
)

type MealPattern string

const (
	FullFast MealPattern = "full-fast"
	OneMeal  MealPattern = "one-meal"
	Regular  MealPattern = "regular"
)

// Common interfaces that all records will implement
type Record interface {
	GetDate() time.Time
	IsCompliant() bool
	Validate() error
}
