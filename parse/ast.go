package parse

import (
	"fmt"
	"strings"

	"github.com/Bokume2/bfc_llvm/lex"
	"github.com/Bokume2/bfc_llvm/llvm"
	"github.com/llir/llvm/ir"
)

type Node interface {
	AddIR(*llvm.BFContext, *ir.Block) *ir.Block
}

type CmdNode struct {
	Cmd lex.Token
}

func (cn CmdNode) String() string {
	return "CmdNode{" + cn.Cmd.ToString(lex.DefaultTokenDef()) + "}"
}

type LoopNode struct {
	Block []Node
}

func (ln LoopNode) String() string {
	return "LoopNode" + fmt.Sprint(ln.Block)
}

type AST = []Node

func FormatAST(t AST) string {
	var formatBlock func([]Node) string
	formatBlock = func(ns []Node) string {
		var sb strings.Builder
		for i, n := range ns {
			nodeStr := ""
			switch v := n.(type) {
			case CmdNode:
				nodeStr = v.String()
			case LoopNode:
				nodeStr = "LoopNode\n" + formatBlock(v.Block)
				if i < len(ns)-1 {
					nodeStr = strings.ReplaceAll(nodeStr, "\n", "\n│ ")
				} else {
					nodeStr = strings.ReplaceAll(nodeStr, "\n", "\n  ")
				}
			}
			if i < len(ns)-1 {
				sb.WriteString("├─" + nodeStr + "\n")
			} else {
				sb.WriteString("└─" + nodeStr + "\n")
			}
		}
		s, _ := strings.CutSuffix(sb.String(), "\n")
		return s
	}
	return "AST\n" + formatBlock(t)
}
