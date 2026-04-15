package main

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/pkg/errors"
)

// tokenize is the entry point for the tokenization process.
// It takes a string like "2d6+3" and builds an AST (Abstract Syntax Tree).
//
// How it works:
// 1. Wraps the input in a tokenExpr to start tokenization
// 2. Passes to tokenizeExpr which recursively splits at operators
// 3. Returns the root of the AST
//
// Example flow for "2d6+3":
//
//	Input: "2d6+3"
//	Step 1: wrap in tokenExpr -> tokenExpr{token: "2d6+3"}
//	Step 2: tokenizeExpr processes in PEMDAS order:
//	  - Try '+': find at position 3, split -> lhs="2d6", rhs="3"
//	  - Recurse on lhs="2d6":
//	    - Try 'd': find at position 1, split -> lhs="2", rhs="6"
//	Step 3: Returns AST:
//	    binaryExpr(
//	      op=AddOp,
//	      lhs=binaryExpr(op=RollOp, lhs=2, rhs=6),
//	      rhs=3
//	    )
func tokenize(s string) (expr, error) {
	expr := &tokenExpr{token: s}
	e, err := expr.Tokenize()
	return e, errors.Wrap(err, "initial tokenization")
}

// tokenExpr.Tokenize splits a token string at a given operator.
// This is the core of the tokenizer - it finds operators in the
// string and builds the AST recursively. The final end product is
// tree of binary operators with literal value token expression leaves.
func (t *tokenExpr) Tokenize() (e expr, err error) {
	for _, opGroup := range pemdas {
		// Split the token on the first instance of the operator.
		ind := strings.LastIndexAny(t.token, opGroup)
		if ind < 0 {
			// Can't split on this opGroup; try to find the next operator.
			continue
		}

		op := ops[[]rune(t.token)[ind]]

		// Split token at operator position
		switch typ := op.(type) {
		case binaryOp:
			lhs := &tokenExpr{t.token[:ind]}
			rhs := &tokenExpr{t.token[ind+1:]}

			bexpr := &binaryExpr{op: typ}

			// Handle RollOp defaults: "d6" -> "1d6"
			switch typ.(type) {
			case *RollOp, *DropHighestOp, *DropLowestOp, *MaxOp, *MinOp:
				if strings.TrimSpace(lhs.token) == "" {
					bexpr.lhs = &litValExpr{val: 1}
				} else if bexpr.lhs, err = lhs.Tokenize(); err != nil {
					return nil, errors.Wrap(err, "tokenizing lhs of binary op")
				}

				if strings.TrimSpace(rhs.token) == "" {
					return nil, errors.New("cannot roll a die with an undefined number of sides")
				}

				if bexpr.rhs, err = rhs.Tokenize(); err != nil {
					return nil, errors.Wrap(err, "tokenizing rhs of binary op")
				}
			default:
				if strings.TrimSpace(lhs.token) == "" {
					return nil, fmt.Errorf("op %T cannot have an empty lhs", typ)
				}
				if bexpr.lhs, err = lhs.Tokenize(); err != nil {
					return nil, errors.Wrap(err, "tokenizing lhs of binary op")
				}

				if strings.TrimSpace(rhs.token) == "" {
					return nil, fmt.Errorf("op %T cannot have an empty rhs", typ)
				}
				if bexpr.rhs, err = rhs.Tokenize(); err != nil {
					return nil, errors.Wrap(err, "tokenizing rhs of binary op")
				}
			}
			return bexpr, nil
		}

		return nil, fmt.Errorf("unknown op type: %T", op)
	}

	// This token can't be broken down any further; this must be a literal value.
	val, err := strconv.Atoi(t.token)
	if err != nil {
		return nil, errors.Wrap(err, "literal values must be integers")
	}

	return &litValExpr{val: val}, nil
}
