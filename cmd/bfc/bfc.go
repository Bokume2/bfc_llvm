package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"

	"github.com/Bokume2/bfc_llvm/lex"
	"github.com/Bokume2/bfc_llvm/llvm"
	"github.com/Bokume2/bfc_llvm/parse"
)

func main() {
	srcFilePath := ""
	optLv := 3
	optRE := regexp.MustCompile(`-O(\d)`)
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}

		if optRE.MatchString(arg) {
			optLv, _ = strconv.Atoi(optRE.FindStringSubmatch(arg)[1])
			if optLv < 0 || optLv > 3 {
				fmt.Fprintln(os.Stderr, "Error:", "Optimization "+arg+" is not supported")
				os.Exit(1)
			}
		} else if srcFilePath == "" {
			srcFilePath = arg
		} else {
			fmt.Fprintln(os.Stderr, "Error:", "Unknown argument "+arg)
			os.Exit(1)
		}
	}
	if srcFilePath == "" {
		fmt.Fprintln(os.Stderr, "Error:", "Source file is required")
		os.Exit(1)
	}

	code, err := os.ReadFile(srcFilePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error:", "Cannot read source file with following error,")
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
	llvm.OptimizeIR(m, uint(optLv))
	fmt.Println(m)
}
