package parse

import (
	"github.com/Bokume2/bfc_llvm/lex"
	"github.com/Bokume2/bfc_llvm/llvm"
	goLLVM "tinygo.org/x/go-llvm"
)

func GenIR(t AST, tapeLen int32) goLLVM.Module {
	cc := llvm.InitIR(tapeLen)
	defer cc.Builder.Dispose()
	addIRFor(t, cc)
	llvm.CloseIR(cc)
	return cc.Module
}

func addIRFor(ns []Node, cc *llvm.CompilerContext) {
	for _, n := range ns {
		n.AddIR(cc)
	}
}

func (cn CmdNode) AddIR(cc *llvm.CompilerContext) {
	switch cn.Cmd {
	case lex.IncToken:
		llvm.IncIR(cc)
	case lex.DecToken:
		llvm.DecIR(cc)
	case lex.PrvToken:
		llvm.PrvIR(cc)
	case lex.NxtToken:
		llvm.NxtIR(cc)
	case lex.GetToken:
		llvm.GetIR(cc)
	case lex.PutToken:
		llvm.PutIR(cc)
	}
}

func (ln LoopNode) AddIR(cc *llvm.CompilerContext) {
	loopStart, next := llvm.LoopStartIR(cc)
	addIRFor(ln.Block, cc)
	llvm.LoopEndIR(cc, loopStart, next)
}
