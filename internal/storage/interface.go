package storage

import (
	"time"

	"github.com/DhirajZope/lazytodo/internal/models"
)

// StorageInterface defines the interface that all storage implementations must satisfy
type StorageInterface interface {
	// Load loads the application data
	Load() (*models.Application, error)

	// Save saves the application data
	Save(app *models.Application) error

	// GetDataPath returns the path to the data storage
	GetDataPath() string

	// Todo List operations
	CreateTodoList(app *models.Application, name, description string) string
	UpdateTodoList(app *models.Application, listID, name, description string) error
	DeleteTodoList(app *models.Application, listID string) error

	// Task operations
	CreateTask(app *models.Application, listID, title, description string, priority models.Priority, deadline *time.Time) (string, error)
	UpdateTask(app *models.Application, listID, taskID, title, description string, priority models.Priority, deadline *time.Time) error
	ToggleTask(app *models.Application, listID, taskID string) error
	DeleteTask(app *models.Application, listID, taskID string) error

	// Close closes any resources (for database connections)
	Close() error
}
