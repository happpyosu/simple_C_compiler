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
