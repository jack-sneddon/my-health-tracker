package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"os/user"
	"path/filepath" // Add this
	"sync"
	"time"

	"github.com/jack-sneddon/my-health-tracker/internal/models"
)

const (
	WeightIDPrefix = "w"
	IDLength       = 5 // number of digits after the prefix
)

const (
	// Default data directory structure
	DefaultDataDir    = "data"
	ProductionDataDir = "production"
	TestDataDir       = "test"
)

// File names
const (
	WeightFileName   = "weight.json"
	ExerciseFileName = "exercise.json"
	FastingFileName  = "fasting.json"
	SodaFileName     = "soda.json"
)

// JSONStorage handles persistence of records to JSON files
type JSONStorage struct {
	rootDir   string // Root directory for all data
	dataDir   string // production or test
	fileLocks map[string]*sync.RWMutex
	mu        sync.RWMutex
}

// Add this helper function to generate IDs
func generateID(prefix string, currentRecords int) string {
	return fmt.Sprintf("%s%05d", prefix, currentRecords+1)
}

// Exercise record implementations
func (s *JSONStorage) AddExercise(record models.ExerciseRecord) error {
	return s.addRecord("exercise", record)
}

func (s *JSONStorage) GetExerciseRange(start, end time.Time, isDefaultRange bool) ([]models.ExerciseRecord, error) {
	filepath := s.getFilePath("exercise")
	lock := s.getLock(filepath)

	lock.RLock()
	defer lock.RUnlock()

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read exercise file: %w", err)
	}

	var records []models.ExerciseRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("failed to parse exercise data: %w", err)
	}

	var filtered []models.ExerciseRecord
	for _, record := range records {
		compareDate := record.Date

		if isDefaultRange {
			// Create dates with same year for comparison
			compareStart := time.Date(2000, start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
			compareEnd := time.Date(2000, end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)
			compareDate = time.Date(2000, record.Date.Month(), record.Date.Day(), 0, 0, 0, 0, time.UTC)

			if (compareDate.Equal(compareStart) || compareDate.After(compareStart)) &&
				(compareDate.Equal(compareEnd) || compareDate.Before(compareEnd)) {
				filtered = append(filtered, record)
			}
		} else {
			if (record.Date.Equal(start) || record.Date.After(start)) &&
				(record.Date.Equal(end) || record.Date.Before(end)) {
				filtered = append(filtered, record)
			}
		}
	}

	return filtered, nil
}

func (s *JSONStorage) GetExercise(date time.Time) (*models.ExerciseRecord, error) {
	records, err := s.GetExerciseRange(date, date, false)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}
	return &records[0], nil
}

// Fasting record implementations
func (s *JSONStorage) AddFasting(record models.FastingRecord) error {
	return s.addRecord("fasting", record)
}

func (s *JSONStorage) GetFastingRange(start, end time.Time, isDefaultRange bool) ([]models.FastingRecord, error) {
	filepath := s.getFilePath("fasting")
	lock := s.getLock(filepath)

	lock.RLock()
	defer lock.RUnlock()

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read fasting file: %w", err)
	}

	var records []models.FastingRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("failed to parse fasting data: %w", err)
	}

	var filtered []models.FastingRecord
	for _, record := range records {
		compareDate := record.Date

		if isDefaultRange {
			// Create dates with same year for comparison
			compareStart := time.Date(2000, start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
			compareEnd := time.Date(2000, end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)
			compareDate = time.Date(2000, record.Date.Month(), record.Date.Day(), 0, 0, 0, 0, time.UTC)

			if (compareDate.Equal(compareStart) || compareDate.After(compareStart)) &&
				(compareDate.Equal(compareEnd) || compareDate.Before(compareEnd)) {
				filtered = append(filtered, record)
			}
		} else {
			if (record.Date.Equal(start) || record.Date.After(start)) &&
				(record.Date.Equal(end) || record.Date.Before(end)) {
				filtered = append(filtered, record)
			}
		}
	}

	return filtered, nil
}

