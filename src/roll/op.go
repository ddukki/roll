package main

import (
	"math/rand"
	"slices"
)

// ===== CONSTANTS =====
// Operator symbols used in dice expressions.
// These runes are used to identify operators when parsing expressions.
const (
	SCE_ADD = '+' // Addition: 2d6+3
	SCE_SUB = '-' // Subtraction: 2d6-1
	SCE_MUL = '*' // Multiplication: 2d6*2 (multiply result by 2)
	SCE_DIV = '/' // Division: 2d6/2 (divide result by 2)

	SCE_ROLL = 'd' // Roll dice: 2d6 (roll 2 six-sided dice)
	SCE_MAX  = 'x' // Max aggregation: 3d6x (take max of 3 dice)
	SCE_MIN  = 'n' // Min aggregation: 3d6n (take min of 3 dice)
	SCE_DHV  = 'h' // Drop highest: 3d6h (drop highest die)
	SCE_DLV  = 'l' // Drop lowest: 3d6l (drop lowest die)
)

// ===== RANDOM INTERFACE =====
// Rand defines the interface for random number generators.
// This interface allows both real rand.Rand and fake/test implementations.
// Go 1.24+ changed how rand.Rand works - Intn() delegation doesn't work properly,
// so we use this interface instead of relying on *rand.Rand directly.
type Rand interface {
	Int63() int64   // Returns next random int64
	Intn(n int) int // Returns random int in range [0, n)
}

// globalRand is the default random number generator.
// Each call creates a new Source with a random seed for true randomness.
// Using rand.Int63() as the seed provides unpredictability.
var globalRand Rand = rand.New(rand.NewSource(rand.Int63()))

// SetGlobalRand replaces the random number generator.
// This is primarily used for testing with deterministic fakeRand.
func SetGlobalRand(r Rand) {
	globalRand = r
}

// ResetGlobalRand restores the default random number generator.
func ResetGlobalRand() {
	globalRand = rand.New(rand.NewSource(rand.Int63()))
}

// getRand returns the current random number generator.
func getRand() Rand {
	return globalRand
}

// ===== OPERATOR INTERFACES =====
// op is the base interface for all operators.
// Every operator must be able to return its rune symbol
// so the tokenizer knows which operator to apply.
type op interface {
	Rune() rune
}

// binaryOp combines op with Apply.
// Operators that combine two values (lhs and rhs) implement this.
// Example: AddOp.Apply(2, 3) = 5.
type binaryOp interface {
	op
	Apply(lhs int, rhs int) *evalInfo
}

// Compile-time check that all operators implement binaryOp.
// This ensures we haven't forgotten to implement Apply() for any operator.
var (
	_ binaryOp = (*AddOp)(nil)
	_ binaryOp = (*SubOp)(nil)
	_ binaryOp = (*MulOp)(nil)
	_ binaryOp = (*DivOp)(nil)
	_ binaryOp = (*RollOp)(nil)
	_ binaryOp = (*MaxOp)(nil)
	_ binaryOp = (*MinOp)(nil)
)

// ===== MATH OPERATORS =====
// AddOp: + (addition)
type AddOp struct{}

// SubOp: - (subtraction)
type SubOp struct{}

// MulOp: * (multiplication)
type MulOp struct{}

// DivOp: / (division)
type DivOp struct{}

// ===== OPERATOR PRECEDENCE =====
// pemdas defines operator precedence for parsing.
// Lower indices = processed first (binding tighter).
// Order: + - * / d x n h l
//
// The Roll/Chain operators (d, x, n, h, l) are processed AFTER math
// so that "3d6+2" splits at '+' before 'd'. This produces:
// (3d6) + 2 rather than 3 * (d6 + 2).
var pemdas = []string{
	string([]rune{SCE_ADD, SCE_SUB}), // + Addition, - Subtraction
	string([]rune{SCE_MUL, SCE_DIV}), // * Multiplication, / Division
	string([]rune{
		SCE_ROLL, // d Roll dice
		SCE_MAX,  // x Max aggregation
		SCE_MIN,  // n Min aggregation
		SCE_DHV,  // h Drop highest
		SCE_DLV,  // l Drop lowest
	}),
}

