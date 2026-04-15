package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testRand struct {
	values []int
	idx    int
}

func (r *testRand) Int63() int64 {
	return int64(r.values[r.idx])
}

func (r *testRand) Intn(n int) int {
	v := r.values[r.idx%len(r.values)] % n
	r.idx++
	return v
}

func TestAddOp(t *testing.T) {
	op := &AddOp{}

	res := op.Apply(2, 3)

	require.NotNil(t, res)
	assert.Equal(t, 5, res.value)
	assert.IsType(t, &AddOp{}, res.op)
}

func TestSubOp(t *testing.T) {
	op := &SubOp{}

	res := op.Apply(5, 3)

	require.NotNil(t, res)
	assert.Equal(t, 2, res.value)
	assert.IsType(t, &SubOp{}, res.op)
}

func TestMulOp(t *testing.T) {
	op := &MulOp{}

	res := op.Apply(2, 3)

	require.NotNil(t, res)
	assert.Equal(t, 6, res.value)
	assert.IsType(t, &MulOp{}, res.op)
}

func TestDivOp(t *testing.T) {
	op := &DivOp{}

	res := op.Apply(6, 2)

	require.NotNil(t, res)
	assert.Equal(t, 3, res.value)
	assert.IsType(t, &DivOp{}, res.op)
}

func TestRollOp(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{0, 0, 0}}) // 1, 1, 1
	defer ResetGlobalRand()

	op := &RollOp{}
	res := op.Apply(3, 6)

	require.NotNil(t, res)
	assert.Equal(t, 3, res.value)
	require.NotNil(t, res.dri)
	assert.Len(t, res.dri.rolls, 3)
}

func TestRollOp_SingleDie(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{0}})
	defer ResetGlobalRand()

	op := &RollOp{}

	res := op.Apply(1, 6)

	require.NotNil(t, res)
	assert.Equal(t, 1, res.value)
	assert.Equal(t, []int{1}, res.dri.rolls)
}

func TestRollOp_ZeroDice(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{0}})
	defer ResetGlobalRand()

	op := &RollOp{}

	res := op.Apply(0, 6)

	require.NotNil(t, res)
	assert.Equal(t, 0, res.value)
	assert.Empty(t, res.dri.rolls)
}

func TestMaxOp(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{5, 4, 3}}) // 6, 5, 4
	defer ResetGlobalRand()

	op := &MaxOp{}
	res := op.Apply(3, 6)

	require.NotNil(t, res)
	assert.Equal(t, 6, res.value)
	require.NotNil(t, res.dri)
	assert.Len(t, res.dri.rolls, 3)
}

func TestMinOp(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{0, 1, 2, 3, 4, 5}})
	defer ResetGlobalRand()

	op := &MinOp{}

	res := op.Apply(3, 6)

	require.NotNil(t, res)
	assert.Equal(t, 1, res.value)
	require.NotNil(t, res.dri)
	assert.Len(t, res.dri.rolls, 3)
}

func TestDropHighestOp(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{0, 1, 2}})
	defer ResetGlobalRand()

	op := &DropHighestOp{}

	res := op.Apply(3, 6)

	require.NotNil(t, res)
	assert.Equal(t, 3, res.value)
	require.NotNil(t, res.dri)
	assert.Equal(t, []int{1, 2, 3}, res.dri.rolls)
}

func TestDropHighestOp_TwoDice(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{0, 1}})
	defer ResetGlobalRand()

	op := &DropHighestOp{}

	res := op.Apply(2, 6)

	require.NotNil(t, res)
	assert.Equal(t, 1, res.value)
	assert.Equal(t, []int{1, 2}, res.dri.rolls)
}

func TestDropHighestOp_Panic_OneDie(t *testing.T) {
	op := &DropHighestOp{}

	assert.Panics(t, func() {
		op.Apply(1, 6)
	})
}

func TestDropLowestOp(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{0, 1, 2}})
	defer ResetGlobalRand()

	op := &DropLowestOp{}

	res := op.Apply(3, 6)

	require.NotNil(t, res)
	assert.Equal(t, 5, res.value)
	require.NotNil(t, res.dri)
	assert.Equal(t, []int{1, 2, 3}, res.dri.rolls)
}

func TestDropLowestOp_TwoDice(t *testing.T) {
	SetGlobalRand(&testRand{values: []int{0, 1}})
	defer ResetGlobalRand()

	op := &DropLowestOp{}

	res := op.Apply(2, 6)

	require.NotNil(t, res)
	assert.Equal(t, 2, res.value)
	assert.Equal(t, []int{1, 2}, res.dri.rolls)
}

func TestDropLowestOp_Panic_OneDie(t *testing.T) {
	op := &DropLowestOp{}

	assert.Panics(t, func() {
		op.Apply(1, 6)
	})
}

func TestOperatorRunes(t *testing.T) {
	tests := []struct {
		op op
		rn rune
	}{
		{&AddOp{}, '+'},
		{&SubOp{}, '-'},
		{&MulOp{}, '*'},
		{&DivOp{}, '/'},
		{&RollOp{}, 'd'},
		{&MaxOp{}, 'x'},
		{&MinOp{}, 'n'},
		{&DropHighestOp{}, 'h'},
		{&DropLowestOp{}, 'l'},
	}

	for _, tt := range tests {
		t.Run(string(tt.rn), func(t *testing.T) {
			assert.Equal(t, tt.rn, tt.op.Rune())
		})
	}
}

func TestPemdasPrecedence(t *testing.T) {
	assert.Equal(t, []string{"+-", "*/", "dxnhl"}, pemdas)
}

func TestOpsMap(t *testing.T) {
	tests := []struct {
		rn   rune
		want string
	}{
		{'+', "AddOp"},
		{'-', "SubOp"},
		{'*', "MulOp"},
		{'/', "DivOp"},
		{'d', "RollOp"},
		{'x', "MaxOp"},
		{'n', "MinOp"},
		{'h', "DropHighestOp"},
		{'l', "DropLowestOp"},
	}

	for _, tt := range tests {
		t.Run(tt.want, func(t *testing.T) {
			op := ops[tt.rn]
			require.NotNil(t, op)
			assert.Equal(t, tt.rn, op.Rune())
		})
	}
}

func TestSetGlobalRand(t *testing.T) {
	fr := &testRand{values: []int{1}}

	SetGlobalRand(fr)
	assert.Same(t, fr, globalRand)

	ResetGlobalRand()
	assert.NotSame(t, fr, globalRand)
}
