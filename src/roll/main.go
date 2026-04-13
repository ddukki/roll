package main

import (
  "flag"
  "fmt"
  "os"
  "strings"
  "time"

  "charm.land/bubbles/v2/spinner"
  tea "charm.land/bubbletea/v2"
)

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
  s.Style = FocusStyle

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
  s += FocusStyle.Render("  Roll expression: ")
  s += InputStyle.Render(m.input)
  s += "\n\n"

  if m.loading {
    s += "  "
    s += m.spinner.View()
    s += FocusStyle.Render("  Rolling...")
    s += "\n\n"
  }

  if m.err != nil {
    s += "  "
    s += ErrorStyle.Render(m.err.Error())
    s += "\n"
  }

  if m.result != "" && !m.loading {
    s += "  "
    s += ResultStyle.Render("Result: " + m.result)
    s += "\n"
  }

  s += "\n"
  s += RenderHistory(m.history, m.width)
  s += "\n\n"

  s += "\n  "
  s += HelpStyle.Render("Enter: roll • Backspace: delete • q: quit")
  s += "\n"

  v = tea.NewView(s)
  v.AltScreen = true

  return v
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
