package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/Bokume2/bfc_llvm/lex"
	"github.com/Bokume2/bfc_llvm/llvm"
	"github.com/Bokume2/bfc_llvm/parse"
	goLLVM "tinygo.org/x/go-llvm"
)

func main() {
	srcFilePath := ""
	maybeOption := true
	optLv := 3
	optRE := regexp.MustCompile(`-O(\d)`)
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}

		if maybeOption && optRE.MatchString(arg) {
			optLv, _ = strconv.Atoi(optRE.FindStringSubmatch(arg)[1])
			if optLv < 0 || optLv > 3 {
				fmt.Fprintln(os.Stderr, "Error:", "Optimization "+arg+" is not supported")
				os.Exit(1)
			}
		} else if arg == "--" {
			maybeOption = false
		} else if srcFilePath == "" && (!maybeOption || !strings.HasPrefix(arg, "-")) {
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
	m, main := parse.GenIR(ast, 4096)
	defer m.Context().Dispose()
	llvm.OptimizeIR(m, uint(optLv))
	ee, err := goLLVM.NewInterpreter(m)
	defer ee.Dispose()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	_ = ee.RunFunction(main, []goLLVM.GenericValue{})
}
