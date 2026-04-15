package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {
	ResetGlobalRand()
	flag.Parse()
	if flag.NArg() > 0 {
		expr := strings.Join(flag.Args(), " ")

		ast, err := tokenize(expr)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Error:", err)
			os.Exit(1)
		}

		ast.Value()
		fmt.Println(RenderResult(expr, ast))
		return
	}
}
