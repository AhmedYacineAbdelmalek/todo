package taskdata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	dataFileName = "tasks.json"
	dateFormat   = "2006-01-02"
)

type Task struct {
	ID          int    `json:"id"`
	Description string `json:"description"`
	DueDate     string `json:"due_date"`
	Priority    string `json:"priority"`
	Completed   bool   `json:"completed"`
}

type TaskStore struct {
	Tasks  []Task `json:"tasks"`
	NextID int    `json:"next_id"`
}

// ValidatePriority checks if the priority is valid
func ValidatePriority(priority string) error {
	validPriorities := []string{"low", "L", "N", "normal", "H", "high"}
	priority = strings.ToLower(priority)

	for _, valid := range validPriorities {
		if priority == valid {
			return nil
		}
	}

	return fmt.Errorf("invalid priority '%s'. Valid priorities are: low, normal, high", priority)
}

// ValidateDate checks if the date format is valid (YYYY-MM-DD)
func ValidateDate(date string) error {
	if date == "" {
		return nil // Empty date is allowed
	}

	// Check format using regex
	dateRegex := regexp.MustCompile(`^\d{4}-\d{2}-\d{2}$`)
	if !dateRegex.MatchString(date) {
		return fmt.Errorf("invalid date format '%s'. Please use YYYY-MM-DD format", date)
	}

	// Check if it's a valid date
	_, err := time.Parse(dateFormat, date)
	if err != nil {
		return fmt.Errorf("invalid date '%s'. Please provide a valid date in YYYY-MM-DD format", date)
	}

	return nil
}

// GetDataFilePath returns the path to the tasks data file
func GetDataFilePath() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		// Fallback to current directory if home directory is not accessible
		return dataFileName
	}
	return filepath.Join(homeDir, ".todo", dataFileName)
}

// LoadTasks loads tasks from the data file
func LoadTasks() (*TaskStore, error) {
	filePath := GetDataFilePath()

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %v", err)
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		// Create new empty task store
		return &TaskStore{
			Tasks:  []Task{},
			NextID: 1,
		}, nil
	}

	// Read file
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tasks file: %v", err)
	}

	// Parse JSON
	var store TaskStore
	if err := json.Unmarshal(data, &store); err != nil {
		return nil, fmt.Errorf("failed to parse tasks file: %v", err)
	}

	return &store, nil
}

// SaveTasks saves tasks to the data file
func (store *TaskStore) SaveTasks() error {
	filePath := GetDataFilePath()

	// Create directory if it doesn't exist
	dir := filepath.Dir(filePath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create data directory: %v", err)
	}

	// Marshal to JSON
	data, err := json.MarshalIndent(store, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal tasks: %v", err)
	}

	// Write to file
	if err := os.WriteFile(filePath, data, 0644); err != nil {
		return fmt.Errorf("failed to write tasks file: %v", err)
	}

	return nil
}

// AddTask adds a new task to the store
func (store *TaskStore) AddTask(description, dueDate, priority string) (*Task, error) {
	// Validate inputs
	if err := ValidateDate(dueDate); err != nil {
		return nil, err
	}

	if err := ValidatePriority(priority); err != nil {
		return nil, err
	}

	// Create new task
	task := Task{
		ID:          store.NextID,
		Description: description,
		DueDate:     dueDate,
		Priority:    strings.ToLower(priority),
		Completed:   false,
	}

	// Add to store
	store.Tasks = append(store.Tasks, task)
	store.NextID++

	return &task, nil
}

// CompleteTask marks a task as completed
func (store *TaskStore) CompleteTask(id int) error {
	for i, task := range store.Tasks {
		if task.ID == id {
			store.Tasks[i].Completed = true
			return nil
		}
	}
	return fmt.Errorf("task with ID %d not found", id)
}
