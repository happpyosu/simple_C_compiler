package parser

import (
	"fmt"
	"testing"
)

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

	stx := NewSyntax(S, []Token{S, N, V}, []Token{s, t, g, w, e, d}, dev)
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

	for k, v := range firstSet {
		t.Log(fmt.Sprintf("token %v -> %v", k, v))
	}
}

func TestAbstractParser_GetDerivationByIndex(t *testing.T) {
	ap := initAbstractParser()
	dNum := ap.Syntax.GetDerivationNum()
	t.Log(fmt.Sprintf("derivation nums is %d", dNum))

	for i := 0; i < dNum; i++ {
		mps := ap.Syntax.GetDerivationByIndex(i)
		t.Log(mps)
	}

}

func TestAbstractParser_FirstSetForSentences(t *testing.T) {
	ap := initAbstractParser()
	firstSetForSentences := ap.FirstSetForSentences()

	dNum := ap.Syntax.GetDerivationNum()

	for i := 0; i < dNum; i++ {
		temp := firstSetForSentences[i]
		tokenList := temp.toTokenList()
		fmt.Println(tokenList)
	}

}

func initNullableTestSyntaxAndInput() (*Syntax, []Token) {
	// ntk
	var Z Token
	var Y Token
	var X Token

	Z = 1
	Y = 2
	X = 3

	// tk
	var d Token
	var c Token
	var a Token

	d = 4
	c = 5
	a = 6

	// derivations
	dev := map[Token][][]Token{
		Z: {
			{d},
			{X, Y, Z},
		},

		Y: {
			{c},
			{Epsilon},
		},

		X: {
			{Y},
			{a},
		},
	}

	input := []Token{d, c, a}
	stx := NewSyntax(Z, []Token{Z, Y, X}, []Token{d, c, a}, dev)

	return stx, input
}

func TestAbstractParser_NullableSet(t *testing.T) {
	stx, input := initNullableTestSyntaxAndInput()
	abstractParser := NewAbstractParser(stx, input)

	nullableSet := abstractParser.NullableSet()
	t.Log(nullableSet)
}
