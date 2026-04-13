package main

import (
  "fmt"
  "math/rand"
)

type rollStep struct {
  expr  *binaryExpr // The AST node for this roll (e.g., "3d6" or "Nd8")
  sides int         // Number of sides on the die (the "d" operand)
}

// buildRollChain walks a nested dice expression from outermost to innermost.
// For "3d6d8", the AST looks like: (3d6) d8 where 3d6 is the left side of the outer d8.
// The chain collects each roll operation in order: [3d6, 11d8].
func buildRollChain(outer *binaryExpr) []rollStep {
  var chain []rollStep
  cur := outer

  for {
    sides := extractSides(cur)
    chain = append(chain, rollStep{expr: cur, sides: sides})

    // Check if the left-hand side is another roll operation (e.g., "3d6" inside "3d6d8")
    if subRoll, ok := cur.lhs.(*binaryExpr); ok {
      if _, isRollOp := subRoll.op.(*RollOp); isRollOp {
        cur = subRoll
        continue
      }
    }
    break
  }
  return chain
}

// extractSides extracts the die faces (sides) from a roll expression.
// Handles both "3d6" (rhs is literal 6) and "XdYdZ" (rhs is another binary expr).
func extractSides(expr *binaryExpr) int {
  sides := 0
  if rhsLit, ok := expr.rhs.(*litValExpr); ok {
    sides = rhsLit.val
  } else if rhsBin, ok := expr.rhs.(*binaryExpr); ok {
    if subRhsLit, ok := rhsBin.rhs.(*litValExpr); ok {
      sides = subRhsLit.val
    }
  }
  if sides == 0 {
    sides = 6
  }
  return sides
}

// evaluateRollChain executes the chain from innermost to outermost.
// For "3d6d8": first roll 3d6 (get sum 11), then roll 11d8 (use previous sum as count).
func evaluateRollChain(chain []rollStep) *RollDetails {
  var finalRoll *RollDetails

  // Process from innermost (last in chain) to outermost (first in chain)
  for i := len(chain) - 1; i >= 0; i-- {
    step := chain[i]
    sides := step.sides
    count := 0

    // Innermost roll: use the literal count from the AST (e.g., "3" from "3d6")
    // Outer rolls: use the total from the previous roll (e.g., sum of 3d6)
    if i == len(chain)-1 {
      count = step.expr.lhs.Value()
    } else {
      count = finalRoll.Total
    }

    rolls := make([]int, count)
    var val int
    for j := range rolls {
      rolls[j] = 1 + rand.Intn(sides)
      val += rolls[j]
    }

    curRoll := &RollDetails{
      Dice:  rolls,
      Sides: sides,
      Total: val,
      Op:    "sum",
      Expr:  fmt.Sprintf("%dd%d", count, sides),
    }
    if finalRoll != nil {
      curRoll.Nested = append(finalRoll.Nested, finalRoll)
    }
    rollStack = append(rollStack, curRoll)
    finalRoll = curRoll
  }

  return finalRoll
}
