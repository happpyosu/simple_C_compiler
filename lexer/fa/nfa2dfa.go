package fa

type NFA2DFAConvert interface {
	ConvertNFA2DFA(nfa *NFA) *DFA
}

// nfa2dfa :nfa to dfa convertor
type nfa2dfa struct {
	nfa *NFA

	dfaStatusNum int
	DStatusList  []*DStatus
	unmarkedNum  int
	DTran        map[int]map[byte]int
}

func NewNFA2DFAConvert() NFA2DFAConvert {
	return &nfa2dfa{}
}

func initNFA2DFA(nfa *NFA) *nfa2dfa {
	return &nfa2dfa{
		nfa:          nfa,
		dfaStatusNum: 0,
		DStatusList:  make([]*DStatus, 0),
		unmarkedNum:  0,
		DTran:        make(map[int]map[byte]int),
	}
}

// epsilonClosure is the impl of ε-closure(T) function, which is defined as: from any statuses in T, the reachable
// statues that can arrive by jumping over one or several ε.
func (N *nfa2dfa) epsilonClosure(nfa *NFA, ss *StateSet) *StateSet {
	dst := NewStateSet()
	N.recurAddEpsilonClosure(nfa, ss, dst)
	return dst
}

func (N *nfa2dfa) recurAddEpsilonClosure(nfa *NFA, src *StateSet, dst *StateSet) {
	all := src.GetAllStates()
	dst.AddStates(all...)
	for _, s := range all {
		stf := nfa.GetStfOnState(s)
		if stf[eps] != nil && len(stf[eps]) > 0 {
			dst.AddStates(stf[eps]...)
		}
	}

	if src.Equals(dst) {
		return
	}

	src = NewStateSet(dst.GetAllStates()...)
	N.recurAddEpsilonClosure(nfa, src, dst)
}

// move is the impl of move(T, a) function, which is defined to get: from the statuses in T, and give the input a, the
// reachable statues
func (N *nfa2dfa) move(nfa *NFA, T *StateSet, a byte) *StateSet {
	set := NewStateSet()
	all := T.GetAllStates()

	for _, s := range all {
		stf := nfa.GetStfOnState(s)
		if stf[a] != nil && len(stf[a]) > 0 {
			set.AddStates(stf[a]...)
		}
	}
	return set
}

func (N *nfa2dfa) appendDStatusList(d *DStatus) int {
	N.DStatusList = append(N.DStatusList, d)
	if !d.marked {
		N.unmarkedNum++
	}
	N.dfaStatusNum++
	return len(N.DStatusList) - 1
}

func (N *nfa2dfa) getNextUnmarkedDStatus() (*DStatus, int) {
	if N.unmarkedNum > 0 {
		for idx, ds := range N.DStatusList {
			if !ds.marked {
				return ds, int(int64(idx + 1))

			}
		}
	}
	return nil, -1
}

func (N *nfa2dfa) containsDStatus(ss *StateSet) bool {
	for _, s := range N.DStatusList {
		if s.Equals(ss) {
			return true
		}
	}
	return false
}

func (N *nfa2dfa) addDTran(T int, ch byte, U int) {
	if N.DTran[T] == nil {
		N.DTran[T] = make(map[byte]int)
	}
	N.DTran[T][ch] = U
}

func (N *nfa2dfa) toDFA() *DFA {

	charLi := N.nfa.chList // the character set of the NFA

	A := N.epsilonClosure(N.nfa, NewStateSet(1)) // get the initial status set for DFA
	N.appendDStatusList(A.ToDStatus())           // append the initial DFA Status to the DStatusList

	for N.unmarkedNum > 0 { // while there exists an unmarked status T
		T, tIdx := N.getNextUnmarkedDStatus() // obtain the next unmarked DStatus T
		T.mark()                              // mark T
		N.unmarkedNum--
		for _, a := range charLi {
			M := N.move(N.nfa, &T.StateSet, a)
			U := N.epsilonClosure(N.nfa, M)
			if !N.containsDStatus(U) {
				uIdx := N.appendDStatusList(U.ToDStatus())
				N.addDTran(tIdx, a, uIdx)
			}
		}
	}

	termStatus := make([]int, 0)
	for idx, ds := range N.DStatusList {
		for _, ts := range N.nfa.TermStates {
			if ds.Contains(ts) {
				termStatus = append(termStatus, idx)
			}
		}
	}

	return &DFA{
		DStatusNum:     N.dfaStatusNum,
		DStatusList:    N.DStatusList,
		DTran:          N.DTran,
		chList:         N.nfa.chList,
		TerminalStatus: termStatus,
	}
}

func (N *nfa2dfa) ConvertNFA2DFA(nfa *NFA) *DFA {
	return initNFA2DFA(nfa).toDFA()
}
