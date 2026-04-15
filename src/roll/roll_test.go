package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type rollExpectation struct {
	dieSize int
	rolls   []int
}

type rollTestCase struct {
	expr     string
	randVals []int
	expected int
	rolls    []rollExpectation
}

func collectDiceRolls(e expr, acc *[]*diceRollInfo) {
	if be, ok := e.(*binaryExpr); ok {
		collectDiceRolls(be.lhs, acc)
		collectDiceRolls(be.rhs, acc)
		if be.evalInfo != nil && be.evalInfo.dri != nil {
			*acc = append(*acc, be.evalInfo.dri)
		}
	}
}

func TestRoll(t *testing.T) {
	cases := []rollTestCase{
		{
			expr:     "2d6",
			randVals: []int{0, 1},
			expected: 3,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2}},
			},
		},
		{
			expr:     "d6",
			randVals: []int{2},
			expected: 3,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{3}},
			},
		},
		{
			expr:     "2d6+3",
			randVals: []int{0, 1},
			expected: 6,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2}},
			},
		},
		{
			expr:     "3+2d6",
			randVals: []int{0, 1},
			expected: 6,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2}},
			},
		},
		{
			expr:     "2d6+2d6",
			randVals: []int{0, 1, 2, 3},
			expected: 10,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2}},
				{dieSize: 6, rolls: []int{3, 4}},
			},
		},
		{
			expr:     "2d6*2",
			randVals: []int{0, 1},
			expected: 6,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2}},
			},
		},
		{
			expr:     "2d6/2",
			randVals: []int{0, 1},
			expected: 1,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2}},
			},
		},
		{
			expr:     "2d6-1",
			randVals: []int{0, 1},
			expected: 2,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2}},
			},
		},
		{
			expr:     "1+2+3",
			randVals: []int{},
			expected: 6,
			rolls:    []rollExpectation{},
		},
		{
			expr:     "2d6+2d6+2d6",
			randVals: []int{0, 1, 2, 3, 4, 5},
			expected: 21,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2}},
				{dieSize: 6, rolls: []int{3, 4}},
				{dieSize: 6, rolls: []int{5, 6}},
			},
		},
		{
			expr:     "2d6+3*2",
			randVals: []int{0, 1},
			expected: 9,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2}},
			},
		},
		{
			expr:     "3d6",
			randVals: []int{0, 1, 2},
			expected: 6,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2, 3}},
			},
		},
		{
			expr:     "4d8",
			randVals: []int{0, 1, 2, 3},
			expected: 10,
			rolls: []rollExpectation{
				{dieSize: 8, rolls: []int{1, 2, 3, 4}},
			},
		},
		{
			expr:     "1d20+1d6",
			randVals: []int{9, 0},
			expected: 11,
			rolls: []rollExpectation{
				{dieSize: 20, rolls: []int{10}},
				{dieSize: 6, rolls: []int{1}},
			},
		},
		{
			expr:     "1d6+1d6+1d6+1d6",
			randVals: []int{0, 1, 2, 3},
			expected: 10,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1}},
				{dieSize: 6, rolls: []int{2}},
				{dieSize: 6, rolls: []int{3}},
				{dieSize: 6, rolls: []int{4}},
			},
		},
		{
			expr:     "5d10-5",
			randVals: []int{0, 1, 2, 3, 4},
			expected: 10,
			rolls: []rollExpectation{
				{dieSize: 10, rolls: []int{1, 2, 3, 4, 5}},
			},
		},
		{
			expr:     "2d6+2d8+2d10",
			randVals: []int{0, 1, 0, 1, 0, 1},
			expected: 9,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2}},
				{dieSize: 8, rolls: []int{1, 2}},
				{dieSize: 10, rolls: []int{1, 2}},
			},
		},
		{
			expr:     "d100",
			randVals: []int{49},
			expected: 50,
			rolls: []rollExpectation{
				{dieSize: 100, rolls: []int{50}},
			},
		},
		{
			expr:     "10d6",
			randVals: []int{0, 1, 2, 3, 4, 5, 0, 1, 2, 3},
			expected: 31,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2, 3, 4, 5, 6, 1, 2, 3, 4}},
			},
		},
		{
			expr:     "2d6*3+1",
			randVals: []int{0, 1},
			expected: 10,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2}},
			},
		},
		{
			expr:     "1+2d6*3",
			randVals: []int{0, 1},
			expected: 10,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2}},
			},
		},
		{
			expr:     "2*2d6+3",
			randVals: []int{0, 1},
			expected: 9,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2}},
			},
		},
		{
			expr:     "2d6+2d6*2",
			randVals: []int{0, 1, 2, 3},
			expected: 17,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2}},
				{dieSize: 6, rolls: []int{3, 4}},
			},
		},
		{
			expr:     "2d6*2+2d6",
			randVals: []int{0, 1, 2, 3},
			expected: 13,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2}},
				{dieSize: 6, rolls: []int{3, 4}},
			},
		},
		{
			expr:     "100d6",
			randVals: []int{0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3, 4, 5, 0, 1, 2, 3},
			expected: 346,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4, 5, 6, 1, 2, 3, 4}},
			},
		},
		{
			expr:     "2d20",
			randVals: []int{9, 9},
			expected: 20,
			rolls: []rollExpectation{
				{dieSize: 20, rolls: []int{10, 10}},
			},
		},
		{
			expr:     "3d6+3d6",
			randVals: []int{0, 1, 2, 3, 4, 5},
			expected: 21,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2, 3}},
				{dieSize: 6, rolls: []int{4, 5, 6}},
			},
		},
		{
			expr:     "4d6+2",
			randVals: []int{0, 1, 2, 3},
			expected: 12,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2, 3, 4}},
			},
		},
		{
			expr:     "5d4-3",
			randVals: []int{0, 1, 2, 3, 0},
			expected: 8,
			rolls: []rollExpectation{
				{dieSize: 4, rolls: []int{1, 2, 3, 4, 1}},
			},
		},
		{
			expr:     "6d8/2",
			randVals: []int{0, 1, 2, 3, 4, 5},
			expected: 10,
			rolls: []rollExpectation{
				{dieSize: 8, rolls: []int{1, 2, 3, 4, 5, 6}},
			},
		},
		{
			expr:     "7d12*2",
			randVals: []int{0, 1, 2, 3, 4, 5, 6},
			expected: 56,
			rolls: []rollExpectation{
				{dieSize: 12, rolls: []int{1, 2, 3, 4, 5, 6, 7}},
			},
		},
		{
			expr:     "8d4+8",
			randVals: []int{0, 1, 2, 3, 0, 1, 2, 3},
			expected: 28,
			rolls: []rollExpectation{
				{dieSize: 4, rolls: []int{1, 2, 3, 4, 1, 2, 3, 4}},
			},
		},
		{
			expr:     "9d6-9",
			randVals: []int{0, 1, 2, 3, 4, 5, 0, 1, 2},
			expected: 18,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2, 3, 4, 5, 6, 1, 2, 3}},
			},
		},
		{
			expr:     "2d6+2d6+2d6+2d6",
			randVals: []int{0, 1, 0, 1, 0, 1, 0, 1},
			expected: 12,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2}},
				{dieSize: 6, rolls: []int{1, 2}},
				{dieSize: 6, rolls: []int{1, 2}},
				{dieSize: 6, rolls: []int{1, 2}},
			},
		},
		{
			expr:     "3d8+2d4",
			randVals: []int{0, 1, 2, 0, 1},
			expected: 9,
			rolls: []rollExpectation{
				{dieSize: 8, rolls: []int{1, 2, 3}},
				{dieSize: 4, rolls: []int{1, 2}},
			},
		},
		{
			expr:     "4d12+3d6",
			randVals: []int{0, 1, 2, 3, 0, 1, 2},
			expected: 16,
			rolls: []rollExpectation{
				{dieSize: 12, rolls: []int{1, 2, 3, 4}},
				{dieSize: 6, rolls: []int{1, 2, 3}},
			},
		},
		{
			expr:     "2d10*2+2",
			randVals: []int{0, 1},
			expected: 8,
			rolls: []rollExpectation{
				{dieSize: 10, rolls: []int{1, 2}},
			},
		},
		{
			expr:     "3*3d6",
			randVals: []int{0, 1, 2},
			expected: 18,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2, 3}},
			},
		},
		{
			expr:     "1d8+1d6+1d4+1d10",
			randVals: []int{0, 0, 0, 0},
			expected: 4,
			rolls: []rollExpectation{
				{dieSize: 8, rolls: []int{1}},
				{dieSize: 6, rolls: []int{1}},
				{dieSize: 4, rolls: []int{1}},
				{dieSize: 10, rolls: []int{1}},
			},
		},
		{
			expr:     "2d6+2d6+2d6+2d6+2d6",
			randVals: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
			expected: 31,
			rolls: []rollExpectation{
				{dieSize: 6, rolls: []int{1, 2}},
				{dieSize: 6, rolls: []int{3, 4}},
				{dieSize: 6, rolls: []int{5, 6}},
				{dieSize: 6, rolls: []int{1, 2}},
				{dieSize: 6, rolls: []int{3, 4}},
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.expr, func(t *testing.T) {
			SetGlobalRand(&testRand{values: tc.randVals})
			defer ResetGlobalRand()

			expr, err := tokenize(tc.expr)
			require.NoError(t, err)

			val := expr.Value()
			assert.Equal(t, tc.expected, val, "value mismatch")

			var collected []*diceRollInfo
			collectDiceRolls(expr, &collected)

			require.Len(t, collected, len(tc.rolls), "roll count mismatch")
			for i, exp := range tc.rolls {
				actual := collected[i]
				assert.Equal(t, exp.dieSize, actual.dieSize, "dieSize mismatch at index %d", i)
				assert.Equal(t, exp.rolls, actual.rolls, "rolls mismatch at index %d", i)
			}
		})
	}
}

