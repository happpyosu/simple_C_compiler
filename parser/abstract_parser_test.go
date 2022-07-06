package parser

import "testing"

func GetTestSyntaxAndInput() (*Syntax, []Token) {
	var S Token
	var N Token
	var V Token

	S = 0
	N = 1
	V = 2

	var s Token
	var t Token
	var g Token
	var w Token
	var e Token
	var d Token

	s = 3
	t = 4
	g = 5
	w = 6
	e = 7
	d = 8

	dev := map[Token][][]Token{
		S: {
			[]Token{N, V, N},
		},

		N: {
			[]Token{s},
			[]Token{t},
			[]Token{g},
			[]Token{w},
		},

		V: {
			[]Token{e},
			[]Token{d},
		},
	}

	input := []Token{s, d, w}

	stx := NewSyntax()
	stx.SetStartSymbol(S).SetNonTermSymbols(S, N, V).SetTermSymbols(s, t, g, w, e, d).SetDerivations(dev)

	return stx, input
}

func initAbstractParser() *AbstractParser {
	stx, input := GetTestSyntaxAndInput()
	abstractParser := NewAbstractParser(stx, input)

	return &abstractParser
}

func TestAbstractParser_FirstSet(t *testing.T) {
	ap := initAbstractParser()
	firstSet := ap.FirstSet()
	t.Log(firstSet)
}
