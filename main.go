package main

import (
	"fmt"
	"os"

	"forge-tui/internal/ai"
	"forge-tui/internal/ui"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// 1. Initialize the Local AI Engine
	aiEngine := ai.NewEngine("llama3.1:8b") // Change to "phi3:mini" if low on RAM

	// 2. Initialize the UI Model
	m := ui.NewModel(aiEngine)

	// 3. Start the Bubbletea Program
	p := tea.NewProgram(
		m,
		tea.WithAltScreen(),       // Use the full terminal screen
		tea.WithMouseCellMotion(), // Enable mouse scrolling
	)

	if _, err := p.Run(); err != nil {
		fmt.Printf("Error running Forge TUI: %v\n", err)
		os.Exit(1)
	}
}