# 🎯 LazyTodo - Smart Todo Application v2.1

A beautiful and feature-rich Terminal User Interface (TUI) todo application built with Go and [Bubble Tea](https://github.com/charmbracelet/bubbletea). **Now featuring an elegant multi-window interface with responsive layout, beautiful borders, and enhanced visual feedback!**

## ✨ Features

### 🎨 Multi-Window Interface (NEW!)
- **📱 Responsive Layout**: Elegant split-pane interface with sidebar and main content areas
- **🖼️ Beautiful Borders**: Multiple Unicode border styles (elegant, double, thick, subtle)
- **🎯 Focus Management**: Seamless navigation between windows with visual focus indicators  
- **🌈 Rich Visual Feedback**: Color-coded status messages and priority indicators
- **📋 Overlay Forms**: Centered form dialogs with proper z-ordering
- **💡 Smart Sizing**: Responsive window sizing that adapts to terminal dimensions

### 📋 Core Functionality
- **📝 Multiple Todo Lists**: Create and manage separate todo lists for different projects or contexts
- **✅ Rich Task Management**: Add, edit, delete, and toggle completion status of tasks
- **⏰ Deadline Support**: Set deadlines for tasks with reminder notifications
- **🎨 Priority Levels**: Assign priority levels (Low, Medium, High, Critical) to tasks
- **🔔 Smart Reminders**: Get notified before task deadlines (configurable reminder window)
- **💾 SQLite Database Storage**: ACID-compliant database storage with automatic backups
- **🔄 Automatic Migration**: Seamlessly migrates from old JSON format to database
- **🎯 Intuitive Navigation**: Vim-like keybindings plus window focus controls
- **📊 Progress Tracking**: Visual progress indicators for each todo list
- **🔍 Visual Status Indicators**: Clear visual cues for task status, priority, and deadlines
- **🖥️ Command Line Interface**: Rich CLI for storage info and management

## 🚀 Installation

### Prerequisites

- Go 1.23 or higher installed on your system
- Windows, macOS, or Linux
- CGO enabled (for SQLite support)

### Build from Source

1. Clone the repository:
```powershell
git clone https://github.com/DhirajZope/lazytodo.git
cd lazytodo
```

2. Install dependencies:
```powershell
go mod tidy
```

3. Build the application:
```powershell
go build -o lazytodo.exe cmd/main.go
```

4. Run the application:
```powershell
.\lazytodo.exe
```

Or run directly with Go:
```powershell
go run cmd/main.go
```

## 🎮 Usage

### Running the Application

```powershell
# Start the TUI application
.\lazytodo.exe

# Show storage information and statistics
.\lazytodo.exe --info

# Show help
.\lazytodo.exe --help

# Show version
.\lazytodo.exe --version

# Manually run migration (if needed)
.\lazytodo.exe --migrate
```

### Navigation

The application uses intuitive keybindings inspired by Vim:

#### Global Keys
- `q` or `Ctrl+C` - Quit application
- `?` - Toggle help menu

#### Todo Lists View
- `↑`/`↓` or `k`/`j` - Navigate between lists
- `Enter` - Open selected list
- `n` - Create new todo list
- `e` - Edit selected list
- `d` - Delete selected list
- `s` - Open settings

#### Tasks View
- `↑`/`↓` or `k`/`j` - Navigate between tasks
- `Space` - Toggle task completion
- `a` - Add new task
- `e` - Edit selected task
- `d` - Delete selected task
- `Esc` - Back to lists view

#### Forms
- `Tab`/`Shift+Tab` - Navigate between form fields
- `Enter` - Save changes
- `Esc` - Cancel and go back

### Visual Indicators

#### Task Status
- `○` - Incomplete task
- `✓` - Complete task

#### Priority Levels
- `⚡` - Medium priority
- `🔥` - High priority
- `🚨` - Critical priority

#### Deadlines
- `⏰` - Task due soon (within 24 hours)
- `⚠️` - Task overdue

## 📁 Data Storage

### Multi-Window Interface (v2.1+)
LazyTodo features a sophisticated multi-window TUI with:

#### 🪟 Window Layout
- **Sidebar**: Todo list navigation with progress indicators
- **Main Window**: Task details, settings, and content views  
- **Status Bar**: Real-time information and keyboard shortcuts
- **Overlay Windows**: Forms and help with proper z-ordering

#### 🎯 Window Navigation
- `Ctrl+→` / `Ctrl+←`: Navigate between windows
- `Ctrl+M`: Focus main window
- `Ctrl+S`: Focus sidebar
- `?`: Toggle help overlay
- Visual focus indicators show active window

#### 🎨 Visual Enhancements  
- **Beautiful Unicode Borders**: Multiple styles for different contexts
- **Color-Coded Messages**: Success (green), warning (orange), error (red), info (blue)
- **Priority Indicators**: Visual priority levels with icons
- **Responsive Design**: Adapts to any terminal size (minimum 80x24)

See [MULTIWINDOW_UI.md](MULTIWINDOW_UI.md) for detailed documentation.

### SQLite Database (v2.0+)
LazyTodo uses SQLite for data storage, providing:
- **ACID compliance** for data integrity
- **Better performance** with indexed queries
- **Concurrent access safety**
- **Automatic schema migrations**

Database location:
- **Windows**: `%USERPROFILE%\.lazytodo\lazytodo.db`
- **macOS/Linux**: `~/.lazytodo/lazytodo.db`

### Migration from JSON (v1.x)
If you're upgrading from v1.x, LazyTodo will automatically:
1. Detect your existing JSON data file
2. Migrate all data to the new SQLite database
3. Create a backup of your JSON file
4. Preserve all your todo lists, tasks, and settings

**No data loss** - your existing data is fully preserved!

### Database Schema
The database includes the following tables:
- `todo_lists` - Stores todo list information
- `tasks` - Stores individual tasks with foreign key references
- `settings` - Stores application settings
- `schema_migrations` - Tracks applied database migrations

## ⚙️ Configuration

Settings are now stored in the database. Default settings:

- **Reminder Window**: 60 minutes before deadline
- **Show Completed Tasks**: Enabled
- **Auto Save**: Enabled (immediate database updates)

## 🎯 Task Deadlines

When creating or editing tasks, you can set deadlines using the format:
```
YYYY-MM-DD HH:MM
```

Examples:
- `2024-12-25 09:00` - Christmas morning at 9 AM
- `2024-07-16 14:30` - Today at 2:30 PM

### Reminders

The application checks for upcoming deadlines every minute and displays notifications for tasks that are:
- Due within your configured reminder window (default: 1 hour)
- Overdue

## 🏗️ Project Structure

```
smart-todo/
├── cmd/
│   └── main.go              # Application entry point with CLI
├── internal/
│   ├── models/
│   │   └── models.go        # Data models and types
│   ├── storage/
│   │   ├── interface.go     # Storage interface definition
│   │   ├── storage.go       # Legacy JSON file storage
│   │   ├── database.go      # SQLite database storage
│   │   └── migration.go     # Data migration utilities
│   └── ui/
│       ├── model.go         # Main TUI model and state management
│       └── views.go         # UI rendering and interactions
├── migrations/
│   ├── 001_initial_schema.up.sql    # Database schema
│   └── 001_initial_schema.down.sql  # Rollback schema
├── go.mod                   # Go module definition
├── go.sum                   # Go dependencies
├── build.ps1               # PowerShell build script
└── README.md               # This file
```

## 🛠️ Development

### Dependencies

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) - TUI framework
- [Lipgloss](https://github.com/charmbracelet/lipgloss) - Styling and layout
- [Bubbles](https://github.com/charmbracelet/bubbles) - Common TUI components
- [go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite database driver

### Building

```powershell
# Build for current platform
go build -o lazytodo.exe cmd/main.go

# Build for different platforms
# Windows
GOOS=windows GOARCH=amd64 go build -o lazytodo-windows.exe cmd/main.go

# macOS
GOOS=darwin GOARCH=amd64 go build -o lazytodo-macos cmd/main.go

# Linux
GOOS=linux GOARCH=amd64 go build -o lazytodo-linux cmd/main.go
```

### Database Management

```powershell
# Check storage information
.\lazytodo.exe --info

# Manually run migration
.\lazytodo.exe --migrate

# View database schema (requires sqlite3 CLI)
sqlite3 %USERPROFILE%\.lazytodo\lazytodo.db ".schema"
```

## 🔄 Version History

### v2.0.0 (Current)
- **New**: SQLite database storage
- **New**: Automatic migration from JSON format
- **New**: Command line interface for management
- **Improved**: Data integrity and performance
- **Enhanced**: ACID-compliant transactions

### v1.x
- JSON file-based storage
- Basic TUI functionality
- File-based persistence

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📝 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🙏 Acknowledgments

- [Charm](https://charm.sh/) team for the amazing TUI libraries
- [mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) for SQLite driver
- Go community for the excellent ecosystem
- All contributors and users of this project

## 📞 Support

If you encounter any issues or have feature requests, please open an issue on GitHub.

### Common Issues

**Migration Issues**: If automatic migration fails, you can:
1. Run `.\lazytodo.exe --migrate` manually
2. Check that your JSON file is valid
3. Ensure you have write permissions to the data directory

**Database Corruption**: SQLite is very reliable, but if issues occur:
1. Check the `--info` command for database status
2. Your JSON backup file is always preserved
3. You can delete the database file to start fresh

---

Made with ❤️ and Go • Enhanced with SQLite 