package main

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
)

type exprRenderer struct {
	colorCursor int
	raw         string
	e           expr
}

func (er *exprRenderer) getPreviousColor() lipgloss.Style {
	return NestedColors[er.colorCursor-1]
}

func (er *exprRenderer) popColor() lipgloss.Style {
	er.colorCursor = (er.colorCursor + 1) % len(NestedColors)
	return NestedColors[er.colorCursor]
}

func (er *exprRenderer) renderResult() string {
	var lines []string
	lines = append(lines, fmt.Sprintf("Roll Expression = %s", er.raw))
	lines = append(lines, fmt.Sprintf("Final Value = %d", er.e.Value()))
	lines = append(lines, er.renderExpr(er.e, make([]bool, 0))...)
	return strings.Join(lines, "\n")
}

// renderEvalInfo renders out the roll details of a dice roll.
func (er *exprRenderer) renderEvalInfo(ei *evalInfo, prefix, connector string, countColor, valueColor lipgloss.Style) string {
	return fmt.Sprintf(
		"%s%s%sd%d: [%s] = %s",
		prefix,
		connector,
		countColor.Render(strconv.Itoa(len(ei.dri.rolls))),
		ei.dri.dieSize,
		formatDiceList(ei.dri.rolls),
		valueColor.Render(strconv.Itoa(ei.value)),
	)
}

func (er *exprRenderer) renderRollExpr(bexpr *binaryExpr, prefix, connector string, isLastChildStack []bool) []string {
	lines := []string{er.renderEvalInfo(bexpr.evalInfo, prefix, connector, er.popColor(), er.getPreviousColor())}

	if isRollExpr(bexpr.lhs) {
		lines = append(lines, er.renderExpr(bexpr.lhs, append(isLastChildStack, true))...)
	}

	return lines
}

func (er *exprRenderer) renderMathExpr(bexpr *binaryExpr, prefix, connector string, isLastChildStack []bool) []string {
	opSymbol := fmt.Sprintf("%c", bexpr.op.Rune())
	lines := []string{fmt.Sprintf("%s%s= %d", prefix, connector, bexpr.evalInfo.value)}

	if bexpr.lhs != nil {
		lines = append(lines, er.renderExpr(bexpr.lhs, append(isLastChildStack, false))...)
	}

	var opConnector string
	if connector != "" {
		if len(isLastChildStack) > 1 {
			opConnector = "│   "
		} else if len(isLastChildStack) > 0 {
			opConnector = "    "
		}
	}

	lines = append(lines, fmt.Sprintf("%s%s%s", prefix, opConnector, opSymbol))

	if bexpr.rhs != nil {
		lines = append(lines, er.renderExpr(bexpr.rhs, append(isLastChildStack, true))...)
	}

	return lines
}

func renderLiteralExpr(lexpr *litValExpr, prefix, connector string) []string {
	return []string{fmt.Sprintf("%s%s%d", prefix, connector, lexpr.val)}
}

func (er *exprRenderer) renderExpr(e expr, isLastChildStack []bool) []string {
	if len(isLastChildStack) == 0 {
		return er.renderExpr(e, append(isLastChildStack, true))
	}

	depth := len(isLastChildStack) - 1

	prefix := buildPrefix(isLastChildStack)
	var connector string
	if isLastChildStack[depth] {
		connector = "└── "
	} else {
		connector = "├── "
	}

	switch v := e.(type) {
	case *binaryExpr:
		switch v.op.(type) {
		case *RollOp, *MaxOp, *MinOp, *DropHighestOp, *DropLowestOp:
			return er.renderRollExpr(v, prefix, connector, isLastChildStack)
		default:
			return er.renderMathExpr(v, prefix, connector, isLastChildStack)
		}
	case *litValExpr:
		return renderLiteralExpr(v, prefix, connector)
	}

	return nil
}

// isRollExpr returns true iff the expr involves a dice roll.
func isRollExpr(e expr) bool {
	if be, ok := e.(*binaryExpr); ok {
		switch be.op.(type) {
		case *RollOp, *MaxOp, *MinOp, *DropHighestOp, *DropLowestOp:
			return true
		}
	}
	return false
}

func buildPrefix(isLastStack []bool) string {
	if len(isLastStack) == 0 {
		return ""
	}
	var sb strings.Builder
	for i := 0; i < len(isLastStack)-1; i++ {
		if isLastStack[i] {
			sb.WriteString("    ")
		} else {
			sb.WriteString("│   ")
		}
	}
	return sb.String()
}

func formatDiceList(dice []int) string {
	var parts []string
	for _, d := range dice {
		parts = append(parts, strconv.Itoa(d))
	}
	return strings.Join(parts, ", ")
}

func RenderHistory(history []rollResult, width int) string {
	if len(history) == 0 {
		return ""
	}

	var lines []string

	for i := len(history) - 1; i >= 0; i-- {
		r := history[i]

		er := &exprRenderer{raw: r.expr, e: r.ast}
		lines = append(lines, r.timestamp.Format("2006-01-02 15:04:05"))
		lines = append(lines, er.renderResult())
		if i > 0 {
			lines = append(lines, "")
		}
	}

	return BaseTableStyle.Width(width).Render(strings.Join(lines, "\n"))
}
