package parser

type Parser interface {
	parse() bool
}

type AbstractParser struct {
	StartSymbol      Token    // the start symbol of the Syntactic Analysis
	TermSymbolSet    TokenSet // the terminal symbol set
	NonTermSymbolSet TokenSet // the non-terminal symbol set

	// the derivations of the syntax, we use a two-dimension array to represent the derivation relations, for example,
	// consider we have the following syntax:
	// ---------------------
	// S -> NVN
	// N -> s | t | g | w
	// V -> e | d
	// ---------------------
	// the derivations can be expressed as the following map:
	// ---------------------
	//{
	//	S: [N, V, N]
	//	N: [s], [t], [g], [w]
	//	V: [e], [d]
	//}
	Derivations map[Token][][]Token
	InputTokens []Token // the input tokens
}

func (A *AbstractParser) parse() bool {
	panic("[AbstractParser]: the parse method should be implemented by its subclass")
	return false
}

func (A *AbstractParser) GetStartSymbol() Token {
	return A.StartSymbol
}

func NewAbstractParser() AbstractParser {
	return AbstractParser{
		StartSymbol:      0,
		TermSymbolSet:    *NewTokenSet(),
		NonTermSymbolSet: *NewTokenSet(),
		Derivations:      nil,
		InputTokens:      nil,
	}
}

func (A *AbstractParser) SetStartSymbol(token Token) *AbstractParser {
	A.StartSymbol = token
	return A
}

func (A *AbstractParser) SetTermSymbols(tokes ...Token) *AbstractParser {
	A.TermSymbolSet.addTokens(tokes...)
	return A
}

func (A *AbstractParser) SetNonTermSymbols(tokes ...Token) *AbstractParser {
	A.NonTermSymbolSet.addTokens(tokes...)
	return A
}

func (A *AbstractParser) SetDerivations(dev map[Token][][]Token) *AbstractParser {
	A.Derivations = dev
	return A
}

func (A *AbstractParser) SetInputTokens(tokens []Token) *AbstractParser {
	A.InputTokens = tokens
	return A
}
