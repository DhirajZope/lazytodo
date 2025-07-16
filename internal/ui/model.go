package ui

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/DhirajZope/lazytodo/internal/models"
	"github.com/DhirajZope/lazytodo/internal/storage"
)

// ViewState represents the current view/state of the application
type ViewState int

const (
	ListsView ViewState = iota
	TasksView
	CreateListView
	EditListView
	CreateTaskView
	EditTaskView
	SettingsView
	HelpView
)

// Model represents the main application model
type Model struct {
	// Application state
	app     *models.Application
	storage storage.StorageInterface

	// Current view state
	state         ViewState
	previousState ViewState

	// Multi-window layout system
	layout       *Layout
	windowStyles map[WindowID]WindowStyle

	// List views
	todoListsList list.Model
	tasksList     list.Model

	// Currently selected list
	currentListID string

	// Form inputs
	titleInput       textinput.Model
	descriptionInput textinput.Model
	deadlineInput    textinput.Model

	// Form states
	formFocusIndex int
	editing        bool
	editingTaskID  string

	// UI dimensions
	width  int
	height int

	// Error and status messages
	message     string
	messageTime time.Time
	messageType string

	// Reminder system
	lastReminderCheck time.Time

	// Key bindings
	keys KeyMap
}

// KeyMap defines the key bindings for the application
type KeyMap struct {
	Up           key.Binding
	Down         key.Binding
	Left         key.Binding
	Right        key.Binding
	Enter        key.Binding
	Back         key.Binding
	Quit         key.Binding
	Help         key.Binding
	NewList      key.Binding
	NewTask      key.Binding
	Edit         key.Binding
	Delete       key.Binding
	Toggle       key.Binding
	Settings     key.Binding
	Tab          key.Binding
	ShiftTab     key.Binding
	NextWindow   key.Binding
	PrevWindow   key.Binding
	FocusMain    key.Binding
	FocusSidebar key.Binding
}

// DefaultKeyMap returns the default key mappings
func DefaultKeyMap() KeyMap {
	return KeyMap{
		Up: key.NewBinding(
			key.WithKeys("up", "k"),
			key.WithHelp("↑/k", "move up"),
		),
		Down: key.NewBinding(
			key.WithKeys("down", "j"),
			key.WithHelp("↓/j", "move down"),
		),
		Left: key.NewBinding(
			key.WithKeys("left", "h"),
			key.WithHelp("←/h", "go back"),
		),
		Right: key.NewBinding(
			key.WithKeys("right", "l"),
			key.WithHelp("→/l", "enter"),
		),
		Enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
		),
		Back: key.NewBinding(
			key.WithKeys("esc"),
			key.WithHelp("esc", "back"),
		),
		Quit: key.NewBinding(
			key.WithKeys("ctrl+c", "q"),
			key.WithHelp("q/ctrl+c", "quit"),
		),
		Help: key.NewBinding(
			key.WithKeys("?"),
			key.WithHelp("?", "help"),
		),
		NewList: key.NewBinding(
			key.WithKeys("n"),
			key.WithHelp("n", "new list"),
		),
		NewTask: key.NewBinding(
			key.WithKeys("a"),
			key.WithHelp("a", "add task"),
		),
		Edit: key.NewBinding(
			key.WithKeys("e"),
			key.WithHelp("e", "edit"),
		),
		Delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
		Toggle: key.NewBinding(
			key.WithKeys(" "),
			key.WithHelp("space", "toggle"),
		),
		Settings: key.NewBinding(
			key.WithKeys("s"),
			key.WithHelp("s", "settings"),
		),
		Tab: key.NewBinding(
			key.WithKeys("tab"),
			key.WithHelp("tab", "next field"),
		),
		ShiftTab: key.NewBinding(
			key.WithKeys("shift+tab"),
			key.WithHelp("shift+tab", "prev field"),
		),
		NextWindow: key.NewBinding(
			key.WithKeys("ctrl+right", "ctrl+l"),
			key.WithHelp("ctrl+→", "next window"),
		),
		PrevWindow: key.NewBinding(
			key.WithKeys("ctrl+left", "ctrl+h"),
			key.WithHelp("ctrl+←", "prev window"),
		),
		FocusMain: key.NewBinding(
			key.WithKeys("ctrl+m"),
			key.WithHelp("ctrl+m", "focus main"),
		),
		FocusSidebar: key.NewBinding(
			key.WithKeys("ctrl+s"),
			key.WithHelp("ctrl+s", "focus sidebar"),
		),
	}
}

