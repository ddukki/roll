package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTokenize_SimpleLiteral(t *testing.T) {
	expr, err := tokenize("42")
	require.NoError(t, err)
	require.IsType(t, &litValExpr{}, expr)
	assert.Equal(t, 42, expr.(*litValExpr).val)
}

func TestTokenize_SimpleAdd(t *testing.T) {
	expr, err := tokenize("2+3")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	be := expr.(*binaryExpr)
	assert.Equal(t, '+', be.op.Rune())
	require.IsType(t, &litValExpr{}, be.lhs)
	require.IsType(t, &litValExpr{}, be.rhs)
	assert.Equal(t, 2, be.lhs.(*litValExpr).val)
	assert.Equal(t, 3, be.rhs.(*litValExpr).val)
}

func TestTokenize_SimpleSub(t *testing.T) {
	expr, err := tokenize("5-3")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	be := expr.(*binaryExpr)
	assert.Equal(t, '-', be.op.Rune())
	assert.Equal(t, 5, be.lhs.(*litValExpr).val)
	assert.Equal(t, 3, be.rhs.(*litValExpr).val)
}

func TestTokenize_SimpleMul(t *testing.T) {
	expr, err := tokenize("2*3")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	be := expr.(*binaryExpr)
	assert.Equal(t, '*', be.op.Rune())
	assert.Equal(t, 2, be.lhs.(*litValExpr).val)
	assert.Equal(t, 3, be.rhs.(*litValExpr).val)
}

func TestTokenize_SimpleDiv(t *testing.T) {
	expr, err := tokenize("6/2")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	be := expr.(*binaryExpr)
	assert.Equal(t, '/', be.op.Rune())
	assert.Equal(t, 6, be.lhs.(*litValExpr).val)
	assert.Equal(t, 2, be.rhs.(*litValExpr).val)
}

func TestTokenize_SimpleRoll(t *testing.T) {
	expr, err := tokenize("2d6")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	be := expr.(*binaryExpr)
	assert.Equal(t, 'd', be.op.Rune())
	assert.Equal(t, 2, be.lhs.(*litValExpr).val)
	assert.Equal(t, 6, be.rhs.(*litValExpr).val)
}

func TestTokenize_DefaultCount(t *testing.T) {
	expr, err := tokenize("d6")
	require.NoError(t, err)

	be := expr.(*binaryExpr)
	assert.Equal(t, 'd', be.op.Rune())
	require.IsType(t, &litValExpr{}, be.lhs)
	assert.Equal(t, 1, be.lhs.(*litValExpr).val)
	require.IsType(t, &litValExpr{}, be.rhs)
	assert.Equal(t, 6, be.rhs.(*litValExpr).val)
}

func TestTokenize_DefaultSides_Error(t *testing.T) {
	_, err := tokenize("2d")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "undefined number of sides")
}

func TestTokenize_RollWithModifier(t *testing.T) {
	expr, err := tokenize("2d6+3")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	be := expr.(*binaryExpr)
	assert.Equal(t, '+', be.op.Rune())
	require.IsType(t, &binaryExpr{}, be.lhs)
	assert.Equal(t, 3, be.rhs.(*litValExpr).val)

	inner := be.lhs.(*binaryExpr)
	assert.Equal(t, 'd', inner.op.Rune())
	assert.Equal(t, 2, inner.lhs.(*litValExpr).val)
	assert.Equal(t, 6, inner.rhs.(*litValExpr).val)
}

func TestTokenize_NestedRoll(t *testing.T) {
	expr, err := tokenize("3d6d8")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	be, ok := expr.(*binaryExpr)
	require.True(t, ok)
	assert.Equal(t, 'd', be.op.Rune())
	inner, ok := be.lhs.(*binaryExpr)
	require.True(t, ok)
	assert.Equal(t, 'd', inner.op.Rune())
	assert.Equal(t, 3, inner.lhs.Value())
	assert.Equal(t, 6, inner.rhs.Value())
	outerRhs, ok := be.rhs.(*litValExpr)
	require.True(t, ok)
	assert.Equal(t, 8, outerRhs.val)
}

func TestTokenize_MultipleAdds(t *testing.T) {
	expr, err := tokenize("1+2+3")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	be := expr.(*binaryExpr)
	assert.Equal(t, '+', be.op.Rune())
	assert.Equal(t, "1+2+3", expressionToString(expr))
}

