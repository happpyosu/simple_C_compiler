package parser

type Token uint64

// TokenSet maintains the collection of the tokens
type TokenSet struct {
	ts map[Token]bool
}

func NewTokenSet() *TokenSet {
	return &TokenSet{
		ts: make(map[Token]bool),
	}
}

func (T *TokenSet) hasToken(tk Token) bool {
	return T.ts[tk]
}

func (T *TokenSet) addTokens(tks ...Token) {
	for _, tk := range tks {
		T.ts[tk] = true
	}
}

func (T *TokenSet) toTokenList() []Token {
	res := make([]Token, 0)
	for tk := range T.ts {
		res = append(res, tk)
	}
	return res
}

func (T *TokenSet) equals(tks *TokenSet) bool {
	for tk := range tks.ts {
		if !T.ts[tk] {
			return false
		}
	}

	for tk := range T.ts {
		if !tks.ts[tk] {
			return false
		}
	}

	return true
}
