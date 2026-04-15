package main

// import (
// 	"fmt"
// 	"strings"

// 	"charm.land/lipgloss/v2"
// 	"github.com/muesli/reflow/wordwrap"
// )

// // ===== HISTORY RENDERING =====
// // These functions render roll history for display in the TUI.
// // History shows all previous rolls with their results.

// // RenderHistory renders a list of roll results into a styled string.
// // Each roll is displayed with its subtotals and total.
// // Parameters:
// //   - history: Previous roll results
// //   - width:   Terminal width for wrapping
// //
// // Returns: Styled multi-line string for display
// func RenderHistory(history []rollResult, width int) string {
// 	if len(history) == 0 {
// 		return ""
// 	}

// 	// Calculate wrap width (leave room for borders)
// 	wrapWidth := width - 4
// 	if wrapWidth <= 0 {
// 		wrapWidth = 60
// 	}

// 	var lines []string

// 	// Process each roll result from oldest to newest
// 	for i := len(history) - 1; i >= 0; i-- {
// 		r := history[i]
// 		// Add timestamp
// 		lines = append(lines, r.timestamp.Format("2006-01-02 15:04:05"))

// 		// Track colors for alternating subtotals
// 		colorIdx := 0
// 		prevStyle := ResultStyle
// 		prevTotal := ""

// 		// Render each subtotal (nested roll results)
// 		// This shows intermediate steps like 3d6 in "3d6d8"
// 		for j := 0; j < len(r.nested); j++ {
// 			n := r.nested[j]
// 			currentTotal := fmt.Sprintf("%d", n.Total)

// 			// Alternate colors for nested results
// 			var countStyle = ResultStyle
// 			if j > 0 {
// 				countStyle = prevStyle
// 			}
// 			totalStyle := NestedColors[colorIdx%len(NestedColors)]

// 			// Format and render this subtotal line
// 			line := formatSubtotalLine(n.Expr, currentTotal, n.Dice, "Subtotal", prevTotal, countStyle, totalStyle)
// 			lines = appendWrappedLines(lines, line, wrapWidth)

// 			prevTotal = currentTotal
// 			prevStyle = totalStyle
// 			colorIdx++
// 		}

// 		// Render the final total line
// 		totalLine := formatTotalLine(r.expr, r.value, r.diceValues, "Total", prevTotal, prevStyle, ResultStyle)
// 		lines = appendWrappedLines(lines, totalLine, wrapWidth)

// 		// Add blank line between rolls (except first)
// 		if i > 0 {
// 			lines = append(lines, "")
// 		}
// 	}

// 	// Join all lines and apply table styling
// 	result := strings.Join(lines, "\n")
// 	return BaseTableStyle.Render(result)
// }

// // appendWrappedLines wraps text at specified width and adds to lines.
// // Handles multi-line wrapped text by indenting continuation lines.
// func appendWrappedLines(lines []string, text string, width int) []string {
// 	wrapped := wordwrap.String(text, width)
// 	parts := strings.Split(wrapped, "\n")
// 	for k, part := range parts {
// 		// Indent continuation lines
// 		if k > 0 {
// 			part = "  " + part
// 		}
// 		lines = append(lines, part)
// 	}
// 	return lines
// }

// // formatSubtotalLine formats a single subtotal line for display.
// // Shows the expression, dice values, and subtotal.
// // Parameters:
// //   - expr:       The expression (e.g., "3d6")
// //   - total:     The total value
// //   - dice:      Individual dice values
// //   - label:     Label to display ("Subtotal")
// //   - countVal:  Previous total for chaining
// //   - countStyle: Style for count display
// //   - totalStyle: Style for total display
// func formatSubtotalLine(expr, total string, dice []int, label, countVal string, countStyle, totalStyle interface{}) string {
// 	diceStr := formatDiceList(dice)
// 	sideStr := extractSidesFromExpr(expr)
// 	countColored := renderStyle(countStyle, countVal)
// 	totalColored := renderStyle(totalStyle, total)

// 	// If no sides or has math operators, show simpler format
// 	if sideStr == "" || strings.Contains(expr, "+") || strings.Contains(expr, "-") {
// 		return fmt.Sprintf("%s %s: [%s] = %s", label, expr, diceStr, totalColored)
// 	}
// 	// Otherwise show dice expression
// 	return fmt.Sprintf("%s (%s%s): [%s] = %s", label, countColored, sideStr, diceStr, totalColored)
// }

// // formatTotalLine formats the final total line.
// // Similar to formatSubtotalLine but uses "Total" label.
// func formatTotalLine(expr, total string, dice []int, label, countVal string, countStyle, totalStyle interface{}) string {
// 	diceStr := formatDiceList(dice)
// 	sideStr := extractSidesFromExpr(expr)
// 	countColored := renderStyle(countStyle, countVal)
// 	totalColored := renderStyle(totalStyle, total)

// 	// If no sides or has math operators, show simpler format
// 	if sideStr == "" || strings.Contains(expr, "+") || strings.Contains(expr, "-") {
// 		return fmt.Sprintf("%s %s: [%s] = %s", label, expr, diceStr, totalColored)
// 	}
// 	// Otherwise show dice expression
// 	return fmt.Sprintf("%s (%s%s): [%s] = %s", label, countColored, sideStr, diceStr, totalColored)
// }

// // renderStyle applies a lipgloss style to text.
// // If the style is not available, returns plain text.
// func renderStyle(style interface{}, text string) string {
// 	if s, ok := style.(lipgloss.Style); ok {
// 		return s.Render(text)
// 	}
// 	return text
// }

// // formatDiceList formats dice values as a comma-separated string.
// // For example: [3, 1, 6] -> "3, 1, 6"
// func formatDiceList(dice []int) string {
// 	var parts []string
// 	for _, d := range dice {
// 		parts = append(parts, fmt.Sprintf("%d", d))
// 	}
// 	return strings.Join(parts, ", ")
// }

// // extractSidesFromExpr extracts "dN" from an expression.
// // For example: "3d6" -> "d6", "2d8" -> "d8"
// func extractSidesFromExpr(expr string) string {
// 	parts := strings.Split(expr, "d")
// 	if len(parts) > 1 {
// 		return "d" + parts[1]
// 	}
// 	return ""
// }
