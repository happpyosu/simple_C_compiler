package parser

import "testing"

func TestLL1Parser(t_ *testing.T) {
	stx, inputs := GetTestSyntaxAndInput()

	ll1Parser := NewLL1Parser(stx, inputs)

	t_.Log(ll1Parser.BuiltLL1ParsingTable())
}
