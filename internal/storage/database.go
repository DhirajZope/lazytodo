package storage

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"

	"github.com/DhirajZope/lazytodo/internal/models"
)

// We'll read migration files from filesystem instead of embedding for now

const (
	DatabaseName = "lazytodo.db"
	DatabaseDir  = ".lazytodo"
)

// DatabaseStorage handles data persistence using SQLite
type DatabaseStorage struct {
	db       *sql.DB
	dataPath string
}

// NewDatabase creates a new database storage instance
func NewDatabase() (*DatabaseStorage, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("failed to get user home directory: %w", err)
	}

	dataDir := filepath.Join(homeDir, DatabaseDir)
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create data directory: %w", err)
	}

	dataPath := filepath.Join(dataDir, DatabaseName)

	// Open database connection
	db, err := sql.Open("sqlite3", dataPath+"?_foreign_keys=on")
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	storage := &DatabaseStorage{
		db:       db,
		dataPath: dataPath,
	}

	// Run migrations
	if err := storage.runMigrations(); err != nil {
		return nil, fmt.Errorf("failed to run migrations: %w", err)
	}

	return storage, nil
}

// Close closes the database connection
func (s *DatabaseStorage) Close() error {
	if s.db != nil {
		return s.db.Close()
	}
	return nil
}

// runMigrations applies database migrations
func (s *DatabaseStorage) runMigrations() error {
	// Create migrations table if it doesn't exist
	createMigrationsTable := `
		CREATE TABLE IF NOT EXISTS schema_migrations (
			version INTEGER PRIMARY KEY,
			applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
		);
	`

	if _, err := s.db.Exec(createMigrationsTable); err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	// Get applied migrations
	appliedMigrations := make(map[int]bool)
	rows, err := s.db.Query("SELECT version FROM schema_migrations")
	if err != nil {
		return fmt.Errorf("failed to query applied migrations: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var version int
		if err := rows.Scan(&version); err != nil {
			return fmt.Errorf("failed to scan migration version: %w", err)
		}
		appliedMigrations[version] = true
	}

	// Apply new migrations from filesystem
	migrationsDir := "migrations"
	migrationEntries, err := os.ReadDir(migrationsDir)
	if err != nil {
		// If migrations directory doesn't exist, create tables directly
		return s.createInitialSchema()
	}

	for _, entry := range migrationEntries {
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".up.sql") {
			// Extract version number from filename (e.g., "001_initial_schema.up.sql" -> 1)
			versionStr := strings.Split(entry.Name(), "_")[0]
			version, err := strconv.Atoi(versionStr)
			if err != nil {
				continue // Skip invalid migration files
			}

			if !appliedMigrations[version] {
				// Read and execute migration
				migrationPath := filepath.Join(migrationsDir, entry.Name())
				migrationSQL, err := os.ReadFile(migrationPath)
				if err != nil {
					return fmt.Errorf("failed to read migration %s: %w", entry.Name(), err)
				}

				if _, err := s.db.Exec(string(migrationSQL)); err != nil {
					return fmt.Errorf("failed to apply migration %s: %w", entry.Name(), err)
				}

				// Record migration as applied
				if _, err := s.db.Exec("INSERT INTO schema_migrations (version) VALUES (?)", version); err != nil {
					return fmt.Errorf("failed to record migration %s: %w", entry.Name(), err)
				}
			}
		}
	}

	return nil
}

// createInitialSchema creates the initial database schema when migrations are not available
func (s *DatabaseStorage) createInitialSchema() error {
	schema := `
-- Create todo_lists table
CREATE TABLE IF NOT EXISTS todo_lists (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    description TEXT DEFAULT '',
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Create tasks table
CREATE TABLE IF NOT EXISTS tasks (
    id TEXT PRIMARY KEY,
    list_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT DEFAULT '',
    completed BOOLEAN NOT NULL DEFAULT FALSE,
    priority INTEGER NOT NULL DEFAULT 0,
    deadline DATETIME NULL,
    created_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (list_id) REFERENCES todo_lists(id) ON DELETE CASCADE
);

-- Create settings table
CREATE TABLE IF NOT EXISTS settings (
    key TEXT PRIMARY KEY,
    value TEXT NOT NULL,
    updated_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Insert default settings
INSERT OR IGNORE INTO settings (key, value) VALUES
    ('reminder_minutes', '60'),
    ('show_completed', 'true'),
    ('auto_save', 'true');

-- Create indexes for better performance
CREATE INDEX IF NOT EXISTS idx_tasks_list_id ON tasks(list_id);
CREATE INDEX IF NOT EXISTS idx_tasks_completed ON tasks(completed);
CREATE INDEX IF NOT EXISTS idx_tasks_deadline ON tasks(deadline);
CREATE INDEX IF NOT EXISTS idx_tasks_priority ON tasks(priority);
	`

	_, err := s.db.Exec(schema)
	if err != nil {
		return fmt.Errorf("failed to create initial schema: %w", err)
	}

	// Record initial migration as applied
	if _, err := s.db.Exec("INSERT INTO schema_migrations (version) VALUES (1)"); err != nil {
		return fmt.Errorf("failed to record initial migration: %w", err)
	}

	return nil
}

