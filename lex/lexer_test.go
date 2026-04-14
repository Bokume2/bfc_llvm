package lex_test

import (
	"slices"
	"testing"

	"github.com/Bokume2/bfc_llvm/lex"
	bfcTest "github.com/Bokume2/bfc_llvm/test"
)

func TestLexSuccess(t *testing.T) {
	result := lex.NewLexer(lex.DefaultTokenDef()).Lex(bfcTest.TestCode)
	if !slices.Equal(result, bfcTest.ExpectedTokens) {
		t.Log("Lex result:", result)
		t.Fatal("Lex test failed")
	}
}
