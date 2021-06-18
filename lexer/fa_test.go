package lexer

import (
	"testing"
)

func TestFA(t *testing.T) {
	nfa := &NFA{
		statesNum: 11,
		chList:    []byte{'a', 'b'},
		stf:       []map[byte][]int{
			//0
			{
				eps: []int{1, 7},
			},

			//1
			{
				eps: []int{2, 4},
			},

			//2
			{
				'a': []int{3},
			},

			//3
			{
				eps: []int{6},
			},

			//4
			{
				'b': []int{5},
			},

			//5
			{
				eps: []int{6},
			},

			//6
			{
				eps: []int{1, 7},
			},

			//7
			{
				'a': []int{8},
			},

			//8
			{
				'b': []int{9},
			},

			//9
			{
				'b': []int{10},
			},

			//10 is the final state
			{
			},
		},
	}

	N := NewNFA2DFA(nfa)
	dfa := N.ToDFA()
	dfa.PrintDFA()
}