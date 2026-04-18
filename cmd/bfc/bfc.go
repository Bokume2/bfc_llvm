package main

import (
	"fmt"
	"os"

	"github.com/Bokume2/bfc_llvm/lex"
	"github.com/Bokume2/bfc_llvm/parse"
)

func main() {
	if len(os.Args) <= 1 {
		fmt.Fprintln(os.Stderr, "Error: Source file is required")
		os.Exit(1)
	}

	code, err := os.ReadFile(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error: Cannot read source file with following error,")
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	tokens := lex.NewLexer(lex.DefaultTokenDef()).Lex(string(code))
	ast, err := parse.NewParser().Parse(tokens)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Parse error:", err)
		os.Exit(1)
	}
	m := parse.GenIR(ast, 4096)
	fmt.Println(m)
}
