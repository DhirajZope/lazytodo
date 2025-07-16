# Multi-Window UI System Documentation

## Overview

LazyTodo v2.1 introduces a sophisticated multi-window Terminal User Interface (TUI) that transforms the application from a single-window design into an elegant, responsive, multi-panel layout with beautiful borders and enhanced visual feedback.

## Architecture

### Core Components

#### 1. Layout System (`internal/ui/layout.go`)
- **Window Management**: Manages multiple windows with focus handling
- **Responsive Sizing**: Automatically calculates window positions and sizes
- **Window Types**: MainWindow, SidebarWindow, StatusWindow, FormWindow, HelpWindow
- **Focus Navigation**: Support for window-to-window navigation

#### 2. Styling System (`internal/ui/styles.go`)
- **Elegant Color Palette**: Purple-green theme with sophisticated grays
- **Multiple Border Styles**: Elegant, Double, Thick, and Subtle borders
- **Component Styles**: Buttons, forms, lists, progress bars, status indicators
- **Visual Feedback**: Success, warning, error, and info styling

#### 3. Multi-Window Views (`internal/ui/multiwindow_views.go`)
- **Sidebar Content**: Todo list navigation
- **Main Content**: Task management and settings
- **Form Overlays**: Create/edit dialogs
- **Status Bar**: Real-time information and key hints

## Window Layout

```
┏━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓
┃ ╭─ 📋 Todo Lists ────────╮ ┏━━━ 📝 Tasks ━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┓ ┃
┃ │                        │ ┃                                               ┃ ┃
┃ │ Sidebar Window         │ ┃ Main Window                                   ┃ ┃
┃ │ - Todo Lists           │ ┃ - Task Details                                ┃ ┃
┃ │ - Progress Indicators  │ ┃ - Settings                                    ┃ ┃
┃ │ - Navigation           │ ┃ - Content Views                               ┃ ┃
┃ │                        │ ┃                                               ┃ ┃
┃ ╰────────────────────────╯ ┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛ ┃
┗━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━┛
Status Bar: Lists: 3 • Tasks: 5/8 • Focus: Main    ? Help  Ctrl+→/← Windows  q Quit
```

### Overlay Windows

#### Form Window (Centered Overlay)
```
              ╔══════════════════════════════════════╗
              ║ ➕ Create Task                       ║
              ║                                      ║
              ║ Title:                               ║
              ║ ┌──────────────────────────────────┐ ║
              ║ │ Task title here...               │ ║
              ║ └──────────────────────────────────┘ ║
              ║                                      ║
              ║ Description:                         ║
              ║ ┌──────────────────────────────────┐ ║
              ║ │ Optional description...          │ ║
              ║ └──────────────────────────────────┘ ║
              ║                                      ║
              ║ Tab/Shift+Tab • Enter Save • Esc Cancel ║
              ╚══════════════════════════════════════╝
```

#### Help Window (Large Overlay)
```
          ╭──────────────────── ❓ Help ────────────────────────╮
          │                                                     │
          │ 🌐 General                                          │
          │   q/Ctrl+C • Quit application                      │
          │   ? • Toggle help                                   │
          │   Ctrl+→/← • Navigate windows                       │
          │                                                     │
          │ 📋 Lists & Tasks                                   │
          │   ↑/↓ • Navigate items                             │
          │   Enter • Select/Open item                          │
          │   n • New todo list                                 │
          │   a • Add task                                      │
          │                                                     │
          │ Press ? or Esc to close help                        │
          ╰─────────────────────────────────────────────────────╯
```

## Key Features

### 🎨 Enhanced Visual Design
- **Beautiful Unicode Borders**: Multiple border styles for different contexts
- **Color-Coded Status**: Success (green), Warning (orange), Error (red), Info (blue)
- **Priority Indicators**: Visual priority levels with icons and colors
- **Progress Visualization**: Task completion progress bars
- **Focus Indicators**: Clear visual feedback for active windows

### 🚀 Improved Navigation
- **Window Focus Management**: `Ctrl+→/←` to navigate between windows
- **Direct Window Access**: `Ctrl+M` (Main), `Ctrl+S` (Sidebar)
- **Context-Aware Controls**: Different keybindings based on focused window
- **Overlay Handling**: Form and help windows with proper z-ordering

### 📱 Responsive Layout
- **Intelligent Sizing**: Automatic window sizing based on terminal dimensions
- **Minimum Size Constraints**: Graceful handling of small terminals
- **Dynamic Content**: Content adapts to available space
- **Overlay Positioning**: Centered overlays that scale with terminal size

