package lexer

import (
	"fmt"
)

const eps byte = 0 // the epsilon means a empty string.

// StateSet refers to a NFA status set
type StateSet struct {
	mp map[int]bool // a map used for deduplication
}

type NFA struct {
	statesNum int // states num of the NFA, the states index are named from 1 ~ stateNum, where 0 indicate no-states
	chList []byte // possible input of the NFA
	stf []map[byte][]int // transfer map of the NFA, the length of the map is equal to the length of chList.
}

type DFA struct {
	DStatusNum int // DFA status num
	DStatusList []*DStatus // DStatus list, indexing in the range [0, DStatusNum)
	DTran map[int]map[byte]int // the transfer function of the DFA
	TerminalStatus []int // terminal status of the DFA
	chList []byte // charset of the DFA
}

// DStatus is a subset of NFA status, we use a status set for representing the status collection.
type DStatus struct {
	StateSet
	marked bool  // whether a DFA status is marked
}

// nfa2dfa convertor
type NFA2DFA struct {
	nfa *NFA

	dfaStatusNum int
	DStatusList []*DStatus
	unmarkedNum int
	DTran map[int]map[byte]int
}

type DFASimplify struct {
	dfa *DFA
	pi []*StateSet
}

func NewDFASimplify(dfa *DFA) *DFASimplify{
	return &DFASimplify{
		dfa: dfa,
		pi:  make([]*StateSet, 0),
	}
}

func (D *DFASimplify) Simplify(){
	D.initialSplit() // do the initial split according to the status belongs to term status or non-term status

	// TODO
}

func (D *DFASimplify) containsSubset(ss *StateSet) bool{
	for _, st := range D.pi{
		if !st.ContainsSubset(ss){
			return false
		}
	}
	return true
}

func (D *DFASimplify) initialSplit(){
	I2 := NewStateSet(D.dfa.TerminalStatus...) // the terminal statuses
	nonTermStatus := make([]int, 0) // the non-term statuses

	for i := 0; i < D.dfa.DStatusNum; i++{
		if !I2.Contains(i){
			nonTermStatus = append(nonTermStatus, i)
		}
	}
	I1 := NewStateSet(nonTermStatus...)
	D.pi = append(D.pi, I1, I2)
}




func NewDStatusFromStateSet(ss *StateSet) *DStatus{
	return &DStatus{
		StateSet: *ss,
		marked:   false,
	}
}

func NewDStatus(state... int) *DStatus{
	return &DStatus{
		StateSet: *NewStateSet(state...),
		marked: false,
	}
}

func (D *DStatus) mark(){
	D.marked = true
}

func (D *DStatus) ToStateSet() *StateSet{
	return &D.StateSet
}


func (N *NFA) GetStfOnState(state int) map[byte][]int {
	if state > N.statesNum{
		panic(fmt.Sprintf("The nfa has only %d states, " +
			"however trying to get the STF on state %d", N.statesNum, state))
	}
	return N.stf[state]
}





func NewNFA2DFA(nfa *NFA) *NFA2DFA{
	return &NFA2DFA{
		nfa:          nfa,
		dfaStatusNum: 0,
		DStatusList:  make([]*DStatus, 0),
		unmarkedNum:  0,
		DTran:        make(map[int]map[byte]int),
	}
}

func NewStateSet(states ...int) *StateSet {
	set := StateSet{mp: make(map[int]bool)}
	for _, s := range states{
		set.AddStates(s)
	}
	return &set
}

func (S *StateSet) ToDStatus() *DStatus{
	return &DStatus{
		StateSet: *S,
		marked:   false,
	}
}

func (S *StateSet) Equals(set *StateSet) bool{
	mp1 := S.mp
	mp2 := set.mp
	for k := range mp1{
		if mp2[k] == false{
			return false
		}
	}
	for k := range mp2{
		if mp1[k] == false{
			return false
		}
	}
	return true
}

func (S *StateSet) Contains(state int) bool {
	return S.mp[state]
}

func (S *StateSet) Size() int{
	return len(S.mp)
}

func (S *StateSet) GetAllStates() []int{
	arr := make([]int, 0, 16)
	for k, _ := range S.mp{
		arr = append(arr, k)
	}
	return arr
}

