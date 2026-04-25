package llvm

import (
	"fmt"

	"tinygo.org/x/go-llvm"
)

func init() {
	llvm.InitializeNativeTarget()
}

func OptimizeIR(m llvm.Module, optLv uint) {
	t, err := llvm.GetTargetFromTriple(llvm.DefaultTargetTriple())
	if err != nil {
		panic(err)
	}
	tm := t.CreateTargetMachine(llvm.DefaultTargetTriple(), "", "", llvm.CodeGenLevelDefault, llvm.RelocDefault, llvm.CodeModelDefault)
	if err = m.RunPasses(fmt.Sprintf("default<O%d>", optLv), tm, llvm.NewPassBuilderOptions()); err != nil {
		panic(err)
	}
}
