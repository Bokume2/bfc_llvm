package parse

import (
	"github.com/Bokume2/bfc_llvm/lex"
)

type Parser struct {
}

func NewParser() Parser {
	return Parser{}
}

func (p Parser) Parse(ts []lex.Token) (AST, error) {
	result := make([]Node, 0, len(ts))
	for i := uint(0); i < uint(len(ts)); i++ {
		switch ts[i] {
		case lex.OpnToken:
			loop, j, err := p.parseLoop(ts, i+1)
			if err != nil {
				return nil, err
			}
			i = j
			result = append(result, loop)
		case lex.ClsToken:
			return nil, NewInvalidLoopError(lex.OpnToken)
		case lex.UnknownToken:
			continue
		default:
			result = append(result, CmdNode{Cmd: ts[i]})
		}
	}
	return result, nil
}

func (p Parser) parseLoop(ts []lex.Token, i uint) (LoopNode, uint, error) {
	buf := make([]Node, 0)
	var empty LoopNode
	for ; i < uint(len(ts)); i++ {
		switch ts[i] {
		case lex.OpnToken:
			loop, j, err := p.parseLoop(ts, i+1)
			if err != nil {
				return empty, 0, err
			}
			i = j
			buf = append(buf, loop)
		case lex.ClsToken:
			return LoopNode{Block: buf}, i, nil
		case lex.UnknownToken:
			continue
		default:
			buf = append(buf, CmdNode{Cmd: ts[i]})
		}
	}
	return empty, 0, NewInvalidLoopError(lex.ClsToken)
}