// ops maps each operator symbol to its operator instance.
// Used during tokenization to look up operators by their rune.
var ops = map[rune]op{
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

// ===== OPERATOR IMPLEMENTATIONS =====
// Apply methods perform the actual operation on two integer values.

func (a *AddOp) Apply(lhs int, rhs int) *evalInfo { return &evalInfo{op: a, value: lhs + rhs} }
func (s *SubOp) Apply(lhs int, rhs int) *evalInfo { return &evalInfo{op: s, value: lhs - rhs} }
func (d *DivOp) Apply(lhs int, rhs int) *evalInfo { return &evalInfo{op: d, value: lhs / rhs} }
func (m *MulOp) Apply(lhs int, rhs int) *evalInfo { return &evalInfo{op: m, value: lhs * rhs} }

// Rune methods return the operator symbol for display/debug.

func (a *AddOp) Rune() rune { return SCE_ADD }
func (s *SubOp) Rune() rune { return SCE_SUB }
func (d *DivOp) Rune() rune { return SCE_DIV }
func (m *MulOp) Rune() rune { return SCE_MUL }

// Compile-time check for RollOp (also a binaryOp).
var (
	_ op = (*RollOp)(nil)
)

// ===== DICE OPERATORS =====
// These operators handle dice-specific operations.

// RollOp: d (roll dice)
// Rolls a number of dice with a number of sides and returns their sum.
// lhs = number of dice (e.g., 2 for "2d6")
// rhs = number of sides (e.g., 6 for "d6")
type RollOp struct{}

// MaxOp: x (maximum aggregation)
// Returns the maximum value between lhs and rhs.
// In a chain like 3d6x3d6, this compares the totals.
type MaxOp struct{}

// MinOp: n (minimum aggregation)
// Returns the minimum value between lhs and rhs.
type MinOp struct{}

// DropHighestOp: h (drop highest)
// Drops the highest die from a rolled set.
// Example: 3d6h rolls 3 dice and drops the highest.
type DropHighestOp struct{}

// DropLowestOp: l (drop lowest)
// Drops the lowest die from a rolled set.
// Example: 3d6l rolls 3 dice and drops the lowest.
type DropLowestOp struct{}

// RollOp.Apply rolls dice and returns their sum.
// This is the core dice rolling logic. It:
// 1. Creates a slice to hold individual dice results
// 2. Rolls each die using d() (helper that returns 1 to sides)
// 3. Sums all dice
// 4. Stores result in rollStack for later access
func (r *RollOp) Apply(lhs, rhs int) *evalInfo {
	ei := &evalInfo{
		op:  r,
		dri: &diceRollInfo{dieSize: rhs},
	}
	rolls := make([]int, lhs)

	for i := range rolls {
		rolls[i] = d(rhs)
		ei.value += rolls[i]
	}

	ei.dri.rolls = rolls
	return ei
}

// MaxOp.Apply compares lhs and rhs, returning the maximum.
func (mx *MaxOp) Apply(lhs, rhs int) *evalInfo {
	ei := &evalInfo{
		op:  mx,
		dri: &diceRollInfo{dieSize: rhs},
	}
	rolls := make([]int, 0, lhs)

	for range lhs {
		rolls = append(rolls, d(rhs))
	}

	ei.value = slices.Max(rolls)
	ei.dri.rolls = rolls
	return ei
}

// MinOp.Apply compares lhs and rhs, returning the minimum.
func (mn *MinOp) Apply(lhs, rhs int) *evalInfo {
	ei := &evalInfo{
		op:  mn,
		dri: &diceRollInfo{dieSize: rhs},
	}
	rolls := make([]int, 0, lhs)

	for range lhs {
		rolls = append(rolls, d(rhs))
	}

	ei.value = slices.Min(rolls)
	ei.dri.rolls = rolls
	return ei
}

// DropHighestOp.Apply rolls dice and drops the highest.
// This is used for "3d6h" style expressions where you roll
// multiple dice and discard the highest result.
func (dh *DropHighestOp) Apply(lhs, rhs int) *evalInfo {
	if lhs <= 1 {
		panic("there must be more than one die rolled for dropping!")
	}
	ei := &evalInfo{
		op:  dh,
		dri: &diceRollInfo{dieSize: rhs},
	}
	rolls := make([]int, 0, lhs)

	for range lhs {
		val := d(rhs)
		rolls = append(rolls, val)
		ei.value += val
	}

	ei.value -= slices.Max(rolls)
	ei.dri.rolls = rolls
	return ei
}

// DropLowestOp.Apply rolls dice and drops the lowest.
func (dl *DropLowestOp) Apply(lhs, rhs int) *evalInfo {
	if lhs <= 1 {
		panic("there must be more than one die rolled for dropping!")
	}
	ei := &evalInfo{
		op:  dl,
		dri: &diceRollInfo{dieSize: rhs},
	}
	rolls := make([]int, 0, lhs)

	for range lhs {
		val := d(rhs)
		rolls = append(rolls, val)
		ei.value += val
	}

	ei.value -= slices.Min(rolls)
	ei.dri.rolls = rolls
	return ei
}

// Rune methods return each operator's symbol.

func (m *RollOp) Rune() rune        { return SCE_ROLL }
func (k *MaxOp) Rune() rune         { return SCE_MAX }
func (k *MinOp) Rune() rune         { return SCE_MIN }
func (k *DropHighestOp) Rune() rune { return SCE_DHV }
func (k *DropLowestOp) Rune() rune  { return SCE_DLV }

// d returns a random value between 1 and i inclusive.
// This is the helper function used by Apply() methods.
func d(i int) int { return 1 + getRand().Intn(i) }
