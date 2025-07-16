package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/DhirajZope/lazytodo/internal/models"
	"github.com/charmbracelet/lipgloss"
)

// renderSidebarContent renders the sidebar content (todo lists)
func (m *Model) renderSidebarContent() string {
	if m.app == nil {
		return BaseSubtitleStyle.Render("Loading application...")
	}

	if len(m.app.TodoLists) == 0 {
		welcomeTitle := BaseTitleStyle.Copy().
			Foreground(PrimaryColor).
			Render("🎯 Welcome to LazyTodo!")

		emptyMsg := BaseSubtitleStyle.Copy().
			Foreground(AccentColor).
			Render("📋 Ready to get organized?")

		instructions := []string{
			"✨ Create your first todo list to get started",
			"",
			"📝 Press 'n' to create a new list",
			"🎯 Organize tasks by projects or contexts",
			"⚡ Set priorities and deadlines",
			"📊 Track your progress",
		}

		instructionText := DescStyle.Render(strings.Join(instructions, "\n"))

		return lipgloss.JoinVertical(
			lipgloss.Left,
			welcomeTitle,
			"",
			emptyMsg,
			"",
			instructionText,
		)
	}

	// Use the bubble tea list for sidebar if we're in list view
	var content string
	if m.todoListsList.Items() != nil {
		content = m.todoListsList.View()
	} else {
		// Fallback rendering
		var lines []string
		for i, todoList := range m.app.TodoLists {
			icon := "📋"
			title := todoList.Name
			subtitle := fmt.Sprintf("%.0f%% complete (%d tasks)",
				todoList.GetProgress(), todoList.GetTotalCount())

			selected := (m.currentListID == todoList.ID)
			item := RenderEnhancedListItem(icon, title, subtitle, selected, false)

			if i == 0 && m.currentListID == "" {
				m.currentListID = todoList.ID
				m.updateTasksList()
			}

			lines = append(lines, item)
		}
		content = lipgloss.JoinVertical(lipgloss.Left, lines...)
	}

	return content
}

// renderMainContent renders the main window content based on current state
func (m *Model) renderMainContent() string {
	switch m.state {
	case ListsView, TasksView:
		return m.renderTasksContent()
	case SettingsView:
		return m.renderSettingsContent()
	default:
		return m.renderTasksContent()
	}
}

