package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/DhirajZope/lazytodo/internal/models"
)

// Styles
var (
	titleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("205")).
			MarginBottom(1)

	headerStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("86")).
			MarginBottom(1)

	formStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("62")).
			Padding(1).
			MarginTop(1)

	helpStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("241")).
			MarginTop(1)

	errorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")).
			Bold(true)

	successStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("46")).
			Bold(true)

	priorityStyles = map[models.Priority]lipgloss.Style{
		models.Low:      lipgloss.NewStyle().Foreground(lipgloss.Color("244")),
		models.Medium:   lipgloss.NewStyle().Foreground(lipgloss.Color("220")),
		models.High:     lipgloss.NewStyle().Foreground(lipgloss.Color("208")),
		models.Critical: lipgloss.NewStyle().Foreground(lipgloss.Color("196")).Bold(true),
	}
)

// List item implementations
type listItem struct {
	id          string
	title       string
	description string
	progress    float64
	taskCount   int
}

func (i listItem) FilterValue() string { return i.title }
func (i listItem) Title() string       { return i.title }
func (i listItem) Description() string {
	if i.taskCount == 0 {
		return i.description
	}
	progress := fmt.Sprintf("%.0f%% complete (%d tasks)", i.progress, i.taskCount)
	if i.description != "" {
		return fmt.Sprintf("%s • %s", i.description, progress)
	}
	return progress
}

type taskItem struct {
	id          string
	title       string
	description string
	completed   bool
	priority    models.Priority
	deadline    *time.Time
	overdue     bool
	dueSoon     bool
}

func (i taskItem) FilterValue() string { return i.title }
func (i taskItem) Title() string {
	prefix := "○"
	if i.completed {
		prefix = "✓"
	}

	title := fmt.Sprintf("%s %s", prefix, i.title)

	// Add priority indicator
	if i.priority != models.Low {
		priorityStr := ""
		switch i.priority {
		case models.Medium:
			priorityStr = "⚡"
		case models.High:
			priorityStr = "🔥"
		case models.Critical:
			priorityStr = "🚨"
		}
		title = fmt.Sprintf("%s %s", title, priorityStr)
	}

	// Add deadline indicator
	if i.deadline != nil {
		if i.overdue {
			title = fmt.Sprintf("%s ⚠️", title)
		} else if i.dueSoon {
			title = fmt.Sprintf("%s ⏰", title)
		}
	}

	return title
}

func (i taskItem) Description() string {
	parts := []string{}

	if i.description != "" {
		parts = append(parts, i.description)
	}

	if i.deadline != nil {
		deadlineStr := i.deadline.Format("2006-01-02 15:04")
		if i.overdue {
			deadlineStr = fmt.Sprintf("Due: %s (OVERDUE)", deadlineStr)
		} else if i.dueSoon {
			deadlineStr = fmt.Sprintf("Due: %s (SOON)", deadlineStr)
		} else {
			deadlineStr = fmt.Sprintf("Due: %s", deadlineStr)
		}
		parts = append(parts, deadlineStr)
	}

	return strings.Join(parts, " • ")
}

