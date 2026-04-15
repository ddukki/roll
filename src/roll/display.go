package main

import (
	"fmt"
	"strings"
)

func RenderResult(exprStr string, e expr) string {
	var lines []string
	lines = append(lines, exprStr)
	renderExpr(&lines, e, "", true)
	lines = append(lines, fmt.Sprintf("= %d", e.Value()))
	return strings.Join(lines, "\n")
}

func renderExpr(lines *[]string, e expr, prefix string, isLast bool) {
	childPrefix := "│   "
	if isLast {
		childPrefix = "    "
	}

	switch v := e.(type) {
	case *binaryExpr:
		switch v.op.(type) {
		case *RollOp, *MaxOp, *MinOp, *DropHighestOp, *DropLowestOp:
			dri := v.evalInfo.dri
			diceStr := formatDiceList(dri.rolls)
			opRune := v.op.Rune()
			count := len(dri.rolls)
			sides := dri.dieSize
			connector := "└── "
			if !isLast {
				connector = "├── "
			}
			*lines = append(*lines, fmt.Sprintf("%s%s%d%c%d: [%s] = %d", prefix, connector, count, opRune, sides, diceStr, v.evalInfo.value))

			newPrefix := prefix + childPrefix
			if v.lhs != nil {
				if be, ok := v.lhs.(*binaryExpr); ok {
					switch be.op.(type) {
					case *RollOp, *MaxOp, *MinOp, *DropHighestOp, *DropLowestOp:
						renderExpr(lines, be, newPrefix, false)
					}
				}
			}
			if v.rhs != nil {
				if be, ok := v.rhs.(*binaryExpr); ok {
					switch be.op.(type) {
					case *RollOp, *MaxOp, *MinOp, *DropHighestOp, *DropLowestOp:
						renderExpr(lines, be, newPrefix, true)
					}
				}
			}
		default:
			renderChild(lines, v.lhs, prefix, childPrefix, false)
			renderChild(lines, v.rhs, prefix, childPrefix, true)
		}
	}
}

func renderChild(lines *[]string, e expr, parentPrefix string, childPrefix string, isLast bool) {
	var connector string
	if isLast {
		connector = "└── "
	} else {
		connector = "├── "
	}

	switch v := e.(type) {
	case *binaryExpr:
		switch v.op.(type) {
		case *RollOp, *MaxOp, *MinOp, *DropHighestOp, *DropLowestOp:
			dri := v.evalInfo.dri
			diceStr := formatDiceList(dri.rolls)
			opRune := v.op.Rune()
			count := len(dri.rolls)
			sides := dri.dieSize
			*lines = append(*lines, fmt.Sprintf("%s%d%c%d: [%s] = %d", parentPrefix+connector, count, opRune, sides, diceStr, v.evalInfo.value))

			newPrefix := parentPrefix + childPrefix

			if v.lhs != nil {
				if be, ok := v.lhs.(*binaryExpr); ok {
					switch be.op.(type) {
					case *RollOp, *MaxOp, *MinOp, *DropHighestOp, *DropLowestOp:
						renderChild(lines, v.lhs, newPrefix, childPrefix, false)
					}
				}
			}

			if v.rhs != nil {
				if be, ok := v.rhs.(*binaryExpr); ok {
					switch be.op.(type) {
					case *RollOp, *MaxOp, *MinOp, *DropHighestOp, *DropLowestOp:
						renderChild(lines, v.rhs, newPrefix, childPrefix, true)
					}
				}
			}
		default:
			*lines = append(*lines, fmt.Sprintf("%s= %d", parentPrefix+connector, v.evalInfo.value))
			renderExpr(lines, v, parentPrefix+childPrefix, isLast)
		}
	case *litValExpr:
		*lines = append(*lines, fmt.Sprintf("%s%d", parentPrefix+connector, v.val))
	}
}

func formatDiceList(dice []int) string {
	var parts []string
	for _, d := range dice {
		parts = append(parts, fmt.Sprintf("%d", d))
	}
	return strings.Join(parts, ", ")
}
