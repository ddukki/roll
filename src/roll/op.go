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
)

type AddOp struct{}
type SubOp struct{}
type MulOp struct{}
type DivOp struct{}

var (
	pemdas = []rune{
		SCE_ADD, SCE_SUB, SCE_MUL, SCE_DIV, SCE_ROLL,
	}

	ops = map[rune]op{
		SCE_ADD:  &AddOp{},
		SCE_SUB:  &SubOp{},
		SCE_DIV:  &DivOp{},
		SCE_MUL:  &MulOp{},
		SCE_ROLL: &RollOp{},
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
	return '+'
}
func (s *SubOp) Rune() rune {
	return '-'
}
func (d *DivOp) Rune() rune {
	return '/'
}
func (m *MulOp) Rune() rune {
	return '*'
}

var (
	_ op = (*RollOp)(nil)
)

type RollOp struct{}

func (r *RollOp) Apply(lhs, rhs int) int {
	rawRoll := d(rhs)
	val := lhs * rawRoll

	rollStr := fmt.Sprintf("d%d", rhs)
	mathStr := fmt.Sprintf("%d", val)
	if lhs > 1 {
		rollStr = fmt.Sprintf("%d%s", lhs, rollStr)
		mathStr = fmt.Sprintf("%d * [%d] = %s", lhs, rawRoll, mathStr)
	} else {
		mathStr = "[" + mathStr + "]"
	}

	fmt.Printf("Rolling %s: %s\n", rollStr, mathStr)
	return val
}

func (m *RollOp) Rune() rune {
	return 'd'
}

// d returns a value between 1 and i inclusive.
func d(i int) int {
	return 1 + rand.Intn(i)
}