// updateTodoListsList updates the todo lists list model
func (m *Model) updateTodoListsList() {
	items := make([]list.Item, len(m.app.TodoLists))
	for i, todoList := range m.app.TodoLists {
		items[i] = listItem{
			id:          todoList.ID,
			title:       todoList.Name,
			description: todoList.Description,
			progress:    todoList.GetProgress(),
			taskCount:   todoList.GetTotalCount(),
		}
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		BorderForeground(lipgloss.Color("62")).
		Foreground(lipgloss.Color("86"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		BorderForeground(lipgloss.Color("62")).
		Foreground(lipgloss.Color("244"))

	// Calculate dimensions first
	listWidth := 25  // Default width
	listHeight := 15 // Default height

	// Set size based on sidebar window dimensions if available
	sidebarWindow := m.layout.GetWindow(SidebarWindow)
	if sidebarWindow != nil && sidebarWindow.Position.Width > 0 && sidebarWindow.Position.Height > 0 {
		listWidth = sidebarWindow.Position.Width - 4   // Account for borders
		listHeight = sidebarWindow.Position.Height - 6 // Account for borders and title
	} else if m.width > 0 && m.height > 0 {
		// Fallback to model dimensions
		listWidth = (m.width / 3) - 4
		listHeight = m.height - 8
	}

	// Ensure minimum dimensions
	if listWidth < 10 {
		listWidth = 10
	}
	if listHeight < 5 {
		listHeight = 5
	}

	// Create list with proper dimensions
	m.todoListsList = list.New(items, delegate, listWidth, listHeight)
	m.todoListsList.Title = "📋 Todo Lists"
	m.todoListsList.SetShowStatusBar(false)
	m.todoListsList.SetShowHelp(false)
}

// updateTasksList updates the tasks list model
func (m *Model) updateTasksList() {
	currentList := m.getCurrentList()
	if currentList == nil {
		return
	}

	var items []list.Item
	for _, task := range currentList.Tasks {
		if !m.app.Settings.ShowCompleted && task.Completed {
			continue
		}

		items = append(items, taskItem{
			id:          task.ID,
			title:       task.Title,
			description: task.Description,
			completed:   task.Completed,
			priority:    task.Priority,
			deadline:    task.Deadline,
			overdue:     task.IsOverdue(),
			dueSoon:     task.IsDueSoon(),
		})
	}

	delegate := list.NewDefaultDelegate()
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.
		BorderForeground(lipgloss.Color("62")).
		Foreground(lipgloss.Color("86"))
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.
		BorderForeground(lipgloss.Color("62")).
		Foreground(lipgloss.Color("244"))

	// Calculate dimensions first
	listWidth := 40  // Default width
	listHeight := 15 // Default height

	// Set size based on main window dimensions if available
	mainWindow := m.layout.GetWindow(MainWindow)
	if mainWindow != nil && mainWindow.Position.Width > 0 && mainWindow.Position.Height > 0 {
		listWidth = mainWindow.Position.Width - 4   // Account for borders
		listHeight = mainWindow.Position.Height - 6 // Account for borders and title
	} else if m.width > 0 && m.height > 0 {
		// Fallback to model dimensions
		listWidth = (m.width * 2 / 3) - 4
		listHeight = m.height - 8
	}

	// Ensure minimum dimensions
	if listWidth < 10 {
		listWidth = 10
	}
	if listHeight < 5 {
		listHeight = 5
	}

	// Create list with proper dimensions
	m.tasksList = list.New(items, delegate, listWidth, listHeight)
	m.tasksList.Title = fmt.Sprintf("📝 %s", currentList.Name)
	m.tasksList.SetShowStatusBar(false)
	m.tasksList.SetShowHelp(false)
}

// Lists view - now handles sidebar interaction
func (m *Model) updateListsView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Only handle if sidebar is focused
	if m.layout.GetFocusedWindowID() != SidebarWindow {
		return m, nil
	}

	switch {
	case key.Matches(msg, m.keys.NewList):
		m.resetForm()
		m.state = CreateListView
		return m, nil

	case key.Matches(msg, m.keys.Settings):
		m.previousState = m.state
		m.state = SettingsView
		m.layout.SetFocus(MainWindow)
		return m, nil

	case key.Matches(msg, m.keys.Enter):
		if selected := m.todoListsList.SelectedItem(); selected != nil {
			if item, ok := selected.(listItem); ok {
				m.currentListID = item.id
				m.updateTasksList()
				m.state = TasksView
				m.layout.SetFocus(MainWindow)
				m.showMessageWithType("Switched to "+item.title, "success")
				return m, nil
			}
		}

	case key.Matches(msg, m.keys.Edit):
		if selected := m.todoListsList.SelectedItem(); selected != nil {
			if item, ok := selected.(listItem); ok {
				m.currentListID = item.id
				m.prepareEditListForm()
				m.state = EditListView
				return m, nil
			}
		}

	case key.Matches(msg, m.keys.Delete):
		if selected := m.todoListsList.SelectedItem(); selected != nil {
			if item, ok := selected.(listItem); ok {
				if err := m.storage.DeleteTodoList(m.app, item.id); err != nil {
					m.showMessageWithType(fmt.Sprintf("Error: %v", err), "error")
				} else {
					m.updateTodoListsList()
					m.showMessageWithType("List deleted successfully", "success")
					return m, m.saveData()
				}
			}
		}
	}

	var cmd tea.Cmd
	m.todoListsList, cmd = m.todoListsList.Update(msg)
	return m, cmd
}

func (m *Model) renderListsView() string {
	var listView string
	if m.todoListsList.Items() != nil {
		listView = m.todoListsList.View()
	} else {
		listView = "Loading..."
	}

	content := []string{
		titleStyle.Render("🎯 LazyTodo - Smart Todo Manager"),
		listView,
		m.getStatusBar(),
	}

	if m.app != nil && len(m.app.TodoLists) == 0 {
		emptyMsg := helpStyle.Render("No todo lists yet. Press 'n' to create your first list!")
		content = append(content[:len(content)-1], emptyMsg, content[len(content)-1])
	}

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

// Tasks view - now handles main window interaction
func (m *Model) updateTasksView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	// Only handle if main window is focused
	if m.layout.GetFocusedWindowID() != MainWindow {
		return m, nil
	}

	switch {
	case key.Matches(msg, m.keys.Back):
		m.layout.SetFocus(SidebarWindow)
		return m, nil

	case key.Matches(msg, m.keys.NewTask):
		m.resetForm()
		m.state = CreateTaskView
		return m, nil

	case key.Matches(msg, m.keys.Edit):
		if selected := m.tasksList.SelectedItem(); selected != nil {
			if item, ok := selected.(taskItem); ok {
				m.editingTaskID = item.id
				m.prepareEditTaskForm()
				m.state = EditTaskView
				return m, nil
			}
		}

	case key.Matches(msg, m.keys.Toggle):
		if selected := m.tasksList.SelectedItem(); selected != nil {
			if item, ok := selected.(taskItem); ok {
				if err := m.storage.ToggleTask(m.app, m.currentListID, item.id); err != nil {
					m.showMessageWithType(fmt.Sprintf("Error: %v", err), "error")
				} else {
					m.updateTasksList()
					status := "completed"
					msgType := "success"
					if item.completed {
						status = "uncompleted"
						msgType = "info"
					}
					m.showMessageWithType(fmt.Sprintf("Task %s", status), msgType)
					return m, m.saveData()
				}
			}
		}

	case key.Matches(msg, m.keys.Delete):
		if selected := m.tasksList.SelectedItem(); selected != nil {
			if item, ok := selected.(taskItem); ok {
				if err := m.storage.DeleteTask(m.app, m.currentListID, item.id); err != nil {
					m.showMessageWithType(fmt.Sprintf("Error: %v", err), "error")
				} else {
					m.updateTasksList()
					m.showMessageWithType("Task deleted successfully", "success")
					return m, m.saveData()
				}
			}
		}
	}

	var cmd tea.Cmd
	// Only update tasksList if it's initialized
	if m.tasksList.Items() != nil {
		m.tasksList, cmd = m.tasksList.Update(msg)
	}
	return m, cmd
}

func (m *Model) renderTasksView() string {
	var taskView string
	if m.tasksList.Items() != nil {
		taskView = m.tasksList.View()
	} else {
		taskView = "Loading..."
	}

	content := []string{
		titleStyle.Render("🎯 LazyTodo - Task View"),
		taskView,
		m.getStatusBar(),
	}

	currentList := m.getCurrentList()
	if currentList != nil && len(currentList.Tasks) == 0 {
		emptyMsg := helpStyle.Render("No tasks yet. Press 'a' to add your first task!")
		content = append(content[:len(content)-1], emptyMsg, content[len(content)-1])
	}

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

// Form handling
func (m *Model) resetForm() {
	m.titleInput.SetValue("")
	m.descriptionInput.SetValue("")
	m.deadlineInput.SetValue("")
	m.formFocusIndex = 0
	m.titleInput.Focus()
	m.descriptionInput.Blur()
	m.deadlineInput.Blur()
	m.editing = false
	m.editingTaskID = ""
}

func (m *Model) prepareEditListForm() {
	if currentList := m.getCurrentList(); currentList != nil {
		m.titleInput.SetValue(currentList.Name)
		m.descriptionInput.SetValue(currentList.Description)
		m.deadlineInput.SetValue("")
		m.formFocusIndex = 0
		m.titleInput.Focus()
		m.descriptionInput.Blur()
		m.deadlineInput.Blur()
		m.editing = true
	}
}

func (m *Model) prepareEditTaskForm() {
	currentList := m.getCurrentList()
	if currentList == nil {
		return
	}

	for _, task := range currentList.Tasks {
		if task.ID == m.editingTaskID {
			m.titleInput.SetValue(task.Title)
			m.descriptionInput.SetValue(task.Description)
			if task.Deadline != nil {
				m.deadlineInput.SetValue(task.Deadline.Format("2006-01-02 15:04"))
			} else {
				m.deadlineInput.SetValue("")
			}
			m.formFocusIndex = 0
			m.titleInput.Focus()
			m.descriptionInput.Blur()
			m.deadlineInput.Blur()
			m.editing = true
			break
		}
	}
}

// List form - now works with form window overlay
func (m *Model) updateListForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		m.state = ListsView
		m.layout.SetFocus(SidebarWindow)
		return m, nil

	case key.Matches(msg, m.keys.Tab):
		m.formFocusIndex = (m.formFocusIndex + 1) % 2
		m.updateFormFocus()
		return m, nil

	case key.Matches(msg, m.keys.ShiftTab):
		m.formFocusIndex = (m.formFocusIndex - 1 + 2) % 2
		m.updateFormFocus()
		return m, nil

	case key.Matches(msg, m.keys.Enter):
		if m.titleInput.Value() == "" {
			m.showMessageWithType("Title is required", "warning")
			return m, nil
		}

		if m.editing {
			// Update existing list
			err := m.storage.UpdateTodoList(m.app, m.currentListID, m.titleInput.Value(), m.descriptionInput.Value())
			if err != nil {
				m.showMessageWithType(fmt.Sprintf("Error: %v", err), "error")
				return m, nil
			}
			m.showMessageWithType("List updated successfully", "success")
		} else {
			// Create new list
			m.storage.CreateTodoList(m.app, m.titleInput.Value(), m.descriptionInput.Value())
			m.showMessageWithType("List created successfully", "success")
		}

		m.updateTodoListsList()
		m.state = ListsView
		m.layout.SetFocus(SidebarWindow)
		return m, m.saveData()
	}

	var cmd tea.Cmd
	switch m.formFocusIndex {
	case 0:
		m.titleInput, cmd = m.titleInput.Update(msg)
	case 1:
		m.descriptionInput, cmd = m.descriptionInput.Update(msg)
	}

	return m, cmd
}

func (m *Model) renderListForm() string {
	title := "Create New List"
	if m.editing {
		title = "Edit List"
	}

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		"Title:",
		m.titleInput.View(),
		"",
		"Description:",
		m.descriptionInput.View(),
		"",
		helpStyle.Render("Tab/Shift+Tab: Navigate • Enter: Save • Esc: Cancel"),
	)

	content := []string{
		titleStyle.Render(title),
		formStyle.Render(form),
		m.getStatusBar(),
	}

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

// Task form - now works with form window overlay
func (m *Model) updateTaskForm(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		m.state = TasksView
		m.layout.SetFocus(MainWindow)
		return m, nil

	case key.Matches(msg, m.keys.Tab):
		m.formFocusIndex = (m.formFocusIndex + 1) % 3
		m.updateFormFocus()
		return m, nil

	case key.Matches(msg, m.keys.ShiftTab):
		m.formFocusIndex = (m.formFocusIndex - 1 + 3) % 3
		m.updateFormFocus()
		return m, nil

	case key.Matches(msg, m.keys.Enter):
		if m.titleInput.Value() == "" {
			m.showMessageWithType("Title is required", "warning")
			return m, nil
		}

		var deadline *time.Time
		if m.deadlineInput.Value() != "" {
			if parsed, err := time.Parse("2006-01-02 15:04", m.deadlineInput.Value()); err != nil {
				m.showMessageWithType("Invalid deadline format (use YYYY-MM-DD HH:MM)", "warning")
				return m, nil
			} else {
				deadline = &parsed
			}
		}

		if m.editing {
			// Update existing task
			err := m.storage.UpdateTask(m.app, m.currentListID, m.editingTaskID,
				m.titleInput.Value(), m.descriptionInput.Value(), models.Medium, deadline)
			if err != nil {
				m.showMessageWithType(fmt.Sprintf("Error: %v", err), "error")
				return m, nil
			}
			m.showMessageWithType("Task updated successfully", "success")
		} else {
			// Create new task
			_, err := m.storage.CreateTask(m.app, m.currentListID,
				m.titleInput.Value(), m.descriptionInput.Value(), models.Medium, deadline)
			if err != nil {
				m.showMessageWithType(fmt.Sprintf("Error: %v", err), "error")
				return m, nil
			}
			m.showMessageWithType("Task created successfully", "success")
		}

		m.updateTasksList()
		m.state = TasksView
		m.layout.SetFocus(MainWindow)
		return m, m.saveData()
	}

	var cmd tea.Cmd
	switch m.formFocusIndex {
	case 0:
		m.titleInput, cmd = m.titleInput.Update(msg)
	case 1:
		m.descriptionInput, cmd = m.descriptionInput.Update(msg)
	case 2:
		m.deadlineInput, cmd = m.deadlineInput.Update(msg)
	}

	return m, cmd
}

