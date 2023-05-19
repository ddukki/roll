package main

import (
	"fmt"
	"math/rand"
)

const (
	SCE_ADD = '+'
	SCE_SUB = '-'
	SCE_MUL = '*'
	SCE_DIV = '/'

	SCE_ROLL = 'd'
	SCE_MAX  = 'x'
	SCE_MIN  = 'n'
	SCE_DHV  = 'h'
	SCE_DLV  = 'l'
)

type op interface {
	Rune() rune
}

type binaryOp interface {
	op
	Apply(lhs int, rhs int) int
}

var (
	_ binaryOp = (*AddOp)(nil)
	_ binaryOp = (*SubOp)(nil)
	_ binaryOp = (*MulOp)(nil)
	_ binaryOp = (*DivOp)(nil)
	_ binaryOp = (*RollOp)(nil)
	_ binaryOp = (*MaxOp)(nil)
	_ binaryOp = (*MinOp)(nil)
)

type AddOp struct{}
type SubOp struct{}
type MulOp struct{}
type DivOp struct{}

var (
	pemdas = []rune{
		SCE_ADD,
		SCE_SUB,
		SCE_MUL,
		SCE_DIV,
		SCE_ROLL,
		SCE_MAX,
		SCE_MIN,
		SCE_DHV,
		SCE_DLV,
	}

	ops = map[rune]op{
		SCE_ADD:  &AddOp{},
		SCE_SUB:  &SubOp{},
		SCE_DIV:  &DivOp{},
		SCE_MUL:  &MulOp{},
		SCE_ROLL: &RollOp{},
		SCE_MAX:  &MaxOp{},
		SCE_MIN:  &MinOp{},
		SCE_DHV:  &DropHighestOp{},
		SCE_DLV:  &DropLowestOp{},
	}
)

func (a *AddOp) Apply(lhs int, rhs int) int { return lhs + rhs }
func (s *SubOp) Apply(lhs int, rhs int) int { return lhs - rhs }
func (d *DivOp) Apply(lhs int, rhs int) int { return lhs / rhs }
func (m *MulOp) Apply(lhs int, rhs int) int { return lhs * rhs }

func (a *AddOp) Rune() rune { return SCE_ADD }
func (s *SubOp) Rune() rune { return SCE_SUB }
func (d *DivOp) Rune() rune { return SCE_DIV }
func (m *MulOp) Rune() rune { return SCE_MUL }

var (
	_ op = (*RollOp)(nil)
)

type RollOp struct{}
type MaxOp struct{}
type MinOp struct{}
type DropHighestOp struct{}
type DropLowestOp struct{}

func (r *RollOp) Apply(lhs, rhs int) int {
	rolls := make([]int, lhs)
	var val int
	for i := range rolls {
		rolls[i] = d(rhs)
		val += rolls[i]
	}

	rollStr := fmt.Sprintf("d%d", rhs)
	mathStr := fmt.Sprintf("%d", val)
	if lhs > 1 {
		rollStr = fmt.Sprintf("%d%s", lhs, rollStr)
		mathStr = fmt.Sprintf("sum%+v = %s", rolls, mathStr)
	} else {
		mathStr = "[" + mathStr + "]"
	}

	fmt.Printf("Rolling %s: %s\n", rollStr, mathStr)
	return val
}

func (mx *MaxOp) Apply(lhs, rhs int) int {
	rolls := make([]int, lhs)
	var val int
	for i := range rolls {
		rolls[i] = d(rhs)
		if val < rolls[i] {
			val = rolls[i]
		}
	}

	rollStr := fmt.Sprintf("%c%d", mx.Rune(), rhs)
	mathStr := fmt.Sprintf("%d", val)
	if lhs > 1 {
		rollStr = fmt.Sprintf("%d%s", lhs, rollStr)
		mathStr = fmt.Sprintf("max%+v = %s", rolls, mathStr)
	} else {
		mathStr = "[" + mathStr + "]"
	}

	fmt.Printf("Rolling %s: %s\n", rollStr, mathStr)
	return val
}

func (mn *MinOp) Apply(lhs, rhs int) int {
	rolls := make([]int, lhs)
	val := rhs + 1
	for i := range rolls {
		rolls[i] = d(rhs)
		if val > rolls[i] {
			val = rolls[i]
		}
	}

	rollStr := fmt.Sprintf("%c%d", mn.Rune(), rhs)
	mathStr := fmt.Sprintf("%d", val)
	if lhs > 1 {
		rollStr = fmt.Sprintf("%d%s", lhs, rollStr)
		mathStr = fmt.Sprintf("min%+v = %s", rolls, mathStr)
	} else {
		mathStr = "[" + mathStr + "]"
	}

	fmt.Printf("Rolling %s: %s\n", rollStr, mathStr)
	return val
}

func (dh *DropHighestOp) Apply(lhs, rhs int) int {
	rolls := make([]int, lhs)
	var max, val int
	for i := range rolls {
		rolls[i] = d(rhs)
		val += rolls[i]
		if max < rolls[i] {
			max = rolls[i]
		}
	}

	val -= max

	rollStr := fmt.Sprintf("%c%d", dh.Rune(), rhs)
	if lhs <= 1 {
		panic("there must be more than one die rolled for dropping!")
	}

	rollStr = fmt.Sprintf("%d%s", lhs, rollStr)
	mathStr := fmt.Sprintf("sum%+v - max::%d = %d", rolls, max, val)

	fmt.Printf("Rolling %s: %s\n", rollStr, mathStr)
	return val
}

func (dl *DropLowestOp) Apply(lhs, rhs int) int {
	rolls := make([]int, lhs)
	min := rhs + 1
	var val int
	for i := range rolls {
		rolls[i] = d(rhs)
		val += rolls[i]
		if min > rolls[i] {
			min = rolls[i]
		}
	}

	val -= min

	rollStr := fmt.Sprintf("%c%d", dl.Rune(), rhs)
	if lhs <= 1 {
		panic("there must be more than one die rolled for dropping!")
	}

	rollStr = fmt.Sprintf("%d%s", lhs, rollStr)
	mathStr := fmt.Sprintf("sum%+v - min::%d = %d", rolls, min, val)

	fmt.Printf("Rolling %s: %s\n", rollStr, mathStr)
	return val
}

func (m *RollOp) Rune() rune        { return SCE_ROLL }
func (k *MaxOp) Rune() rune         { return SCE_MAX }
func (k *MinOp) Rune() rune         { return SCE_MIN }
func (k *DropHighestOp) Rune() rune { return SCE_DHV }
func (k *DropLowestOp) Rune() rune  { return SCE_DLV }

// d returns a value between 1 and i inclusive.
func d(i int) int { return 1 + rand.Intn(i) }