// NewModel creates a new application model
func NewModel() (*Model, error) {
	storage, err := storage.NewWithMigration()
	if err != nil {
		return nil, fmt.Errorf("failed to create storage: %w", err)
	}

	app, err := storage.Load()
	if err != nil {
		return nil, fmt.Errorf("failed to load application data: %w", err)
	}

	// Create text inputs
	titleInput := textinput.New()
	titleInput.Placeholder = "Enter title..."
	titleInput.Focus()

	descriptionInput := textinput.New()
	descriptionInput.Placeholder = "Enter description (optional)..."

	deadlineInput := textinput.New()
	deadlineInput.Placeholder = "Enter deadline (YYYY-MM-DD HH:MM) (optional)..."

	// Create layout and window styles
	layout := NewLayout()
	windowStyles := CreateWindowStyles()

	model := &Model{
		app:               app,
		storage:           storage,
		state:             ListsView,
		layout:            layout,
		windowStyles:      windowStyles,
		titleInput:        titleInput,
		descriptionInput:  descriptionInput,
		deadlineInput:     deadlineInput,
		keys:              DefaultKeyMap(),
		lastReminderCheck: time.Now(),
		width:             80, // Default width
		height:            24, // Default height
		messageType:       "info",
	}

	// Set initial layout dimensions
	model.layout.SetScreenSize(model.width, model.height)

	// Initialize layout windows
	model.initializeWindows()

	// Initialize lists
	model.updateTodoListsList()

	// Auto-select first list if available
	if len(model.app.TodoLists) > 0 && model.currentListID == "" {
		model.currentListID = model.app.TodoLists[0].ID
		model.updateTasksList()
	}

	return model, nil
}

// initializeWindows sets up the initial windows in the layout
func (m *Model) initializeWindows() {
	// Create sidebar window
	sidebarWindow := &Window{
		ID:      SidebarWindow,
		Title:   "📋 Todo Lists",
		Content: "",
		Focused: false,
		Visible: true,
		Border:  true,
		Style:   m.windowStyles[SidebarWindow],
	}
	m.layout.AddWindow(sidebarWindow)

	// Create main window
	mainWindow := &Window{
		ID:      MainWindow,
		Title:   "📝 Tasks",
		Content: "",
		Focused: true,
		Visible: true,
		Border:  true,
		Style:   m.windowStyles[MainWindow],
	}
	m.layout.AddWindow(mainWindow)

	// Create status window
	statusWindow := &Window{
		ID:      StatusWindow,
		Title:   "",
		Content: "",
		Focused: false,
		Visible: true,
		Border:  false,
		Style:   m.windowStyles[StatusWindow],
	}
	m.layout.AddWindow(statusWindow)

	// Create form window (initially hidden)
	formWindow := &Window{
		ID:      FormWindow,
		Title:   "Form",
		Content: "",
		Focused: false,
		Visible: false,
		Border:  true,
		Style:   m.windowStyles[FormWindow],
	}
	m.layout.AddWindow(formWindow)

	// Create help window (initially hidden)
	helpWindow := &Window{
		ID:      HelpWindow,
		Title:   "❓ Help",
		Content: "",
		Focused: false,
		Visible: false,
		Border:  true,
		Style:   m.windowStyles[HelpWindow],
	}
	m.layout.AddWindow(helpWindow)
}

// updateListDimensions updates the list component dimensions based on window sizes
func (m *Model) updateListDimensions() {
	// Get sidebar window dimensions for lists
	sidebarWindow := m.layout.GetWindow(SidebarWindow)
	if sidebarWindow != nil {
		listWidth := sidebarWindow.Position.Width - 4   // Account for borders and padding
		listHeight := sidebarWindow.Position.Height - 6 // Account for borders and title

		// Ensure minimum dimensions
		if listWidth < 10 {
			listWidth = 10
		}
		if listHeight < 5 {
			listHeight = 5
		}

		if m.todoListsList.Items() != nil {
			m.todoListsList.SetSize(listWidth, listHeight)
		}
	}

	// Get main window dimensions for task list
	mainWindow := m.layout.GetWindow(MainWindow)
	if mainWindow != nil {
		listWidth := mainWindow.Position.Width - 4   // Account for borders and padding
		listHeight := mainWindow.Position.Height - 6 // Account for borders and title

		// Ensure minimum dimensions
		if listWidth < 10 {
			listWidth = 10
		}
		if listHeight < 5 {
			listHeight = 5
		}

		if m.tasksList.Items() != nil {
			m.tasksList.SetSize(listWidth, listHeight)
		}
	}
}

