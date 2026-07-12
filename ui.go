package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// --- CATPPUCCIN MOCHA THEME ---
var (
	BaseBg    = lipgloss.Color("#1E1E2E")
	TextFg    = lipgloss.Color("#CDD6F4")
	Blue      = lipgloss.Color("#89B4FA")
	Pink      = lipgloss.Color("#F38BA8")
	Green     = lipgloss.Color("#A6E3A1")
	Subtext   = lipgloss.Color("#6C7086")

	TitleStyle = lipgloss.NewStyle().
			Bold(true).
			Foreground(BaseBg).
			Background(Blue).
			Padding(0, 1)

	BorderStyle = lipgloss.NewStyle().
			BorderStyle(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("#45475A")).
			Padding(1, 2).
			Background(BaseBg)

	PromptStyle = lipgloss.NewStyle().Foreground(Pink).Bold(true)
	AIStyle     = lipgloss.NewStyle().Foreground(Green)
	SubtextStyle = lipgloss.NewStyle().Foreground(Subtext).Italic(true)
)

// --- BUBBLETEA MESSAGES ---
type aiTokenMsg string
type aiDoneMsg struct{}

// --- THE MODEL ---
type Model struct {
	aiEngine  *AIEngine
	viewport  viewport.Model
	textInput textinput.Model
	spinner   spinner.Model
	
	messages  []string // Stores the chat history
	query     string   // Current input
	
	aiRunning bool
	ctx       context.Context
	cancel    context.CancelFunc
	
	width     int
	height    int
}

func NewModel(ai *AIEngine) Model {
	ti := textinput.New()
	ti.Placeholder = "Ask Forge to analyze your code..."
	ti.PlaceholderStyle = lipgloss.NewStyle().Foreground(Subtext)
	ti.TextStyle = lipgloss.NewStyle().Foreground(TextFg)
	ti.CharLimit = 256
	ti.Width = 50
	ti.Focus()

	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(Blue)

	ctx, cancel := context.WithCancel(context.Background())

	return Model{
		aiEngine:  ai,
		textInput: ti,
		spinner:   s,
		messages:  []string{AIStyle.Render("⚡ Forge TUI initialized. Local AI is ready.")},
		ctx:       ctx,
		cancel:    cancel,
	}
}

func (m Model) Init() tea.Cmd {
	return tea.Batch(m.spinner.Tick, m.listenForAI())
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		
		// Calculate viewport dimensions
		verticalMargins := 6 // Title + borders + input area
		m.viewport.Width = msg.Width - 6
		m.viewport.Height = msg.Height - verticalMargins
		
		m.textInput.Width = msg.Width - 10

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			m.cancel() // Cancel any running AI context
			return m, tea.Quit
		case "enter":
			if !m.aiRunning && m.textInput.Value() != "" {
				m.query = m.textInput.Value()
				m.messages = append(m.messages, PromptStyle.Render("❯ ")+m.query)
				m.textInput.SetValue("")
				m.aiRunning = true
				
				// Create a new context for this specific AI request
				m.ctx, m.cancel = context.WithCancel(context.Background())
				cmds = append(cmds, m.startAI(m.query))
			}
		case "esc":
			if m.aiRunning {
				m.cancel() // Stop the AI mid-generation
				m.aiRunning = false
				m.messages = append(m.messages, SubtextStyle.Render("[Generation cancelled]"))
			}
		}

	case aiTokenMsg:
		// Append token to the last message or create a new one
		if len(m.messages) > 0 && !m.isLastMessageAI() {
			m.messages = append(m.messages, AIStyle.Render(string(msg)))
		} else {
			m.messages[len(m.messages)-1] += string(msg)
		}
		cmds = append(cmds, m.listenForAI())

	case aiDoneMsg:
		m.aiRunning = false
		m.messages = append(m.messages, "") // Add a blank line for spacing

	case spinner.TickMsg:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		cmds = append(cmds, cmd)
	}

	// Update text input
	var tiCmd tea.Cmd
	m.textInput, tiCmd = m.textInput.Update(msg)
	cmds = append(cmds, tiCmd)

	// Update viewport content and auto-scroll
	m.viewport.SetContent(strings.Join(m.messages, "\n"))
	m.viewport.GotoBottom()

	return m, tea.Batch(cmds...)
}

func (m Model) View() string {
	if m.width == 0 {
		return "Initializing Forge TUI..."
	}

	// Header
	title := TitleStyle.Render("⚡ FORGE TUI: Local AI Agent")
	
	// Main Content Area (Viewport)
	content := BorderStyle.Width(m.width - 4).Height(m.height - 6).Render(m.viewport.View())

	// Footer / Input Area
	var footer string
	if m.aiRunning {
		footer = fmt.Sprintf("%s %s", m.spinner.View(), SubtextStyle.Render("Press [Esc] to cancel"))
	} else {
		footer = PromptStyle.Render("❯ ") + m.textInput.View()
	}
	
	footerStyled := lipgloss.NewStyle().
		Foreground(TextFg).
		Background(BaseBg).
		Padding(0, 2).
		Width(m.width - 4).
		Render(footer)

	return lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		content,
		footerStyled,
	)
}

// --- HELPER METHODS ---

func (m Model) isLastMessageAI() bool {
	if len(m.messages) == 0 { return false }
	last := m.messages[len(m.messages)-1]
	return strings.Contains(last, AIStyle.Render("")) // Hacky but works for simple streaming
}

// Channels to bridge the async Ollama stream with Bubbletea's sync UI
var tokenChan = make(chan string)
var doneChan = make(chan bool)

func (m Model) startAI(prompt string) tea.Cmd {
	return func() tea.Msg {
		go m.aiEngine.StreamPrompt(m.ctx, prompt, tokenChan, doneChan)
		return nil
	}
}

func (m Model) listenForAI() tea.Cmd {
	return func() tea.Msg {
		select {
		case token := <-tokenChan:
			return aiTokenMsg(token)
		case <-doneChan:
			return aiDoneMsg{}
		case <-m.ctx.Done():
			return aiDoneMsg{}
		}
	}
}