package ui

import (
	"github.com/charmbracelet/lipgloss"
)

// Color palette for the elegant theme
var (
	// Primary colors
	PrimaryColor    = lipgloss.Color("#7C3AED") // Purple
	AccentColor     = lipgloss.Color("#10B981") // Green
	BackgroundColor = lipgloss.Color("#1F2937") // Dark gray
	SurfaceColor    = lipgloss.Color("#374151") // Medium gray

	// Text colors
	TextPrimary   = lipgloss.Color("#F9FAFB") // Light gray
	TextSecondary = lipgloss.Color("#D1D5DB") // Medium light gray
	TextMuted     = lipgloss.Color("#9CA3AF") // Muted gray

	// Status colors
	SuccessColor = lipgloss.Color("#10B981") // Green
	WarningColor = lipgloss.Color("#F59E0B") // Orange
	ErrorColor   = lipgloss.Color("#EF4444") // Red
	InfoColor    = lipgloss.Color("#3B82F6") // Blue

	// Border colors
	BorderPrimary   = lipgloss.Color("#7C3AED") // Purple
	BorderSecondary = lipgloss.Color("#6B7280") // Gray
	BorderFocused   = lipgloss.Color("#10B981") // Green
	BorderUnfocused = lipgloss.Color("#4B5563") // Dark gray
)

// Border styles
var (
	// Elegant rounded border
	ElegantBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "╭",
		TopRight:    "╮",
		BottomLeft:  "╰",
		BottomRight: "╯",
	}

	// Double line border for emphasis
	DoubleBorder = lipgloss.Border{
		Top:         "═",
		Bottom:      "═",
		Left:        "║",
		Right:       "║",
		TopLeft:     "╔",
		TopRight:    "╗",
		BottomLeft:  "╚",
		BottomRight: "╝",
	}

	// Thick border for main windows
	ThickBorder = lipgloss.Border{
		Top:         "━",
		Bottom:      "━",
		Left:        "┃",
		Right:       "┃",
		TopLeft:     "┏",
		TopRight:    "┓",
		BottomLeft:  "┗",
		BottomRight: "┛",
	}

	// Subtle border for secondary content
	SubtleBorder = lipgloss.Border{
		Top:         "─",
		Bottom:      "─",
		Left:        "│",
		Right:       "│",
		TopLeft:     "┌",
		TopRight:    "┐",
		BottomLeft:  "└",
		BottomRight: "┘",
	}
)

// Base styles
var (
	// Base content style
	BaseContentStyle = lipgloss.NewStyle().
				Foreground(TextPrimary)

	// Title style
	BaseTitleStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true).
			Padding(0, 1)

	// Subtitle style
	BaseSubtitleStyle = lipgloss.NewStyle().
				Foreground(TextSecondary).
				Italic(true)
)

