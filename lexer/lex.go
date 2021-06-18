package lexer

import (
	"io"
	"os"
)

const EOF = 0

type Lexer struct {
	fs *FileStream  // a lexer will accept byte from a filestream
	lineNum int64
	ch	*byte
}

func NewLexer(path string) (*Lexer, error){
	fs, err := NewFileStream(path)
	if nil != err{
		return nil, err
	}
	return &Lexer{
		fs:      fs,
		lineNum: 0,
		ch:      &fs.Ch,
	}, nil
}

func (L *Lexer) skipWhiteSpace()  {
	for ; L.fs.Ch == ' ' || L.fs.Ch == '\t' || L.fs.Ch == '\r' || L.fs.Ch == '\n'; L.fs.GetChar(){
		if L.fs.Ch == '\n' || L.fs.Ch == '\r'{
			L.lineNum += 1
		}
	}
}

func (L *Lexer) Preprocess()  {
	L.fs.GetChar()
	for true{
		if *L.ch == ' ' || *L.ch == '\t' || *L.ch == '\r' || *L.ch == '\n'{
			L.skipWhiteSpace()
		}else if *L.ch == '/'{
			L.fs.GetChar()
			if *L.ch == '*'{
				L.parseComment()
			}else {
				L.fs.UnGetChar() // if not the pair "/*" we need to put back a char.
			}
		}else {
			break
		}
	}
}

func (L *Lexer) parseComment(){
	L.fs.GetChar()
	for true{
		for true{
			if *L.ch == '\n' || *L.ch == '\t' || *L.ch == '*' || *L.ch == EOF{
				break
			}
			L.fs.GetChar()
		}

		switch *L.ch {
		case '\n', '\t':
			L.lineNum += 1
		case '*':
			L.fs.GetChar()
			if *L.ch == '/'{
				return
			}
		case EOF:
			panic("Not able to Find the paired comment symbols until the EOF")
			return
		}
	}
}



type FileStream struct {
	f *os.File
	Ch byte
}

func NewFileStream(path string) (*FileStream, error) {
	fs := FileStream{}
	err := fs.open(path)
	if nil != err{
		return nil, err
	}
	return &fs, nil
}

type StreamOp interface {
	open(path string) error
	Close() error
	GetChar() byte
	UnGetChar() byte
}

func (F *FileStream) open(path string) error{
	file, err := os.OpenFile(path, os.O_RDONLY, 0600)
	if nil != err{
		return err
	}
	F.f = file
	return nil
}

func (F *FileStream) Close() error{
	if F.f == nil{
		return nil
	}
	return F.f.Close()
}

func (F *FileStream) GetChar() byte {
	if F.f == nil{
		panic("no file is opened, open a file first.")
	}
	buf := make([]byte, 1)
	_, err := F.f.Read(buf)
	if nil != err{
		if err == io.EOF{
			return EOF
		}
		panic(err)
	}
	F.Ch = buf[0]
	return F.Ch
}

func (F *FileStream) UnGetChar() byte{
	// get the current cursor
	cur, err := F.f.Seek(0, 1)
	if nil != err{
		panic("failed to unGetChar from the file. Failed to Seek the file.")
	}
	if cur == 0 {
		return F.Ch
	}else if cur == 1{
		// the current cursor is at the second character, if we unGet a character, the global ch should be reset to 0
		F.Ch = 0
		_, err := F.f.Seek(-1, 1)
		if nil != err {
			panic(err)
		}
	}else {
		_, err := F.f.Seek(-2, 1)
		if nil != err{
			panic(err)
		}
		F.GetChar()
	}
	return F.Ch
}
