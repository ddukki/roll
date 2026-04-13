package main

import (
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

type RollDetails struct {
	Expr   string
	Dice   []int
	Sides  int
	Total  int
	Op     string
	Nested []*RollDetails
}

var rollStack []*RollDetails

func GetRollDetails() *RollDetails {
	if len(rollStack) == 0 {
		return nil
	}
	r := rollStack[len(rollStack)-1]
	rollStack = rollStack[:len(rollStack)-1]
	return r
}

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
	prevRoll := GetRollDetails()

	rolls := make([]int, lhs)
	var val int
	for i := range rolls {
		rolls[i] = d(rhs)
		val += rolls[i]
	}

	curRoll := &RollDetails{
		Dice:  rolls,
		Sides: rhs,
		Total: val,
		Op:    "sum",
	}

	if prevRoll != nil {
		curRoll.Nested = append(curRoll.Nested, prevRoll)
		curRoll.Expr = itoa(lhs) + "d[" + prevRoll.Expr + "]"
	} else {
		curRoll.Expr = itoa(lhs) + "d" + itoa(rhs)
	}

	rollStack = append(rollStack, curRoll)
	return val
}

func (mx *MaxOp) Apply(lhs, rhs int) int {
	prevRoll := GetRollDetails()

	rolls := make([]int, lhs)
	var val int
	for i := range rolls {
		rolls[i] = d(rhs)
		if val < rolls[i] {
			val = rolls[i]
		}
	}
	curRoll := &RollDetails{
		Expr:  itoa(lhs) + "x" + itoa(rhs),
		Dice:  rolls,
		Sides: rhs,
		Total: val,
		Op:    "max",
	}
	if prevRoll != nil {
		curRoll.Nested = append(curRoll.Nested, prevRoll)
	}

	rollStack = append(rollStack, curRoll)
	return val
}

func (mn *MinOp) Apply(lhs, rhs int) int {
	prevRoll := GetRollDetails()

	rolls := make([]int, lhs)
	val := rhs + 1
	for i := range rolls {
		rolls[i] = d(rhs)
		if val > rolls[i] {
			val = rolls[i]
		}
	}
	curRoll := &RollDetails{
		Expr:  itoa(lhs) + "n" + itoa(rhs),
		Dice:  rolls,
		Sides: rhs,
		Total: val,
		Op:    "min",
	}
	if prevRoll != nil {
		curRoll.Nested = append(curRoll.Nested, prevRoll)
	}

	rollStack = append(rollStack, curRoll)
	return val
}

func (dh *DropHighestOp) Apply(lhs, rhs int) int {
	if lhs <= 1 {
		panic("there must be more than one die rolled for dropping!")
	}
	prevRoll := GetRollDetails()

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
	curRoll := &RollDetails{
		Expr:  itoa(lhs) + "h" + itoa(rhs),
		Dice:  rolls,
		Sides: rhs,
		Total: val,
		Op:    "drop-h",
	}
	if prevRoll != nil {
		curRoll.Nested = append(curRoll.Nested, prevRoll)
	}

	rollStack = append(rollStack, curRoll)
	return val
}

func (dl *DropLowestOp) Apply(lhs, rhs int) int {
	if lhs <= 1 {
		panic("there must be more than one die rolled for dropping!")
	}
	prevRoll := GetRollDetails()

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
	curRoll := &RollDetails{
		Expr:  itoa(lhs) + "l" + itoa(rhs),
		Dice:  rolls,
		Sides: rhs,
		Total: val,
		Op:    "drop-l",
	}
	if prevRoll != nil {
		curRoll.Nested = append(curRoll.Nested, prevRoll)
	}

	rollStack = append(rollStack, curRoll)
	return val
}

func (m *RollOp) Rune() rune        { return SCE_ROLL }
func (k *MaxOp) Rune() rune         { return SCE_MAX }
func (k *MinOp) Rune() rune         { return SCE_MIN }
func (k *DropHighestOp) Rune() rune { return SCE_DHV }
func (k *DropLowestOp) Rune() rune  { return SCE_DLV }

// d returns a value between 1 and i inclusive.
func d(i int) int { return 1 + rand.Intn(i) }

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	s := ""
	for n > 0 {
		s = string(rune('0'+n%10)) + s
		n /= 10
	}
	return s
}
