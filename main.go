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
)

type state int

const (
	menu state = iota
	editor
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
}

type exerciseResultMsg struct {
	output string
}
type setProgramMsg struct{ program *tea.Program }

type runnerDebugMsg struct{ Line string }
type runnerOutputMsg struct {
	Line      string
	Assertion bool
}
type runnerDoneMsg struct{ Err error }

// errorLinePattern matches e.g. "typecheck.ts(10,44): error ..."
var errorLinePattern = regexp.MustCompile(`typecheck\.ts\((\d+),\d+\): error`)

func printHelpfulTSCErrors(tscOutput, harnessPath string, p *tea.Program) {
	harnessBytes, _ := os.ReadFile(harnessPath)
	harnessLines := strings.Split(string(harnessBytes), "\n")
	scanner := bufio.NewScanner(strings.NewReader(tscOutput))
	errorLinePattern := regexp.MustCompile(`typecheck\.ts\((\d+),\d+\): error`)
	for scanner.Scan() {
		line := scanner.Text()
		if m := errorLinePattern.FindStringSubmatch(line); m != nil {
			lineNum, _ := strconv.Atoi(m[1])
			// Print the comment line (if any) above the assertion
			if lineNum-2 >= 0 && lineNum-2 < len(harnessLines) {
				p.Send(runnerOutputMsg{Line: harnessLines[lineNum-2], Assertion: true})
			}
			// Print the assertion line itself
			// if lineNum-1 >= 0 && lineNum-1 < len(harnessLines) {
			// 	p.Send(runnerOutputMsg{Line: harnessLines[lineNum-1]})
			// }
		}
		// Always print the error itself
		p.Send(runnerOutputMsg{Line: line})
	}
}