func TestTokenize_PEMDAS_MulBeforeAdd(t *testing.T) {
	expr, err := tokenize("1+2*3")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	be := expr.(*binaryExpr)
	assert.Equal(t, '+', be.op.Rune())

	assert.Equal(t, 1, be.lhs.(*litValExpr).val)
	require.IsType(t, &binaryExpr{}, be.rhs)
	mul := be.rhs.(*binaryExpr)
	assert.Equal(t, '*', mul.op.Rune())
	assert.Equal(t, 2, mul.lhs.(*litValExpr).val)
	assert.Equal(t, 3, mul.rhs.(*litValExpr).val)
}

func TestTokenize_PEMDAS_ParenLike(t *testing.T) {
	expr, err := tokenize("2d6*2")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	be := expr.(*binaryExpr)
	assert.Equal(t, '*', be.op.Rune())

	require.IsType(t, &binaryExpr{}, be.lhs)
	dice := be.lhs.(*binaryExpr)
	assert.Equal(t, 'd', dice.op.Rune())
	assert.Equal(t, 2, be.rhs.(*litValExpr).val)
}

func TestTokenize_MaxOp(t *testing.T) {
	_, err := tokenize("4x")
	require.Error(t, err)

	expr, err := tokenize("3x6")
	require.NoError(t, err)
	be, ok := expr.(*binaryExpr)
	require.True(t, ok)
	lhs, ok := be.lhs.(*litValExpr)
	require.True(t, ok)
	rhs, ok := be.rhs.(*litValExpr)
	require.True(t, ok)
	assert.Equal(t, 3, lhs.Value())
	assert.Equal(t, 6, rhs.Value())
	assert.Equal(t, &MaxOp{}, be.op)
}

func TestTokenize_MinOp(t *testing.T) {
	_, err := tokenize("4n")
	require.Error(t, err)

	expr, err := tokenize("3n6")
	require.NoError(t, err)
	be, ok := expr.(*binaryExpr)
	require.True(t, ok)
	lhs, ok := be.lhs.(*litValExpr)
	require.True(t, ok)
	rhs, ok := be.rhs.(*litValExpr)
	require.True(t, ok)
	assert.Equal(t, 3, lhs.Value())
	assert.Equal(t, 6, rhs.Value())
	assert.Equal(t, &MinOp{}, be.op)
}

func TestTokenize_DropHighest(t *testing.T) {
	_, err := tokenize("3d6h")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "undefined number of sides")
}

func TestTokenize_DropHighest_WithRhs(t *testing.T) {
	expr, err := tokenize("3d6h1")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	assert.Equal(t, "3d6h1", expressionToString(expr))
}

func TestTokenize_DropLowest(t *testing.T) {
	_, err := tokenize("3d6l")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "undefined number of sides")
}

func TestTokenize_DropLowest_WithRhs(t *testing.T) {
	expr, err := tokenize("3d6l1")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	assert.Equal(t, "3d6l1", expressionToString(expr))
}

func TestTokenize_Complex_NestedWithModifier(t *testing.T) {
	expr, err := tokenize("3d6d8+2")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	be := expr.(*binaryExpr)
	assert.Equal(t, '+', be.op.Rune())

	require.IsType(t, &binaryExpr{}, be.lhs)
	outer := be.lhs.(*binaryExpr)
	assert.Equal(t, 'd', outer.op.Rune())
	assert.Equal(t, 2, be.rhs.(*litValExpr).val)
}

func TestTokenize_Complex_RollChain(t *testing.T) {
	expr, err := tokenize("3d6d8d10")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	be := expr.(*binaryExpr)
	assert.Equal(t, 'd', be.op.Rune())
}

func TestTokenize_Deep_Nesting(t *testing.T) {
	expr, err := tokenize("1d20+1d6")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	be := expr.(*binaryExpr)
	assert.Equal(t, '+', be.op.Rune())
}

func TestTokenize_PEMDAS_DivBeforeAdd(t *testing.T) {
	expr, err := tokenize("1+6/2")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	be := expr.(*binaryExpr)
	assert.Equal(t, '+', be.op.Rune())
	assert.Equal(t, 1, be.lhs.(*litValExpr).val)
	require.IsType(t, &binaryExpr{}, be.rhs)
	div := be.rhs.(*binaryExpr)
	assert.Equal(t, '/', div.op.Rune())
}