func TestRoll_OrderDeterministic(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{0, 1, 2, 3}})
	defer ResetGlobalRand()

	expr, err := tokenize("2d6+2d6")
	require.NoError(t, err)

	val := expr.Value()
	assert.Equal(t, 10, val)

	var collected []*diceRollInfo
	collectDiceRolls(expr, &collected)

	require.Len(t, collected, 2)
	assert.Equal(t, []int{1, 2}, collected[0].rolls)
	assert.Equal(t, []int{3, 4}, collected[1].rolls)
}

func TestRoll_Caching(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{0, 1}})
	defer ResetGlobalRand()

	expr, err := tokenize("2d6")
	require.NoError(t, err)

	val1 := expr.Value()
	val2 := expr.Value()

	assert.Equal(t, val1, val2)

	var collected []*diceRollInfo
	collectDiceRolls(expr, &collected)
	require.Len(t, collected, 1)
	assert.Equal(t, []int{1, 2}, collected[0].rolls)
}

func TestRoll_MultipleEvaluations(t *testing.T) {
	tests := []struct {
		expr     string
		randVals []int
		expected int
	}{
		{"2d6", []int{0, 1, 2, 3, 4, 5}, 3},
		{"4d6", []int{0, 1, 2, 3, 4, 5, 0, 1}, 10},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			SetGlobalRand(&testRand{values: tt.randVals})
			defer ResetGlobalRand()

			expr, err := tokenize(tt.expr)
			require.NoError(t, err)

			for i := 0; i < 3; i++ {
				val := expr.Value()
				assert.Equal(t, tt.expected, val, "evaluation %d mismatch", i)
			}
		})
	}
}