func (m *model) appendDebug(msg string) {
	if m.debugMode {
		m.debugLog = append(m.debugLog, msg)
		if len(m.debugLog) > 100 {
			m.debugLog = m.debugLog[len(m.debugLog)-100:]
		}
	}
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

	//l := list.New(items, list.NewDefaultDelegate(), 30, 14)
	l := list.New(makeListItems(exs, state.Completed), list.NewDefaultDelegate(), 30, 14)
	l.Title = "Select an Exercise"
	l.SetShowHelp(false)

	t := textarea.New()
	t.Placeholder = "Edit your code here..."
	t.SetHeight(16)
	t.SetWidth(80)
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

func runExerciseStreamed(userCode, testScript string, program *tea.Program, ex internal.Exercise) tea.Cmd {
	return func() tea.Msg {
		// Create temp dir
		tmpDir, err := os.MkdirTemp("", "tskoans-*")
		if err != nil {
			program.Send(runnerOutputMsg{Line: fmt.Sprintf("Failed to create temp dir: %v", err)})
			program.Send(runnerDoneMsg{Err: err})
			return nil
		}
		copyVersionFilesToTempDir(tmpDir)

		defer os.RemoveAll(tmpDir)

		typecheckPath := filepath.Join(tmpDir, "typecheck.ts")

		typeHarness := `// -- Type Assertion Utilities --

// Produces a type error if 'T' is not 'true'
type Assert<T extends true> = T;

// Checks if two types are the same (structurally)
type IsType<A, B> = (<T>() => T extends A ? 1 : 2) extends
                    (<T>() => T extends B ? 1 : 2) ? true : false;

// Checks if type A is not assignable to B
type IsNotType<A, B> = IsType<A, B> extends true ? false : true;

type IsNotReadonly<T, K extends keyof T> =
  // Compare: If { [P in K]: T[P] } is assignable to { -readonly [P in K]: T[P] }
  // then K is not readonly
  { [P in K]: T[P] } extends { -readonly [P in K]: T[P] }
    ? true
    : false;


// Checks if type A is assignable to B
type IsAssignable<A, B> = A extends B ? true : false;

`

		fullTypecheck := string(userCode) + "\n\n" + typeHarness + "\n" + ex.TypeAssertions + "\n"

		if err := os.WriteFile(typecheckPath, []byte(fullTypecheck), 0644); err != nil {
			program.Send(runnerOutputMsg{Line: fmt.Sprintf("Failed to write typecheck.ts: %v", err)})
			program.Send(runnerDoneMsg{Err: err})
			return nil
		}

		tscCmd := exec.Command("tsc", typecheckPath, "--target", "es2020", "--module", "commonjs", "--outDir", tmpDir)
		tscCmd.Dir = tmpDir

		var tscStderrBuf bytes.Buffer
		tscCmd.Stderr = &tscStderrBuf

		var tscStdoutBuf bytes.Buffer
		tscCmd.Stdout = &tscStdoutBuf

		if err := tscCmd.Run(); err != nil {
			program.Send(runnerDebugMsg{Line: fmt.Sprintf("tsc exit error: %v", err)})

			program.Send(runnerDebugMsg{Line: "STDERR: " + tscStderrBuf.String()})
			program.Send(runnerDebugMsg{Line: "STDOUT: " + tscStdoutBuf.String()})

			// Print annotated errors to your TUI/debug panel/logs
			printHelpfulTSCErrors(tscStdoutBuf.String(), typecheckPath, program)

			program.Send(runnerOutputMsg{Line: fmt.Sprintf("[tsc] Compilation failed: %v", err)})
			program.Send(runnerDoneMsg{Err: err})
			return nil
		}

		jsBytes, err := os.ReadFile(filepath.Join(tmpDir, "typecheck.js"))
		if err != nil { /* handle error */
		}

		combined := fmt.Sprintf("%s\n\n%s\n", string(jsBytes), ex.TestScript)

		// Write to a single file, say, run.js
		runPath := filepath.Join(tmpDir, "run.js")
		os.WriteFile(runPath, []byte(combined), 0644)
		// Write runner.mjs
		runnerPath := filepath.Join(tmpDir, "runner.mjs")
		runner := fmt.Sprintf(`
import { readFileSync } from "fs";
import vm from "vm";

const combined = readFileSync("./run.js", "utf8");

// Optionally: set up a basic sandbox (exports/global, etc.)
const sandbox = { exports: {}, module: { exports: {}}, console };
const context = vm.createContext(sandbox);

try {
  vm.runInContext(combined, context, { timeout: 1000 });
  console.log("✅ All tests passed!")
} catch (err) {
  // Print clean error message
  console.log("❌ Test failed:", err && err.message ? err.message : err);
  process.exit(1);
}

`)
		if err := os.WriteFile(runnerPath, []byte(runner), 0644); err != nil {
			program.Send(runnerOutputMsg{Line: fmt.Sprintf("Failed to write runner.mjs: %v", err)})
			program.Send(runnerDoneMsg{Err: err})
			return nil
		}

		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
		defer cancel()

		nodeCmd := exec.CommandContext(ctx, "node", runnerPath)
		nodeCmd.Dir = tmpDir

		var nodeStdoutBuf, nodeStderrBuf bytes.Buffer
		nodeCmd.Stdout = &nodeStdoutBuf
		nodeCmd.Stderr = &nodeStderrBuf

		err = nodeCmd.Run()

		if ctx.Err() == context.DeadlineExceeded {
			program.Send(runnerOutputMsg{Line: "[node] Execution timed out!"})
			program.Send(runnerDoneMsg{Err: ctx.Err()})
			return nil
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
		program.Send(runnerDoneMsg{Err: err})
		return nil

	}
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
		switch msg := msg.(type) {
		case setProgramMsg:
			m.program = msg.program

		case tea.WindowSizeMsg:
			m.width = msg.Width
			m.height = msg.Height
			m.list.SetSize(msg.Width, msg.Height)

		case tea.KeyMsg:
			switch msg.String() {
			case "q", "ctrl+c":
				// Save before quitting
				m.persistentState.SelectedIndex = m.selected
				internal.SaveState(m.persistentState)
				return m, tea.Quit
			case "enter":
				m.outputLines = nil
				m.persistentState.SelectedIndex = m.selected
				internal.SaveState(m.persistentState)
				i := m.list.Index()
				m.selected = i
				selectedEx := m.exercises[i]
				// Restore user's previous solution if it exists
				if code, ok := m.persistentState.Solutions[i]; ok && code != "" {
					m.textarea.SetValue(code)
				} else {
					m.textarea.SetValue(selectedEx.StarterCode)
				}
				m.textarea.Focus()
				m.state = editor
			}
		}
		var cmd tea.Cmd
		m.list, cmd = m.list.Update(msg)
		return m, cmd

	case editor:
		switch msg := msg.(type) {
		case runnerOutputMsg:
			m.outputLines = append(m.outputLines, msg)
			if len(m.outputLines) > 100 {
				m.outputLines = m.outputLines[len(m.outputLines)-100:]
			}

		case runnerDebugMsg:
			if m.debugMode {
				m.appendDebug(msg.Line)
			}

			return m, nil

		case runnerDoneMsg:
			m.running = false

			output := outputLinesToString(m.outputLines)
			if strings.Contains(output, "✅") {
				if m.persistentState.Completed == nil {
					m.persistentState.Completed = make(map[int]bool)
				}
				m.persistentState.Completed[m.selected] = true
				m.list.SetItems(makeListItems(m.exercises, m.persistentState.Completed))
				internal.SaveState(m.persistentState)
			}
			return m, nil
		case tea.WindowSizeMsg:

			m.width = msg.Width
			m.height = msg.Height

			descLines := len(strings.Split(m.exercises[m.selected].Description(), "\n"))
			headerLines := 1
			blankAbove := 2
			blankBelow := 1
			outputPanelHeight := 10
			debugPanelHeight := 0
			if m.debugMode {
				debugPanelHeight = 5
			}
			blankBelowPanels := 1
			helpLines := 1

			linesAboveEditor := headerLines + descLines + blankAbove
			linesBelowEditor := blankBelow + outputPanelHeight + debugPanelHeight + blankBelowPanels + helpLines

			editorHeight := m.height - linesAboveEditor - linesBelowEditor
			if editorHeight < 3 {
				editorHeight = 3
			}

			m.textarea.SetHeight(editorHeight)
			m.textarea.SetWidth(m.width)
			m.textarea.SetCursor(0)

			return m, nil
		case tea.KeyMsg:
			switch msg.String() {
			case "esc":
				m.outputLines = nil
				m.persistentState.SelectedIndex = m.selected
				m.persistentState.Solutions[m.selected] = m.textarea.Value()
				internal.SaveState(m.persistentState)
				m.state = menu
				m.textarea.Blur()
				return m, nil
			case "tab":
				m.textarea = insertSpacesAtCursor(m.textarea, 2)
				return m, nil
			case "f5":
				m.persistentState.Solutions[m.selected] = m.textarea.Value()
				internal.SaveState(m.persistentState)
				userCode := m.textarea.Value()
				testScript := m.exercises[m.selected].TestScript
				m.outputLines = nil
				m.running = true
				return m, runExerciseStreamed(userCode, testScript, m.program, m.exercises[m.selected])
			}
		}
		var cmd tea.Cmd
		m.textarea, cmd = m.textarea.Update(msg)
		return m, cmd
	}

	if m.running {
		var spinCmd tea.Cmd
		m.spinner, spinCmd = m.spinner.Update(msg)
		return m, spinCmd
	}
	return m, nil
}

var outputStyle = lipgloss.NewStyle().
	Border(lipgloss.NormalBorder()).
	Padding(0, 1).
	MarginTop(1).
	MarginLeft(1).
	MarginRight(1).
	Width(80).
	Faint(true)

func (m model) renderOutputPanel() string {
	boxHeight := 10

	if m.running {
		spin := m.spinner.View()
		lines := []string{spin + " Running..."}
		for len(lines) < boxHeight {
			lines = append(lines, "")
		}
		return outputStyle.Render(strings.Join(lines, "\n"))
	}

	lines := m.outputLines
	if len(lines) > boxHeight {
		lines = lines[len(lines)-boxHeight:]
	}
	for len(lines) < boxHeight {
		lines = append(lines, runnerOutputMsg{Line: ""})
	}

	// Style for assertion lines
	assertionStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("9")).Bold(true)
	renderedLines := make([]string, len(lines))
	for i, msg := range lines {
		if msg.Assertion {
			renderedLines[i] = assertionStyle.Render(msg.Line)
		} else {
			renderedLines[i] = msg.Line
		}
	}

	return outputStyle.Copy().Width(m.width).Render(strings.Join(renderedLines, "\n"))
}

