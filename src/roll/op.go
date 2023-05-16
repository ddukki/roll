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
	SCE_KHT  = 'h'
	SCE_KLT  = 'l'
	SCE_DHT  = 'H'
	SCE_DLT  = 'L'
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
	_ binaryOp = (*KeepHighestOp)(nil)
	_ binaryOp = (*KeepLowestOp)(nil)
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
		SCE_KHT,
		SCE_KLT,
		SCE_DHT,
		SCE_DLT,
	}

	ops = map[rune]op{
		SCE_ADD:  &AddOp{},
		SCE_SUB:  &SubOp{},
		SCE_DIV:  &DivOp{},
		SCE_MUL:  &MulOp{},
		SCE_ROLL: &RollOp{},
		SCE_KHT:  &KeepHighestOp{},
		SCE_KLT:  &KeepLowestOp{},
		SCE_DHT:  &DropHighestOp{},
		SCE_DLT:  &DropLowestOp{},
	}
)

func (a *AddOp) Apply(lhs int, rhs int) int {
	return lhs + rhs
}
func (s *SubOp) Apply(lhs int, rhs int) int {
	return lhs - rhs
}
func (d *DivOp) Apply(lhs int, rhs int) int {
	return lhs / rhs
}
func (m *MulOp) Apply(lhs int, rhs int) int {
	return lhs * rhs
}

func (a *AddOp) Rune() rune {
	return SCE_ADD
}
func (s *SubOp) Rune() rune {
	return SCE_SUB
}
func (d *DivOp) Rune() rune {
	return SCE_DIV
}
func (m *MulOp) Rune() rune {
	return SCE_MUL
}

var (
	_ op = (*RollOp)(nil)
)

type RollOp struct{}

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

func (m *RollOp) Rune() rune {
	return SCE_ROLL
}

type KeepHighestOp struct{}

func (r *KeepHighestOp) Apply(lhs, rhs int) int {
	rolls := make([]int, lhs)
	var val int
	for i := range rolls {
		rolls[i] = d(rhs)
		if val < rolls[i] {
			val = rolls[i]
		}
	}

	rollStr := fmt.Sprintf("h%d", rhs)
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

func (k *KeepHighestOp) Rune() rune {
	return SCE_KHT
}

type KeepLowestOp struct{}

func (r *KeepLowestOp) Apply(lhs, rhs int) int {
	rolls := make([]int, lhs)
	val := rhs + 1
	for i := range rolls {
		rolls[i] = d(rhs)
		if val > rolls[i] {
			val = rolls[i]
		}
	}

	rollStr := fmt.Sprintf("l%d", rhs)
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

func (k *KeepLowestOp) Rune() rune {
	return SCE_KLT
}

type DropHighestOp struct{}

func (r *DropHighestOp) Apply(lhs, rhs int) int {
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

	rollStr := fmt.Sprintf("H%d", rhs)
	if lhs <= 1 {
		panic("there must be more than one die rolled for dropping!")
	}

	rollStr = fmt.Sprintf("%d%s", lhs, rollStr)
	mathStr := fmt.Sprintf("sum%+v - max<%d> = %d", rolls, max, val)

	fmt.Printf("Rolling %s: %s\n", rollStr, mathStr)
	return val
}

func (k *DropHighestOp) Rune() rune {
	return SCE_DHT
}

type DropLowestOp struct{}

func (r *DropLowestOp) Apply(lhs, rhs int) int {
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

	rollStr := fmt.Sprintf("L%d", rhs)
	if lhs <= 1 {
		panic("there must be more than one die rolled for dropping!")
	}

	rollStr = fmt.Sprintf("%d%s", lhs, rollStr)
	mathStr := fmt.Sprintf("sum%+v - max<%d> = %d", rolls, min, val)

	fmt.Printf("Rolling %s: %s\n", rollStr, mathStr)
	return val
}

func (k *DropLowestOp) Rune() rune {
	return SCE_DLT
}

// d returns a value between 1 and i inclusive.
func d(i int) int {
	return 1 + rand.Intn(i)
}
