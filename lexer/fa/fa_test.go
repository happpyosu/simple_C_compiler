package fa

import "testing"

func TestFA(t *testing.T) {
	nfa := &NFA{
		statesNum:  11,
		chList:     []byte{'a', 'b'},
		TermStates: []int{10},
		stf: map[int]map[byte][]int{

			//1
			1: {
				eps: []int{2, 8},
			},

			//2
			2: {
				eps: []int{3, 5},
			},

			//3
			3: {
				'a': []int{4},
			},

			//4
			4: {
				eps: []int{7},
			},

			//5
			5: {
				'b': []int{6},
			},

			//6
			6: {
				eps: []int{7},
			},

			//7
			7: {
				eps: []int{2, 8},
			},

			//8
			8: {
				'a': []int{9},
			},

			//9
			9: {
				'b': []int{10},
			},

			//10
			10: {
				'b': []int{11},
			},

			//11 is the final state
			11: {},
		},
	}

	d := NewNFA2DFAConvert().ConvertNFA2DFA(nfa)

	d.PrintDFA()
}

func TestDFASimplify(t *testing.T) {
	dfa := &DFA{
		DStatusNum: 6,
		DTran: map[int]map[byte]int{
			1: {
				'f': 2,
			},

			2: {
				'e': 3,
				'i': 5,
			},

			3: {
				'e': 4,
			},

			5: {
				'e': 6,
			},
		},
		TerminalStatus: []int{4, 6},
		chList:         []byte{'e', 'f', 'i'},
	}

	NewDFASimplify().Simplify(dfa).Print()

}