// renderTasksContent renders the task list in the main window
func (m *Model) renderTasksContent() string {
	if m.app == nil {
		return BaseSubtitleStyle.Render("Loading...")
	}

	if len(m.app.TodoLists) == 0 {
		welcomeTitle := BaseTitleStyle.Copy().
			Foreground(PrimaryColor).
			Render("🚀 Let's Get Started!")

		emptyMsg := BaseSubtitleStyle.Copy().
			Foreground(AccentColor).
			Render("📝 No todo lists yet")

		hint := DescStyle.Render("Create your first list in the sidebar (Press Ctrl+S to focus sidebar)")

		features := []string{
			"🎯 LazyTodo Features:",
			"",
			"📋 Multiple todo lists for different projects",
			"⚡ Priority levels (Low, Medium, High, Critical)",
			"📅 Deadline tracking with smart reminders",
			"✅ Progress tracking and completion stats",
			"💾 Automatic SQLite database storage",
			"🎨 Beautiful multi-window interface",
		}

		featureText := DescStyle.Render(strings.Join(features, "\n"))

		return lipgloss.JoinVertical(
			lipgloss.Left,
			welcomeTitle,
			"",
			emptyMsg,
			"",
			hint,
			"",
			featureText,
		)
	}

	currentList := m.getCurrentList()
	if currentList == nil {
		emptyMsg := BaseSubtitleStyle.Render("📋 Select a todo list from the sidebar")
		hint := DescStyle.Render("Use Ctrl+S to focus the sidebar")
		return lipgloss.JoinVertical(lipgloss.Left, emptyMsg, "", hint)
	}

	// Update main window title
	m.layout.SetWindowTitle(MainWindow, fmt.Sprintf("📝 %s", currentList.Name))

	if len(currentList.Tasks) == 0 {
		emptyMsg := BaseSubtitleStyle.Render("No tasks yet")
		hint := DescStyle.Render("Press 'a' to add your first task")
		return lipgloss.JoinVertical(lipgloss.Center, emptyMsg, "", hint)
	}

	// Use the bubble tea list for main content if available
	if m.tasksList.Items() != nil {
		return m.tasksList.View()
	}

	// Fallback rendering
	var lines []string
	for _, task := range currentList.Tasks {
		if !m.app.Settings.ShowCompleted && task.Completed {
			continue
		}

		icon := "○"
		if task.Completed {
			icon = "✓"
		}

		title := task.Title
		subtitle := ""

		// Add description if present
		if task.Description != "" {
			subtitle = task.Description
		}

		// Add deadline info
		if task.Deadline != nil {
			deadlineStr := task.Deadline.Format("2006-01-02 15:04")
			if task.IsOverdue() {
				deadlineStr = "⚠️ Due: " + deadlineStr + " (OVERDUE)"
			} else if task.IsDueSoon() {
				deadlineStr = "⏰ Due: " + deadlineStr + " (SOON)"
			} else {
				deadlineStr = "📅 Due: " + deadlineStr
			}

			if subtitle != "" {
				subtitle += " • " + deadlineStr
			} else {
				subtitle = deadlineStr
			}
		}

		// Add priority indicator
		if task.Priority != models.Low {
			priorityStr := ""
			switch task.Priority {
			case models.Medium:
				priorityStr = "⚡ Medium"
			case models.High:
				priorityStr = "🔥 High"
			case models.Critical:
				priorityStr = "🚨 Critical"
			}

			if subtitle != "" {
				subtitle += " • " + priorityStr
			} else {
				subtitle = priorityStr
			}
		}

		item := RenderEnhancedListItem(icon, title, subtitle, false, task.Completed)
		lines = append(lines, item)
	}

	if len(lines) == 0 {
		emptyMsg := BaseSubtitleStyle.Render("All tasks completed!")
		hint := DescStyle.Render("Toggle 'Show Completed' in settings to see completed tasks")
		return lipgloss.JoinVertical(lipgloss.Center, emptyMsg, "", hint)
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// renderSettingsContent renders the settings view
func (m *Model) renderSettingsContent() string {
	m.layout.SetWindowTitle(MainWindow, "⚙️ Settings")

	var lines []string
	lines = append(lines, BaseTitleStyle.Render("📝 Application Settings"))
	lines = append(lines, "")

	// Settings display
	settings := []string{
		fmt.Sprintf("Reminder Minutes: %d", m.app.Settings.ReminderMinutes),
		fmt.Sprintf("Show Completed: %v", m.app.Settings.ShowCompleted),
		fmt.Sprintf("Auto Save: %v", m.app.Settings.AutoSave),
	}

	for _, setting := range settings {
		lines = append(lines, "  "+DescStyle.Render(setting))
	}

	lines = append(lines, "")
	lines = append(lines, BaseSubtitleStyle.Render("Settings can be modified by editing the database directly"))
	lines = append(lines, DescStyle.Render("Press Esc to go back to task view"))

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// renderStatusContent renders the status bar content
func (m *Model) renderStatusContent() string {
	// Show message if recent
	if time.Since(m.messageTime) < 3*time.Second && m.message != "" {
		return StyleStatusMessage(m.message, m.messageType)
	}

	// Build status information
	var statusParts []string

	// Current state info
	if m.app != nil {
		switch m.state {
		case ListsView, TasksView:
			statusParts = append(statusParts,
				fmt.Sprintf("Lists: %d", len(m.app.TodoLists)))

			if currentList := m.getCurrentList(); currentList != nil {
				statusParts = append(statusParts,
					fmt.Sprintf("Tasks: %d/%d",
						currentList.GetCompletedCount(),
						currentList.GetTotalCount()))
			}
		case SettingsView:
			statusParts = append(statusParts, "Settings")
		}
	}

	// Window focus indicator
	focusedWindow := m.layout.GetFocusedWindowID()
	switch focusedWindow {
	case SidebarWindow:
		statusParts = append(statusParts, "Focus: Sidebar")
	case MainWindow:
		statusParts = append(statusParts, "Focus: Main")
	case FormWindow:
		statusParts = append(statusParts, "Focus: Form")
	case HelpWindow:
		statusParts = append(statusParts, "Focus: Help")
	}

	// Key hints
	keyHints := []string{
		KeyStyle.Render("?") + " Help",
		KeyStyle.Render("Ctrl+→/←") + " Windows",
		KeyStyle.Render("q") + " Quit",
	}

	// Combine status parts and key hints
	status := strings.Join(statusParts, " • ")
	hints := strings.Join(keyHints, "  ")

	// Use available width to balance status and hints
	totalContent := status + "    " + hints

	return BaseContentStyle.Render(totalContent)
}

// renderFormContent renders form content for overlays
func (m *Model) renderFormContent() string {
	switch m.state {
	case CreateListView, EditListView:
		return m.renderListFormContent()
	case CreateTaskView, EditTaskView:
		return m.renderTaskFormContent()
	default:
		return ""
	}
}

// renderListFormContent renders the todo list form
func (m *Model) renderListFormContent() string {
	title := "Create New List"
	if m.editing {
		title = "Edit List"
		m.layout.SetWindowTitle(FormWindow, "✏️ Edit List")
	} else {
		m.layout.SetWindowTitle(FormWindow, "➕ Create List")
	}

	var lines []string
	lines = append(lines, BaseTitleStyle.Render(title))
	lines = append(lines, "")

	// Title field
	titleLabel := FormLabel.Render("Title:")
	var titleField string
	if m.formFocusIndex == 0 {
		titleField = FormFieldFocused.Render(m.titleInput.View())
	} else {
		titleField = FormFieldUnfocused.Render(m.titleInput.View())
	}
	lines = append(lines, titleLabel)
	lines = append(lines, titleField)
	lines = append(lines, "")

	// Description field
	descLabel := FormLabel.Render("Description:")
	var descField string
	if m.formFocusIndex == 1 {
		descField = FormFieldFocused.Render(m.descriptionInput.View())
	} else {
		descField = FormFieldUnfocused.Render(m.descriptionInput.View())
	}
	lines = append(lines, descLabel)
	lines = append(lines, descField)
	lines = append(lines, "")

	// Help text
	helpText := CreateHelpSection("Form Controls", map[string]string{
		"Tab/Shift+Tab": "Navigate fields",
		"Enter":         "Save",
		"Esc":           "Cancel",
	})
	lines = append(lines, helpText)

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// renderTaskFormContent renders the task form
func (m *Model) renderTaskFormContent() string {
	title := "Create New Task"
	if m.editing {
		title = "Edit Task"
		m.layout.SetWindowTitle(FormWindow, "✏️ Edit Task")
	} else {
		m.layout.SetWindowTitle(FormWindow, "➕ Create Task")
	}

	var lines []string
	lines = append(lines, BaseTitleStyle.Render(title))
	lines = append(lines, "")

	// Title field
	titleLabel := FormLabel.Render("Title:")
	var titleField string
	if m.formFocusIndex == 0 {
		titleField = FormFieldFocused.Render(m.titleInput.View())
	} else {
		titleField = FormFieldUnfocused.Render(m.titleInput.View())
	}
	lines = append(lines, titleLabel)
	lines = append(lines, titleField)
	lines = append(lines, "")

	// Description field
	descLabel := FormLabel.Render("Description:")
	var descField string
	if m.formFocusIndex == 1 {
		descField = FormFieldFocused.Render(m.descriptionInput.View())
	} else {
		descField = FormFieldUnfocused.Render(m.descriptionInput.View())
	}
	lines = append(lines, descLabel)
	lines = append(lines, descField)
	lines = append(lines, "")

	// Deadline field
	deadlineLabel := FormLabel.Render("Deadline (YYYY-MM-DD HH:MM):")
	var deadlineField string
	if m.formFocusIndex == 2 {
		deadlineField = FormFieldFocused.Render(m.deadlineInput.View())
	} else {
		deadlineField = FormFieldUnfocused.Render(m.deadlineInput.View())
	}
	lines = append(lines, deadlineLabel)
	lines = append(lines, deadlineField)
	lines = append(lines, "")

	// Help text
	helpText := CreateHelpSection("Form Controls", map[string]string{
		"Tab/Shift+Tab": "Navigate fields",
		"Enter":         "Save",
		"Esc":           "Cancel",
	})
	lines = append(lines, helpText)

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// isInFormState checks if we're currently in a form state
func (m *Model) isInFormState() bool {
	switch m.state {
	case CreateListView, EditListView, CreateTaskView, EditTaskView:
		return true
	default:
		return false
	}
}

// showMessage displays a status message with type
func (m *Model) showMessageWithType(msg, msgType string) {
	m.message = msg
	m.messageType = msgType
	m.messageTime = time.Now()
}