func (m model) View() string {
	switch m.state {
	case menu:
		return m.list.View() + "\n\n[enter] Start | [q] Quit"
	case editor:
		header := lipgloss.NewStyle().Bold(true).Render(m.exercises[m.selected].Title())
		desc := m.exercises[m.selected].Description()
		editor := m.textarea.View()
		output := m.renderOutputPanel()

		debugPanel := ""
		if m.debugMode {
			debugPanelHeight := 5
			n := len(m.debugLog)
			start := 0
			if n > debugPanelHeight {
				start = n - debugPanelHeight
			}
			logs := m.debugLog[start:]
			debugPanel = lipgloss.NewStyle().
				Border(lipgloss.NormalBorder()).
				Padding(0, 1).
				MarginTop(1).
				Width(m.width).
				Faint(true).
				Foreground(lipgloss.Color("8")).
				Render(strings.Join(logs, "\n"))
		}

		if m.debugMode {
			m.appendDebug(fmt.Sprintf("header:\n%s\ndesc:\n%s\n", header, desc))
		}

		return fmt.Sprintf("%s\n%s\n\n%s\n\n%s\n%s\n\n[esc] Back | [F5] Run",
			header, desc, editor, output, debugPanel,
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

	debug := flag.Bool("debug", false, "enable debug mode")
	flag.Parse()

	state, err := internal.LoadState()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Warning: Could not load previous state:", err)
	}
	m := initialModel(state)
	m.debugMode = *debug
	m.debugLog = append(m.debugLog, "Debug panel is working!")

	p := tea.NewProgram(m, tea.WithAltScreen())
	go func() {
		p.Send(setProgramMsg{program: p})
	}()
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
