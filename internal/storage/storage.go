// internal/storage/storage.go
package storage

import (
	"time"

	"github.com/jack-sneddon/my-health-tracker/internal/models"
)

// StorageManager defines the interface for all storage operations
type StorageManager interface {
	Init() error
	IsTestMode() bool
	GetDataDir() string

	// Weight records
	AddWeight(models.WeightRecord) (models.WeightRecord, error)
	GetWeight(time.Time) (*models.WeightRecord, error)
	GetWeightRange(start, end time.Time, isDefaultRange bool) ([]models.WeightRecord, error)
	GetWeightByID(id string) (*models.WeightRecord, error)
	GetLastWeightRecord() (*models.WeightRecord, error)
	GetPreviousWeightRecord(date time.Time) (*models.WeightRecord, error)
	GetNextWeightRecord(date time.Time) (*models.WeightRecord, error)
	UpdateWeight(id string, record models.WeightRecord) error
	DeleteWeight(id string) error

	// Exercise records
	AddExercise(models.ExerciseRecord) error
	GetExercise(time.Time) (*models.ExerciseRecord, error)
	GetExerciseRange(start, end time.Time, isDefaultRange bool) ([]models.ExerciseRecord, error)
	UpdateExercise(date time.Time, record models.ExerciseRecord) error
	DeleteExercise(date time.Time) error

	// Fasting records
	AddFasting(models.FastingRecord) error
	GetFasting(time.Time) (*models.FastingRecord, error)
	GetFastingRange(start, end time.Time, isDefaultRange bool) ([]models.FastingRecord, error)

	// Soda records
	AddSoda(models.SodaRecord) error
	GetSoda(time.Time) (*models.SodaRecord, error)
	GetSodaRange(start, end time.Time, isDefaultRange bool) ([]models.SodaRecord, error)
}