### 💬 Enhanced Feedback
- **Typed Messages**: Success, warning, error, and info messages with icons
- **Status Indicators**: Real-time focus and state information
- **Progress Tracking**: Visual completion indicators
- **Contextual Hints**: Dynamic key binding hints based on current state

## Usage Guide

### Basic Navigation

1. **Starting the Application**
   ```bash
   ./lazytodo.exe
   ```

2. **Window Navigation**
   - `Ctrl+→` or `Ctrl+L`: Next window
   - `Ctrl+←` or `Ctrl+H`: Previous window
   - `Ctrl+M`: Focus main window
   - `Ctrl+S`: Focus sidebar

3. **Sidebar Operations** (Focus: Sidebar)
   - `↑/↓` or `j/k`: Navigate todo lists
   - `Enter`: Select and view list tasks
   - `n`: Create new todo list
   - `e`: Edit selected list
   - `d`: Delete selected list

4. **Main Window Operations** (Focus: Main)
   - `↑/↓` or `j/k`: Navigate tasks
   - `a`: Add new task
   - `Space`: Toggle task completion
   - `e`: Edit selected task
   - `d`: Delete selected task
   - `Esc`: Return focus to sidebar

5. **Form Operations** (Overlay Mode)
   - `Tab`: Next field
   - `Shift+Tab`: Previous field
   - `Enter`: Save
   - `Esc`: Cancel

6. **Global Operations**
   - `?`: Toggle help window
   - `q` or `Ctrl+C`: Quit application

### Advanced Features

#### Multi-List Workflow
1. Use sidebar to browse and select todo lists
2. Main window shows tasks for selected list
3. Create tasks in focused list context
4. Switch between lists without losing task focus

#### Form Overlays
- Forms appear as centered overlays
- Background content remains visible but inactive
- Focus automatically moves to form fields
- Cancel returns to previous window focus

#### Status Bar Information
- Current list and task counts
- Active window indicator
- Dynamic key hints
- Recent operation feedback

## Technical Details

### Window System Architecture

```go
type WindowID int
const (
    MainWindow WindowID = iota
    SidebarWindow
    StatusWindow
    HelpWindow
    FormWindow
)

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
```

### Layout Calculation
- **Sidebar**: 30-50 chars wide (responsive)
- **Main**: Remaining width after sidebar
- **Status**: 2-3 lines height at bottom
- **Overlays**: Centered with intelligent sizing

### Styling System
- **Primary Color**: Purple (#7C3AED)
- **Accent Color**: Green (#10B981)
- **Background**: Dark Gray (#1F2937)
- **Surface**: Medium Gray (#374151)
- **Text**: Light Gray hierarchy

### Performance Optimizations
- **Efficient Rendering**: Only update changed windows
- **Smart Sizing**: Calculate layouts once per resize
- **Content Caching**: Reuse rendered content when possible
- **Focus Management**: Minimal redraws on focus changes

## Migration from Single-Window

The multi-window system maintains full backward compatibility:

- All existing functionality is preserved
- Data format remains unchanged
- Keyboard shortcuts are enhanced, not replaced
- Settings and configuration work identically

### Key Differences
- **Old**: Single full-screen view with modal dialogs
- **New**: Split-pane layout with overlay forms
- **Old**: Linear navigation between views
- **New**: Window-based navigation with focus management
- **Old**: Basic text styling
- **New**: Rich visual feedback with icons and colors

## Future Enhancements

- **Themes**: Multiple color schemes
- **Customizable Layout**: User-defined window sizes
- **Plugin System**: Extensible window types
- **Mouse Support**: Click-to-focus and drag-to-resize
- **Animation**: Smooth transitions and effects

## Troubleshooting

### Common Issues

1. **Layout Issues on Small Terminals**
   - Minimum supported size: 80x24
   - Windows auto-resize with constraints
   - Content adapts to available space

2. **Unicode Characters Not Displaying**
   - Ensure terminal supports UTF-8
   - Check font has Unicode box-drawing characters
   - Fall back to ASCII borders if needed

3. **Key Bindings Not Working**
   - Check terminal key mapping
   - Some terminals may not support all Ctrl combinations
   - Alternative keys provided for compatibility

4. **Focus Issues**
   - Use `Ctrl+M` or `Ctrl+S` to explicitly set focus
   - Help window (`?`) shows current focus state
   - Status bar indicates active window

This multi-window system represents a significant evolution in LazyTodo's user experience, providing a modern, efficient, and visually appealing interface for task management. 