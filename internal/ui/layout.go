package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// WindowID represents different window types
type WindowID int

const (
	MainWindow WindowID = iota
	SidebarWindow
	StatusWindow
	HelpWindow
	FormWindow
)

// Window represents a UI window/panel
type Window struct {
	ID       WindowID
	Title    string
	Content  string
	Focused  bool
	Visible  bool
	Position Rect
	Style    WindowStyle
	Border   bool
}

// Rect represents window position and dimensions
type Rect struct {
	X      int
	Y      int
	Width  int
	Height int
}

// WindowStyle contains styling information for a window
type WindowStyle struct {
	Border    lipgloss.Style
	Title     lipgloss.Style
	Content   lipgloss.Style
	Focused   lipgloss.Style
	Unfocused lipgloss.Style
}

// Layout manages the arrangement of windows
type Layout struct {
	windows      map[WindowID]*Window
	focusOrder   []WindowID
	currentFocus int
	screenWidth  int
	screenHeight int
}

// NewLayout creates a new layout manager
func NewLayout() *Layout {
	return &Layout{
		windows:      make(map[WindowID]*Window),
		focusOrder:   []WindowID{MainWindow, SidebarWindow, StatusWindow},
		currentFocus: 0,
	}
}

// SetScreenSize updates the screen dimensions and recalculates layout
func (l *Layout) SetScreenSize(width, height int) {
	l.screenWidth = width
	l.screenHeight = height
	l.calculateLayout()
}

// AddWindow adds a window to the layout
func (l *Layout) AddWindow(window *Window) {
	l.windows[window.ID] = window
	l.calculateLayout()
}

// GetWindow returns a window by ID
func (l *Layout) GetWindow(id WindowID) *Window {
	return l.windows[id]
}

// SetWindowVisible sets the visibility of a window
func (l *Layout) SetWindowVisible(id WindowID, visible bool) {
	if window := l.windows[id]; window != nil {
		window.Visible = visible
		l.calculateLayout()
	}
}

// SetWindowContent updates the content of a window
func (l *Layout) SetWindowContent(id WindowID, content string) {
	if window := l.windows[id]; window != nil {
		window.Content = content
	}
}

// SetWindowTitle updates the title of a window
func (l *Layout) SetWindowTitle(id WindowID, title string) {
	if window := l.windows[id]; window != nil {
		window.Title = title
	}
}

// NextFocus moves focus to the next window
func (l *Layout) NextFocus() {
	if len(l.focusOrder) == 0 {
		return
	}

	// Clear current focus
	if current := l.GetFocusedWindow(); current != nil {
		current.Focused = false
	}

	// Move to next visible window
	for i := 0; i < len(l.focusOrder); i++ {
		l.currentFocus = (l.currentFocus + 1) % len(l.focusOrder)
		windowID := l.focusOrder[l.currentFocus]
		if window := l.windows[windowID]; window != nil && window.Visible {
			window.Focused = true
			break
		}
	}
}

// PrevFocus moves focus to the previous window
func (l *Layout) PrevFocus() {
	if len(l.focusOrder) == 0 {
		return
	}

	// Clear current focus
	if current := l.GetFocusedWindow(); current != nil {
		current.Focused = false
	}

	// Move to previous visible window
	for i := 0; i < len(l.focusOrder); i++ {
		l.currentFocus = (l.currentFocus - 1 + len(l.focusOrder)) % len(l.focusOrder)
		windowID := l.focusOrder[l.currentFocus]
		if window := l.windows[windowID]; window != nil && window.Visible {
			window.Focused = true
			break
		}
	}
}

// GetFocusedWindow returns the currently focused window
func (l *Layout) GetFocusedWindow() *Window {
	for _, window := range l.windows {
		if window.Focused {
			return window
		}
	}
	return nil
}

// GetFocusedWindowID returns the ID of the currently focused window
func (l *Layout) GetFocusedWindowID() WindowID {
	if window := l.GetFocusedWindow(); window != nil {
		return window.ID
	}
	return MainWindow
}

// SetFocus sets focus to a specific window
func (l *Layout) SetFocus(id WindowID) {
	// Clear all focus
	for _, window := range l.windows {
		window.Focused = false
	}

	// Set focus to target window
	if window := l.windows[id]; window != nil && window.Visible {
		window.Focused = true
		// Update current focus index
		for i, focusID := range l.focusOrder {
			if focusID == id {
				l.currentFocus = i
				break
			}
		}
	}
}

