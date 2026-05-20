package llvm

import (
	"tinygo.org/x/go-llvm"
)

type CompilerContext struct {
	LLVMContext llvm.Context
	Module      llvm.Module
	Builder     llvm.Builder
	Main        llvm.Value
	CellType    llvm.Type
	HeadType    llvm.Type
	Head        llvm.Value
	Tape        llvm.Value
	TapeLen     int32
	GetType     llvm.Type
	Get         llvm.Value
	PutType     llvm.Type
	Put         llvm.Value
}

func (cc *CompilerContext) Load() (head, vPtr, v llvm.Value) {
	head = cc.Builder.CreateLoad(cc.HeadType, cc.Head, "")
	vPtr = cc.Builder.CreateInBoundsGEP(llvm.PointerType(cc.CellType, 0), cc.Tape, []llvm.Value{head}, "")
	v = cc.Builder.CreateLoad(cc.CellType, vPtr, "")
	return
}

func InitIR(tapeLen int32) *CompilerContext {
	if tapeLen <= 0 {
		panic("Length of tape cannot be negative nor zero")
	}
	c := llvm.NewContext()
	b := c.NewBuilder()
	m := c.NewModule("")
	m.SetTarget(llvm.DefaultTargetTriple())
	main := llvm.AddFunction(m, "main", llvm.FunctionType(c.Int32Type(), []llvm.Type{}, false))
	entry := llvm.AddBasicBlock(main, "entry")
	b.SetInsertPoint(entry, entry.FirstInstruction())
	cellType := c.Int8Type()
	headType := c.Int32Type()
	head := b.CreateAlloca(headType, "head")
	b.CreateStore(llvm.ConstInt(headType, 0, false), head)
	callocType := llvm.FunctionType(llvm.PointerType(cellType, 0), []llvm.Type{
		c.Int64Type(),
		c.Int64Type(),
	}, false)
	calloc := llvm.AddFunction(m, "calloc", callocType)
	tape := b.CreateCall(callocType, calloc, []llvm.Value{
		llvm.ConstInt(c.Int64Type(), uint64(tapeLen), false),
		llvm.ConstInt(c.Int64Type(), 1, false),
	}, "tape")
	getcharType := llvm.FunctionType(c.Int32Type(), []llvm.Type{}, false)
	getchar := llvm.AddFunction(m, "getchar", getcharType)
	putcharType := llvm.FunctionType(c.Int32Type(), []llvm.Type{c.Int32Type()}, false)
	putchar := llvm.AddFunction(m, "putchar", putcharType)
	next := llvm.AddBasicBlock(main, "")
	b.CreateBr(next)
	b.SetInsertPoint(next, next.FirstInstruction())
	cc := &CompilerContext{
		LLVMContext: c,
		Module:      m,
		Builder:     b,
		Main:        main,
		CellType:    cellType,
		HeadType:    headType,
		Head:        head,
		Tape:        tape,
		TapeLen:     tapeLen,
		GetType:     getcharType,
		Get:         getchar,
		PutType:     putcharType,
		Put:         putchar,
	}
	return cc
}

func CloseIR(cc *CompilerContext) {
	cc.Builder.CreateFree(cc.Tape)
	cc.Builder.CreateRet(llvm.ConstInt(cc.LLVMContext.Int32Type(), 0, false))
	if err := llvm.VerifyModule(cc.Module, llvm.ReturnStatusAction); err != nil {
		panic(err)
	}
}

func IncIR(cc *CompilerContext) {
	_, vPtr, v := cc.Load()
	tmp := cc.Builder.CreateAdd(v, llvm.ConstInt(cc.CellType, 1, false), "")
	cc.Builder.CreateStore(tmp, vPtr)
	next := llvm.AddBasicBlock(cc.Builder.GetInsertBlock().Parent(), "")
	cc.Builder.CreateBr(next)
	cc.Builder.SetInsertPoint(next, next.FirstInstruction())
}

func DecIR(cc *CompilerContext) {
	_, vPtr, v := cc.Load()
	tmp := cc.Builder.CreateSub(v, llvm.ConstInt(cc.CellType, 1, false), "")
	cc.Builder.CreateStore(tmp, vPtr)
	next := llvm.AddBasicBlock(cc.Builder.GetInsertBlock().Parent(), "")
	cc.Builder.CreateBr(next)
	cc.Builder.SetInsertPoint(next, next.FirstInstruction())
}

