package main

import (
	"fmt"
	"log"
	"os"

	"github.com/DhirajZope/lazytodo/internal/storage"
	"github.com/DhirajZope/lazytodo/internal/ui"
	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Check for command line arguments
	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "--info", "-i":
			showStorageInfo()
			return
		case "--migrate", "-m":
			runMigration()
			return
		case "--help", "-h":
			showHelp()
			return
		case "--version", "-v":
			showVersion()
			return
		default:
			fmt.Printf("Unknown option: %s\n", os.Args[1])
			showHelp()
			os.Exit(1)
		}
	}

	// Initialize the model
	model, err := ui.NewModel()
	if err != nil {
		log.Fatalf("Failed to initialize application: %v", err)
	}

	// Create the program
	program := tea.NewProgram(
		model,
		tea.WithAltScreen(),
		tea.WithMouseCellMotion(),
	)

	// Run the program
	if _, err := program.Run(); err != nil {
		fmt.Printf("Error running program: %v\n", err)
		os.Exit(1)
	}
}

func showStorageInfo() {
	fmt.Println("🎯 LazyTodo - Storage Information")
	fmt.Println("===============================")

	// Try to initialize storage to get info
	storageInstance, err := storage.NewWithMigration()
	if err != nil {
		fmt.Printf("Error initializing storage: %v\n", err)
		os.Exit(1)
	}
	defer storageInstance.Close()

	fmt.Printf("Storage Backend: %s\n", storage.GetStorageInfo(storageInstance))

	// Load data to show statistics
	app, err := storageInstance.Load()
	if err != nil {
		fmt.Printf("Error loading data: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Todo Lists: %d\n", len(app.TodoLists))

	totalTasks := 0
	completedTasks := 0
	for _, list := range app.TodoLists {
		totalTasks += len(list.Tasks)
		completedTasks += list.GetCompletedCount()
	}

	fmt.Printf("Total Tasks: %d\n", totalTasks)
	fmt.Printf("Completed Tasks: %d\n", completedTasks)
	if totalTasks > 0 {
		fmt.Printf("Completion Rate: %.1f%%\n", float64(completedTasks)/float64(totalTasks)*100)
	}

	fmt.Printf("\nSettings:\n")
	fmt.Printf("  Reminder Minutes: %d\n", app.Settings.ReminderMinutes)
	fmt.Printf("  Show Completed: %v\n", app.Settings.ShowCompleted)
	fmt.Printf("  Auto Save: %v\n", app.Settings.AutoSave)
}

func runMigration() {
	fmt.Println("🎯 LazyTodo - Manual Migration")
	fmt.Println("=============================")

	dbStorage, err := storage.NewDatabase()
	if err != nil {
		fmt.Printf("Error creating database storage: %v\n", err)
		os.Exit(1)
	}
	defer dbStorage.Close()

	if err := storage.MigrateFromJSON(dbStorage); err != nil {
		fmt.Printf("Migration failed: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Migration completed successfully!")
}

func showHelp() {
	fmt.Println("🎯 LazyTodo - Smart Todo Application")
	fmt.Println("===================================")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  lazytodo                Run the TUI application")
	fmt.Println("  lazytodo --info, -i     Show storage information and statistics")
	fmt.Println("  lazytodo --migrate, -m  Manually run JSON to database migration")
	fmt.Println("  lazytodo --help, -h     Show this help message")
	fmt.Println("  lazytodo --version, -v  Show version information")
	fmt.Println()
	fmt.Println("Storage:")
	fmt.Println("  LazyTodo now uses SQLite database for improved reliability and performance.")
	fmt.Println("  Data is stored in: %USERPROFILE%\\.lazytodo\\lazytodo.db")
	fmt.Println("  Old JSON data will be automatically migrated on first run.")
	fmt.Println()
	fmt.Println("For more information, visit: https://github.com/DhirajZope/lazytodo")
}

func showVersion() {
	fmt.Println("🎯 LazyTodo v2.0.0")
	fmt.Println("Enhanced with SQLite database storage")
	fmt.Println("Built with Go and Bubble Tea")
}
