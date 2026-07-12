# Forge-TUI
# Forge TUI

A blazing-fast, 100% local, dark-mode Terminal User Interface (TUI) for interacting with AI. 

Built with Go, Bubbletea, and Ollama. **Zero API costs. Zero cloud dependencies. Your data never leaves your machine.**

![Go](https://img.shields.io/badge/Go-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![Ollama](https://img.shields.io/badge/Ollama-Local_AI-black?style=for-the-badge)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)

## Features

- Gorgeous Dark Mode UI:** Built with the Catppuccin Mocha theme using `lipgloss`.
- Real-Time Streaming:** Watch the AI think and type token-by-token directly in your terminal.
- 100% Private & Free:** Runs entirely on your local machine via Ollama. No API keys required.
- Context Cancellation:** Press `Esc` at any time to instantly stop the AI generation.
- High Performance:** Built in Go, leveraging Goroutines for non-blocking UI and AI streaming.

## Quick Start

### Prerequisites
1. Install [Go](https://go.dev/dl/)
2. Install [Ollama](https://ollama.com/) and pull a model:
   ```bash
   ollama pull llama3.1:8b