func (m *Model) renderTaskForm() string {
	title := "Create New Task"
	if m.editing {
		title = "Edit Task"
	}

	form := lipgloss.JoinVertical(
		lipgloss.Left,
		"Title:",
		m.titleInput.View(),
		"",
		"Description:",
		m.descriptionInput.View(),
		"",
		"Deadline (YYYY-MM-DD HH:MM):",
		m.deadlineInput.View(),
		"",
		helpStyle.Render("Tab/Shift+Tab: Navigate • Enter: Save • Esc: Cancel"),
	)

	content := []string{
		titleStyle.Render(title),
		formStyle.Render(form),
		m.getStatusBar(),
	}

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

func (m *Model) updateFormFocus() {
	switch m.formFocusIndex {
	case 0:
		m.titleInput.Focus()
		m.descriptionInput.Blur()
		m.deadlineInput.Blur()
	case 1:
		m.titleInput.Blur()
		m.descriptionInput.Focus()
		m.deadlineInput.Blur()
	case 2:
		m.titleInput.Blur()
		m.descriptionInput.Blur()
		m.deadlineInput.Focus()
	}
}

// Settings view
func (m *Model) updateSettingsView(msg tea.KeyMsg) (tea.Model, tea.Cmd) {
	switch {
	case key.Matches(msg, m.keys.Back):
		m.state = TasksView
		m.layout.SetFocus(MainWindow)
		return m, nil
	}

	return m, nil
}

func (m *Model) renderSettingsView() string {
	settings := []string{
		fmt.Sprintf("Reminder Minutes: %d", m.app.Settings.ReminderMinutes),
		fmt.Sprintf("Show Completed: %v", m.app.Settings.ShowCompleted),
		fmt.Sprintf("Auto Save: %v", m.app.Settings.AutoSave),
	}

	content := []string{
		titleStyle.Render("⚙️ Settings"),
		headerStyle.Render("Current Settings:"),
		strings.Join(settings, "\n"),
		"",
		helpStyle.Render("Settings can be modified by editing the data file directly"),
		helpStyle.Render("Press Esc to go back"),
		m.getStatusBar(),
	}

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}

// Help view
func (m *Model) renderHelpView() string {
	keyHelp := []string{
		"📋 Todo Lists View:",
		"  n - Create new list",
		"  Enter - Open list",
		"  e - Edit list",
		"  d - Delete list",
		"  s - Settings",
		"",
		"📝 Tasks View:",
		"  a - Add new task",
		"  Space - Toggle task completion",
		"  e - Edit task",
		"  d - Delete task",
		"  Esc - Back to lists",
		"",
		"📝 Forms:",
		"  Tab/Shift+Tab - Navigate fields",
		"  Enter - Save",
		"  Esc - Cancel",
		"",
		"🌐 Global:",
		"  ? - Toggle help",
		"  q/Ctrl+C - Quit",
		"",
		"🎨 Visual Indicators:",
		"  ○ - Incomplete task",
		"  ✓ - Complete task",
		"  ⚡ - Medium priority",
		"  🔥 - High priority",
		"  🚨 - Critical priority",
		"  ⏰ - Task due soon",
		"  ⚠️ - Task overdue",
	}

	content := []string{
		titleStyle.Render("❓ Help"),
		strings.Join(keyHelp, "\n"),
		"",
		helpStyle.Render("Press ? or Esc to close help"),
		m.getStatusBar(),
	}

	return lipgloss.JoinVertical(lipgloss.Left, content...)
}
