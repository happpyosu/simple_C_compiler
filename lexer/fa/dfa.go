package fa

import "fmt"

// IDFA DFA drive interface
type IDFA interface {
	NextToken(s string, offset int) (string, int)
}

// DStatus is a subset of NFA status, we use a status set for representing the status collection.
type DStatus struct {
	StateSet
	marked bool // whether a DFA status is marked
}

type DFA struct {
	DStatusNum     int                  // DFA status num, the states index are named from 1 ~ stateNum, where 0 indicate no-states
	DStatusList    []*DStatus           // DStatus list, indexing in the range [0, DStatusNum), indicating the origin status inherited from a NFA
	DTran          map[int]map[byte]int // the transfer function of the DFA
	TerminalStatus []int                // terminal status of the DFA
	chList         []byte               // charset of the DFA
}

func (D *DStatus) mark() {
	D.marked = true
}

func (D *DStatus) ToStateSet() *StateSet {
	return &D.StateSet
}

func (D *DFA) PrintDFA() {
	fmt.Printf("Summary of the DFA \n")
	fmt.Printf("Total DFA Status: %d \n", D.DStatusNum)
	fmt.Println()

	fmt.Println("----------------------DFA Status-------------------------")
	for idx, d := range D.DStatusList {
		fmt.Printf("State %d contains NFA states: %v \n", idx+1, d.GetAllStates())
	}
	fmt.Println("----------------------DFA Status-------------------------")
	fmt.Println()

	fmt.Println("----------------------Trans function---------------------")
	fmt.Println(D.DTran)
	fmt.Println("----------------------Trans function---------------------")
	fmt.Println()

	fmt.Println("----------------------Term Status------------------------")
	fmt.Println(D.TerminalStatus)
	fmt.Println("----------------------Term Status------------------------")
}

// NextToken DFA最长匹配算法
func (D *DFA) NextToken(s string, offset int) (string, int) {
	state := 1 // the initial state
	stack := make([]int, 0)

	cur := offset

	for state != 0 {
		c := s[cur]
		cur++

		if D.isTermState(state) {
			stack = make([]int, 0) //clear stack
		}
		stack = append(stack, state)
		state = D.DTran[state][c]
	}

	for !D.isTermState(state) {
		state = stack[len(stack)-1]
		stack = stack[0 : len(stack)-1]
		cur--
	}

	return s[0 : cur+1], cur

}

func (D *DFA) isTermState(state int) bool {
	for _, ts := range D.TerminalStatus {
		if ts == state {
			return true
		}
	}
	return false
}
