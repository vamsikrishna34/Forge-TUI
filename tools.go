package main

import (
	"context"
	"fmt"
	"os/exec"

	"github.com/tmc/langchaingo/tools"
)

// --- Tool 1: Read Local Files ---
type ReadFileTool struct{}

func (t ReadFileTool) Description() string {
	return "Use this tool to read the contents of a local file. Input should be the file path (e.g., 'main.go')."
}

func (t ReadFileTool) Name() string {
	return "read_file"
}

func (t ReadFileTool) Call(ctx context.Context, input string) (string, error) {
	// In a real app, you'd sanitize this path to prevent directory traversal!
	content, err := exec.CommandContext(ctx, "cat", input).Output()
	if err != nil {
		return fmt.Sprintf("Error reading file: %v", err), nil
	}
	return string(content), nil
}

// --- Tool 2: Execute Bash Commands ---
type BashTool struct{}

func (t BashTool) Description() string {
	return "Use this tool to execute a bash command on the local machine. Input should be the exact command (e.g., 'ls -la')."
}

func (t BashTool) Name() string {
	return "run_bash"
}

func (t BashTool) Call(ctx context.Context, input string) (string, error) {
	output, err := exec.CommandContext(ctx, "bash", "-c", input).Output()
	if err != nil {
		return fmt.Sprintf("Error executing command: %v", err), nil
	}
	return string(output), nil
}

// GetTools returns all available tools for the AI
func GetTools() []tools.Tool {
	return []tools.Tool{
		ReadFileTool{},
		BashTool{},
	}
}