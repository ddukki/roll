package main

import (
	"fmt"
	"strings"
)

type expr interface {
	Value() int
}

var _ expr = (*binaryExpr)(nil)

type binaryExpr struct {
	lhs expr
	rhs expr
	op  binaryOp
}

func (b *binaryExpr) Value() int {
	return b.op.Apply(b.lhs.Value(), b.rhs.Value())
}

var _ expr = (*unaryExpr)(nil)

type unaryExpr struct {
	val expr
	op  unaryOp
}

func (u *unaryExpr) Value() int {
	return u.op.Apply(u.val.Value())
}

var _ expr = (*litValExpr)(nil)

type litValExpr struct {
	val int
}

func (lv *litValExpr) Value() int {
	return lv.val
}

var _ expr = (*tokenExpr)(nil)

type tokenExpr struct {
	token string
}

func (t *tokenExpr) Value() int {
	return 0
}

// Tokenize tokenizes the stored token once, breaking the token with the given
// op. If the given op is not present in the stored token, the token expression
// is returned with no changes. Otherwise, a new expression is returned,
// according to the op.
func (t *tokenExpr) Tokenize(o op) expr {
	ind := strings.Index(t.token, string(o.Rune()))
	if ind < 0 {
		return t
	}

	switch typ := o.(type) {
	case binaryOp:
		lhs := &tokenExpr{t.token[:ind]}
		rhs := &tokenExpr{t.token[ind+1:]}

		return &binaryExpr{lhs: lhs.Tokenize(o), rhs: rhs.Tokenize(o), op: typ}
	case unaryOp:
		val := &tokenExpr{t.token[ind+1:]}

		return &unaryExpr{val: val, op: typ}
	}

	panic(fmt.Sprintf("unknown op type: %T", o))
}
