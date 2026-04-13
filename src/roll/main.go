package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"charm.land/bubbles/v2/spinner"
	"charm.land/bubbles/v2/table"
	tea "charm.land/bubbletea/v2"
	"charm.land/lipgloss/v2"
)

var focusStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
var inputStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
var resultStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("86")).Bold(true)
var errorStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("196"))
var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("245"))

var baseTableStyle = lipgloss.NewStyle().
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("240"))

const maxHistory = 20

type rollResult struct {
	expr       string
	value      string
	diceValues []int
	nested     []string
}

type model struct {
	input    string
	spinner  spinner.Model
	loading  bool
	result   string
	err      error
	history  []rollResult
	table    table.Model
	quitting bool
}

func initialModel() model {
	s := spinner.New()
	s.Spinner = spinner.Dot
	s.Style = lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	columns := []table.Column{
		{Title: "Expression", Width: 50},
		{Title: "Total", Width: 8},
	}
	tbl := table.New(
		table.WithColumns(columns),
		table.WithHeight(12),
		table.WithWidth(60),
	)
	styles := table.DefaultStyles()
	styles.Header = styles.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	tbl.SetStyles(styles)

	return model{spinner: s, table: tbl}
}

func (m model) Init() tea.Cmd {
	return tea.ClearScreen
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
			})
			if len(m.history) > maxHistory {
				m.history = m.history[len(m.history)-maxHistory:]
			}
			m.table.SetRows(buildTableRows(m.history))
		}
		return m, nil
	default:
		var cmd tea.Cmd
		m.spinner, cmd = m.spinner.Update(msg)
		m.table, _ = m.table.Update(msg)
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
	s += baseTableStyle.Render(m.table.View()) + "\n"

	s += "\n  "
	s += helpStyle.Render("Enter: roll • Backspace: delete • q: quit")
	s += "\n"

	v = tea.NewView(s)
	v.AltScreen = true

	return v
}

func buildTableRows(history []rollResult) []table.Row {
	rows := make([]table.Row, 0, len(history))
	for i := len(history) - 1; i >= 0; i-- {
		r := history[i]
		display := formatRollResult(r)
		rows = append(rows, table.Row{display, r.value})
	}
	return rows
}

func formatDice(dice []int) string {
	if len(dice) == 0 {
		return ""
	}
	s := ""
	for i, d := range dice {
		if i > 0 {
			s += ", "
		}
		s += fmt.Sprintf("%d", d)
	}
	return s
}

func formatRollResult(r rollResult) string {
	var parts []string
	parts = append(parts, formatDice(r.diceValues))
	for _, n := range r.nested {
		parts = append(parts, n)
	}
	return r.expr + " = " + strings.Join(parts, " <- ")
}

type doneMsg struct {
	expr   string
	value  string
	dice   []int
	nested []string
	err    error
}

func runRoll(expr string) tea.Cmd {
	return func() tea.Msg {
		time.Sleep(2 * time.Second)

		var err error
		var dice []int
		var nested []string
		var val int
		var resultExpr string

		func() {
			defer func() {
				if r := recover(); r != nil {
					err = fmt.Errorf("%v", r)
					return
				}
			}()

			tkn := tokenize(expr)
			ast := parse(tkn)
			ast.Validate()
			val = ast.Value()

			details := GetRollDetails()
			if details != nil {
				dice = details.Dice
				resultExpr = details.Expr
				for _, n := range details.Nested {
					nested = append(nested, fmt.Sprintf("%s:%v", n.Expr, n.Dice))
				}
			}
		}()

		displayExpr := resultExpr
		if displayExpr == "" {
			displayExpr = expr
		}

		return doneMsg{expr: displayExpr, value: fmt.Sprintf("%d", val), dice: dice, nested: nested, err: err}
	}
}

func main() {
	flag.Parse()
	if flag.NArg() > 0 {
		expr := strings.Join(flag.Args(), " ")

		tkn := tokenize(expr)
		ast := parse(tkn)
		ast.Validate()
		val := ast.Value()

		details := GetRollDetails()
		var display string
		if details != nil {
			parts := []string{fmt.Sprintf("%v", details.Dice)}
			for _, n := range details.Nested {
				parts = append(parts, fmt.Sprintf("%s:%v", n.Expr, n.Dice))
			}
			display = details.Expr + " = " + strings.Join(parts, " <- ")
		} else {
			display = expr + " = error"
		}

		fmt.Printf("%-40s %s\n", display, fmt.Sprintf("%d", val))
		return
	}

	p := tea.NewProgram(initialModel())
	if _, err := p.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
