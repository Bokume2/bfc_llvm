package test

import "github.com/Bokume2/bfc_llvm/lex"

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

const ExpectedOutput = "0"
