package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/DhirajZope/lazytodo/internal/models"
)

const (
	DataFileName = "lazytodo.json"
	DataDir      = ".lazytodo"
)

// Storage handles data persistence
type Storage struct {
	dataPath string
}

// New creates a new Storage instance
func New() (*Storage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	dataDir := filepath.Join(homeDir, DataDir)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	dataPath := filepath.Join(dataDir, DataFileName)

	return &Storage{
		dataPath: dataPath,
	}, nil
}

// Load loads the application data from file
func (s *Storage) Load() (*models.Application, error) {
	// Check if file exists
	if _, err := os.Stat(s.dataPath); os.IsNotExist(err) {
		// Return default application if file doesn't exist
		return &models.Application{
			TodoLists: []models.TodoList{},
			Settings:  models.DefaultSettings(),
		}, nil
	}

	data, err := os.ReadFile(s.dataPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read data file: %w", err)
	}

	var app models.Application
	if err := json.Unmarshal(data, &app); err != nil {
		return nil, fmt.Errorf("failed to unmarshal data: %w", err)
	}

	// Ensure settings have default values if missing
	if app.Settings.ReminderMinutes == 0 {
		app.Settings = models.DefaultSettings()
	}

	return &app, nil
}

// Save saves the application data to file
func (s *Storage) Save(app *models.Application) error {
	data, err := json.MarshalIndent(app, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal data: %w", err)
	}

	// Create backup if file exists
	if _, err := os.Stat(s.dataPath); err == nil {
		backupPath := s.dataPath + ".backup"
		if err := s.copyFile(s.dataPath, backupPath); err != nil {
			// Log warning but don't fail the save operation
			fmt.Printf("Warning: failed to create backup: %v\n", err)
		}
	}

	if err := os.WriteFile(s.dataPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write data file: %w", err)
	}

	return nil
}

// copyFile creates a copy of the source file
func (s *Storage) copyFile(src, dst string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	return os.WriteFile(dst, data, 0644)
}

// GetDataPath returns the path to the data file
func (s *Storage) GetDataPath() string {
	return s.dataPath
}

// CreateTodoList creates a new todo list
func (s *Storage) CreateTodoList(app *models.Application, name, description string) string {
	id := generateID()
	newList := models.TodoList{
		ID:          id,
		Name:        name,
		Description: description,
		Tasks:       []models.Task{},
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	app.TodoLists = append(app.TodoLists, newList)
	return id
}

// UpdateTodoList updates an existing todo list
func (s *Storage) UpdateTodoList(app *models.Application, listID, name, description string) error {
	for i := range app.TodoLists {
		if app.TodoLists[i].ID == listID {
			app.TodoLists[i].Name = name
			app.TodoLists[i].Description = description
			app.TodoLists[i].UpdatedAt = time.Now()
			return nil
		}
	}
	return fmt.Errorf("todo list with ID %s not found", listID)
}

// DeleteTodoList deletes a todo list
func (s *Storage) DeleteTodoList(app *models.Application, listID string) error {
	for i, list := range app.TodoLists {
		if list.ID == listID {
			app.TodoLists = append(app.TodoLists[:i], app.TodoLists[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("todo list with ID %s not found", listID)
}

// CreateTask creates a new task in a todo list
func (s *Storage) CreateTask(app *models.Application, listID, title, description string, priority models.Priority, deadline *time.Time) (string, error) {
	for i := range app.TodoLists {
		if app.TodoLists[i].ID == listID {
			taskID := generateID()
			newTask := models.Task{
				ID:          taskID,
				Title:       title,
				Description: description,
				Completed:   false,
				Priority:    priority,
				CreatedAt:   time.Now(),
				UpdatedAt:   time.Now(),
				Deadline:    deadline,
			}

			app.TodoLists[i].Tasks = append(app.TodoLists[i].Tasks, newTask)
			app.TodoLists[i].UpdatedAt = time.Now()
			return taskID, nil
		}
	}
	return "", fmt.Errorf("todo list with ID %s not found", listID)
}

// UpdateTask updates an existing task
func (s *Storage) UpdateTask(app *models.Application, listID, taskID, title, description string, priority models.Priority, deadline *time.Time) error {
	for i := range app.TodoLists {
		if app.TodoLists[i].ID == listID {
			for j := range app.TodoLists[i].Tasks {
				if app.TodoLists[i].Tasks[j].ID == taskID {
					app.TodoLists[i].Tasks[j].Title = title
					app.TodoLists[i].Tasks[j].Description = description
					app.TodoLists[i].Tasks[j].Priority = priority
					app.TodoLists[i].Tasks[j].Deadline = deadline
					app.TodoLists[i].Tasks[j].UpdatedAt = time.Now()
					app.TodoLists[i].UpdatedAt = time.Now()
					return nil
				}
			}
			return fmt.Errorf("task with ID %s not found in list %s", taskID, listID)
		}
	}
	return fmt.Errorf("todo list with ID %s not found", listID)
}

// ToggleTask toggles the completion status of a task
func (s *Storage) ToggleTask(app *models.Application, listID, taskID string) error {
	for i := range app.TodoLists {
		if app.TodoLists[i].ID == listID {
			for j := range app.TodoLists[i].Tasks {
				if app.TodoLists[i].Tasks[j].ID == taskID {
					app.TodoLists[i].Tasks[j].Completed = !app.TodoLists[i].Tasks[j].Completed
					app.TodoLists[i].Tasks[j].UpdatedAt = time.Now()
					app.TodoLists[i].UpdatedAt = time.Now()
					return nil
				}
			}
			return fmt.Errorf("task with ID %s not found in list %s", taskID, listID)
		}
	}
	return fmt.Errorf("todo list with ID %s not found", listID)
}

// DeleteTask deletes a task from a todo list
func (s *Storage) DeleteTask(app *models.Application, listID, taskID string) error {
	for i := range app.TodoLists {
		if app.TodoLists[i].ID == listID {
			for j, task := range app.TodoLists[i].Tasks {
				if task.ID == taskID {
					app.TodoLists[i].Tasks = append(app.TodoLists[i].Tasks[:j], app.TodoLists[i].Tasks[j+1:]...)
					app.TodoLists[i].UpdatedAt = time.Now()
					return nil
				}
			}
			return fmt.Errorf("task with ID %s not found in list %s", taskID, listID)
		}
	}
	return fmt.Errorf("todo list with ID %s not found", listID)
}

// Close is a no-op for file storage (satisfies StorageInterface)
func (s *Storage) Close() error {
	return nil
}

// generateID generates a simple unique ID
func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