// Load loads the application data from database
func (s *DatabaseStorage) Load() (*models.Application, error) {
	app := &models.Application{
		TodoLists: []models.TodoList{},
	}

	// Load settings
	settings, err := s.loadSettings()
	if err != nil {
		return nil, fmt.Errorf("failed to load settings: %w", err)
	}
	app.Settings = settings

	// Load todo lists
	todoLists, err := s.loadTodoLists()
	if err != nil {
		return nil, fmt.Errorf("failed to load todo lists: %w", err)
	}
	app.TodoLists = todoLists

	return app, nil
}

// loadSettings loads application settings from database
func (s *DatabaseStorage) loadSettings() (models.Settings, error) {
	settings := models.DefaultSettings()

	rows, err := s.db.Query("SELECT key, value FROM settings")
	if err != nil {
		return settings, fmt.Errorf("failed to query settings: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var key, value string
		if err := rows.Scan(&key, &value); err != nil {
			continue // Skip invalid settings
		}

		switch key {
		case "reminder_minutes":
			if minutes, err := strconv.Atoi(value); err == nil {
				settings.ReminderMinutes = minutes
			}
		case "show_completed":
			settings.ShowCompleted = value == "true"
		case "auto_save":
			settings.AutoSave = value == "true"
		}
	}

	return settings, nil
}

// loadTodoLists loads all todo lists with their tasks
func (s *DatabaseStorage) loadTodoLists() ([]models.TodoList, error) {
	var todoLists []models.TodoList

	rows, err := s.db.Query(`
		SELECT id, name, description, created_at, updated_at 
		FROM todo_lists 
		ORDER BY created_at ASC
	`)
	if err != nil {
		return nil, fmt.Errorf("failed to query todo lists: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var list models.TodoList
		var createdAt, updatedAt string

		if err := rows.Scan(&list.ID, &list.Name, &list.Description, &createdAt, &updatedAt); err != nil {
			continue // Skip invalid lists
		}

		// Parse timestamps
		if ct, err := time.Parse("2006-01-02 15:04:05", createdAt); err == nil {
			list.CreatedAt = ct
		}
		if ut, err := time.Parse("2006-01-02 15:04:05", updatedAt); err == nil {
			list.UpdatedAt = ut
		}

		// Load tasks for this list
		tasks, err := s.loadTasksForList(list.ID)
		if err != nil {
			return nil, fmt.Errorf("failed to load tasks for list %s: %w", list.ID, err)
		}
		list.Tasks = tasks

		todoLists = append(todoLists, list)
	}

	return todoLists, nil
}

// loadTasksForList loads all tasks for a specific todo list
func (s *DatabaseStorage) loadTasksForList(listID string) ([]models.Task, error) {
	var tasks []models.Task

	rows, err := s.db.Query(`
		SELECT id, title, description, completed, priority, deadline, created_at, updated_at
		FROM tasks 
		WHERE list_id = ? 
		ORDER BY created_at ASC
	`, listID)
	if err != nil {
		return nil, fmt.Errorf("failed to query tasks: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var task models.Task
		var deadline sql.NullString
		var createdAt, updatedAt string

		if err := rows.Scan(
			&task.ID, &task.Title, &task.Description, &task.Completed,
			&task.Priority, &deadline, &createdAt, &updatedAt,
		); err != nil {
			continue // Skip invalid tasks
		}

		// Parse deadline
		if deadline.Valid {
			if dl, err := time.Parse("2006-01-02 15:04:05", deadline.String); err == nil {
				task.Deadline = &dl
			}
		}

		// Parse timestamps
		if ct, err := time.Parse("2006-01-02 15:04:05", createdAt); err == nil {
			task.CreatedAt = ct
		}
		if ut, err := time.Parse("2006-01-02 15:04:05", updatedAt); err == nil {
			task.UpdatedAt = ut
		}

		tasks = append(tasks, task)
	}

	return tasks, nil
}

// Save saves the application data to database
func (s *DatabaseStorage) Save(app *models.Application) error {
	tx, err := s.db.Begin()
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback()

	// This method is used for auto-save, but with database we save immediately
	// on each operation, so this can be a no-op or just ensure data consistency

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

// GetDataPath returns the path to the database file
func (s *DatabaseStorage) GetDataPath() string {
	return s.dataPath
}

// CreateTodoList creates a new todo list
func (s *DatabaseStorage) CreateTodoList(app *models.Application, name, description string) string {
	id := generateDatabaseID()

	_, err := s.db.Exec(`
		INSERT INTO todo_lists (id, name, description) 
		VALUES (?, ?, ?)
	`, id, name, description)

	if err != nil {
		return "" // Return empty string on error
	}

	// Add to in-memory structure for consistency
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
func (s *DatabaseStorage) UpdateTodoList(app *models.Application, listID, name, description string) error {
	_, err := s.db.Exec(`
		UPDATE todo_lists 
		SET name = ?, description = ? 
		WHERE id = ?
	`, name, description, listID)

	if err != nil {
		return fmt.Errorf("failed to update todo list: %w", err)
	}

	// Update in-memory structure
	for i := range app.TodoLists {
		if app.TodoLists[i].ID == listID {
			app.TodoLists[i].Name = name
			app.TodoLists[i].Description = description
			app.TodoLists[i].UpdatedAt = time.Now()
			break
		}
	}

	return nil
}

// DeleteTodoList deletes a todo list and all its tasks
func (s *DatabaseStorage) DeleteTodoList(app *models.Application, listID string) error {
	_, err := s.db.Exec("DELETE FROM todo_lists WHERE id = ?", listID)
	if err != nil {
		return fmt.Errorf("failed to delete todo list: %w", err)
	}

	// Remove from in-memory structure
	for i, list := range app.TodoLists {
		if list.ID == listID {
			app.TodoLists = append(app.TodoLists[:i], app.TodoLists[i+1:]...)
			break
		}
	}

	return nil
}

// CreateTask creates a new task in a todo list
func (s *DatabaseStorage) CreateTask(app *models.Application, listID, title, description string, priority models.Priority, deadline *time.Time) (string, error) {
	taskID := generateDatabaseID()

	var deadlineStr sql.NullString
	if deadline != nil {
		deadlineStr = sql.NullString{String: deadline.Format("2006-01-02 15:04:05"), Valid: true}
	}

	_, err := s.db.Exec(`
		INSERT INTO tasks (id, list_id, title, description, priority, deadline) 
		VALUES (?, ?, ?, ?, ?, ?)
	`, taskID, listID, title, description, int(priority), deadlineStr)

	if err != nil {
		return "", fmt.Errorf("failed to create task: %w", err)
	}

	// Add to in-memory structure
	for i := range app.TodoLists {
		if app.TodoLists[i].ID == listID {
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
			break
		}
	}

	return taskID, nil
}

// UpdateTask updates an existing task
func (s *DatabaseStorage) UpdateTask(app *models.Application, listID, taskID, title, description string, priority models.Priority, deadline *time.Time) error {
	var deadlineStr sql.NullString
	if deadline != nil {
		deadlineStr = sql.NullString{String: deadline.Format("2006-01-02 15:04:05"), Valid: true}
	}

	_, err := s.db.Exec(`
		UPDATE tasks 
		SET title = ?, description = ?, priority = ?, deadline = ? 
		WHERE id = ? AND list_id = ?
	`, title, description, int(priority), deadlineStr, taskID, listID)

	if err != nil {
		return fmt.Errorf("failed to update task: %w", err)
	}

	// Update in-memory structure
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
		}
	}

	return fmt.Errorf("task not found in memory")
}

// ToggleTask toggles the completion status of a task
func (s *DatabaseStorage) ToggleTask(app *models.Application, listID, taskID string) error {
	// First get current status
	var completed bool
	err := s.db.QueryRow("SELECT completed FROM tasks WHERE id = ? AND list_id = ?", taskID, listID).Scan(&completed)
	if err != nil {
		return fmt.Errorf("failed to get task status: %w", err)
	}

	// Toggle it
	newCompleted := !completed
	_, err = s.db.Exec("UPDATE tasks SET completed = ? WHERE id = ? AND list_id = ?", newCompleted, taskID, listID)
	if err != nil {
		return fmt.Errorf("failed to toggle task: %w", err)
	}

	// Update in-memory structure
	for i := range app.TodoLists {
		if app.TodoLists[i].ID == listID {
			for j := range app.TodoLists[i].Tasks {
				if app.TodoLists[i].Tasks[j].ID == taskID {
					app.TodoLists[i].Tasks[j].Completed = newCompleted
					app.TodoLists[i].Tasks[j].UpdatedAt = time.Now()
					app.TodoLists[i].UpdatedAt = time.Now()
					return nil
				}
			}
		}
	}

	return fmt.Errorf("task not found in memory")
}

// DeleteTask deletes a task from a todo list
func (s *DatabaseStorage) DeleteTask(app *models.Application, listID, taskID string) error {
	_, err := s.db.Exec("DELETE FROM tasks WHERE id = ? AND list_id = ?", taskID, listID)
	if err != nil {
		return fmt.Errorf("failed to delete task: %w", err)
	}

	// Remove from in-memory structure
	for i := range app.TodoLists {
		if app.TodoLists[i].ID == listID {
			for j, task := range app.TodoLists[i].Tasks {
				if task.ID == taskID {
					app.TodoLists[i].Tasks = append(app.TodoLists[i].Tasks[:j], app.TodoLists[i].Tasks[j+1:]...)
					app.TodoLists[i].UpdatedAt = time.Now()
					return nil
				}
			}
		}
	}

	return fmt.Errorf("task not found in memory")
}

// generateDatabaseID generates a simple unique ID for database records
func generateDatabaseID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}
