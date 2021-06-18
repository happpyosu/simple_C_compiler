package lexer

const(
	// operator and separator
	TK_PLUS = iota
	TK_MINUS
	TK_STAR
	TK_DIVIDE
	TK_MOD
	TK_EQ
	TK_NEQ
	TK_LT
	TK_LEQ
	TK_GT
	TK_GEQ
	TK_ASSIGN
	TK_POINTSTO
	TK_DOT
	TK_AND
	TK_OPENPAs
	TK_CLOSEPA
	TK_OPENBR
	TK_CLOSEBR
	TK_BEGIN
	TK_END
	TK_SEMICOLON
	TK_COMMA
	TK_ELLIPSIS
	TK_EOF

	// constant
	TK_CINT
	TK_CCHAR
	TK_CSTR

	// Key word
	KW_CHAR
	KW_SHORT
	KW_INT
	KW_VOID
	KW_STRUCT
	KW_IF
	KW_ELSE
	KW_FOR
	KW_CONTINUE
	KW_BREAK
	KW_RETURN
	KW_SIZEOF
	KW_ALIGN
	KW_CDECL
	KW_STDCALL

	TK_IDENT
)

// token word
type TkWord struct {
	tkCode int // token index
	next *TkWord
	spelling string
}

func elfHash(key string) int {
	h := 0
	var g int
	for _, ch := range key{
		h = (h << 4) + int(ch)
		g = h & 0xf0000000
		if g > 0{
			h ^= g >> 24
		}
		h &= ^g
	}
	return h & MAXKEY
}

const MAXKEY = 1024 // capacity of the word map

// token table operation
type TkTableOperation interface {
	tkWordDirectInsert(word *TkWord) *TkWord	// insert operator, keyword, constant to word map directly.
	tkWorkFind(word string, keyNo int) *TkWord	// find a token word in word map.
	tkWordInsert(word string) *TkWord	// insert indentifier into word map if not exist
}

// token table
type TkTable struct {
	tkHashTable []*TkWord	// word hashmap
	tkTable []*TkWord
}

func InitLex(){

}


func (T *TkTable) tkWordDirectInsert(word *TkWord) *TkWord{
	T.tkTable = append(T.tkTable, word)
	key := elfHash(word.spelling)
	word.next = T.tkHashTable[key]
	T.tkHashTable[key] = word
	return word
}

func (T *TkTable) tkWorkFind(word string, keyNo int) *TkWord{
	var tp *TkWord = nil
	for tp1 := T.tkHashTable[keyNo]; tp1 != nil; tp1 = tp1.next{
		if word == tp1.spelling{
			tp = tp1
		}
	}
	return tp
}

func (T *TkTable) tkWordInsert(p string) *TkWord{
	keyNo := elfHash(p)
	tp := T.tkWorkFind(p, keyNo)
	if nil == tp{
		tp = &TkWord{
			tkCode:   0,
			next:     T.tkHashTable[keyNo],
			spelling: "",
		}
		T.tkHashTable[keyNo] = tp
		T.tkTable = append(T.tkTable, tp)
		tp.tkCode = len(T.tkTable) - 1
		tp.spelling = p
	}
	return tp
}



