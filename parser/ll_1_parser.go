package parser

type LL1Parser struct {
	AbstractParser

	parsingTable map[Token]map[Token]int
}

func NewLL1Parser(syntax *Syntax, inputTks []Token) *LL1Parser {
	return &LL1Parser{
		AbstractParser: NewAbstractParser(syntax, inputTks),
		parsingTable:   nil,
	}

}

func (L *LL1Parser) BuiltLL1ParsingTable() map[Token]map[Token]int {
	// the parsing table
	pTable := make(map[Token]map[Token]int)

	// build the first set
	firstSet := L.FirstSet()

	// derivation index
	dIndex := 0

	// traversal the derivations
	for ntk, devs := range L.Syntax.Derivations {

		// init the pTable for the non-term symbol
		pTable[ntk] = make(map[Token]int)

		// the first set for the given non-term symbol
		firstSetOfNtk := firstSet[ntk]

		for range devs {
			for _, ttk := range firstSetOfNtk.toTokenList() {
				pTable[ntk][ttk] = dIndex
			}
			dIndex++
		}

	}

	return pTable
}
