package parser

import "sort"

type Token uint64

// TokenSet maintains the collection of the tokens
type TokenSet struct {
	orderedList []Token // used for save the ordered token

	ts map[Token]bool // used for fast indexing
}

func NewEmptyTokenSet() *TokenSet {
	return &TokenSet{
		orderedList: make([]Token, 0),
		ts:          make(map[Token]bool),
	}
}

func NewTokenSet(tks ...Token) *TokenSet {
	set := NewEmptyTokenSet()
	set.addTokens(tks...)
	return set
}

func (T *TokenSet) hasToken(tk Token) bool {
	return T.ts[tk]
}

func (T *TokenSet) addTokens(tks ...Token) {
	//todo: maintain the orderedList
	for _, tk := range tks {
		if T.ts[tk] {
			continue
		}
		T.ts[tk] = true
		T.orderedList = append(T.orderedList, tk)
	}

	// maintain the token set as asc ranking
	sort.Slice(T.orderedList, func(i, j int) bool {
		return T.orderedList[i] < T.orderedList[j]
	})
}

func (T *TokenSet) toTokenList() []Token {
	return T.orderedList
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
