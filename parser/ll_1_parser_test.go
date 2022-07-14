package parser

import "testing"

func TestLL1Parser(t_ *testing.T) {
	stx, inputs := GetTestSyntaxAndInput()

	ll1Parser := NewLL1Parser(stx, inputs)

	parsingTable := ll1Parser.BuildLL1ParsingTable()
	t_.Log(parsingTable)
}
