package parse

import (
	"github.com/Bokume2/bfc_llvm/lex"
	"github.com/Bokume2/bfc_llvm/llvm"
	"github.com/llir/llvm/ir"
)

func GenIR(t AST, tapeLen int32) *ir.Module {
	module, ctx, next := llvm.InitIR(tapeLen)
	next = addIRFor(t, ctx, next)
	llvm.CloseIR(ctx, next)
	return module
}

func addIRFor(ns []Node, ctx *llvm.BFContext, block *ir.Block) *ir.Block {
	next := block
	for _, n := range ns {
		next = n.AddIR(ctx, next)
	}
	return next
}

func (cn CmdNode) AddIR(ctx *llvm.BFContext, block *ir.Block) *ir.Block {
	switch cn.Cmd {
	case lex.IncToken:
		return llvm.IncIR(ctx, block)
	case lex.DecToken:
		return llvm.DecIR(ctx, block)
	case lex.PrvToken:
		return llvm.PrvIR(ctx, block)
	case lex.NxtToken:
		return llvm.NxtIR(ctx, block)
	case lex.GetToken:
		return llvm.GetIR(ctx, block)
	case lex.PutToken:
		return llvm.PutIR(ctx, block)
	default:
		return block
	}
}

func (ln LoopNode) AddIR(ctx *llvm.BFContext, block *ir.Block) *ir.Block {
	loop, next := llvm.LoopStartIR(ctx, block)
	loopEnd := addIRFor(ln.Block, ctx, loop)
	llvm.LoopEndIR(block, loopEnd)
	return next
}
