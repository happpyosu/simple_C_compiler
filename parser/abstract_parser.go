package parser

type Parser interface {
	parse() bool
}

var Epsilon Token = 0

type Syntax struct {
	StartSymbol      Token     // the start symbol of the Syntactic Analysis
	TermSymbolSet    *TokenSet // the terminal symbol set
	NonTermSymbolSet *TokenSet // the non-terminal symbol set

	SymbolPrintMap map[Token]string // todo: used for print a number-represented token with string

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

	// the derivations of the syntax, which is accessed by using the index
	DerivationsIndex []map[Token][]Token

	// empty str, which will directly initialized as zero
	Epsilon Token
}

func NewEmptySyntax() *Syntax {
	return &Syntax{
		StartSymbol:      0,
		TermSymbolSet:    NewEmptyTokenSet(),
		NonTermSymbolSet: NewEmptyTokenSet(),
		Derivations:      nil,
		DerivationsIndex: nil,
		Epsilon:          0,
	}
}

func NewSyntax(startSymbol Token, nonTerms, terms []Token, derivations map[Token][][]Token) *Syntax {
	syntax := &Syntax{
		StartSymbol:      startSymbol,
		TermSymbolSet:    NewTokenSet(terms...),
		NonTermSymbolSet: NewTokenSet(nonTerms...),
		Derivations:      derivations,
		Epsilon:          0,
	}
	syntax.DerivationsIndex = syntax.initDerivationsIndex()
	return syntax
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

func (S *Syntax) initDerivationsIndex() []map[Token][]Token {
	devsIndex := make([]map[Token][]Token, 0)
	nonTermTokenList := S.NonTermSymbolSet.toTokenList()
	for _, ntk := range nonTermTokenList {
		devs := S.Derivations[ntk]
		for _, dev := range devs {
			devsIndex = append(devsIndex, map[Token][]Token{
				ntk: dev,
			})
		}
	}
	return devsIndex
}

func (S *Syntax) GetDerivationByIndex(index int) map[Token][]Token {
	if nil == S.DerivationsIndex {
		panic("[GetDerivationByIndex]: DerivationsIndex has not been initialized")
	}

	return S.DerivationsIndex[index]
}

func (S *Syntax) GetDerivationNum() int {
	return len(S.DerivationsIndex)
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

// FirstSet compute the first set for the given Syntax
func (A *AbstractParser) FirstSet() map[Token]*TokenSet {
	a := make(map[Token]*TokenSet)
	b := make(map[Token]*TokenSet)

	nullable := A.NullableSet()

	for _, item := range A.Syntax.NonTermSymbolSet.toTokenList() {
		a[item] = NewEmptyTokenSet()
		b[item] = NewEmptyTokenSet()
	}

	for !A.doFirstSetOneStep(a, b, nullable) {
		// deep copy b to a
		a = make(map[Token]*TokenSet)
		for tk, tks := range b {
			a[tk] = tks
		}
	}

	return a
}

// unchanged point method for computing the first set, the function will return true if the token set a equals b.
func (A *AbstractParser) doFirstSetOneStep(a, b map[Token]*TokenSet, nullable *TokenSet) bool {
	for leftToken, dev := range A.Syntax.Derivations {
		tkSet := b[leftToken]
		for _, beta := range dev {

			for _, r := range beta {
				if A.Syntax.TermSymbolSet.hasToken(r) {
					tkSet.addTokens(r)
				} else if A.Syntax.NonTermSymbolSet.hasToken(r) {
					fistSetOfStartTk := b[r]
					tkSet.addTokens(fistSetOfStartTk.toTokenList()...)
					if !nullable.hasToken(r) {
						break
					}
				}
			}

			//startTk := rightToken[0]
			//
			//if A.Syntax.TermSymbolSet.hasToken(startTk) {
			//	tkSet.addTokens(startTk)
			//} else if A.Syntax.NonTermSymbolSet.hasToken(startTk) {
			//	fistSetOfStartTk := b[startTk]
			//	tkSet.addTokens(fistSetOfStartTk.toTokenList()...)
			//}
		}
	}

	for tk, tks := range a {
		tksOfb := b[tk]
		if !tks.equals(tksOfb) {
			return false
		}
	}

	for tk, tks := range b {
		tksOfa := a[tk]
		if !tks.equals(tksOfa) {
			return false
		}
	}

	return true
}

// FirstSetForSentences build the first set for each derivation
func (A *AbstractParser) FirstSetForSentences() map[int]*TokenSet {
	// the derivation nums
	dNum := A.Syntax.GetDerivationNum()

	// build the firstSet for the parser
	fs := A.FirstSet()

	firstS := make(map[int]*TokenSet)
	for i := 0; i < dNum; i++ {
		// get one derivation
		derivation := A.Syntax.GetDerivationByIndex(i)
		var rightDev []Token

		for _, v := range derivation {
			rightDev = v
			break
		}

		// get the first token for this derivation
		firstToken := rightDev[0]

		// if the first token is a non-terminal symbol
		if A.Syntax.NonTermSymbolSet.hasToken(firstToken) {
			temp := fs[firstToken]
			firstS[i] = temp
		} else {
			// otherwise, this is a terminal symbol
			firstS[i] = NewTokenSet(firstToken)
		}
	}

	return firstS
}

// NullableSet computes the nullable token set for the given Syntax
func (A *AbstractParser) NullableSet() *TokenSet {
	a := NewEmptyTokenSet()
	b := NewEmptyTokenSet()

	// while b is still changing...
	for !A.doNullableSetOneStep(a, b) {
		// deep copy b to a
		a = NewEmptyTokenSet()
		a.addTokens(b.toTokenList()...)
	}

	return a
}

func (A *AbstractParser) doNullableSetOneStep(a, b *TokenSet) bool {
	for ntk, devs := range A.Syntax.Derivations {

		for _, dev := range devs {
			// if we directly encounter a null string
			if len(dev) == 1 && dev[0] == A.Syntax.Epsilon {
				b.addTokens(ntk)
			} else {
				// otherwise, we have to check if the right side of the derivation are all belongs to the nullable set
				flag := true
				for _, tk := range dev {
					if !b.hasToken(tk) {
						flag = false
						break
					}
				}

				if flag {
					b.addTokens(ntk)
				}
			}

		}
	}

	if !a.equals(b) {
		return false
	}
	return true
}
