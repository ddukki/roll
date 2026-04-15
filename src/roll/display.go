package main

import (
	"fmt"
	"strings"
)

func RenderResult(exprStr string, e expr) string {
	var lines []string
	lines = append(lines, exprStr)
	lines = append(lines, fmt.Sprintf("= %d", e.Value()))
	renderExpr(&lines, e, []bool{true}, true)
	return strings.Join(lines, "\n")
}

func renderExpr(lines *[]string, e expr, isLastStack []bool, isRoot bool) {
	if isRoot {
		renderExprBody(lines, e, isLastStack)
		return
	}

	depth := len(isLastStack) - 1

	prefix := buildPrefix(isLastStack)
	var connector string
	if isLastStack[depth] {
		connector = "└── "
	} else {
		connector = "├── "
	}

	renderExprBodyWithPrefix(lines, e, isLastStack, prefix, connector)
}

func renderExprBody(lines *[]string, e expr, isLastStack []bool) {
	switch v := e.(type) {
	case *binaryExpr:
		switch v.op.(type) {
		case *RollOp, *MaxOp, *MinOp, *DropHighestOp, *DropLowestOp:
			dri := v.evalInfo.dri
			diceStr := formatDiceList(dri.rolls)
			count := len(dri.rolls)
			sides := dri.dieSize
			*lines = append(*lines, fmt.Sprintf("%dd%d: [%s] = %d", count, sides, diceStr, v.evalInfo.value))

			hasLhsRoll := hasRollChild(v.lhs)
			hasRhsRoll := hasRollChild(v.rhs)

			if hasLhsRoll || hasRhsRoll {
				newStack := append(isLastStack, !hasRhsRoll)
				if hasLhsRoll {
					renderExpr(lines, v.lhs, newStack, false)
				}
				if hasRhsRoll {
					renderExpr(lines, v.rhs, newStack, false)
				}
			}
		default:
			newStack := append(isLastStack, false)
			if v.lhs != nil {
				renderChild(lines, v.lhs, newStack)
			}
			if v.rhs != nil {
				newStack[len(newStack)-1] = true
				renderChild(lines, v.rhs, newStack)
			}
		}
	}
}

func renderExprBodyWithPrefix(lines *[]string, e expr, isLastStack []bool, prefix string, connector string) {
	switch v := e.(type) {
	case *binaryExpr:
		switch v.op.(type) {
		case *RollOp, *MaxOp, *MinOp, *DropHighestOp, *DropLowestOp:
			dri := v.evalInfo.dri
			diceStr := formatDiceList(dri.rolls)
			count := len(dri.rolls)
			sides := dri.dieSize
			*lines = append(*lines, prefix+connector+fmt.Sprintf("%dd%d: [%s] = %d", count, sides, diceStr, v.evalInfo.value))

			hasLhsRoll := hasRollChild(v.lhs)
			hasRhsRoll := hasRollChild(v.rhs)

			if hasLhsRoll || hasRhsRoll {
				newStack := append(isLastStack, !hasRhsRoll)
				if hasLhsRoll {
					renderExpr(lines, v.lhs, newStack, false)
				}
				if hasRhsRoll {
					renderExpr(lines, v.rhs, newStack, false)
				}
			}
		default:
			*lines = append(*lines, prefix+connector+fmt.Sprintf("= %d", v.evalInfo.value))

			newStack := append(isLastStack, false)
			if v.lhs != nil {
				renderChild(lines, v.lhs, newStack)
			}
			if v.rhs != nil {
				newStack[len(newStack)-1] = true
				renderChild(lines, v.rhs, newStack)
			}
		}
	}
}

func renderChild(lines *[]string, e expr, isLastStack []bool) {
	depth := len(isLastStack) - 1

	prefix := buildPrefix(isLastStack)
	connector := "└── "
	if !isLastStack[depth] {
		connector = "├── "
	}

	switch v := e.(type) {
	case *binaryExpr:
		switch v.op.(type) {
		case *RollOp, *MaxOp, *MinOp, *DropHighestOp, *DropLowestOp:
			dri := v.evalInfo.dri
			diceStr := formatDiceList(dri.rolls)
			count := len(dri.rolls)
			sides := dri.dieSize
			*lines = append(*lines, prefix+connector+fmt.Sprintf("%dd%d: [%s] = %d", count, sides, diceStr, v.evalInfo.value))

			hasLhsRoll := hasRollChild(v.lhs)
			hasRhsRoll := hasRollChild(v.rhs)

			if hasLhsRoll || hasRhsRoll {
				newStack := append(isLastStack, !hasRhsRoll)
				if hasLhsRoll {
					renderExpr(lines, v.lhs, newStack, false)
				}
				if hasRhsRoll {
					renderExpr(lines, v.rhs, newStack, false)
				}
			}
		default:
			*lines = append(*lines, prefix+connector+fmt.Sprintf("= %d", v.evalInfo.value))

			newStack := append(isLastStack, false)
			if v.lhs != nil {
				renderChild(lines, v.lhs, newStack)
			}
			if v.rhs != nil {
				newStack[len(newStack)-1] = true
				renderChild(lines, v.rhs, newStack)
			}
		}
	case *litValExpr:
		*lines = append(*lines, prefix+connector+fmt.Sprintf("%d", v.val))
	}
}

func hasRollChild(e expr) bool {
	if be, ok := e.(*binaryExpr); ok {
		switch be.op.(type) {
		case *RollOp, *MaxOp, *MinOp, *DropHighestOp, *DropLowestOp:
			return true
		}
	}
	return false
}

func buildPrefix(isLastStack []bool) string {
	if len(isLastStack) <= 1 {
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
		parts = append(parts, fmt.Sprintf("%d", d))
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
