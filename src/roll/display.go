package main

import (
	"fmt"
	"strconv"
	"strings"

	"charm.land/lipgloss/v2"
)

func RenderResult(exprStr string, e expr) string {
	var lines []string
	lines = append(lines, exprStr)
	lines = append(lines, fmt.Sprintf("= %d", e.Value()))
	renderExpr(&lines, e, make([]bool, 0))
	return strings.Join(lines, "\n")
}

// renderEvalInfo renders out the roll details of a dice roll.
func renderEvalInfo(prefix, connector string, ei *evalInfo, countColor, valueColor lipgloss.Style) string {
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

func renderRollExpr(bexpr *binaryExpr, lines *[]string, prefix, connector string, isLastChildStack []bool) {
	var countColor, valueColor lipgloss.Style
	*lines = append(*lines, renderEvalInfo(prefix, connector, bexpr.evalInfo, countColor, valueColor))

	hasLhsRoll := isRollExpr(bexpr.lhs)
	hasRhsRoll := isRollExpr(bexpr.rhs)

	if hasLhsRoll || hasRhsRoll {
		newStack := append(isLastChildStack, !hasRhsRoll)
		if hasLhsRoll {
			renderExpr(lines, bexpr.lhs, newStack)
		}
		if hasRhsRoll {
			renderExpr(lines, bexpr.rhs, newStack)
		}
	}
}

func renderMathExpr(bexpr *binaryExpr, lines *[]string, prefix, connector string, isLastStack []bool) {
	opSymbol := fmt.Sprintf("%c", bexpr.op.Rune())
	if len(isLastStack) > 1 {
		*lines = append(*lines, fmt.Sprintf("%s%s= %d", prefix, connector, bexpr.evalInfo.value))
	}

	if bexpr.lhs != nil {
		renderExpr(lines, bexpr.lhs, append(isLastStack, false))
	}

	var opConnector string
	if connector != "" && len(isLastStack) > 1 {
		opConnector = "│   "
	}

	*lines = append(*lines, fmt.Sprintf("%s%s%s", prefix, opConnector, opSymbol))

	if bexpr.rhs != nil {
		renderExpr(lines, bexpr.rhs, append(isLastStack, true))
	}
}

func renderLiteralExpr(lexpr *litValExpr, lines *[]string, prefix, connector string) {
	*lines = append(*lines, fmt.Sprintf("%s%s%d", prefix, connector, lexpr.val))
}

func renderExpr(lines *[]string, e expr, isLastChildStack []bool) {
	if len(isLastChildStack) == 0 {
		renderExpr(lines, e, append(isLastChildStack, false))
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
			renderRollExpr(v, lines, prefix, connector, isLastChildStack)
		default:
			renderMathExpr(v, lines, prefix, connector, isLastChildStack)
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
	if len(isLastStack) < 2 {
		return ""
	}
	var sb strings.Builder
	for i := 1; i < len(isLastStack)-1; i++ {
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
		lines = append(lines, r.timestamp.Format("2006-01-02 15:04:05"))
		lines = append(lines, RenderResult(r.expr, r.ast))
		if i > 0 {
			lines = append(lines, "")
		}
	}

	return BaseTableStyle.Render(strings.Join(lines, "\n"))
}
