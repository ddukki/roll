package main

import "math/rand"

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

type unaryOp interface {
	op
	Apply(val int) int
	PrePost() bool
}

type RollOp struct{}

func (r *RollOp) Apply(val int) int {
	return d(val)
}

func (r *RollOp) PrePost() bool {
	return true
}

func (m *RollOp) Rune() rune {
	return 'd'
}

// d returns a value between 1 and i inclusive.
func d(i int) int {
	return 1 + rand.Intn(i)
}
