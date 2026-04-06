package main

import (
	"fmt"
	"strings"

	"github.com/alecthomas/chroma/v2"
	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/styles"
	"github.com/charmbracelet/lipgloss"
)

// --- Syntax highlighting ---
//
// The textarea widget still handles all editing (keystrokes, cursor, buffer),
// but we replace its View() output with our own highlighted rendering.

var (
	// chromaTheme provides the color palette for syntax tokens (keywords, strings, etc.)
	chromaTheme = styles.Get("monokai")

	// tsLexer is the TypeScript tokenizer, initialized once at startup.
	// chroma.Coalesce merges adjacent tokens of the same type for efficiency.
	// IIFE to handle potential nil lexer gracefully.
	// the ts lexer should never be nil in reality, but... whatever
	tsLexer = func() chroma.Lexer {
		l := lexers.Get("typescript")
		if l == nil {
			l = lexers.Fallback
		}
		return chroma.Coalesce(l)
	}()
)

// styledSpan is a fragment of text on a single line with a single style.
// A line of code is represented as []styledSpan — one span per token/color change.
type styledSpan struct {
	text  string
	style lipgloss.Style
}

// chromaToLipgloss converts a chroma style entry (color, bold, italic)
// into a lipgloss style that the terminal can render.
func chromaToLipgloss(entry chroma.StyleEntry) lipgloss.Style {
	s := lipgloss.NewStyle()
	if entry.Colour.IsSet() {
		s = s.Foreground(lipgloss.Color(entry.Colour.String()))
	}
	if entry.Bold == chroma.Yes {
		s = s.Bold(true)
	}
	if entry.Italic == chroma.Yes {
		s = s.Italic(true)
	}
	return s
}

// highlightLines tokenizes the full code string and returns a 2D slice:
// result[lineIndex] = []styledSpan for that line.
//
// Chroma produces a flat stream of tokens — some tokens span multiple lines
// (e.g. multi-line strings). We split those on "\n" to get per-line spans.
func highlightLines(code string) [][]styledSpan {
	iter, err := tsLexer.Tokenise(nil, code)
	if err != nil {
		// Fallback: no highlighting, just plain text per line
		raw := strings.Split(code, "\n")
		result := make([][]styledSpan, len(raw))
		for i, l := range raw {
			result[i] = []styledSpan{{text: l, style: lipgloss.NewStyle()}}
		}
		return result
	}

	var result [][]styledSpan
	var cur []styledSpan // spans accumulated for the current line

	for _, tok := range iter.Tokens() {
		s := chromaToLipgloss(chromaTheme.Get(tok.Type))

		// A single token may contain newlines (e.g. template literals).
		// Split it so each piece lands on the correct output line.
		parts := strings.Split(tok.Value, "\n")
		for i, part := range parts {
			if i > 0 {
				// Newline boundary — flush current line and start a new one
				result = append(result, cur)
				cur = nil
			}
			if part != "" {
				cur = append(cur, styledSpan{text: part, style: s})
			}
		}
	}
	// Don't forget the last line (no trailing newline to trigger a flush)
	result = append(result, cur)
	return result
}

// renderStyledLine renders one line of highlighted spans into a string.
// If showCursor is true, it draws a reverse-video block cursor at cursorCol
// by splitting the span that contains the cursor position into three parts:
// before-cursor, cursor-char (reversed), after-cursor.
func renderStyledLine(spans []styledSpan, cursorCol int, showCursor bool) string {
	var b strings.Builder
	pos := 0            // character position across all spans
	cursorDone := false // have we already rendered the cursor?

	for _, sp := range spans {
		runes := []rune(sp.text)
		rLen := len(runes)

		// Check if the cursor falls inside this span
		if showCursor && !cursorDone && pos+rLen > cursorCol {
			off := cursorCol - pos // offset within this span

			// Render text before the cursor in normal style
			if off > 0 {
				b.WriteString(sp.style.Render(string(runes[:off])))
			}
			// Render the cursor character with reverse video
			b.WriteString(sp.style.Reverse(true).Render(string(runes[off : off+1])))
			// Render text after the cursor in normal style
			if off+1 < rLen {
				b.WriteString(sp.style.Render(string(runes[off+1:])))
			}
			cursorDone = true
		} else {
			b.WriteString(sp.style.Render(sp.text))
		}
		pos += rLen
	}

	// If cursor is past the end of all text (e.g. at end of line),
	// draw a floating block cursor in empty space
	if showCursor && !cursorDone {
		b.WriteString(lipgloss.NewStyle().Reverse(true).Render(" "))
	}

	return b.String()
}

// renderHighlightedCode produces the full editor panel content:
// line numbers + syntax-highlighted code + cursor, scrolled to keep
// the cursor row visible within the given viewHeight.
func (m model) renderHighlightedCode(viewHeight int) string {
	code := m.textarea.Value()
	styledLines := highlightLines(code)
	totalLines := len(styledLines)

	// Get cursor position from the textarea (it still tracks editing state)
	curRow := m.textarea.Line()
	curCol := m.textarea.LineInfo().CharOffset

	// Line number gutter width — at least 3 chars so it doesn't jump around
	numWidth := len(fmt.Sprintf("%d", totalLines))
	if numWidth < minGutterWidth {
		numWidth = minGutterWidth
	}

	var lines []string
	for i := 0; i < viewHeight; i++ {
		idx := m.start + i
		if idx < totalLines {
			isCursor := idx == curRow

			// Highlight the current line number more brightly
			nStyle := lineNumStyle
			if isCursor {
				nStyle = cursorLineNumStyle
			}
			num := nStyle.Render(fmt.Sprintf("%*d ", numWidth, idx+1))
			content := renderStyledLine(styledLines[idx], curCol, isCursor)
			lines = append(lines, num+content)
		} else {
			// Past end of file — show tilde like vim
			lines = append(lines, lineNumStyle.Render(fmt.Sprintf("%*s ", numWidth, "~")))
		}
	}

	return strings.Join(lines, "\n")
}
