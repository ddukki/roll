package main

import (
	"fmt"
)

// tokenize splits the token into parts by op and builds a tokenized graph.
func tokenize(s string) expr {
	return tokenizeExpr(&tokenExpr{token: s})
}

// tokenizeExpr recursively tokenizes the values at each node.
func tokenizeExpr(e expr) expr {
	for _, p := range pemdas {
		switch t := e.(type) {
		case *binaryExpr:
			t.lhs = tokenizeExpr(t.lhs)
			t.rhs = tokenizeExpr(t.rhs)
			t.Validate()
		case *tokenExpr:
			e = t.Tokenize(ops[p])
			e.Validate()
		case *litValExpr:
		default:
			panic(fmt.Sprintf("unknown expr type: %T", e))
		}
	}

	return e
}
