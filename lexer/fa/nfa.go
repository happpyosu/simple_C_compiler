package fa

import "fmt"

type NFA struct {
	statesNum  int                    // states num of the NFA, the states index are named from 1 ~ stateNum, where 0 indicate no-states
	chList     []byte                 // possible input of the NFA
	stf        map[int]map[byte][]int // transfer map of the NFA, the length of the map is equal to the length of chList.
	TermStates []int                  // terminal status of the NFA
}

// GetStfOnState : get the status transfer fucntion on a given status
func (N *NFA) GetStfOnState(state int) map[byte][]int {
	if state > N.statesNum {
		panic(fmt.Sprintf("The nfa has only %d states, "+
			"however trying to get the STF on state %d", N.statesNum, state))
	}
	return N.stf[state]
}
