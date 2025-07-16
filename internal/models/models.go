package models

import (
	"time"
)

// Priority represents the priority level of a task
type Priority int

const (
	Low Priority = iota
	Medium
	High
	Critical
)

func (p Priority) String() string {
	switch p {
	case Low:
		return "Low"
	case Medium:
		return "Medium"
	case High:
		return "High"
	case Critical:
		return "Critical"
	default:
		return "Unknown"
	}
}

// Task represents a single todo task
type Task struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Description string     `json:"description"`
	Completed   bool       `json:"completed"`
	Priority    Priority   `json:"priority"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	Deadline    *time.Time `json:"deadline,omitempty"`
}

// IsOverdue checks if the task is overdue
func (t *Task) IsOverdue() bool {
	if t.Deadline == nil || t.Completed {
		return false
	}
	return time.Now().After(*t.Deadline)
}

// IsDueSoon checks if the task is due within the next 24 hours
func (t *Task) IsDueSoon() bool {
	if t.Deadline == nil || t.Completed {
		return false
	}
	return time.Now().Add(24*time.Hour).After(*t.Deadline) && !t.IsOverdue()
}

// TodoList represents a collection of tasks
type TodoList struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Tasks       []Task    `json:"tasks"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetCompletedCount returns the number of completed tasks
func (tl *TodoList) GetCompletedCount() int {
	count := 0
	for _, task := range tl.Tasks {
		if task.Completed {
			count++
		}
	}
	return count
}

// GetTotalCount returns the total number of tasks
func (tl *TodoList) GetTotalCount() int {
	return len(tl.Tasks)
}

// GetProgress returns the completion percentage
func (tl *TodoList) GetProgress() float64 {
	if len(tl.Tasks) == 0 {
		return 0
	}
	return float64(tl.GetCompletedCount()) / float64(len(tl.Tasks)) * 100
}

// Application represents the entire application state
type Application struct {
	TodoLists []TodoList `json:"todo_lists"`
	Settings  Settings   `json:"settings"`
}

// Settings represents application settings
type Settings struct {
	ReminderMinutes int  `json:"reminder_minutes"` // Minutes before deadline to remind
	ShowCompleted   bool `json:"show_completed"`   // Whether to show completed tasks
	AutoSave        bool `json:"auto_save"`        // Whether to auto-save changes
}

// DefaultSettings returns default application settings
func DefaultSettings() Settings {
	return Settings{
		ReminderMinutes: 60, // 1 hour before deadline
		ShowCompleted:   true,
		AutoSave:        true,
	}
}
