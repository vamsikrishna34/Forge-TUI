package main

import (
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func main() {
	// Initialize the Local AI Engine (Using the model you just pulled!)
	aiEngine := NewEngine("llama3.1:8b") 

	// Initialize the UI Model
	m := NewModel(aiEngine)

	// Start the Bubbletea Program
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