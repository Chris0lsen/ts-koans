package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/chris0lsen/ts-koans/internal"
)

type state int

const (
	menu state = iota
	editor
)

type model struct {
	state           state
	list            list.Model
	textarea        textarea.Model
	output          string
	selected        int
	exercises       []internal.Exercise
	width           int
	height          int
	running         bool
	spinner         spinner.Model
	persistentState internal.PersistentState
}

type exerciseResultMsg struct {
	output string
}

func initialModel(state internal.PersistentState) model {
	exs := internal.Exercises()
	items := make([]list.Item, len(exs))
	for i, ex := range exs {
		items[i] = ex
	}

	l := list.New(items, list.NewDefaultDelegate(), 30, 14)
	l.Title = "Select an Exercise"
	l.SetShowHelp(false)

	t := textarea.New()
	t.Placeholder = "Edit your code here..."
	t.SetHeight(8)
	t.SetWidth(50)
	t.ShowLineNumbers = true

	s := spinner.New()
	s.Spinner = spinner.Dot

	m := model{
		persistentState: state,
		selected:        state.SelectedIndex,
		state:           menu,
		list:            l,
		textarea:        t,
		exercises:       exs,
		spinner:         s,
	}

	// If user has a saved solution for this exercise, load it into textarea
	if code, ok := state.Solutions[state.SelectedIndex]; ok && code != "" {
		m.textarea.SetValue(code)
	} else {
		m.textarea.SetValue(m.exercises[m.selected].StarterCode)
	}
	// ... any additional model init ...
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func insertSpacesAtCursor(t textarea.Model, n int) textarea.Model {
	t.InsertString(strings.Repeat(" ", n))
	return t
}

func (m model) runExercise() tea.Cmd {
	userCode := m.textarea.Value()
	testScript := m.exercises[m.selected].TestScript
	return func() tea.Msg {
		out, _ := internal.RunTypeScript(userCode, testScript)
		return exerciseResultMsg{output: out}
	}
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.running {
		var spinCmd tea.Cmd
		m.spinner, spinCmd = m.spinner.Update(msg)
		return m, spinCmd
	}

	switch m.state {
	case menu:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
			m.list.SetSize(msg.Width, msg.Height)

		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c":
				// Save before quitting
				m.persistentState.SelectedIndex = m.selected
				m.persistentState.Solutions[m.selected] = m.textarea.Value()
				internal.SaveState(m.persistentState)
				return m, tea.Quit
			case "enter":
				m.persistentState.SelectedIndex = m.selected
				m.persistentState.Solutions[m.selected] = m.textarea.Value()
				internal.SaveState(m.persistentState)
				i := m.list.Index()
				m.selected = i
				selectedEx := m.exercises[i]
				m.textarea.SetValue(selectedEx.StarterCode)
				m.textarea.Focus()
				m.state = editor
			}
		}
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd

	case editor:
		switch msg := msg.(type) {
		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
		case exerciseResultMsg:
			m.running = false
			m.output = msg.output
			return m, nil
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.state = menu
				m.textarea.Blur()
				return m, nil
			case "tab":
				m.textarea = insertSpacesAtCursor(m.textarea, 2)
				return m, nil
			case "f5":
				m.persistentState.Solutions[m.selected] = m.textarea.Value()
				internal.SaveState(m.persistentState)
				m.output = ""    // Clear output
				m.running = true // Start spinner
				return m, m.runExercise()
			}
		}
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	return m, nil
}

var outputStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	Padding(0, 1).
	MarginTop(1).
	Width(50).
	Faint(true)

func (m model) outputColor() lipgloss.Style {
	if strings.Contains(m.output, "✅") {
		return outputStyle.Copy().Foreground(lipgloss.Color("10")) // green
	}
	if strings.Contains(m.output, "❌") {
		return outputStyle.Copy().Foreground(lipgloss.Color("9")) // red
	}
	return outputStyle
}

func (m model) renderOutputPanel() string {
	boxHeight := 10
	lines := strings.Split(m.output, "\n")
	if len(lines) > boxHeight {
		lines = lines[:boxHeight]
	}
	// Pad if fewer than 10 lines
	for len(lines) < boxHeight {
		lines = append(lines, "")
	}
	// Show spinner if running
	outStr := strings.Join(lines, "\n")
	if m.running {
		outStr = m.spinner.View() + " Running..." + "\n" + strings.Repeat(" ", m.width)
	}
	// Colorize: red for compiler error, green for success, yellow for test failure
	style := outputStyle
	switch {
	case strings.Contains(m.output, "error") || strings.Contains(m.output, "Error"):
		style = style.Foreground(lipgloss.Color("9")) // Red
	case strings.Contains(m.output, "✅"):
		style = style.Foreground(lipgloss.Color("10")) // Green
	case strings.Contains(m.output, "❌"):
		style = style.Foreground(lipgloss.Color("11")) // Yellow
	}
	return style.Render(outStr)
}

func (m model) View() string {
	switch m.state {
	case menu:
		// Render the exercise menu (list.Select or similar)
		// and any menu help text
		return m.list.View() + "\n\n[enter] Start | [q] Quit"
	case editor:
		// Render the header, desc, editor, output panel, etc.
		header := lipgloss.NewStyle().Bold(true).Render(m.exercises[m.selected].Title())
		desc := m.exercises[m.selected].Description()
		editor := m.textarea.View()
		output := m.renderOutputPanel()
		return fmt.Sprintf("%s\n%s\n\n%s\n\n%s\n\n[esc] Back | [F5] Run",
			header, desc, editor, output,
		)
	default:
		return "Loading..."
	}
}

func nodeAvailable() bool {
	cmd := exec.Command("node", "--version")
	err := cmd.Run()
	return err == nil
}

func tscAvailable() bool {
	cmd := exec.Command("tsc", "--version")
	err := cmd.Run()
	return err == nil
}

func main() {
	if !nodeAvailable() {
		fmt.Fprintln(os.Stderr, "❌ Node.js not found in PATH. Please install Node.js (https://nodejs.org/) and try again.")
		os.Exit(1)
	}

	if !tscAvailable() {
		fmt.Fprintln(os.Stderr, "❌ tsc not found in PATH. Please install with `npm install -g typescript`.")
		os.Exit(1)
	}

	state, err := internal.LoadState()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Warning: Could not load previous state:", err)
		// You may want to continue with a fresh state, or exit if this is critical
	}

	p := tea.NewProgram(initialModel(state), tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
