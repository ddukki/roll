package main

import (
	"fmt"
	"os"
	"strings"
)

const (
	SCE_ADD = '+'
	SCE_SUB = '-'
	SCE_MUL = '*'
	SCE_DIV = '/'

	SCE_ROLL = 'd'
)

func main() {
	args := os.Args[1:]
	rollStr := strings.Join(args, " ")
	tkn := tokenize(rollStr)
	ast := parse(tkn)
	fmt.Println(ast.Value())
}