func (s *JSONStorage) GetFasting(date time.Time) (*models.FastingRecord, error) {
	records, err := s.GetFastingRange(date, date, false)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}
	return &records[0], nil
}

// Soda record implementations
func (s *JSONStorage) AddSoda(record models.SodaRecord) error {
	return s.addRecord("soda", record)
}

func (s *JSONStorage) GetSodaRange(start, end time.Time, isDefaultRange bool) ([]models.SodaRecord, error) {
	filepath := s.getFilePath("soda")
	lock := s.getLock(filepath)

	lock.RLock()
	defer lock.RUnlock()

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read soda file: %w", err)
	}

	var records []models.SodaRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("failed to parse soda data: %w", err)
	}

	var filtered []models.SodaRecord
	for _, record := range records {
		compareDate := record.Date

		if isDefaultRange {
			// Create dates with same year for comparison
			compareStart := time.Date(2000, start.Month(), start.Day(), 0, 0, 0, 0, time.UTC)
			compareEnd := time.Date(2000, end.Month(), end.Day(), 0, 0, 0, 0, time.UTC)
			compareDate = time.Date(2000, record.Date.Month(), record.Date.Day(), 0, 0, 0, 0, time.UTC)

			if (compareDate.Equal(compareStart) || compareDate.After(compareStart)) &&
				(compareDate.Equal(compareEnd) || compareDate.Before(compareEnd)) {
				filtered = append(filtered, record)
			}
		} else {
			if (record.Date.Equal(start) || record.Date.After(start)) &&
				(record.Date.Equal(end) || record.Date.Before(end)) {
				filtered = append(filtered, record)
			}
		}
	}

	return filtered, nil
}

func (s *JSONStorage) GetSoda(date time.Time) (*models.SodaRecord, error) {
	records, err := s.GetSodaRange(date, date, false)
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}
	return &records[0], nil
}

func (s *JSONStorage) getFilePath(recordType string) string {
	filename := fmt.Sprintf("%s.json", recordType)
	return filepath.Join(s.rootDir, s.dataDir, filename)
}

func (s *JSONStorage) getLock(filepath string) *sync.RWMutex {
	s.mu.Lock()
	defer s.mu.Unlock()

	if lock, exists := s.fileLocks[filepath]; exists {
		return lock
	}

	lock := &sync.RWMutex{}
	s.fileLocks[filepath] = lock
	return lock
}

// Generic function to add a record to a JSON file
func (s *JSONStorage) addRecord(recordType string, record interface{}) error {
	filepath := s.getFilePath(recordType)
	lock := s.getLock(filepath)

	lock.Lock()
	defer lock.Unlock()

	// Read existing records
	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read %s file: %w", recordType, err)
	}

	var records []interface{}
	if err := json.Unmarshal(data, &records); err != nil {
		return fmt.Errorf("failed to parse %s data: %w", recordType, err)
	}

	// Add new record
	records = append(records, record)

	// Write back to file
	updatedData, err := json.MarshalIndent(records, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal %s data: %w", recordType, err)
	}

	if err := os.WriteFile(filepath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write %s file: %w", recordType, err)
	}

	return nil
}

func (s *JSONStorage) Init() error {
	// Create full directory path if it doesn't exist
	fullPath := filepath.Join(s.rootDir, s.dataDir)
	if err := os.MkdirAll(fullPath, 0755); err != nil {
		return fmt.Errorf("failed to create storage directory %s: %w", fullPath, err)
	}

	// Initialize empty files if they don't exist
	files := map[string]string{
		"weight":   WeightFileName,
		"exercise": ExerciseFileName,
		"fasting":  FastingFileName,
		"soda":     SodaFileName,
	}

	for _, filename := range files {
		filepath := filepath.Join(fullPath, filename)
		if _, err := os.Stat(filepath); os.IsNotExist(err) {
			if err := os.WriteFile(filepath, []byte("[]"), 0644); err != nil {
				return fmt.Errorf("failed to create file %s: %w", filepath, err)
			}
		}
	}

	return nil
}

func (s *JSONStorage) IsTestMode() bool {
	return s.dataDir == TestDataDir
}

