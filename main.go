package main

import (
	"bufio"
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/spinner"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/chris0lsen/ts-koans/internal"

	"github.com/lithammer/dedent"
)

type state int

const (
	menu state = iota
	editor
)

// Layout constants
const (
	panelHorizChrome = 8 // Total horizontal chrome: MarginLeft(4) + Border(1) + Padding(1) each side
	panelVertChrome  = 8 // Vertical chrome for editor + output panels (borders, margins)
	editorLeftX      = 6 // MarginLeft(4) + Border(1) + Padding(1)
	editorMarginTop  = 2 // Editor MarginTop(1) + top border line

	maxOutputHeight = 10
	minEditorHeight = 3
	minOutputHeight = 3
	debugPanelLines = 5
	minGutterWidth  = 3

	listItemHeight = 3 // Default delegate Height(2) + Spacing(1)

	tabWidth       = 2
	maxBufferLines = 100 // Max retained output/debug lines
)

type model struct {
	program         *tea.Program
	state           state
	list            list.Model
	textarea        textarea.Model
	outputLines     []runnerOutputMsg
	selected        int
	exercises       []internal.Exercise
	width           int
	height          int
	running         bool
	spinner         spinner.Model
	persistentState internal.PersistentState
	debugMode       bool
	debugLog        []string
	start           int
	gutterWidth     int
	editorTopY      int
	editorHeight    int
	outputHeight    int
}

type setProgramMsg struct{ program *tea.Program }

type runnerDebugMsg struct{ Line string }
type runnerOutputMsg struct {
	Line      string
	Assertion bool
}
type runnerDoneMsg struct{ Err error }

var (
	headerStyle = lipgloss.NewStyle().
			Bold(true).
			PaddingTop(2).PaddingLeft(4)

	descStyle = lipgloss.NewStyle().
			PaddingTop(0).PaddingBottom(1).PaddingLeft(4)

	editorStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(0, 1).
			MarginTop(1).
			MarginLeft(4)

	outputStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(0, 1).
			MarginTop(1).
			MarginLeft(4).
			Faint(true)

	debugStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder()).
			Padding(0, 1).
			MarginTop(1).
			MarginLeft(4).
			Faint(true).
			Foreground(lipgloss.Color("8"))

	helpStyle = lipgloss.NewStyle().
			MarginTop(1).
			PaddingLeft(4)

	infoStyle = lipgloss.NewStyle().
			MarginLeft(2).
			MarginTop(1).
			PaddingLeft(1).
			Border(lipgloss.NormalBorder())

	lineNumStyle       = lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	cursorLineNumStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	assertionStyle     = lipgloss.NewStyle().Foreground(lipgloss.Color("12")).Bold(true)
)

func printHelpfulTSCErrors(tscOutput, harnessPath string, p *tea.Program) {
	harnessBytes, _ := os.ReadFile(harnessPath)
	harnessLines := strings.Split(string(harnessBytes), "\n")
	scanner := bufio.NewScanner(strings.NewReader(tscOutput))
	// errorLinePattern matches e.g. "typecheck.ts(10,44): error ..."
	errorLinePattern := regexp.MustCompile(`typecheck\.ts\((\d+),\d+\): error`)
	for scanner.Scan() {
		line := scanner.Text()
		if m := errorLinePattern.FindStringSubmatch(line); m != nil {
			lineNum, _ := strconv.Atoi(m[1])
			// Print the comment line (if any) above the assertion
			if lineNum-2 >= 0 && lineNum-2 < len(harnessLines) {
				p.Send(runnerOutputMsg{Line: harnessLines[lineNum-2], Assertion: true})
			}
		}
		// Always print the error itself
		p.Send(runnerOutputMsg{Line: line})
	}
}

func (m *model) appendDebug(msg string) {
	if m.debugMode {
		m.debugLog = append(m.debugLog, msg)
		if len(m.debugLog) > maxBufferLines {
			m.debugLog = m.debugLog[len(m.debugLog)-maxBufferLines:]
		}
	}
}

