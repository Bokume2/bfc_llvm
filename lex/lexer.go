package lex

import "strings"

type Lexer struct {
	TokenDef TokenDef
}

func NewLexer(td TokenDef) Lexer {
	return Lexer{
		TokenDef: td,
	}
}

func (l Lexer) Lex(code string) []Token {
	rest := code
	result := make([]Token, 0, 256)
	for len(rest) > 0 {
		found := false
		for _, v := range l.TokenDef.Orderd() {
			if rest, found = strings.CutPrefix(rest, v); found {
				result = append(result, l.TokenDef.ToToken(v))
				break
			}
		}
		if !found {
			rest = rest[1:]
		}
	}
	return result
}
