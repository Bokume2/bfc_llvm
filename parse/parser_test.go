package parse_test

import (
	"errors"
	"slices"
	"testing"

	"github.com/Bokume2/bfc_llvm/lex"
	"github.com/Bokume2/bfc_llvm/parse"
	"github.com/Bokume2/bfc_llvm/test"
)

func TestParseSuccess(t *testing.T) {
	result, err := parse.NewParser().Parse(test.ExpectedTokens())
	if err != nil {
		t.Fatal("Parse test faild:", err)
	}
	ok := slices.EqualFunc(result, test.ExpectedAST(), func(n1, n2 parse.Node) bool {
		switch v1 := n1.(type) {
		case parse.CmdNode:
			if v2, ok := n2.(parse.CmdNode); ok {
				return v1 == v2
			}
			return false
		case parse.LoopNode:
			if v2, ok := n2.(parse.LoopNode); ok {
				return slices.Equal(v1.Block, v2.Block)
			}
			return false
		default:
			return false
		}
	})
	if !ok {
		t.Log(result)
		t.Fatal("Parse test failed")
	}
}

func TestOpnLackFailure(t *testing.T) {
	_, err := parse.NewParser().Parse(test.OpnLackTokens())
	if err == nil {
		t.Fatal("Parse with lack of Opn test failed")
	}
	if !errors.Is(err, parse.NewInvalidLoopError(lex.OpnToken)) {
		t.Fatal("Parse with lack of Opn test failed with wrong error:", err)
	}
}

func TestClsLackFailure(t *testing.T) {
	_, err := parse.NewParser().Parse(test.ClsLackTokens())
	if err == nil {
		t.Fatal("Parse with lack of Cls test failed")
	}
	if !errors.Is(err, parse.NewInvalidLoopError(lex.ClsToken)) {
		t.Fatal("Parse with lack of Cls test failed with wrong error:", err)
	}
}
