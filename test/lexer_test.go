package test

import (
	"fmt"
	"testing"
)
import "simple_C_compiler/lexer"

func TestFileStream(t *testing.T) {
	fs, err := lexer.NewFileStream("./src.txt")
	if err == nil{
		fmt.Println(fs.GetChar())
		fmt.Println(fs.GetChar())
		fmt.Println(fs.GetChar())
		fmt.Println(fs.GetChar())
		fmt.Println(fs.GetChar())
	}

}
