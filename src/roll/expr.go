package main

import (
	"fmt"
)

// ===== EXPR INTERFACE =====
// All expression types must implement this interface.
// This enables polymorphic handling in the parser and evaluator.
type expr interface {
	// Value returns the value of the expression.
	Value() int
}

// evalInfo holds all the information about what happened during expression evaluation.
type evalInfo struct {
	dri   *diceRollInfo
	op    op
	value int
}

// diceRollInfo stores information about a particular roll of dice.
type diceRollInfo struct {
	dieSize int
	rolls   []int
}

// Compile-time check that binaryExpr implements expr.
var _ expr = (*binaryExpr)(nil)

// ===== BINARY EXPRESSION =====
// binaryExpr represents a binary operation with a left side, right side, and operator.
// The AST represents mathematical expressions as trees. For example:
//
//	"2d6+3"  ->  binaryExpr(lhs=binaryExpr(op=d, lhs=2, rhs=6), op=AddOp, rhs=3)
//	"2d6*2"  ->  binaryExpr(lhs=binaryExpr(op=d, lhs=2, rhs=6), op=MulOp, rhs=2)
//
// The tree is built recursively:
//
//	root:          AddOp
//	             /     \
//	        RollOp       litValExpr(3)
//	       /     \
//	 litValExpr(2)  litValExpr(6)
type binaryExpr struct {
	lhs      expr      // Left side (e.g., "2d6" in "2d6+3")
	rhs      expr      // Right side (e.g., "3" in "2d6+3")
	op       binaryOp  // Operator (e.g., AddOp for "+")
	evalInfo *evalInfo // nil if not yet evaluated
}

// binaryExpr.Value evaluates the expression tree.
// It computes the result by recursively evaluating lhs and rhs,
// then applying the operator to combine them.
// For "2d6+3": Value() calls lhs.Value() + rhs.Value()
func (b *binaryExpr) Value() int {
	if b.evalInfo == nil {
		// Evaluate both sides and apply operator
		b.evalInfo = b.op.Apply(b.lhs.Value(), b.rhs.Value())
	}
	return b.evalInfo.value
}

// expressionToString converts an expression tree back to a string.
// This is used for caching and debugging.
// For example: binaryExpr(lhs=2d6, op=+, rhs=3) -> "2d6+3"
func expressionToString(e expr) string {
	switch v := e.(type) {
	case *binaryExpr:
		return expressionToString(v.lhs) + string(v.op.Rune()) + expressionToString(v.rhs)
	case *litValExpr:
		return fmt.Sprintf("%d", v.val)
	case *tokenExpr:
		return v.token
	default:
		return ""
	}
}

// getLitValue extracts a literal value from an expression.
// If the expression is already a literal, returns it directly.
// Otherwise, evaluates the expression (for cases like "2d6+3" in a larger expression).
func getLitValue(e expr) int {
	if lit, ok := e.(*litValExpr); ok {
		return lit.val
	}
	return e.Value()
}

// ===== LITERAL VALUE EXPRESSION =====
// litValExpr represents a numeric literal in the AST.
// This is the leaf node of the expression tree.
var _ expr = (*litValExpr)(nil)

type litValExpr struct {
	val int
}

// Value returns the stored value.
func (lv *litValExpr) Value() int { return lv.val }

// ===== TOKEN EXPRESSION =====
// tokenExpr represents an unparsed token string.
// During tokenization, this holds the raw string that will
// eventually be split into operators and operands.
var _ expr = (*tokenExpr)(nil)

type tokenExpr struct {
	token string
}

// Value() returns 0 for unparsed tokens.
func (t *tokenExpr) Value() int { return 0 }
