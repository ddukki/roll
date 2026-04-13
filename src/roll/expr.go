package main

import (
  "fmt"
  "strings"
)

type expr interface {
  Value() int
  Validate()
  Roll() *RollDetails
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

func (b *binaryExpr) Roll() *RollDetails {
  if _, ok := b.op.(*RollOp); ok {
    oldMode := inRollMode
    inRollMode = true
    defer func() { inRollMode = oldMode }()

    chain := buildRollChain(b)
    return evaluateRollChain(chain)
  }

  b.rhs.Value()
  b.lhs.Value()
  return nil
}

func (b *binaryExpr) Validate() {
  if _, ok := b.op.(*RollOp); ok {
    if texpr, ok := b.lhs.(*tokenExpr); ok {
      if strings.TrimSpace(texpr.token) == "" {
        b.lhs = &litValExpr{val: 1}
      }
    }
  }
}

var _ expr = (*litValExpr)(nil)

type litValExpr struct {
  val int
}

func (lv *litValExpr) Value() int {
  return lv.val
}

func (lv *litValExpr) Validate()          {}
func (lv *litValExpr) Roll() *RollDetails { return nil }

var _ expr = (*tokenExpr)(nil)

type tokenExpr struct {
  token string
}

func (t *tokenExpr) Value() int {
  return 0
}

func (t *tokenExpr) Validate() {}

func (t *tokenExpr) Roll() *RollDetails {
  return nil
}

func (t *tokenExpr) Tokenize(o op) expr {
  ind := strings.LastIndex(t.token, string(o.Rune()))
  if ind < 0 {
    return t
  }

  switch typ := o.(type) {
  case binaryOp:
    lhs := &tokenExpr{t.token[:ind]}
    rhs := &tokenExpr{t.token[ind+1:]}

    return &binaryExpr{lhs: lhs.Tokenize(o), rhs: rhs.Tokenize(o), op: typ}
  }

  panic(fmt.Sprintf("unknown op type: %T", o))
}
