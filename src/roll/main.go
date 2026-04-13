package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
	"github.com/muesli/reflow/wordwrap"
)

var focusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
var inputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
var resultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

var nestedColors = []lipgloss.Style{
	lipgloss.NewStyle().Foreground(lipgloss.Color("205")), // pink
	lipgloss.NewStyle().Foreground(lipgloss.Color("75")),  // cyan
	lipgloss.NewStyle().Foreground(lipgloss.Color("197")), // magenta
	lipgloss.NewStyle().Foreground(lipgloss.Color("33")),  // blue
	lipgloss.NewStyle().Foreground(lipgloss.Color("130")), // orange
}

var baseTableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

const maxHistory = 20

type rollResult struct {
	expr       string
	value      string
	diceValues []int
	nested     []*RollDetails
	timestamp  time.Time
}

type model struct {
	input    string
	spinner  spinner.Model
	loading  bool
	result   string
	err      error
	history  []rollResult
	width    int
	quitting bool
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	return model{spinner: s, width: 80}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyPressMsg:
		switch msg.String() {
		case "ctrl+c", "q", "esc":
			m.quitting = true
			return m, tea.Quit
		case "enter":
			if !m.loading && m.input != "" {
				m.loading = true
				m.result = ""
				m.err = nil
				return m, tea.Batch(m.spinner.Tick, runRoll(m.input))
			}
		case "backspace":
			if len(m.input) > 0 {
				m.input = m.input[:len(m.input)-1]
			}
		default:
			if len(msg.String()) == 1 {
				m.input += msg.String()
			}
			return m, nil
		}
	case doneMsg:
		m.loading = false
		m.result = msg.value
		if msg.err != nil {
			m.err = msg.err
		} else {
			m.history = append(m.history, rollResult{
				expr:       msg.expr,
				value:      msg.value,
				diceValues: msg.dice,
				nested:     msg.nested,
				timestamp:  msg.timestamp,
			})
			if len(m.history) > maxHistory {
				m.history = m.history[len(m.history)-maxHistory:]
			}
		}
		return m, nil
	case tea.WindowSizeMsg:
		m.width = msg.Width
		return m, nil
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		return m, cmd
	}
	return m, nil
}

func (m model) View() tea.View {
	v := tea.NewView("")

	if m.quitting {
		return v
	}

	var s string
	s += "\n"
	s += focusStyle.Render("  Roll expression: ")
	s += inputStyle.Render(m.input)
	s += "\n\n"

	if m.loading {
		s += "  "
		s += m.spinner.View()
		s += focusStyle.Render("  Rolling...")
		s += "\n\n"
	}

	if m.err != nil {
		s += "  "
		s += errorStyle.Render(m.err.Error())
		s += "\n"
	}

	if m.result != "" && !m.loading {
		s += "  "
		s += resultStyle.Render("Result: " + m.result)
		s += "\n"
	}

	s += "\n"
	s += renderHistory(m.history, m.width)
	s += "\n\n"

	s += "\n  "
	s += helpStyle.Render("Enter: roll • Backspace: delete • q: quit")
	s += "\n"

	v = tea.NewView(s)
	v.AltScreen = true

	return v
}

