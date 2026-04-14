package lex

import (
	"slices"
)

type Token int

const (
	IncToken Token = iota
	DecToken
	PrvToken
	NxtToken
	GetToken
	PutToken
	OpnToken
	ClsToken
	UnknownToken
)

func (t Token) String() string {
	return "Token{" + t.ToString(DefaultTokenDef()) + "}"
}

func (t Token) ToString(td TokenDef) string {
	switch t {
	case IncToken:
		return td.Inc
	case DecToken:
		return td.Dec
	case PrvToken:
		return td.Prv
	case NxtToken:
		return td.Nxt
	case GetToken:
		return td.Get
	case PutToken:
		return td.Put
	case OpnToken:
		return td.Opn
	case ClsToken:
		return td.Cls
	default:
		return ""
	}
}

type TokenDef struct {
	Inc      string
	Dec      string
	Prv      string
	Nxt      string
	Get      string
	Put      string
	Opn      string
	Cls      string
	orderd   [8]string
	isOrderd bool
}

func DefaultTokenDef() TokenDef {
	return TokenDef{
		Inc:      "+",
		Dec:      "-",
		Prv:      "<",
		Nxt:      ">",
		Get:      ",",
		Put:      ".",
		Opn:      "[",
		Cls:      "]",
		orderd:   [8]string{"+", "-", "<", ">", ",", ".", "[", "]"},
		isOrderd: true,
	}
}

func (td *TokenDef) ToToken(t string) Token {
	switch t {
	case td.Inc:
		return IncToken
	case td.Dec:
		return DecToken
	case td.Prv:
		return PrvToken
	case td.Nxt:
		return NxtToken
	case td.Get:
		return GetToken
	case td.Put:
		return PutToken
	case td.Opn:
		return OpnToken
	case td.Cls:
		return ClsToken
	default:
		return UnknownToken
	}
}

func (td *TokenDef) Orderd() [8]string {
	if !td.isOrderd {
		sl := []string{
			td.Inc,
			td.Dec,
			td.Prv,
			td.Nxt,
			td.Get,
			td.Put,
			td.Opn,
			td.Cls,
		}
		slices.SortFunc(sl, func(a, b string) int {
			la := len(a)
			lb := len(b)
			switch {
			case la < lb:
				return -1
			case la == lb:
				return 0
			default:
				return 1
			}
		})
		td.orderd = [8]string(sl[0:8])
		td.isOrderd = true
	}
	return td.orderd
}
