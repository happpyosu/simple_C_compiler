package parser

type TopDownParser struct {
	AbstractParser

	//stack []Token
}

func NewTopDownParser(syntax *Syntax, inputTks []Token) *TopDownParser {
	return &TopDownParser{
		AbstractParser: NewAbstractParser(syntax, inputTks),
	}

}

func (T *TopDownParser) parse() bool {
	stack := make([]Token, 0)
	startSymbol := T.Syntax.GetStartSymbol()

	stack = append(stack, startSymbol)

	return T.recurParse(stack)
}

func (T *TopDownParser) recurParse(stack []Token) bool {
	if len(stack) > len(T.InputTokens) {
		return false
	}

	// firstly check weather there is non-term symbols
	hasNonTermSymbol := false
	nonTermIndex := 0
	for index, token := range stack {
		if T.Syntax.NonTermSymbolSet.hasToken(token) {
			hasNonTermSymbol = true
			nonTermIndex = index
			break
		}
	}

	// there still exists non-term symbols, so we have to do the derivations
	if hasNonTermSymbol {
		// check the length, if the len of the stack is larger than the input tokens, we need to directly return
		if len(stack) > len(T.InputTokens) {
			return false
		}

		for i := nonTermIndex; i < len(stack); i++ {
			if T.Syntax.NonTermSymbolSet.hasToken(stack[i]) {
				// do the derivation
				nonTermSymbol := stack[i]
				possibleRightHands := T.Syntax.Derivations[nonTermSymbol]

				for _, derivation := range possibleRightHands {
					newStack := make([]Token, 0)

					newStack = append(newStack, stack[0:i]...)
					newStack = append(newStack, derivation...)
					if i+1 < len(stack) {
						newStack = append(newStack, stack[(i+1):]...)
					}

					if T.recurParse(newStack) {
						return true
					}
				}

			}
		}
	} else {
		// there is no more non-term symbols, we have to direct compare the input symbols and the stack
		if len(stack) != len(T.InputTokens) {
			return false
		}

		for i := 0; i < len(stack); i++ {
			if stack[i] != T.InputTokens[i] {
				return false
			}
		}

		return true
	}

	// should never run here
	return false
}