func TestRoll_SumOfRolls(t *testing.T) {
	tests := []struct {
		expr     string
		randVals []int
		expected int
	}{
		{"3d6", []int{0, 1, 2}, 6},
		{"4d6", []int{0, 1, 2, 3}, 10},
		{"5d6", []int{0, 1, 2, 3, 4}, 15},
		{"6d6", []int{0, 1, 2, 3, 4, 5}, 21},
	}

	for _, tt := range tests {
		t.Run(tt.expr, func(t *testing.T) {
			SetGlobalRand(&testRand{values: tt.randVals})
			defer ResetGlobalRand()

			expr, err := tokenize(tt.expr)
			require.NoError(t, err)

			val := expr.Value()
			assert.Equal(t, tt.expected, val)
		})
	}
}

func TestRoll_LargeDice(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{99}})
	defer ResetGlobalRand()

	expr, err := tokenize("d100")
	require.NoError(t, err)

	val := expr.Value()
	assert.Equal(t, 100, val)

	var collected []*diceRollInfo
	collectDiceRolls(expr, &collected)
	require.Len(t, collected, 1)
	assert.Equal(t, 100, collected[0].dieSize)
	assert.Equal(t, []int{100}, collected[0].rolls)
}

func TestRoll_ZeroDice(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{0}})
	defer ResetGlobalRand()

	expr, err := tokenize("0d6")
	require.NoError(t, err)

	val := expr.Value()
	assert.Equal(t, 0, val)

	var collected []*diceRollInfo
	collectDiceRolls(expr, &collected)
	require.Len(t, collected, 1)
	assert.Empty(t, collected[0].rolls)
}

