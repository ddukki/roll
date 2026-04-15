package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestBinaryExpr_Value_Add(t *testing.T) {
	expr := &binaryExpr{
		lhs: &litValExpr{val: 2},
		rhs: &litValExpr{val: 3},
		op:  &AddOp{},
	}

	assert.Equal(t, 5, expr.Value())
}

func TestBinaryExpr_Value_Sub(t *testing.T) {
	expr := &binaryExpr{
		lhs: &litValExpr{val: 5},
		rhs: &litValExpr{val: 3},
		op:  &SubOp{},
	}

	assert.Equal(t, 2, expr.Value())
}

func TestBinaryExpr_Value_Mul(t *testing.T) {
	expr := &binaryExpr{
		lhs: &litValExpr{val: 2},
		rhs: &litValExpr{val: 3},
		op:  &MulOp{},
	}

	assert.Equal(t, 6, expr.Value())
}

func TestBinaryExpr_Value_Div(t *testing.T) {
	expr := &binaryExpr{
		lhs: &litValExpr{val: 6},
		rhs: &litValExpr{val: 2},
		op:  &DivOp{},
	}

	assert.Equal(t, 3, expr.Value())
}

func TestBinaryExpr_Value_Roll(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{0, 1}})
	defer ResetGlobalRand()

	expr := &binaryExpr{
		lhs: &litValExpr{val: 2},
		rhs: &litValExpr{val: 6},
		op:  &RollOp{},
	}

	assert.Equal(t, 3, expr.Value())
}

func TestBinaryExpr_Value_Caches(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{0, 1}})
	defer ResetGlobalRand()

	expr := &binaryExpr{
		lhs: &litValExpr{val: 2},
		rhs: &litValExpr{val: 6},
		op:  &RollOp{},
	}

	first := expr.Value()
	second := expr.Value()

	assert.Equal(t, first, second)
	assert.NotNil(t, expr.evalInfo)
}

func TestLitValExpr_Value(t *testing.T) {
	expr := &litValExpr{val: 42}

	assert.Equal(t, 42, expr.Value())
}

func TestTokenExpr_Value(t *testing.T) {
	expr := &tokenExpr{token: "2d6"}

	assert.Equal(t, 0, expr.Value())
}

func TestExpressionToString_BinaryExpr(t *testing.T) {
	expr := &binaryExpr{
		lhs: &litValExpr{val: 2},
		rhs: &litValExpr{val: 6},
		op:  &RollOp{},
	}

	assert.Equal(t, "2d6", expressionToString(expr))
}

func TestExpressionToString_Nested(t *testing.T) {
	expr := &binaryExpr{
		lhs: &binaryExpr{
			lhs: &litValExpr{val: 2},
			rhs: &litValExpr{val: 6},
			op:  &RollOp{},
		},
		rhs: &litValExpr{val: 3},
		op:  &AddOp{},
	}

	assert.Equal(t, "2d6+3", expressionToString(expr))
}

func TestExpressionToString_LitValExpr(t *testing.T) {
	expr := &litValExpr{val: 42}

	assert.Equal(t, "42", expressionToString(expr))
}

func TestExpressionToString_TokenExpr(t *testing.T) {
	expr := &tokenExpr{token: "2d6"}

	assert.Equal(t, "2d6", expressionToString(expr))
}

func TestGetLitValue_LitValExpr(t *testing.T) {
	expr := &litValExpr{val: 5}

	assert.Equal(t, 5, getLitValue(expr))
}

func TestGetLitValue_BinaryExpr(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{0, 1}})
	defer ResetGlobalRand()

	expr := &binaryExpr{
		lhs: &litValExpr{val: 2},
		rhs: &litValExpr{val: 6},
		op:  &RollOp{},
	}

	assert.Equal(t, 3, getLitValue(expr))
}