func (s *JSONStorage) GetDataDir() string {
	return filepath.Join(s.rootDir, s.dataDir)
}

func NewJSONStorage(rootDir string, testMode bool) StorageManager {
	if rootDir == "" {
		// Get user's home directory
		usr, err := user.Current()
		if err != nil {
			// Fallback to local directory if we can't get home dir
			rootDir = DefaultDataDir
		} else {
			rootDir = filepath.Join(usr.HomeDir, ".health-tracker", "data")
		}
	}

	dataDir := ProductionDataDir
	if testMode {
		dataDir = TestDataDir
	}

	return &JSONStorage{
		rootDir:   rootDir,
		dataDir:   dataDir,
		fileLocks: make(map[string]*sync.RWMutex),
	}
}

// Weight record implementations
// internal/storage/json.go
func (s *JSONStorage) AddWeight(record models.WeightRecord) (models.WeightRecord, error) {
	filepath := s.getFilePath("weight")
	lock := s.getLock(filepath)

	lock.Lock()
	defer lock.Unlock()

	var records []models.WeightRecord
	data, err := os.ReadFile(filepath)
	if err != nil {
		return record, fmt.Errorf("failed to read weight file: %w", err)
	}

	if err := json.Unmarshal(data, &records); err != nil {
		return record, fmt.Errorf("failed to parse weight data: %w", err)
	}

	// Check for existing record on same date
	recordIndex := -1
	for i, r := range records {
		if r.Date.Equal(record.Date) {
			recordIndex = i
			break
		}
	}

	// Generate ID if it's a new record
	if recordIndex == -1 {
		record.ID = generateID(WeightIDPrefix, len(records))
		records = append(records, record)
	} else {
		// Update existing record but keep its ID
		record.ID = records[recordIndex].ID
		records[recordIndex] = record
	}

	updatedData, err := json.MarshalIndent(records, "", "    ")
	if err != nil {
		return record, fmt.Errorf("failed to marshal weight data: %w", err)
	}

	if err := os.WriteFile(filepath, updatedData, 0644); err != nil {
		return record, fmt.Errorf("failed to write weight file: %w", err)
	}

	return record, nil
}

func (s *JSONStorage) GetWeight(date time.Time) (*models.WeightRecord, error) {
	records, err := s.GetWeightRange(date, date, false) // false since this is a specific date query
	if err != nil {
		return nil, err
	}
	if len(records) == 0 {
		return nil, nil
	}
	return &records[0], nil
}

func (s *JSONStorage) GetWeightRange(start, end time.Time, isDefaultRange bool) ([]models.WeightRecord, error) {
    filepath := s.getFilePath("weight")
    lock := s.getLock(filepath)

    lock.RLock()
    defer lock.RUnlock()

    data, err := os.ReadFile(filepath)
    if err != nil {
        return nil, fmt.Errorf("failed to read weight file: %w", err)
    }

    var records []models.WeightRecord
    if err := json.Unmarshal(data, &records); err != nil {
        return nil, fmt.Errorf("failed to parse weight data: %w", err)
    }

    var filtered []models.WeightRecord
    for _, record := range records {
        if isDefaultRange {
            // For default range (last 30 days), normalize year to the target year
            // and compare only month and day
            startDay := start.YearDay()
            endDay := end.YearDay()
            recordDay := record.Date.YearDay()
            
            // Handle year wrap around (December to January)
            if endDay < startDay {
                // If the end day is less than start day, we're spanning a year boundary
                if recordDay >= startDay || recordDay <= endDay {
                    filtered = append(filtered, record)
                }
            } else {
                // Normal case within the same year
                if recordDay >= startDay && recordDay <= endDay {
                    filtered = append(filtered, record)
                }
            }
        } else {
            // For explicit date ranges, compare full dates
            if (record.Date.Equal(start) || record.Date.After(start)) &&
               (record.Date.Equal(end) || record.Date.Before(end)) {
                filtered = append(filtered, record)
            }
        }
    }

    return filtered, nil
}