func renderHistory(history []rollResult, width int) string {
	if len(history) == 0 {
		return ""
	}

	wrapWidth := width - 4
	if wrapWidth <= 0 {
		wrapWidth = 60
	}

	var wrappedLines []string

	for i := len(history) - 1; i >= 0; i-- {
		r := history[i]

		ts := r.timestamp.Format("2006-01-02 15:04:05")
		wrappedLines = append(wrappedLines, ts)

		colorIdx := 0
		prevColor := resultStyle // for first line, use resultStyle as "previous"
		prevTotal := ""

		// Subtotals
		for j := 0; j < len(r.nested); j++ {
			n := r.nested[j]
			currentTotal := fmt.Sprintf("%d", n.Total)

			// For first subtotal: no previous color, use plain. For others: use prevColor
			var countStyle lipgloss.Style
			if j == 0 {
				countStyle = lipgloss.NewStyle() // plain, no color
			} else {
				countStyle = prevColor
			}
			colorForTotal := nestedColors[colorIdx%len(nestedColors)]

			line := formatRollLineTwoColors(n.Expr, currentTotal, n.Dice, "Subtotal", prevTotal, countStyle, colorForTotal)
			wrapped := wordwrap.String(line, wrapWidth)
			parts := strings.Split(wrapped, "\n")
			for k, part := range parts {
				if k > 0 {
					part = "  " + part
				}
				wrappedLines = append(wrappedLines, part)
			}

			prevTotal = currentTotal
			prevColor = colorForTotal
			colorIdx++
		}

		// Total line - use prevColor for count, resultStyle for final total
		totalLine := formatRollLineTwoColors(r.expr, r.value, r.diceValues, "Total", prevTotal, prevColor, resultStyle)
		wrapped := wordwrap.String(totalLine, wrapWidth)
		parts := strings.Split(wrapped, "\n")
		for k, part := range parts {
			if k > 0 {
				// Wrapped continuation lines: no extra styling
				wrappedLines = append(wrappedLines, "  "+part)
			} else {
				wrappedLines = append(wrappedLines, part)
			}
		}

		if i > 0 {
			wrappedLines = append(wrappedLines, "")
		}
	}

	result := lipgloss.JoinVertical(
		lipgloss.Left,
		wrappedLines...,
	)

	return baseTableStyle.Render(result)
}

func formatRollLineTwoColors(expr, total string, dice []int, label, countVal string, countStyle, totalStyle lipgloss.Style) string {
	var diceParts []string
	for _, d := range dice {
		diceParts = append(diceParts, fmt.Sprintf("%d", d))
	}
	diceStr := strings.Join(diceParts, ", ")

	// Use countVal if provided, otherwise parse from expression
	countStr := countVal
	if countStr == "" {
		exprParts := strings.Split(expr, "d")
		countStr = exprParts[0]
	}

	// Parse expression for sides
	exprParts := strings.Split(expr, "d")
	sideStr := ""
	if len(exprParts) > 1 {
		sideStr = "d" + exprParts[1]
	}

	countColored := countStyle.Render(countStr)
	totalColored := totalStyle.Render(total)

	return fmt.Sprintf("%s (%s%s): [%s] = %s", label, countColored, sideStr, diceStr, totalColored)
}

type doneMsg struct {
	expr      string
	value     string
	dice      []int
	nested    []*RollDetails
	timestamp time.Time
	err       error
}

func runRoll(expr string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)

		var err error
		var dice []int
		var nested []*RollDetails
		var val int
		var resultExpr string

		defer func() {
			if r := recover(); r != nil {
				err = fmt.Errorf("%v", r)
			}
		}()

		tkn := tokenize(expr)
		ast := parse(tkn)
		ast.Validate()

		rollDetails := ast.Roll()
		if rollDetails != nil {
			val = rollDetails.Total
			dice = rollDetails.Dice
			resultExpr = rollDetails.Expr
			nested = rollDetails.Nested
			rollStack = append(rollStack, rollDetails)
		} else {
			val = ast.Value()
		}

		displayExpr := resultExpr
		if displayExpr == "" {
			displayExpr = expr
		}

		return doneMsg{expr: displayExpr, value: fmt.Sprintf("%d", val), dice: dice, nested: nested, timestamp: time.Now(), err: err}
	}
}

func main() {
	flag.Parse()
	if flag.NArg() > 0 {
		expr := strings.Join(flag.Args(), " ")

		tkn := tokenize(expr)
		ast := parse(tkn)
		ast.Validate()

		rollDetails := ast.Roll()
		val := 0
		if rollDetails != nil {
			val = rollDetails.Total
		} else {
			val = ast.Value()
		}

		var display string
		if rollDetails != nil {
			lines := []string{}
			for i := 0; i < len(rollDetails.Nested); i++ {
				n := rollDetails.Nested[i]
				lines = append(lines, fmt.Sprintf("Subtotal %s: [%v] = %d", n.Expr, n.Dice, n.Total))
			}
			lines = append(lines, fmt.Sprintf("Total %s: [%v] = %d", rollDetails.Expr, rollDetails.Dice, rollDetails.Total))
			display = strings.Join(lines, "\n")
			fmt.Println(display)
		} else {
			display = expr + " = " + fmt.Sprintf("%d", val)
			fmt.Println(display)
		}
		return
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