// CreateWindowStyles creates comprehensive styling for all window types
func CreateWindowStyles() map[WindowID]WindowStyle {
	styles := make(map[WindowID]WindowStyle)

	// Main window styles
	styles[MainWindow] = WindowStyle{
		Border: lipgloss.NewStyle().
			Border(ThickBorder).
			BorderForeground(BorderFocused),
		Title: BaseTitleStyle.Copy().
			Foreground(PrimaryColor).
			Bold(true),
		Content: BaseContentStyle.Copy().
			Padding(1),
		Focused: lipgloss.NewStyle().
			Border(ThickBorder).
			BorderForeground(BorderFocused),
		Unfocused: lipgloss.NewStyle().
			Border(ThickBorder).
			BorderForeground(BorderUnfocused),
	}

	// Sidebar window styles
	styles[SidebarWindow] = WindowStyle{
		Border: lipgloss.NewStyle().
			Border(ElegantBorder).
			BorderForeground(BorderSecondary),
		Title: BaseTitleStyle.Copy().
			Foreground(AccentColor),
		Content: BaseContentStyle.Copy().
			Padding(1),
		Focused: lipgloss.NewStyle().
			Border(ElegantBorder).
			BorderForeground(BorderFocused),
		Unfocused: lipgloss.NewStyle().
			Border(ElegantBorder).
			BorderForeground(BorderSecondary),
	}

	// Status window styles
	styles[StatusWindow] = WindowStyle{
		Border: lipgloss.NewStyle().
			Border(SubtleBorder).
			BorderForeground(BorderSecondary),
		Title: BaseTitleStyle.Copy().
			Foreground(InfoColor),
		Content: BaseContentStyle.Copy().
			Padding(0, 1),
		Focused: lipgloss.NewStyle().
			Border(SubtleBorder).
			BorderForeground(BorderFocused),
		Unfocused: lipgloss.NewStyle().
			Border(SubtleBorder).
			BorderForeground(BorderSecondary),
	}

	// Form window styles (overlay)
	styles[FormWindow] = WindowStyle{
		Border: lipgloss.NewStyle().
			Border(DoubleBorder).
			BorderForeground(PrimaryColor),
		Title: BaseTitleStyle.Copy().
			Foreground(PrimaryColor).
			Bold(true),
		Content: BaseContentStyle.Copy().
			Padding(2),
		Focused: lipgloss.NewStyle().
			Border(DoubleBorder).
			BorderForeground(PrimaryColor),
		Unfocused: lipgloss.NewStyle().
			Border(DoubleBorder).
			BorderForeground(BorderUnfocused),
	}

	// Help window styles (overlay)
	styles[HelpWindow] = WindowStyle{
		Border: lipgloss.NewStyle().
			Border(ElegantBorder).
			BorderForeground(InfoColor),
		Title: BaseTitleStyle.Copy().
			Foreground(InfoColor).
			Bold(true),
		Content: BaseContentStyle.Copy().
			Padding(2),
		Focused: lipgloss.NewStyle().
			Border(ElegantBorder).
			BorderForeground(InfoColor),
		Unfocused: lipgloss.NewStyle().
			Border(ElegantBorder).
			BorderForeground(BorderSecondary),
	}

	return styles
}

// Component styles for various UI elements
var (
	// List item styles
	ListItemNormal = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Padding(0, 1)

	ListItemSelected = lipgloss.NewStyle().
				Foreground(BackgroundColor).
				Background(PrimaryColor).
				Bold(true).
				Padding(0, 1)

	ListItemCompleted = lipgloss.NewStyle().
				Foreground(TextMuted).
				Strikethrough(true).
				Padding(0, 1)

	// Progress indicators
	ProgressBarFull = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Background(SuccessColor)

	ProgressBarEmpty = lipgloss.NewStyle().
				Foreground(BorderSecondary).
				Background(BorderSecondary)

	// Status indicators
	StatusSuccess = lipgloss.NewStyle().
			Foreground(SuccessColor).
			Bold(true)

	StatusWarning = lipgloss.NewStyle().
			Foreground(WarningColor).
			Bold(true)

	StatusError = lipgloss.NewStyle().
			Foreground(ErrorColor).
			Bold(true)

	StatusInfo = lipgloss.NewStyle().
			Foreground(InfoColor).
			Bold(true)

	// Form element styles
	FormFieldFocused = lipgloss.NewStyle().
				Border(SubtleBorder).
				BorderForeground(BorderFocused).
				Padding(0, 1)

	FormFieldUnfocused = lipgloss.NewStyle().
				Border(SubtleBorder).
				BorderForeground(BorderSecondary).
				Padding(0, 1)

	FormLabel = lipgloss.NewStyle().
			Foreground(TextSecondary).
			Bold(true)

	// Button styles
	ButtonPrimary = lipgloss.NewStyle().
			Foreground(BackgroundColor).
			Background(PrimaryColor).
			Bold(true).
			Padding(0, 2).
			Margin(0, 1)

	ButtonSecondary = lipgloss.NewStyle().
			Foreground(TextPrimary).
			Background(SurfaceColor).
			Border(SubtleBorder).
			BorderForeground(BorderSecondary).
			Padding(0, 2).
			Margin(0, 1)

	// Key binding help styles
	KeyStyle = lipgloss.NewStyle().
			Foreground(PrimaryColor).
			Bold(true)

	DescStyle = lipgloss.NewStyle().
			Foreground(TextSecondary)

	SeparatorStyle = lipgloss.NewStyle().
			Foreground(TextMuted)
)

