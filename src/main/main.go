package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	args := os.Args[1:]
	rollStr := strings.Join(args, " ")
	tkn := tokenize(rollStr)
	ast := parse(tkn)
	fmt.Println(ast.Value())
}
