package main

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"regexp"
	"strconv"

	"github.com/Bokume2/bfc_llvm/lex"
	"github.com/Bokume2/bfc_llvm/llvm"
	"github.com/Bokume2/bfc_llvm/parse"
)

func main() {
	const llcCmd = "llc"
	const ccCmd = "cc"

	srcFilePath := ""
	optLv := 3
	optRE := regexp.MustCompile(`-O(\d)`)
	toCompile := false
	outFilePath := ""
	switchOutFile := false
	for i, arg := range os.Args {
		if i == 0 {
			continue
		}

		if switchOutFile {
			outFilePath = arg
			switchOutFile = false
			continue
		}

		if optRE.MatchString(arg) {
			optLv, _ = strconv.Atoi(optRE.FindStringSubmatch(arg)[1])
			if optLv < 0 || optLv > 3 {
				fmt.Fprintln(os.Stderr, "Error:", "Optimization "+arg+" is not supported")
				os.Exit(1)
			}
		} else if arg == "-c" || arg == "--compile-to-exe" {
			toCompile = true
		} else if arg == "-o" {
			switchOutFile = true
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
	m, _ := parse.GenIR(ast, 4096)
	defer m.Context().Dispose()
	defer m.Dispose()
	llvm.OptimizeIR(m, uint(optLv))

	if !toCompile {
		if outFilePath != "" {
			irFile, err := os.Create(outFilePath)
			if err != nil {
				fmt.Fprintln(os.Stderr, "Error:", "Cannot write to "+outFilePath)
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
			if _, err := irFile.WriteString(m.String()); err != nil {
				fmt.Fprintln(os.Stderr, "Error:", "Cannot write to "+outFilePath)
				fmt.Fprintln(os.Stderr, err)
				os.Exit(1)
			}
		} else {
			fmt.Println(m)
		}
		return
	}

	exeFilePath := regexp.MustCompile(`\.[\dA-Za-z]*?$`).ReplaceAllString(srcFilePath, "")
	asmFilePath := exeFilePath + ".s"

	if outFilePath != "" {
		exeFilePath = outFilePath
	}

	compileIR := exec.Command(llcCmd, "-o", asmFilePath)
	compileIRStdin, err := compileIR.StdinPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer compileIRStdin.Close()
	io.WriteString(compileIRStdin, m.String())
	compileIRStdin.Close()
	if err := compileIR.Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := exec.Command(ccCmd, "-o", exeFilePath, asmFilePath).Run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if err := os.Remove(asmFilePath); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
