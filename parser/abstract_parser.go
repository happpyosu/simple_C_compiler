package parser

type Parser interface {
	parse() bool
}

type Syntax struct {
	StartSymbol      Token     // the start symbol of the Syntactic Analysis
	TermSymbolSet    *TokenSet // the terminal symbol set
	NonTermSymbolSet *TokenSet // the non-terminal symbol set

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
}

func NewSyntax() *Syntax {
	return &Syntax{
		StartSymbol:      0,
		TermSymbolSet:    NewTokenSet(),
		NonTermSymbolSet: NewTokenSet(),
		Derivations:      nil,
	}

}

func (S *Syntax) SetStartSymbol(token Token) *Syntax {
	S.StartSymbol = token
	return S
}

func (S *Syntax) SetTermSymbols(tokes ...Token) *Syntax {
	S.TermSymbolSet.addTokens(tokes...)
	return S
}

func (S *Syntax) SetNonTermSymbols(tokes ...Token) *Syntax {
	S.NonTermSymbolSet.addTokens(tokes...)
	return S
}

func (S *Syntax) SetDerivations(dev map[Token][][]Token) *Syntax {
	S.Derivations = dev
	return S
}

func (S *Syntax) GetStartSymbol() Token {
	return S.StartSymbol
}

type AbstractParser struct {
	Syntax      *Syntax
	InputTokens []Token // the input tokens
}

func (A *AbstractParser) parse() bool {
	panic("[AbstractParser]: the parse method should be implemented by its subclass")
	return false
}

func NewAbstractParser(syntax *Syntax, inputTks []Token) AbstractParser {
	return AbstractParser{
		Syntax:      syntax,
		InputTokens: inputTks,
	}
}

func (A *AbstractParser) SetInputTokens(tokens []Token) *AbstractParser {
	A.InputTokens = tokens
	return A
}

func (A *AbstractParser) FirstSet() map[Token]TokenSet {
	a := make(map[Token]TokenSet)
	b := make(map[Token]TokenSet)
	for _, item := range A.Syntax.NonTermSymbolSet.toTokenList() {
		a[item] = *NewTokenSet()
		b[item] = *NewTokenSet()
	}

	for !A.doFirstSetOneStep(a, b) {
		// deep copy b to a
		a = make(map[Token]TokenSet)
		for tk, tks := range b {
			a[tk] = tks
		}
	}

	return a
}

// unchanged point method for computing the first set, the function will return true if the token set a equals b.
func (A *AbstractParser) doFirstSetOneStep(a, b map[Token]TokenSet) bool {
	for leftToken, dev := range A.Syntax.Derivations {
		tkSet := b[leftToken]
		for _, rightToken := range dev {
			startTk := rightToken[0]

			if A.Syntax.TermSymbolSet.hasToken(startTk) {
				tkSet.addTokens(startTk)
			} else if A.Syntax.NonTermSymbolSet.hasToken(startTk) {
				fistSetOfStartTk := b[startTk]
				tkSet.addTokens(fistSetOfStartTk.toTokenList()...)
			}
		}
	}

	for tk, tks := range a {
		tksOfb := b[tk]
		if !tks.equals(&tksOfb) {
			return false
		}
	}

	for tk, tks := range b {
		tksOfa := a[tk]
		if !tks.equals(&tksOfa) {
			return false
		}
	}

	return true
}

func (A *AbstractParser) getDerivationsNum() int {
	num := 0
	for _, dev := range A.Syntax.Derivations {
		num += len(dev)
	}
	return num
}

func (A *AbstractParser) GetDerivationByIndex(index int) []Token {
	idx := 0
	for _, devs := range A.Syntax.Derivations {
		for _, d := range devs {
			if idx == index {
				return d
			}
			idx++
		}
	}
	panic("[GetDerivationByIndex]: out of range")
	return nil
}