func TestRoll_NestedRolls(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		randVals []int
		expected int
	}{
		{
			name:     "2d6d8",
			expr:     "2d6d8",
			randVals: []int{0, 1, 0, 1, 2},
			expected: 6,
		},
		{
			name:     "3d6d8",
			expr:     "3d6d8",
			randVals: []int{0, 1, 2, 0, 1, 2, 3, 4, 5},
			expected: 21,
		},
		{
			name:     "1d20d4",
			expr:     "1d20d4",
			randVals: []int{0, 0},
			expected: 1,
		},
		{
			name:     "4d6d6",
			expr:     "4d6d6",
			randVals: []int{0, 1, 2, 3, 0, 1, 2, 3, 4, 5, 6, 7},
			expected: 27,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetGlobalRand(&testRand{values: tt.randVals})
			defer ResetGlobalRand()

			expr, err := tokenize(tt.expr)
			require.NoError(t, err)

			val := expr.Value()
			assert.Equal(t, tt.expected, val)
		})
	}
}

func TestRoll_NestedRollsWithMath(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		randVals []int
		expected int
	}{
		{
			name:     "2d6d8+3",
			expr:     "2d6d8+3",
			randVals: []int{0, 1, 0, 1, 2},
			expected: 9,
		},
		{
			name:     "3d6d8+2",
			expr:     "3d6d8+2",
			randVals: []int{0, 1, 2, 0, 1, 2, 3, 4, 5},
			expected: 23,
		},
		{
			name:     "2d6d8*2",
			expr:     "2d6d8*2",
			randVals: []int{0, 1, 0, 1, 2},
			expected: 12,
		},
		{
			name:     "3d6d8*3",
			expr:     "3d6d8*3",
			randVals: []int{0, 1, 2, 0, 1, 2, 3, 4, 5},
			expected: 63,
		},
		{
			name:     "2d6d8-1",
			expr:     "2d6d8-1",
			randVals: []int{0, 1, 0, 1, 2},
			expected: 5,
		},
		{
			name:     "3d6d8-5",
			expr:     "3d6d8-5",
			randVals: []int{0, 1, 2, 0, 1, 2, 3, 4, 5},
			expected: 16,
		},
		{
			name:     "2d6d8/2",
			expr:     "2d6d8/2",
			randVals: []int{0, 1, 0, 1, 2},
			expected: 3,
		},
		{
			name:     "4d6d6/2",
			expr:     "4d6d6/2",
			randVals: []int{0, 1, 2, 3, 0, 1, 2, 3, 4, 5, 6, 7},
			expected: 13,
		},
		{
			name:     "3+2d6d8",
			expr:     "3+2d6d8",
			randVals: []int{0, 1, 0, 1, 2},
			expected: 9,
		},
		{
			name:     "5*2d6d8",
			expr:     "5*2d6d8",
			randVals: []int{0, 1, 0, 1, 2},
			expected: 30,
		},
		{
			name:     "10-2d6d8",
			expr:     "10-2d6d8",
			randVals: []int{0, 1, 0, 1, 2},
			expected: 4,
		},
		{
			name:     "12/2d6d8",
			expr:     "12/2d6d8",
			randVals: []int{0, 1, 0, 1, 2},
			expected: 2,
		},
		{
			name:     "2d6d8+2d6",
			expr:     "2d6d8+2d6",
			randVals: []int{0, 1, 0, 1, 2, 0, 1},
			expected: 9,
		},
		{
			name:     "2d6+2d6d8",
			expr:     "2d6+2d6d8",
			randVals: []int{0, 1, 0, 1, 2, 3, 4},
			expected: 15,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetGlobalRand(&testRand{values: tt.randVals})
			defer ResetGlobalRand()

			expr, err := tokenize(tt.expr)
			require.NoError(t, err)

			val := expr.Value()
			assert.Equal(t, tt.expected, val)
		})
	}
}