// Priority styling
func GetPriorityStyle(priority string) lipgloss.Style {
	switch priority {
	case "low":
		return lipgloss.NewStyle().Foreground(TextMuted)
	case "medium":
		return lipgloss.NewStyle().Foreground(WarningColor)
	case "high":
		return lipgloss.NewStyle().Foreground(ErrorColor).Bold(true)
	case "critical":
		return lipgloss.NewStyle().Foreground(ErrorColor).Bold(true).Blink(true)
	default:
		return lipgloss.NewStyle().Foreground(TextPrimary)
	}
}

// Deadline styling based on urgency
func GetDeadlineStyle(isOverdue, isDueSoon bool) lipgloss.Style {
	if isOverdue {
		return lipgloss.NewStyle().Foreground(ErrorColor).Bold(true)
	} else if isDueSoon {
		return lipgloss.NewStyle().Foreground(WarningColor).Bold(true)
	}
	return lipgloss.NewStyle().Foreground(InfoColor)
}

// Progress bar rendering
func RenderProgressBar(current, total int, width int) string {
	if total == 0 {
		return ProgressBarEmpty.Width(width).Render("")
	}

	progress := float64(current) / float64(total)
	filled := int(float64(width) * progress)
	empty := width - filled

	filledBar := ProgressBarFull.Width(filled).Render("")
	emptyBar := ProgressBarEmpty.Width(empty).Render("")

	return lipgloss.JoinHorizontal(lipgloss.Left, filledBar, emptyBar)
}

// Enhanced status message styling
func StyleStatusMessage(message string, messageType string) string {
	var icon string
	var style lipgloss.Style

	switch messageType {
	case "success":
		icon = "✓"
		style = StatusSuccess
	case "warning":
		icon = "⚠"
		style = StatusWarning
	case "error":
		icon = "✗"
		style = StatusError
	case "info":
		icon = "ℹ"
		style = StatusInfo
	default:
		icon = "•"
		style = lipgloss.NewStyle().Foreground(TextPrimary)
	}

	return style.Render(icon + " " + message)
}

// Create a help section with key bindings
func CreateHelpSection(title string, bindings map[string]string) string {
	var lines []string

	// Add section title
	lines = append(lines, BaseTitleStyle.Render(title))
	lines = append(lines, "")

	// Add key bindings
	for key, desc := range bindings {
		keyStr := KeyStyle.Render(key)
		descStr := DescStyle.Render(desc)
		separator := SeparatorStyle.Render(" • ")
		lines = append(lines, lipgloss.JoinHorizontal(lipgloss.Left, keyStr, separator, descStr))
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

// Create a visual separator
func CreateSeparator(width int, style string) string {
	var char string
	var color lipgloss.Color

	switch style {
	case "thick":
		char = "━"
		color = BorderPrimary
	case "double":
		char = "═"
		color = BorderSecondary
	case "dotted":
		char = "┄"
		color = BorderSecondary
	default:
		char = "─"
		color = BorderSecondary
	}

	return lipgloss.NewStyle().
		Foreground(color).
		Width(width).
		Render(char)
}

// Enhanced list rendering with icons and styling
func RenderEnhancedListItem(icon, title, subtitle string, selected, completed bool) string {
	var style lipgloss.Style
	var itemIcon string

	if completed {
		style = ListItemCompleted
		itemIcon = "✓"
	} else if selected {
		style = ListItemSelected
		itemIcon = icon
	} else {
		style = ListItemNormal
		itemIcon = icon
	}

	// Build the item
	content := itemIcon + " " + title
	if subtitle != "" {
		subtitleStyle := lipgloss.NewStyle().
			Foreground(TextMuted).
			Italic(true)
		content += "\n  " + subtitleStyle.Render(subtitle)
	}

	return style.Render(content)
}