// toggleHelp shows or hides the help window
func (m *Model) toggleHelp() {
	helpWindow := m.layout.GetWindow(HelpWindow)
	if helpWindow != nil {
		if helpWindow.Visible {
			// Hide help window
			m.layout.SetWindowVisible(HelpWindow, false)
			m.layout.SetFocus(MainWindow)
		} else {
			// Show help window
			m.layout.SetWindowVisible(HelpWindow, true)
			m.layout.SetFocus(HelpWindow)
			// Update help content
			m.updateHelpContent()
		}
	}
}

// updateHelpContent updates the help window content
func (m *Model) updateHelpContent() {
	generalBindings := map[string]string{
		"q/Ctrl+C": "Quit application",
		"?":        "Toggle help",
		"Ctrl+→/←": "Navigate windows",
		"Ctrl+m":   "Focus main window",
		"Ctrl+s":   "Focus sidebar",
	}

	listBindings := map[string]string{
		"↑/↓":   "Navigate items",
		"Enter": "Select/Open item",
		"n":     "New todo list",
		"a":     "Add task",
		"e":     "Edit item",
		"d":     "Delete item",
		"Space": "Toggle task completion",
		"Esc":   "Go back",
	}

	formBindings := map[string]string{
		"Tab":       "Next field",
		"Shift+Tab": "Previous field",
		"Enter":     "Save",
		"Esc":       "Cancel",
	}

	content := CreateHelpSection("🌐 General", generalBindings) + "\n\n" +
		CreateHelpSection("📋 Lists & Tasks", listBindings) + "\n\n" +
		CreateHelpSection("📝 Forms", formBindings) + "\n\n" +
		DescStyle.Render("Press ? or Esc to close help")

	m.layout.SetWindowContent(HelpWindow, content)
}

// Init initializes the model
func (m *Model) Init() tea.Cmd {
	return tea.Batch(
		textinput.Blink,
		m.checkReminders(),
	)
}

// Update handles messages and updates the model
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Update layout dimensions
		m.layout.SetScreenSize(msg.Width, msg.Height)

		// Update list dimensions based on window sizes
		m.updateListDimensions()

	case tea.KeyMsg:
		// Global keys
		switch {
		case key.Matches(msg, m.keys.Quit):
			return m, tea.Quit
		case key.Matches(msg, m.keys.Help):
			m.toggleHelp()
			return m, nil
		case key.Matches(msg, m.keys.NextWindow):
			m.layout.NextFocus()
			return m, nil
		case key.Matches(msg, m.keys.PrevWindow):
			m.layout.PrevFocus()
			return m, nil
		case key.Matches(msg, m.keys.FocusMain):
			m.layout.SetFocus(MainWindow)
			return m, nil
		case key.Matches(msg, m.keys.FocusSidebar):
			m.layout.SetFocus(SidebarWindow)
			return m, nil
		}

		// State-specific handling based on focus and current state
		focusedWindow := m.layout.GetFocusedWindowID()

		// Handle form states first (overlay windows)
		if m.isInFormState() {
			switch m.state {
			case CreateListView, EditListView:
				return m.updateListForm(msg)
			case CreateTaskView, EditTaskView:
				return m.updateTaskForm(msg)
			}
		}

		// Handle help window
		if focusedWindow == HelpWindow {
			if key.Matches(msg, m.keys.Back) || key.Matches(msg, m.keys.Help) {
				m.toggleHelp()
				return m, nil
			}
			return m, nil
		}

		// Route to appropriate handler based on focus and state
		switch focusedWindow {
		case SidebarWindow:
			return m.updateListsView(msg)
		case MainWindow:
			switch m.state {
			case SettingsView:
				return m.updateSettingsView(msg)
			default:
				return m.updateTasksView(msg)
			}
		}

	case reminderMsg:
		m.checkForDueReminders()
		return m, m.checkReminders()

	case errorMsg:
		m.showMessage(string(msg))
		return m, nil
	}

	return m, tea.Batch(cmds...)
}