func TestRoll_NestedRollsDetailed(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{0, 1, 2, 0, 1, 2, 3, 4, 5}})
	defer ResetGlobalRand()

	expr, err := tokenize("3d6d8")
	require.NoError(t, err)

	val := expr.Value()
	assert.Equal(t, 21, val)

	var collected []*diceRollInfo
	collectDiceRolls(expr, &collected)

	require.Len(t, collected, 2)
	assert.Equal(t, 6, collected[0].dieSize)
	assert.Equal(t, []int{1, 2, 3}, collected[0].rolls)
	assert.Equal(t, 8, collected[1].dieSize)
	assert.Equal(t, []int{1, 2, 3, 4, 5, 6}, collected[1].rolls)
}

func TestRoll_NestedRollsWithMathDetailed(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{0, 1, 0, 1, 2}})
	defer ResetGlobalRand()

	expr, err := tokenize("2d6d8+3")
	require.NoError(t, err)

	val := expr.Value()
	assert.Equal(t, 9, val)

	var collected []*diceRollInfo
	collectDiceRolls(expr, &collected)

	require.Len(t, collected, 2)
	assert.Equal(t, 6, collected[0].dieSize)
	assert.Equal(t, []int{1, 2}, collected[0].rolls)
	assert.Equal(t, 8, collected[1].dieSize)
	assert.Equal(t, []int{1, 2, 3}, collected[1].rolls)
}

func TestRoll_MultipleNestedRolls(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11}})
	defer ResetGlobalRand()

	expr, err := tokenize("2d6d8+3d6d4")
	require.NoError(t, err)

	val := expr.Value()
	assert.Equal(t, 33, val)
}

func TestRoll_DeeplyNestedRolls(t *testing.T) {
	tests := []struct {
		name     string
		expr     string
		randVals []int
		expected int
	}{
		{
			name:     "2d4d6d8",
			expr:     "2d4d6d8",
			randVals: []int{0, 1, 0, 1, 2, 0, 1, 2, 3, 4, 5, 6, 7},
			expected: 21,
		},
		{
			name:     "2d6d8+1",
			expr:     "2d6d8+1",
			randVals: []int{0, 1, 0, 1, 2},
			expected: 7,
		},
		{
			name:     "2d6d8*2+1",
			expr:     "2d6d8*2+1",
			randVals: []int{0, 1, 0, 1, 2},
			expected: 13,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SetGlobalRand(&testRand{values: tt.randVals})
			defer ResetGlobalRand()

			expr, err := tokenize(tt.expr)
			require.NoError(t, err)

			val := expr.Value()
			assert.Equal(t, tt.expected, val)
		})
	}
}
