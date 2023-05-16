package main

import (
	"fmt"
	"strconv"
	"strings"
)

// parse recursively parses leaf token expressions into evaluable expressions.
func parse(e expr) expr {
	switch t := e.(type) {
	case *binaryExpr:
		t.lhs = parse(t.lhs)
		t.rhs = parse(t.rhs)
	case *litValExpr:
	case *tokenExpr:
		valStr := strings.TrimSpace(t.token)
		val64, err := strconv.ParseInt(valStr, 10, 32)
		if err != nil {
			panic(err)
		}
		e = &litValExpr{val: int(val64)}
	default:
		panic(fmt.Sprintf("unknown expr type: %T", e))
	}

	return e
}