// View renders the multi-window layout
func (m *Model) View() string {
	if m.app == nil {
		return "Loading application..."
	}

	// Update window contents based on current state
	m.updateWindowContents()

	// Render the complete layout
	return m.layout.Render()
}

// updateWindowContents updates all window contents based on current state
func (m *Model) updateWindowContents() {
	// Update sidebar content (todo lists)
	sidebarContent := m.renderSidebarContent()
	m.layout.SetWindowContent(SidebarWindow, sidebarContent)

	// Update main content (tasks or other views)
	mainContent := m.renderMainContent()
	m.layout.SetWindowContent(MainWindow, mainContent)

	// Update status bar
	statusContent := m.renderStatusContent()
	m.layout.SetWindowContent(StatusWindow, statusContent)

	// Handle form overlay
	if m.isInFormState() {
		formContent := m.renderFormContent()
		m.layout.SetWindowContent(FormWindow, formContent)
		m.layout.SetWindowVisible(FormWindow, true)
	} else {
		m.layout.SetWindowVisible(FormWindow, false)
	}
}

// showMessage displays a status message
func (m *Model) showMessage(msg string) {
	m.showMessageWithType(msg, "info")
}

// getStatusBar returns the status bar content
func (m *Model) getStatusBar() string {
	style := lipgloss.NewStyle().
		Background(lipgloss.Color("62")).
		Foreground(lipgloss.Color("230")).
		Padding(0, 1)

	// Show message if recent
	if time.Since(m.messageTime) < 3*time.Second && m.message != "" {
		return style.Render(m.message)
	}

	// Show current state info
	var status string
	if m.app == nil {
		status = "Loading..."
	} else {
		switch m.state {
		case ListsView:
			status = fmt.Sprintf("Todo Lists (%d) • Press 'n' to create new • '?' for help", len(m.app.TodoLists))
		case TasksView:
			if list := m.getCurrentList(); list != nil {
				status = fmt.Sprintf("%s (%d/%d tasks) • Press 'a' to add task • 'esc' to go back",
					list.Name, list.GetCompletedCount(), list.GetTotalCount())
			} else {
				status = "Task View • 'esc' to go back"
			}
		default:
			status = "LazyTodo - Smart Todo Application"
		}
	}

	return style.Render(status)
}

// getCurrentList returns the currently selected todo list
func (m *Model) getCurrentList() *models.TodoList {
	for i := range m.app.TodoLists {
		if m.app.TodoLists[i].ID == m.currentListID {
			return &m.app.TodoLists[i]
		}
	}
	return nil
}

// saveData saves the application data
func (m *Model) saveData() tea.Cmd {
	return func() tea.Msg {
		if err := m.storage.Save(m.app); err != nil {
			return errorMsg(fmt.Sprintf("Failed to save: %v", err))
		}
		return nil
	}
}

// Message types
type reminderMsg struct{}
type errorMsg string

// checkReminders returns a command to check for reminders periodically
func (m *Model) checkReminders() tea.Cmd {
	return tea.Tick(time.Minute, func(t time.Time) tea.Msg {
		return reminderMsg{}
	})
}

// checkForDueReminders checks for tasks that need reminders
func (m *Model) checkForDueReminders() {
	if time.Since(m.lastReminderCheck) < time.Minute {
		return
	}

	m.lastReminderCheck = time.Now()
	reminderWindow := time.Duration(m.app.Settings.ReminderMinutes) * time.Minute

	for _, list := range m.app.TodoLists {
		for _, task := range list.Tasks {
			if task.Deadline != nil && !task.Completed {
				timeUntilDeadline := time.Until(*task.Deadline)
				if timeUntilDeadline > 0 && timeUntilDeadline <= reminderWindow {
					m.showMessage(fmt.Sprintf("⏰ Task '%s' is due in %s!", task.Title, timeUntilDeadline.Round(time.Minute)))
					return
				}
			}
		}
	}
}