// calculateLayout calculates window positions and sizes based on screen size
func (l *Layout) calculateLayout() {
	if l.screenWidth <= 0 || l.screenHeight <= 0 {
		return
	}

	// Calculate dimensions with minimum sizes
	minSidebarWidth := 35
	minMainWidth := 45
	statusHeight := 3

	// Ensure we have enough space
	if l.screenWidth < minSidebarWidth+minMainWidth {
		minSidebarWidth = l.screenWidth / 3
		minMainWidth = l.screenWidth - minSidebarWidth
	}

	if l.screenHeight < statusHeight+10 {
		statusHeight = 2
	}

	// Calculate sidebar width (1/3 of screen, but with limits)
	sidebarWidth := l.screenWidth / 3
	if sidebarWidth < minSidebarWidth {
		sidebarWidth = minSidebarWidth
	}
	if sidebarWidth > 50 {
		sidebarWidth = 50
	}

	mainWidth := l.screenWidth - sidebarWidth
	mainHeight := l.screenHeight - statusHeight

	// Update window positions and sizes
	if sidebar := l.windows[SidebarWindow]; sidebar != nil {
		sidebar.Position = Rect{
			X:      0,
			Y:      0,
			Width:  sidebarWidth,
			Height: mainHeight,
		}
	}

	if main := l.windows[MainWindow]; main != nil {
		main.Position = Rect{
			X:      sidebarWidth,
			Y:      0,
			Width:  mainWidth,
			Height: mainHeight,
		}
	}

	if status := l.windows[StatusWindow]; status != nil {
		status.Position = Rect{
			X:      0,
			Y:      mainHeight,
			Width:  l.screenWidth,
			Height: statusHeight,
		}
	}

	// Form window (overlay, centered)
	if form := l.windows[FormWindow]; form != nil {
		formWidth := 60
		formHeight := 12
		if l.screenWidth < 70 {
			formWidth = l.screenWidth - 4
		}
		if l.screenHeight < 15 {
			formHeight = l.screenHeight - 3
		}

		form.Position = Rect{
			X:      (l.screenWidth - formWidth) / 2,
			Y:      (l.screenHeight - formHeight) / 2,
			Width:  formWidth,
			Height: formHeight,
		}
	}

	// Help window (overlay, large)
	if help := l.windows[HelpWindow]; help != nil {
		helpWidth := l.screenWidth - 6
		helpHeight := l.screenHeight - 4
		if helpWidth < 60 {
			helpWidth = l.screenWidth
		}
		if helpHeight < 20 {
			helpHeight = l.screenHeight
		}

		help.Position = Rect{
			X:      (l.screenWidth - helpWidth) / 2,
			Y:      (l.screenHeight - helpHeight) / 2,
			Width:  helpWidth,
			Height: helpHeight,
		}
	}
}

// Render renders all visible windows and returns the complete view
func (l *Layout) Render() string {
	if l.screenWidth <= 0 || l.screenHeight <= 0 {
		return "Loading..."
	}

	// Get window content
	sidebarContent := ""
	mainContent := ""
	statusContent := ""

	if sidebar := l.windows[SidebarWindow]; sidebar != nil && sidebar.Visible {
		sidebarContent = l.renderWindow(sidebar)
	}

	if main := l.windows[MainWindow]; main != nil && main.Visible {
		mainContent = l.renderWindow(main)
	}

	if status := l.windows[StatusWindow]; status != nil && status.Visible {
		statusContent = l.renderWindow(status)
	}

	// Create horizontal layout for sidebar and main content
	topRow := lipgloss.JoinHorizontal(
		lipgloss.Top,
		sidebarContent,
		mainContent,
	)

	// Add status bar at bottom
	fullLayout := lipgloss.JoinVertical(
		lipgloss.Left,
		topRow,
		statusContent,
	)

	// Handle overlay windows (forms, help)
	if form := l.windows[FormWindow]; form != nil && form.Visible {
		formContent := l.renderWindow(form)
		// Center the form over the main layout
		formWithBackground := lipgloss.Place(
			l.screenWidth, l.screenHeight,
			lipgloss.Center, lipgloss.Center,
			formContent,
		)
		return formWithBackground
	}

	if help := l.windows[HelpWindow]; help != nil && help.Visible {
		helpContent := l.renderWindow(help)
		// Center the help over the main layout
		helpWithBackground := lipgloss.Place(
			l.screenWidth, l.screenHeight,
			lipgloss.Center, lipgloss.Center,
			helpContent,
		)
		return helpWithBackground
	}

	return fullLayout
}

// renderWindow renders a single window with its styling
func (l *Layout) renderWindow(window *Window) string {
	if window == nil {
		return ""
	}

	// Choose style based on focus
	var borderStyle lipgloss.Style
	if window.Focused {
		borderStyle = window.Style.Focused
	} else {
		borderStyle = window.Style.Unfocused
	}

	content := window.Content

	// Apply content styling if content exists
	if content != "" && window.Style.Content.GetForeground() != lipgloss.Color("") {
		content = window.Style.Content.Render(content)
	}

	// Calculate available dimensions
	availableWidth := window.Position.Width
	availableHeight := window.Position.Height

	// Apply border and sizing
	if window.Border {
		content = borderStyle.
			Width(availableWidth - 2).   // Account for border
			Height(availableHeight - 2). // Account for border
			Render(content)
	} else {
		content = lipgloss.NewStyle().
			Width(availableWidth).
			Height(availableHeight).
			Render(content)
	}

	return content
}
