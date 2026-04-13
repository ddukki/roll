package main

import (
  "fmt"
  "strings"

  "charm.land/lipgloss/v2"
  "github.com/muesli/reflow/wordwrap"
)

func RenderHistory(history []rollResult, width int) string {
  if len(history) == 0 {
    return ""
  }

  wrapWidth := width - 4
  if wrapWidth <= 0 {
    wrapWidth = 60
  }

  var lines []string

  for i := len(history) - 1; i >= 0; i-- {
    r := history[i]
    lines = append(lines, r.timestamp.Format("2006-01-02 15:04:05"))

    colorIdx := 0
    prevStyle := ResultStyle
    prevTotal := ""

    // Subtotals
    for j := 0; j < len(r.nested); j++ {
      n := r.nested[j]
      currentTotal := fmt.Sprintf("%d", n.Total)

      var countStyle = ResultStyle
      if j > 0 {
        countStyle = prevStyle
      }
      totalStyle := NestedColors[colorIdx%len(NestedColors)]

      line := formatSubtotalLine(n.Expr, currentTotal, n.Dice, "Subtotal", prevTotal, countStyle, totalStyle)
      lines = appendWrappedLines(lines, line, wrapWidth)

      prevTotal = currentTotal
      prevStyle = totalStyle
      colorIdx++
    }

    // Total line
    totalLine := formatTotalLine(r.expr, r.value, r.diceValues, "Total", prevTotal, prevStyle, ResultStyle)
    lines = appendWrappedLines(lines, totalLine, wrapWidth)

    if i > 0 {
      lines = append(lines, "")
    }
  }

  result := strings.Join(lines, "\n")
  return BaseTableStyle.Render(result)
}

func appendWrappedLines(lines []string, text string, width int) []string {
  wrapped := wordwrap.String(text, width)
  parts := strings.Split(wrapped, "\n")
  for k, part := range parts {
    if k > 0 {
      part = "  " + part
    }
    lines = append(lines, part)
  }
  return lines
}

func formatSubtotalLine(expr, total string, dice []int, label, countVal string, countStyle, totalStyle interface{}) string {
  diceStr := formatDiceList(dice)
  sideStr := extractSidesFromExpr(expr)
  countColored := renderStyle(countStyle, countVal)
  totalColored := renderStyle(totalStyle, total)
  return fmt.Sprintf("%s (%s%s): [%s] = %s", label, countColored, sideStr, diceStr, totalColored)
}

func formatTotalLine(expr, total string, dice []int, label, countVal string, countStyle, totalStyle interface{}) string {
  diceStr := formatDiceList(dice)
  sideStr := extractSidesFromExpr(expr)
  countColored := renderStyle(countStyle, countVal)
  totalColored := renderStyle(totalStyle, total)
  return fmt.Sprintf("%s (%s%s): [%s] = %s", label, countColored, sideStr, diceStr, totalColored)
}

func renderStyle(style interface{}, text string) string {
  if s, ok := style.(lipgloss.Style); ok {
    return s.Render(text)
  }
  return text
}

func formatDiceList(dice []int) string {
  var parts []string
  for _, d := range dice {
    parts = append(parts, fmt.Sprintf("%d", d))
  }
  return strings.Join(parts, ", ")
}

func extractSidesFromExpr(expr string) string {
  parts := strings.Split(expr, "d")
  if len(parts) > 1 {
    return "d" + parts[1]
  }
  return ""
}