func (m *model) recalcEditorHeight() {
	if m.width == 0 || m.height == 0 {
		return
	}
	header := headerStyle.Render(m.exercises[m.selected].Title())
	desc := descStyle.Render(m.exercises[m.selected].Description())
	m.editorTopY = lipgloss.Height(header) + lipgloss.Height(desc) + editorMarginTop

	debugPanelHeight := 0
	if m.debugMode {
		debugPanelHeight = lipgloss.Height(debugStyle.Width(m.width - panelHorizChrome).Render(strings.Repeat("\n", debugPanelLines-1)))
	}

	help := helpStyle.Render("[esc] Back | [F5] Run | [shift + ← / → ] Prev/Next Exercise")

	fixedHeight := lipgloss.Height(header) +
		lipgloss.Height(desc) +
		debugPanelHeight +
		lipgloss.Height(help)

	available := m.height - fixedHeight - panelVertChrome

	m.outputHeight = maxOutputHeight
	m.editorHeight = available - m.outputHeight
	if m.editorHeight < minEditorHeight {
		m.editorHeight = minEditorHeight
		m.outputHeight = available - m.editorHeight
		if m.outputHeight < minOutputHeight {
			m.outputHeight = minOutputHeight
		}
	}

	m.textarea.SetWidth(m.width - panelHorizChrome)
	m.textarea.SetHeight(m.editorHeight)
}

func makeListItems(exs []internal.Exercise, completed map[int]bool) []list.Item {
	items := make([]list.Item, len(exs))
	for i, ex := range exs {
		label := ex.Title()
		if completed[i] {
			label = "✅ " + label
		}
		ex2 := ex
		ex2.Label = label
		items[i] = ex2
	}
	return items
}

func initialModel(state internal.PersistentState) model {
	exs := internal.Exercises()
	items := make([]list.Item, len(exs))
	for i, ex := range exs {
		items[i] = ex
	}

	l := list.New(makeListItems(exs, state.Completed), list.NewDefaultDelegate(), 30, 14)
	l.Title = "Select an Exercise"
	l.SetShowHelp(false)

	t := textarea.New()
	t.Placeholder = "Edit your code here..."
	t.SetHeight(16)
	t.SetWidth(80)
	t.ShowLineNumbers = false

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
	return m
}

func (m model) Init() tea.Cmd {
	return nil
}

func insertSpacesAtCursor(t textarea.Model, n int) textarea.Model {
	t.InsertString(strings.Repeat(" ", n))
	return t
}

func copyVersionFilesToTempDir(tempDir string) {
	for _, filename := range []string{".tool-versions", ".nvmrc", ".node-version"} {
		if data, err := os.ReadFile(filename); err == nil {
			os.WriteFile(filepath.Join(tempDir, filename), data, 0600)
		}
	}
}

func runExerciseStreamed(userCode string, program *tea.Program, ex internal.Exercise) tea.Cmd {
	return func() tea.Msg {
		tmpDir, err := os.MkdirTemp("", "tskoans-*")
		if err != nil {
			program.Send(runnerOutputMsg{Line: fmt.Sprintf("Failed to create temp dir: %v", err)})
			program.Send(runnerDoneMsg{Err: err})
			return nil
		}
		copyVersionFilesToTempDir(tmpDir)
		defer os.RemoveAll(tmpDir)

		if err := compileTypeScript(tmpDir, userCode, ex, program); err != nil {
			program.Send(runnerDoneMsg{Err: err})
			return nil
		}

		if err := writeTestBundle(tmpDir, ex.TestScript); err != nil {
			program.Send(runnerOutputMsg{Line: fmt.Sprintf("Failed to write test bundle: %v", err)})
			program.Send(runnerDoneMsg{Err: err})
			return nil
		}

		err = runNodeTests(tmpDir, program)
		program.Send(runnerDoneMsg{Err: err})
		return nil
	}
}

