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
	er.renderExpr(&lines, er.e, make([]bool, 0))
	return strings.Join(lines, "\n")
}

// renderEvalInfo renders out the roll details of a dice roll.
func (er *exprRenderer) renderEvalInfo(prefix, connector string, ei *evalInfo, countColor, valueColor lipgloss.Style) string {
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

func (er *exprRenderer) renderRollExpr(bexpr *binaryExpr, lines *[]string, prefix, connector string, isLastChildStack []bool) {
	*lines = append(*lines, er.renderEvalInfo(prefix, connector, bexpr.evalInfo, er.popColor(), er.getPreviousColor()))

	if isRollExpr(bexpr.lhs) {
		er.renderExpr(lines, bexpr.lhs, append(isLastChildStack, true))
	}
}

func (er *exprRenderer) renderMathExpr(bexpr *binaryExpr, lines *[]string, prefix, connector string, isLastChildStack []bool) {
	opSymbol := fmt.Sprintf("%c", bexpr.op.Rune())
	*lines = append(*lines, fmt.Sprintf("%s%s= %d", prefix, connector, bexpr.evalInfo.value))

	if bexpr.lhs != nil {
		er.renderExpr(lines, bexpr.lhs, append(isLastChildStack, false))
	}

	var opConnector string
	if connector != "" {
		if len(isLastChildStack) > 1 {
			opConnector = "│   "
		} else if len(isLastChildStack) > 0 {
			opConnector = "    "
		}
	}

	*lines = append(*lines, fmt.Sprintf("%s%s%s", prefix, opConnector, opSymbol))

	if bexpr.rhs != nil {
		er.renderExpr(lines, bexpr.rhs, append(isLastChildStack, true))
	}
}

func renderLiteralExpr(lexpr *litValExpr, lines *[]string, prefix, connector string) {
	*lines = append(*lines, fmt.Sprintf("%s%s%d", prefix, connector, lexpr.val))
}

func (er *exprRenderer) renderExpr(lines *[]string, e expr, isLastChildStack []bool) {
	if len(isLastChildStack) == 0 {
		er.renderExpr(lines, e, append(isLastChildStack, true))
		return
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
			er.renderRollExpr(v, lines, prefix, connector, isLastChildStack)
		default:
			er.renderMathExpr(v, lines, prefix, connector, isLastChildStack)
		}
	case *litValExpr:
		renderLiteralExpr(v, lines, prefix, connector)
	}
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