func TestTokenize_ExpressionToString(t *testing.T) {
	tests := []string{
		"2+3",
		"2d6",
		"2d6+3",
		"3d6d8",
		"1+2*3",
		"2d6*2",
	}

	for _, s := range tests {
		t.Run(s, func(t *testing.T) {
			expr, err := tokenize(s)
			require.NoError(t, err)
			assert.Equal(t, s, expressionToString(expr))
		})
	}
}

func TestTokenize_ExpressionToString_WithDefaults(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"d6", "1d6"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			expr, err := tokenize(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, expressionToString(expr))
		})
	}
}

func TestTokenize_Value(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{1, 1}})
	defer ResetGlobalRand()

	tests := []struct {
		input    string
		expected int
	}{
		{"2+3", 5},
		{"2d6", 4},
		{"2d6+3", 7},
		{"1+2*3", 7},
		{"2d6*2", 8},
		{"d6", 2},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			expr, err := tokenize(tt.input)
			require.NoError(t, err)
			assert.Equal(t, tt.expected, expr.Value())
		})
	}
}

func TestTokenize_NestedRoll_Value(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{0, 0, 0, 0, 0}})
	defer ResetGlobalRand()

	expr, err := tokenize("3d6d8")
	require.NoError(t, err)

	val := expr.Value()
	assert.GreaterOrEqual(t, val, 3)
	assert.LessOrEqual(t, val, 24)
}

func TestTokenize_EmptyLhs_Error(t *testing.T) {
	_, err := tokenize("+3")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot have an empty lhs")
}

func TestTokenize_EmptyRhs_Error(t *testing.T) {
	_, err := tokenize("2+")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "cannot have an empty rhs")
}

func TestTokenize_InvalidLiteral(t *testing.T) {
	_, err := tokenize("abc")
	require.Error(t, err)
	assert.Contains(t, err.Error(), "literal values must be integers")
}

func TestTokenize_MultipleAdds_Structure(t *testing.T) {
	expr, err := tokenize("1+2+3")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	be := expr.(*binaryExpr)
	assert.Equal(t, '+', be.op.Rune())

	require.IsType(t, &binaryExpr{}, be.lhs)
	left := be.lhs.(*binaryExpr)
	assert.Equal(t, '+', left.op.Rune())
	assert.Equal(t, 1, left.lhs.(*litValExpr).val)
	assert.Equal(t, 2, left.rhs.(*litValExpr).val)

	require.IsType(t, &litValExpr{}, be.rhs)
	assert.Equal(t, 3, be.rhs.(*litValExpr).val)
}

func TestTokenize_ChainedRolls_Structure(t *testing.T) {
	expr, err := tokenize("1d6+2d8+3")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	outer := expr.(*binaryExpr)
	assert.Equal(t, '+', outer.op.Rune())

	require.IsType(t, &binaryExpr{}, outer.lhs)
	left := outer.lhs.(*binaryExpr)
	assert.Equal(t, '+', left.op.Rune())
	require.IsType(t, &binaryExpr{}, left.lhs)
	require.IsType(t, &binaryExpr{}, left.rhs)

	require.IsType(t, &litValExpr{}, outer.rhs)
	assert.Equal(t, 3, outer.rhs.(*litValExpr).val)
}

func TestTokenize_MaxOp_Structure(t *testing.T) {
	expr, err := tokenize("3x6")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	be := expr.(*binaryExpr)
	require.IsType(t, &MaxOp{}, be.op)
	assert.Equal(t, 'x', be.op.Rune())
	require.IsType(t, &litValExpr{}, be.lhs)
	require.IsType(t, &litValExpr{}, be.rhs)
	assert.Equal(t, 3, be.lhs.(*litValExpr).val)
	assert.Equal(t, 6, be.rhs.(*litValExpr).val)
}

func TestTokenize_MinOp_Structure(t *testing.T) {
	expr, err := tokenize("3n6")
	require.NoError(t, err)
	require.IsType(t, &binaryExpr{}, expr)
	be := expr.(*binaryExpr)
	require.IsType(t, &MinOp{}, be.op)
	assert.Equal(t, 'n', be.op.Rune())
	require.IsType(t, &litValExpr{}, be.lhs)
	require.IsType(t, &litValExpr{}, be.rhs)
	assert.Equal(t, 3, be.lhs.(*litValExpr).val)
	assert.Equal(t, 6, be.rhs.(*litValExpr).val)
}