// compileTypeScript writes the user code + type harness + assertions to a .ts file,
// then runs tsc. Returns nil on success, or the tsc error (after sending output messages).
func compileTypeScript(tmpDir, userCode string, ex internal.Exercise, program *tea.Program) error {
	typecheckPath := filepath.Join(tmpDir, "typecheck.ts")
	fullTypecheck := userCode + "\n\n" + internal.TypeHarness + "\n" + ex.TypeAssertions + "\n"

	if err := os.WriteFile(typecheckPath, []byte(fullTypecheck), 0644); err != nil {
		program.Send(runnerOutputMsg{Line: fmt.Sprintf("Failed to write typecheck.ts: %v", err)})
		return err
	}

	tscCmd := exec.Command("tsc", typecheckPath, "--target", "es2020", "--module", "commonjs", "--outDir", tmpDir)
	tscCmd.Dir = tmpDir

	var tscStderrBuf, tscStdoutBuf bytes.Buffer
	tscCmd.Stderr = &tscStderrBuf
	tscCmd.Stdout = &tscStdoutBuf

	if err := tscCmd.Run(); err != nil {
		program.Send(runnerDebugMsg{Line: fmt.Sprintf("tsc exit error: %v", err)})
		program.Send(runnerDebugMsg{Line: "STDERR: " + tscStderrBuf.String()})
		program.Send(runnerDebugMsg{Line: "STDOUT: " + tscStdoutBuf.String()})
		printHelpfulTSCErrors(tscStdoutBuf.String(), typecheckPath, program)
		program.Send(runnerOutputMsg{Line: fmt.Sprintf("[tsc] Compilation failed: %v", err)})
		return err
	}
	return nil
}

// writeTestBundle reads the compiled JS, combines it with the test script,
// and writes both run.js and runner.mjs into tmpDir.
func writeTestBundle(tmpDir, testScript string) error {
	jsBytes, err := os.ReadFile(filepath.Join(tmpDir, "typecheck.js"))
	if err != nil {
		return fmt.Errorf("read compiled JS: %w", err)
	}

	combined := fmt.Sprintf("%s\n\n%s\n", string(jsBytes), testScript)
	if err := os.WriteFile(filepath.Join(tmpDir, "run.js"), []byte(combined), 0644); err != nil {
		return fmt.Errorf("write run.js: %w", err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "runner.mjs"), []byte(internal.RunnerMJS), 0644); err != nil {
		return fmt.Errorf("write runner.mjs: %w", err)
	}
	return nil
}

// runNodeTests executes runner.mjs with a timeout and sends stdout/stderr as output messages.
func runNodeTests(tmpDir string, program *tea.Program) error {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	nodeCmd := exec.CommandContext(ctx, "node", filepath.Join(tmpDir, "runner.mjs"))
	nodeCmd.Dir = tmpDir

	var nodeStdoutBuf, nodeStderrBuf bytes.Buffer
	nodeCmd.Stdout = &nodeStdoutBuf
	nodeCmd.Stderr = &nodeStderrBuf

	err := nodeCmd.Run()

	if ctx.Err() == context.DeadlineExceeded {
		program.Send(runnerOutputMsg{Line: "[node] Execution timed out!"})
		return ctx.Err()
	}
	if nodeStdoutBuf.Len() > 0 {
		program.Send(runnerOutputMsg{Line: "[node stdout] " + nodeStdoutBuf.String()})
	}
	if nodeStderrBuf.Len() > 0 {
		program.Send(runnerOutputMsg{Line: "[node stderr] " + nodeStderrBuf.String()})
	}
	if err != nil {
		program.Send(runnerOutputMsg{Line: fmt.Sprintf("[node] Test runner failed: %v", err)})
	}
	return err
}

func outputLinesToString(lines []runnerOutputMsg) string {
	strs := make([]string, len(lines))
	for i, msg := range lines {
		strs[i] = msg.Line
	}
	return strings.Join(strs, "\n")
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case menu:
		return m.updateMenu(msg)
	case editor:
		return m.updateEditor(msg)
	}
	return m, nil
}

// listHeaderHeight returns the number of terminal rows the list widget
// renders above the first item (title/filter bar + status bar).
func (m model) listHeaderHeight() int {
	h := 0
	if m.list.ShowTitle() || m.list.ShowFilter() {
		// Title bar: 1 line of text + 1 line of bottom padding (default style)
		h += 2
	}
	if m.list.ShowStatusBar() {
		// Status bar: 1 line of text + 1 line of bottom padding (default style)
		h += 2
	}
	return h
}

