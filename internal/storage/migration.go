package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/DhirajZope/lazytodo/internal/models"
)

// MigrateFromJSON migrates data from the old JSON file format to the database
func MigrateFromJSON(dbStorage *DatabaseStorage) error {
	// Check if JSON file exists
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return fmt.Errorf("failed to get user home directory: %w", err)
	}

	jsonPath := filepath.Join(homeDir, ".lazytodo", "lazytodo.json")
	if _, err := os.Stat(jsonPath); os.IsNotExist(err) {
		// No JSON file to migrate
		return nil
	}

	fmt.Printf("Found existing JSON data file. Migrating to database...\n")

	// Read JSON file
	data, err := os.ReadFile(jsonPath)
	if err != nil {
		return fmt.Errorf("failed to read JSON file: %w", err)
	}

	// Parse JSON data
	var jsonApp models.Application
	if err := json.Unmarshal(data, &jsonApp); err != nil {
		return fmt.Errorf("failed to parse JSON data: %w", err)
	}

	// Start database transaction
	tx, err := dbStorage.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// Migrate settings
	for key, value := range map[string]string{
		"reminder_minutes": fmt.Sprintf("%d", jsonApp.Settings.ReminderMinutes),
		"show_completed":   fmt.Sprintf("%t", jsonApp.Settings.ShowCompleted),
		"auto_save":        fmt.Sprintf("%t", jsonApp.Settings.AutoSave),
	} {
		_, err := tx.Exec("INSERT OR REPLACE INTO settings (key, value) VALUES (?, ?)", key, value)
		if err != nil {
			return fmt.Errorf("failed to migrate setting %s: %w", key, err)
		}
	}

	// Migrate todo lists and tasks
	for _, list := range jsonApp.TodoLists {
		// Insert todo list
		_, err := tx.Exec(`
			INSERT OR REPLACE INTO todo_lists (id, name, description, created_at, updated_at) 
			VALUES (?, ?, ?, ?, ?)
		`, list.ID, list.Name, list.Description,
			list.CreatedAt.Format("2006-01-02 15:04:05"),
			list.UpdatedAt.Format("2006-01-02 15:04:05"))

		if err != nil {
			return fmt.Errorf("failed to migrate todo list %s: %w", list.Name, err)
		}

		// Insert tasks for this list
		for _, task := range list.Tasks {
			var deadlineStr *string
			if task.Deadline != nil {
				dl := task.Deadline.Format("2006-01-02 15:04:05")
				deadlineStr = &dl
			}

			_, err := tx.Exec(`
				INSERT OR REPLACE INTO tasks 
				(id, list_id, title, description, completed, priority, deadline, created_at, updated_at) 
				VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
			`, task.ID, list.ID, task.Title, task.Description, task.Completed,
				int(task.Priority), deadlineStr,
				task.CreatedAt.Format("2006-01-02 15:04:05"),
				task.UpdatedAt.Format("2006-01-02 15:04:05"))

			if err != nil {
				return fmt.Errorf("failed to migrate task %s: %w", task.Title, err)
			}
		}
	}

	// Commit transaction
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit migration: %w", err)
	}

	// Create backup of JSON file and remove original
	backupPath := jsonPath + ".backup." + time.Now().Format("20060102-150405")
	if err := os.Rename(jsonPath, backupPath); err != nil {
		fmt.Printf("Warning: failed to backup JSON file: %v\n", err)
	} else {
		fmt.Printf("Migration completed! JSON file backed up to: %s\n", backupPath)
	}

	fmt.Printf("Successfully migrated %d todo lists to database.\n", len(jsonApp.TodoLists))
	return nil
}

// NewWithMigration creates a new database storage and automatically migrates from JSON if needed
func NewWithMigration() (StorageInterface, error) {
	dbStorage, err := NewDatabase()
	if err != nil {
		return nil, fmt.Errorf("failed to create database storage: %w", err)
	}

	// Check if we need to migrate from JSON
	if err := MigrateFromJSON(dbStorage); err != nil {
		dbStorage.Close()
		return nil, fmt.Errorf("failed to migrate from JSON: %w", err)
	}

	return dbStorage, nil
}

// GetStorageInfo returns information about the current storage backend
func GetStorageInfo(storage StorageInterface) string {
	switch s := storage.(type) {
	case *DatabaseStorage:
		return fmt.Sprintf("Database: %s", s.GetDataPath())
	case *Storage:
		return fmt.Sprintf("JSON File: %s", s.GetDataPath())
	default:
		return "Unknown storage backend"
	}
}
