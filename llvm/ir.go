package llvm

import (
	"github.com/llir/llvm/ir"
	"github.com/llir/llvm/ir/constant"
	"github.com/llir/llvm/ir/enum"
	"github.com/llir/llvm/ir/types"
)

var (
	HeadType = types.I32
	CellType = types.I8
)

var (
	CellMin = constant.NewInt(CellType, 0)
	CellMax = constant.NewInt(CellType, 0xFF)
)

type BFContext struct {
	Head    *ir.InstAlloca
	Tape    *ir.InstCall
	TapeLen int32
	Get     *ir.Func
	Put     *ir.Func
}

func (ctx BFContext) Load(block *ir.Block) (head *ir.InstLoad, vPtr *ir.InstGetElementPtr, v *ir.InstLoad) {
	head = block.NewLoad(HeadType, ctx.Head)
	vPtr = block.NewGetElementPtr(CellType, ctx.Tape, head)
	v = block.NewLoad(CellType, vPtr)
	return
}

func InitIR(tapeLen int32) (*ir.Module, *BFContext, *ir.Block) {
	if tapeLen <= 0 {
		panic("Length of tape cannot be negative nor zero")
	}
	module := ir.NewModule()
	main := module.NewFunc("main", types.I32)
	entry := main.NewBlock("")
	head := entry.NewAlloca(HeadType)
	head.SetName("head")
	entry.NewStore(constant.NewInt(HeadType, 0), head)
	calloc := module.NewFunc("calloc", types.NewPointer(CellType), ir.NewParam("", types.I64), ir.NewParam("", types.I64))
	tape := entry.NewCall(calloc, constant.NewInt(types.I64, int64(tapeLen)), constant.NewInt(types.I64, int64(CellType.BitSize/8)))
	tape.SetName("tape")
	getchar := module.NewFunc("getchar", types.I32)
	putchar := module.NewFunc("putchar", types.I32, ir.NewParam("", types.I32))
	ctx := &BFContext{
		Head:    head,
		Tape:    tape,
		TapeLen: tapeLen,
		Get:     getchar,
		Put:     putchar,
	}
	next := main.NewBlock("")
	entry.NewBr(next)
	return module, ctx, next
}

func CloseIR(ctx *BFContext, block *ir.Block) {
	free := block.Parent.Parent.NewFunc("free", types.Void, ir.NewParam("", types.NewPointer(CellType)))
	block.NewCall(free, ctx.Tape)
	block.NewRet(constant.NewInt(types.I32, 0))
}

func IncIR(ctx *BFContext, block *ir.Block) *ir.Block {
	_, vPtr, v := ctx.Load(block)
	newV := block.NewAdd(v, constant.NewInt(CellType, 1))
	block.NewStore(newV, vPtr)
	next := block.Parent.NewBlock("")
	block.NewBr(next)
	return next
}

func DecIR(ctx *BFContext, block *ir.Block) *ir.Block {
	_, vPtr, v := ctx.Load(block)
	newV := block.NewSub(v, constant.NewInt(CellType, 1))
	block.NewStore(newV, vPtr)
	next := block.Parent.NewBlock("")
	block.NewBr(next)
	return next
}

func PrvIR(ctx *BFContext, block *ir.Block) *ir.Block {
	head, _, _ := ctx.Load(block)
	newHead := block.NewSub(head, constant.NewInt(CellType, 1))
	cond := block.NewICmp(enum.IPredSLT, newHead, constant.NewInt(HeadType, 0))
	then := block.Parent.NewBlock("")
	then.NewStore(constant.NewInt(HeadType, int64(ctx.TapeLen-1)), ctx.Head)
	els := block.Parent.NewBlock("")
	els.NewStore(newHead, ctx.Head)
	block.NewCondBr(cond, then, els)
	next := block.Parent.NewBlock("")
	then.NewBr(next)
	els.NewBr(next)
	return next
}

func NxtIR(ctx *BFContext, block *ir.Block) *ir.Block {
	head, _, _ := ctx.Load(block)
	newHead := block.NewAdd(head, constant.NewInt(CellType, 1))
	cond := block.NewICmp(enum.IPredSGE, newHead, constant.NewInt(HeadType, int64(ctx.TapeLen)))
	then := block.Parent.NewBlock("")
	then.NewStore(constant.NewInt(HeadType, 0), ctx.Head)
	els := block.Parent.NewBlock("")
	els.NewStore(newHead, ctx.Head)
	block.NewCondBr(cond, then, els)
	next := block.Parent.NewBlock("")
	then.NewBr(next)
	els.NewBr(next)
	return next
}

func GetIR(ctx *BFContext, block *ir.Block) *ir.Block {
	got32 := block.NewCall(ctx.Get)
	got := block.NewTrunc(got32, CellType)
	_, vPtr, _ := ctx.Load(block)
	block.NewStore(got, vPtr)
	next := block.Parent.NewBlock("")
	block.NewBr(next)
	return next
}

func PutIR(ctx *BFContext, block *ir.Block) *ir.Block {
	_, _, v := ctx.Load(block)
	v32 := block.NewZExt(v, types.I32)
	block.NewCall(ctx.Put, v32)
	next := block.Parent.NewBlock("")
	block.NewBr(next)
	return next
}

func LoopStartIR(ctx *BFContext, block *ir.Block) (loop *ir.Block, next *ir.Block) {
	_, _, v := ctx.Load(block)
	cond := block.NewICmp(enum.IPredNE, v, constant.NewInt(CellType, 0))
	loop = block.Parent.NewBlock("")
	next = block.Parent.NewBlock("")
	block.NewCondBr(cond, loop, next)
	return
}

func LoopEndIR(loopStart *ir.Block, loopEnd *ir.Block) {
	loopEnd.NewBr(loopStart)
}