func (m model) updateMenu(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case setProgramMsg:
		m.program = msg.program

	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		// Reserve 2 lines for the help text rendered below the list
		m.list.SetSize(msg.Width, msg.Height-2)

	case tea.MouseMsg:
		// Don't intercept clicks while the filter input is focused
		if !m.list.SettingFilter() {
			if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
				slot := (msg.Y - m.listHeaderHeight()) / listItemHeight
				page := m.list.Paginator.Page
				perPage := m.list.Paginator.PerPage
				idx := page*perPage + slot
				if idx >= 0 && idx < len(m.list.VisibleItems()) {
					m.list.Select(idx)
					return m, nil
				}
			}
		}

	case tea.KeyMsg:
		// While filtering, let the list widget handle all keys
		// (Enter accepts the filter, Escape cancels it)
		if !m.list.SettingFilter() {
			switch msg.String() {
			case "q", "ctrl+c":
				m.saveState()
				return m, tea.Quit
			case "enter":
				m.outputLines = nil
				m.switchToExercise(m.list.GlobalIndex())
				m.textarea.Focus()
				m.state = editor
			}
		}
	}
	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m model) updateEditor(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case runnerOutputMsg:
		m.outputLines = append(m.outputLines, msg)
		if len(m.outputLines) > maxBufferLines {
			m.outputLines = m.outputLines[len(m.outputLines)-maxBufferLines:]
		}
		m.recalcEditorHeight()

	case runnerDebugMsg:
		if m.debugMode {
			m.appendDebug(msg.Line)
		}

		return m, nil

	case runnerDoneMsg:
		m.running = false
		m.recalcEditorHeight()

		output := outputLinesToString(m.outputLines)
		if strings.Contains(output, "✅") {
			if m.persistentState.Completed == nil {
				m.persistentState.Completed = make(map[int]bool)
			}
			m.persistentState.Completed[m.selected] = true
			m.list.SetItems(makeListItems(m.exercises, m.persistentState.Completed))
			m.saveState()
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		m.recalcEditorHeight()
		m.textarea.SetCursor(0)
		return m, nil
	case tea.MouseMsg:
		// Click to move cursor
		if msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft {
			editorBottomY := m.editorTopY + m.textarea.Height() - 1
			contentLeftX := editorLeftX + m.gutterWidth

			// Get code coordinates
			targetLine := msg.Y - m.editorTopY + m.start
			targetCol := msg.X - contentLeftX

			// Bounds checks (ignore clicks outside editor)
			totalLines := len(strings.Split(m.textarea.Value(), "\n"))
			if msg.Y < m.editorTopY || msg.Y >= editorBottomY || msg.X < contentLeftX || targetLine >= totalLines {
				return m, nil
			}
			if targetCol < 0 {
				targetCol = 0
			}

			// Move cursor in textarea
			currentLine := m.textarea.Line()
			for currentLine < targetLine {
				m.textarea, _ = m.textarea.Update(tea.KeyMsg{Type: tea.KeyDown})
				currentLine++
			}
			for currentLine > targetLine {
				m.textarea, _ = m.textarea.Update(tea.KeyMsg{Type: tea.KeyUp})
				currentLine--
			}

			// Clamp col to len of target line, set col
			lines := strings.Split(m.textarea.Value(), "\n")
			lineLen := len([]rune(lines[targetLine]))
			if targetCol > lineLen {
				targetCol = lineLen
			}

			m.textarea.SetCursor(targetCol)

			m.calculateCursorCoordinates()

			return m, nil
		}
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			// Save before quitting
			m.saveState()
			return m, tea.Quit
		case "esc":
			m.outputLines = nil
			m.saveState()
			m.state = menu
			m.textarea.Blur()
			return m, nil
		case "tab":
			m.textarea = insertSpacesAtCursor(m.textarea, tabWidth)
			return m, nil
		case "f5":
			m.saveState()
			userCode := m.textarea.Value()
			m.outputLines = nil
			m.running = true
			m.recalcEditorHeight()
			return m, tea.Batch(m.spinner.Tick, runExerciseStreamed(userCode, m.program, m.exercises[m.selected]))
		case "shift+right":
			m.switchToExercise((m.selected + 1) % len(m.exercises))
			return m, nil
		case "shift+left":
			m.switchToExercise((m.selected - 1 + len(m.exercises)) % len(m.exercises))
			return m, nil
		}
	}
	if m.running {
		var spinCmd tea.Cmd
		m.spinner, spinCmd = m.spinner.Update(msg)
		return m, spinCmd
	}
	var cmd tea.Cmd
	m.textarea, cmd = m.textarea.Update(msg)
	// Calculate cursor start
	m.calculateCursorCoordinates()
	return m, cmd
}

