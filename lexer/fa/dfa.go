package fa

import "fmt"

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
	fmt.Printf("Total DFA Status: %d \n", D.DStatusNum)
	fmt.Println("----------------------DFA Status------------------------")
	for idx, d := range D.DStatusList {
		fmt.Printf("State %d contains NFA states: %v \n", idx+1, d.GetAllStates())
	}
	fmt.Println("----------------------Trans function------------------------")
	fmt.Println(D.DTran)
	fmt.Println("----------------------Term Status------------------------")
	fmt.Println(D.TerminalStatus)
}
