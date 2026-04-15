package test

import (
	"slices"

	"github.com/Bokume2/bfc_llvm/lex"
	"github.com/Bokume2/bfc_llvm/parse"
)

const TestCode = "++++++[>++++++++<-]>.,"

func ExpectedTokens() []lex.Token {
	return []lex.Token{
		lex.IncToken,
		lex.IncToken,
		lex.IncToken,
		lex.IncToken,
		lex.IncToken,
		lex.IncToken,
		lex.OpnToken,
		lex.NxtToken,
		lex.IncToken,
		lex.IncToken,
		lex.IncToken,
		lex.IncToken,
		lex.IncToken,
		lex.IncToken,
		lex.IncToken,
		lex.IncToken,
		lex.PrvToken,
		lex.DecToken,
		lex.ClsToken,
		lex.NxtToken,
		lex.PutToken,
		lex.GetToken,
	}
}

func ExpectedAST() parse.AST {
	return parse.AST{
		parse.CmdNode{Cmd: lex.IncToken},
		parse.CmdNode{Cmd: lex.IncToken},
		parse.CmdNode{Cmd: lex.IncToken},
		parse.CmdNode{Cmd: lex.IncToken},
		parse.CmdNode{Cmd: lex.IncToken},
		parse.CmdNode{Cmd: lex.IncToken},
		parse.LoopNode{Block: []parse.Node{
			parse.CmdNode{Cmd: lex.NxtToken},
			parse.CmdNode{Cmd: lex.IncToken},
			parse.CmdNode{Cmd: lex.IncToken},
			parse.CmdNode{Cmd: lex.IncToken},
			parse.CmdNode{Cmd: lex.IncToken},
			parse.CmdNode{Cmd: lex.IncToken},
			parse.CmdNode{Cmd: lex.IncToken},
			parse.CmdNode{Cmd: lex.IncToken},
			parse.CmdNode{Cmd: lex.IncToken},
			parse.CmdNode{Cmd: lex.PrvToken},
			parse.CmdNode{Cmd: lex.DecToken},
		}},
		parse.CmdNode{Cmd: lex.NxtToken},
		parse.CmdNode{Cmd: lex.PutToken},
		parse.CmdNode{Cmd: lex.GetToken},
	}
}

func OpnLackTokens() []lex.Token {
	return slices.DeleteFunc(ExpectedTokens(), func(t lex.Token) bool {
		return t == lex.OpnToken
	})
}

func ClsLackTokens() []lex.Token {
	return slices.DeleteFunc(ExpectedTokens(), func(t lex.Token) bool {
		return t == lex.ClsToken
	})
}

const ExpectedOutput = "0"
