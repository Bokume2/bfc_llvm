package parse

import "github.com/Bokume2/bfc_llvm/lex"

type InvalidLoopError struct {
	lack lex.Token
}

func NewInvalidLoopError(lack lex.Token) InvalidLoopError {
	return InvalidLoopError{lack: lack}
}

func (ile InvalidLoopError) Error() string {
	return "Not enough " + ile.lack.String()
}