func (m *model) saveState() {
	m.persistentState.SelectedIndex = m.selected
	m.persistentState.Solutions[m.selected] = m.textarea.Value()
	internal.SaveState(m.persistentState)
}

func (m *model) switchToExercise(i int) {
	m.saveState()
	m.selected = i
	if code, ok := m.persistentState.Solutions[i]; ok && code != "" {
		m.textarea.SetValue(code)
	} else {
		m.textarea.SetValue(m.exercises[i].StarterCode)
	}
	m.recalcEditorHeight()
}

func (m *model) calculateCursorCoordinates() {
	curRow := m.textarea.Line()
	viewHeight := m.textarea.Height()
	totalLines := len(strings.Split(m.textarea.Value(), "\n"))

	m.start = 0
	if curRow >= m.start+viewHeight {
		m.start = curRow - viewHeight + 1
	}
	if totalLines <= viewHeight {
		m.start = 0
	}

	numWidth := len(fmt.Sprintf("%d", totalLines))
	if numWidth < minGutterWidth {
		numWidth = minGutterWidth
	}
	m.gutterWidth = numWidth + 1 // line numbers + space
}

func (m model) renderOutputPanel(boxHeight int) string {
	style := outputStyle.Width(m.width - panelHorizChrome)

	if m.running {
		spin := m.spinner.View()
		lines := []string{spin + " Running..."}
		for len(lines) < boxHeight {
			lines = append(lines, "")
		}
		return style.Render(strings.Join(lines, "\n"))
	}

	lines := m.outputLines
	if len(lines) > boxHeight {
		lines = lines[len(lines)-boxHeight:]
	}
	for len(lines) < boxHeight {
		lines = append(lines, runnerOutputMsg{Line: ""})
	}

	renderedLines := make([]string, len(lines))
	for i, msg := range lines {
		if msg.Assertion {
			renderedLines[i] = assertionStyle.Render(msg.Line)
		} else {
			renderedLines[i] = msg.Line
		}
	}

	return style.Render(strings.Join(renderedLines, "\n"))
}

func (m model) View() string {
	switch m.state {
	case menu:
		return m.list.View() + "\n\n[enter] Start | [q] Quit | [ ← / → ] Prev/Next Page | [/] Filter"
	case editor:
		header := headerStyle.Render(m.exercises[m.selected].Title())
		desc := descStyle.Render(m.exercises[m.selected].Description())

		debugPanel := ""
		if m.debugMode {
			n := len(m.debugLog)
			start := 0
			if n > debugPanelLines {
				start = n - debugPanelLines
			}
			logs := m.debugLog[start:]
			debugPanel = debugStyle.Width(m.width - panelHorizChrome).Render(strings.Join(logs, "\n"))
		}

		help := helpStyle.Render("[esc] Back | [F5] Run | [shift + ← / → ] Prev/Next Exercise")

		output := m.renderOutputPanel(m.outputHeight)

		// Set editor size
		infoChrome := 6 // infoStyle MarginLeft(2) + Border(1) + PaddingLeft(1) + right border(1) + safety(1)
		editorWidth := (m.width - panelHorizChrome) * 60 / 100
		editor := editorStyle.Width(editorWidth).Height(m.editorHeight).Render(m.renderHighlightedCode(m.editorHeight))

		// Join help text panel horizontally with editor (when enough width)
		if m.width > 80 && m.exercises[m.selected].Info() != "" {
			infoWidth := m.width - panelHorizChrome - editorWidth - infoChrome
			if infoWidth > 10 {
				info := infoStyle.Width(infoWidth).Height(m.editorHeight).Render(dedent.Dedent(m.exercises[m.selected].Info()))
				editor = lipgloss.JoinHorizontal(lipgloss.Top, editor, info)
			}
		}

		// Compose
		panels := []string{header, desc, editor, output}
		if m.debugMode {
			panels = append(panels, debugPanel)
		}
		panels = append(panels, help)

		return lipgloss.JoinVertical(lipgloss.Left, panels...)
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

	debug := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	state, err := internal.LoadState()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Warning: Could not load previous state:", err)
	}
	m := initialModel(state)
	m.debugMode = *debug
	m.debugLog = append(m.debugLog, "Debug panel is working!")

	p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
	go func() {
		p.Send(setProgramMsg{program: p})
	}()
	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