func (S *StateSet) AddStates(states ...int) {
	for _, s := range states{
		S.mp[s] = true
	}
}

func (S *StateSet) RmStates(states ...int){
	for _, s := range states{
		S.mp[s] = false
	}
}

func (S *StateSet) ContainsSubset(ss *StateSet) bool{
	for k := range ss.mp{
		if !S.Contains(k){
			return false
		}
	}
	return true
}

// epsilonClosure is the impl of ε-closure(T) function, which is defined as: from any statuses in T, the reachable
// statues that can arrive by jumping over one or several ε.
func (N *NFA2DFA) epsilonClosure(nfa *NFA, ss *StateSet) *StateSet {
	dst := NewStateSet()
	N.recurAddEpsilonClosure(nfa, ss, dst)
	return dst
}

func (N *NFA2DFA) recurAddEpsilonClosure(nfa *NFA, src *StateSet, dst *StateSet){
	all := src.GetAllStates()
	dst.AddStates(all...)
	for _, s := range all{
		stf := nfa.GetStfOnState(s)
		if stf[eps] != nil && len(stf[eps]) > 0{
			dst.AddStates(stf[eps]...)
		}
	}

	if src.Equals(dst){
		return
	}

	src = NewStateSet(dst.GetAllStates()...)
	N.recurAddEpsilonClosure(nfa, src, dst)
}

// move is the impl of move(T, a) function, which is defined to get: from the statuses in T, and give the input a, the
// reachable statues
func (N *NFA2DFA) move(nfa *NFA, T *StateSet, a byte) *StateSet{
	set := NewStateSet()
	all := T.GetAllStates()

	for _, s := range all{
		stf := nfa.GetStfOnState(s)
		if stf[a] != nil && len(stf[a]) > 0{
			set.AddStates(stf[a]...)
		}
	}
	return set
}

func (N *NFA2DFA) appendDStatusList (d *DStatus) int {
	N.DStatusList = append(N.DStatusList, d)
	if !d.marked {
		N.unmarkedNum++
	}
	N.dfaStatusNum++
	return len(N.DStatusList) - 1
}

func (N *NFA2DFA) getNextUnmarkedDStatus() (*DStatus, int) {
	if N.unmarkedNum > 0{
		for idx, ds := range N.DStatusList{
			if !ds.marked{
				return ds, int(int64(idx))

			}
		}
	}
	return nil, -1
}

func (N *NFA2DFA) containsDStatus(ss *StateSet) bool {
	for _, s := range N.DStatusList{
		if s.Equals(ss){
			return true
		}
	}
	return false
}

func (N *NFA2DFA) AddDTran(T int, ch byte, U int){
	if N.DTran[T] == nil{
		N.DTran[T] = make(map[byte]int)
	}
	N.DTran[T][ch] = U
}


func (N *NFA2DFA) ToDFA() *DFA{
	charLi := N.nfa.chList // the character set of the NFA

	A := N.epsilonClosure(N.nfa, NewStateSet(0)) // get the initial status set for DFA
	N.appendDStatusList(A.ToDStatus()) // append the initial DFA Status to the DStatusList


	for N.unmarkedNum > 0{	// while there exists an unmarked status T
		T, tIdx := N.getNextUnmarkedDStatus() // obtain the next unmarked DStatus T
		T.mark() // mark T
		N.unmarkedNum--
		for _, a := range charLi{
			M := N.move(N.nfa, &T.StateSet, a)
			U := N.epsilonClosure(N.nfa, M)
			if !N.containsDStatus(U){
				uIdx := N.appendDStatusList(U.ToDStatus())
				N.AddDTran(tIdx, a, uIdx)
			}
		}
	}

	return &DFA{
		DStatusNum:  N.dfaStatusNum,
		DStatusList: N.DStatusList,
		DTran:       N.DTran,
		chList: 	 N.nfa.chList,
	}
}

func (D *DFA) PrintDFA(){
	fmt.Printf("Total DFA Status: %d \n", D.DStatusNum)
	for idx, d := range D.DStatusList{
		fmt.Printf("State %d contains NFA states: %v \n", idx, d.GetAllStates())
	}
}