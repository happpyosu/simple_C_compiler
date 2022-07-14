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

func (L *LL1Parser) BuildLL1ParsingTable() map[Token]map[Token]int {
	// the parsing table
	pTable := make(map[Token]map[Token]int)

	// firstSet for
	firstSetForSentences := L.FirstSetForSentences()

	// derivation num
	dNum := L.Syntax.GetDerivationNum()

	for i := 0; i < dNum; i++ {
		// get the first_S for ith dev
		oneS := firstSetForSentences[i]

		// get the ith dev
		oneDev := L.Syntax.GetDerivationByIndex(i)

		var ntk Token
		for k, _ := range oneDev {
			ntk = k
			break
		}

		if nil == pTable[ntk] {
			pTable[ntk] = make(map[Token]int)
		}

		for _, tk := range oneS.toTokenList() {
			pTable[ntk][tk] = i
		}

	}

	return pTable

}