func (s *JSONStorage) GetLastWeightRecord() (*models.WeightRecord, error) {
	filepath := s.getFilePath("weight")
	lock := s.getLock(filepath)

	lock.RLock()
	defer lock.RUnlock()

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read weight file: %w", err)
	}

	var records []models.WeightRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("failed to parse weight data: %w", err)
	}

	if len(records) == 0 {
		return nil, nil
	}

	// Return the most recent record
	lastRecord := &records[len(records)-1]
	return lastRecord, nil
}

func (s *JSONStorage) GetWeightByID(id string) (*models.WeightRecord, error) {
	filepath := s.getFilePath("weight")
	lock := s.getLock(filepath)

	lock.RLock()
	defer lock.RUnlock()

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read weight file: %w", err)
	}

	var records []models.WeightRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("failed to parse weight data: %w", err)
	}

	// Find record by ID
	for _, record := range records {
		if record.ID == id {
			return &record, nil
		}
	}

	return nil, nil
}

func (s *JSONStorage) GetPreviousWeightRecord(date time.Time) (*models.WeightRecord, error) {
	filepath := s.getFilePath("weight")
	lock := s.getLock(filepath)

	lock.RLock()
	defer lock.RUnlock()

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read weight file: %w", err)
	}

	var records []models.WeightRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("failed to parse weight data: %w", err)
	}

	var prevRecord *models.WeightRecord
	for i := len(records) - 1; i >= 0; i-- {
		if records[i].Date.Before(date) {
			prevRecord = &records[i]
			break
		}
	}

	return prevRecord, nil
}

func (s *JSONStorage) GetNextWeightRecord(date time.Time) (*models.WeightRecord, error) {
	filepath := s.getFilePath("weight")
	lock := s.getLock(filepath)

	lock.RLock()
	defer lock.RUnlock()

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read weight file: %w", err)
	}

	var records []models.WeightRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return nil, fmt.Errorf("failed to parse weight data: %w", err)
	}

	var nextRecord *models.WeightRecord
	for _, record := range records {
		if record.Date.After(date) {
			nextRecord = &record
			break
		}
	}

	return nextRecord, nil
}

// Add to internal/storage/json.go
func (s *JSONStorage) UpdateWeight(id string, record models.WeightRecord) error {
	filepath := s.getFilePath("weight")
	lock := s.getLock(filepath)

	lock.Lock()
	defer lock.Unlock()

	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read weight file: %w", err)
	}

	var records []models.WeightRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return fmt.Errorf("failed to parse weight data: %w", err)
	}

	updated := false
	for i := range records {
		if records[i].ID == id {
			records[i] = record
			updated = true
			break
		}
	}

	if !updated {
		return fmt.Errorf("record not found: %s", id)
	}

	updatedData, err := json.MarshalIndent(records, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal weight data: %w", err)
	}

	if err := os.WriteFile(filepath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write weight file: %w", err)
	}

	return nil
}

func (s *JSONStorage) DeleteWeight(id string) error {
	filepath := s.getFilePath("weight")
	lock := s.getLock(filepath)

	lock.Lock()
	defer lock.Unlock()

	data, err := os.ReadFile(filepath)
	if err != nil {
		return fmt.Errorf("failed to read weight file: %w", err)
	}

	var records []models.WeightRecord
	if err := json.Unmarshal(data, &records); err != nil {
		return fmt.Errorf("failed to parse weight data: %w", err)
	}

	newRecords := make([]models.WeightRecord, 0, len(records))
	found := false
	for _, record := range records {
		if record.ID == id {
			found = true
			continue
		}
		newRecords = append(newRecords, record)
	}

	if !found {
		return fmt.Errorf("record not found: %s", id)
	}

	updatedData, err := json.MarshalIndent(newRecords, "", "    ")
	if err != nil {
		return fmt.Errorf("failed to marshal weight data: %w", err)
	}

	if err := os.WriteFile(filepath, updatedData, 0644); err != nil {
		return fmt.Errorf("failed to write weight file: %w", err)
	}

	return nil
}