func PrvIR(cc *CompilerContext) {
	head, _, _ := cc.Load()
	tmp := cc.Builder.CreateSub(head, llvm.ConstInt(cc.HeadType, 1, false), "")
	cond := cc.Builder.CreateICmp(llvm.IntSLT, tmp, llvm.ConstInt(cc.HeadType, 0, false), "")
	then := llvm.AddBasicBlock(cc.Builder.GetInsertBlock().Parent(), "")
	els := llvm.AddBasicBlock(cc.Builder.GetInsertBlock().Parent(), "")
	next := llvm.AddBasicBlock(cc.Builder.GetInsertBlock().Parent(), "")
	cc.Builder.CreateCondBr(cond, then, els)
	cc.Builder.SetInsertPoint(then, then.FirstInstruction())
	cc.Builder.CreateStore(llvm.ConstInt(cc.HeadType, uint64(cc.TapeLen-1), false), cc.Head)
	cc.Builder.CreateBr(next)
	cc.Builder.SetInsertPoint(els, els.FirstInstruction())
	cc.Builder.CreateStore(tmp, cc.Head)
	cc.Builder.CreateBr(next)
	cc.Builder.SetInsertPoint(next, next.FirstInstruction())
}

func NxtIR(cc *CompilerContext) {
	head, _, _ := cc.Load()
	tmp := cc.Builder.CreateAdd(head, llvm.ConstInt(cc.HeadType, 1, false), "")
	cond := cc.Builder.CreateICmp(llvm.IntSGE, tmp, llvm.ConstInt(cc.HeadType, uint64(cc.TapeLen), false), "")
	then := llvm.AddBasicBlock(cc.Builder.GetInsertBlock().Parent(), "")
	els := llvm.AddBasicBlock(cc.Builder.GetInsertBlock().Parent(), "")
	next := llvm.AddBasicBlock(cc.Builder.GetInsertBlock().Parent(), "")
	cc.Builder.CreateCondBr(cond, then, els)
	cc.Builder.SetInsertPoint(then, then.FirstInstruction())
	cc.Builder.CreateStore(llvm.ConstInt(cc.HeadType, 0, false), cc.Head)
	cc.Builder.CreateBr(next)
	cc.Builder.SetInsertPoint(els, els.FirstInstruction())
	cc.Builder.CreateStore(tmp, cc.Head)
	cc.Builder.CreateBr(next)
	cc.Builder.SetInsertPoint(next, next.FirstInstruction())
}

func GetIR(cc *CompilerContext) {
	got32 := cc.Builder.CreateCall(cc.GetType, cc.Get, []llvm.Value{}, "")
	got := cc.Builder.CreateTrunc(got32, cc.CellType, "")
	_, vPtr, _ := cc.Load()
	cc.Builder.CreateStore(got, vPtr)
	next := llvm.AddBasicBlock(cc.Builder.GetInsertBlock().Parent(), "")
	cc.Builder.CreateBr(next)
	cc.Builder.SetInsertPoint(next, next.FirstInstruction())
}

func PutIR(cc *CompilerContext) {
	_, _, v := cc.Load()
	v32 := cc.Builder.CreateZExt(v, cc.LLVMContext.Int32Type(), "")
	cc.Builder.CreateCall(cc.PutType, cc.Put, []llvm.Value{v32}, "")
	next := llvm.AddBasicBlock(cc.Builder.GetInsertBlock().Parent(), "")
	cc.Builder.CreateBr(next)
	cc.Builder.SetInsertPoint(next, next.FirstInstruction())
}

func LoopStartIR(cc *CompilerContext) (loopStart, next llvm.BasicBlock) {
	loopStart = cc.Builder.GetInsertBlock()
	_, _, v := cc.Load()
	cond := cc.Builder.CreateICmp(llvm.IntNE, v, llvm.ConstInt(cc.CellType, 0, false), "")
	loop := llvm.AddBasicBlock(cc.Builder.GetInsertBlock().Parent(), "")
	next = llvm.AddBasicBlock(cc.Builder.GetInsertBlock().Parent(), "")
	cc.Builder.CreateCondBr(cond, loop, next)
	cc.Builder.SetInsertPoint(loop, loop.FirstInstruction())
	return
}

func LoopEndIR(cc *CompilerContext, loopStart llvm.BasicBlock, next llvm.BasicBlock) {
	cc.Builder.CreateBr(loopStart)
	cc.Builder.SetInsertPoint(next, next.FirstInstruction())
}
