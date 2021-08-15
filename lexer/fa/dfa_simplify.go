package fa

import "fmt"

type DFASimplify interface {
	Simplify(dfa *DFA) *SimplifiedDFA
}

type SimplifiedDFA struct {
	*DFA

	SimplifiedStatusNum  int
	SimplifiedStatusList []*StateSet
	SimplifiedTrans      map[int]map[byte]int
}

func NewDFASimplify() DFASimplify {
	return &dfaSimplifier{}
}

type dfaSimplifier struct {
	dfa  *DFA        // DFA to simplify
	pi   []*StateSet // StatusSet list used to store the status has been split
	flag bool        // is the list pi has been changed
}

func initDFASimplify(dfa *DFA) *dfaSimplifier {
	return &dfaSimplifier{
		dfa: dfa,
		pi:  make([]*StateSet, 0),
	}
}

func (D *dfaSimplifier) Simplify(dfa *DFA) *SimplifiedDFA {
	return initDFASimplify(dfa).split().rebuildDFA()
}

func (D *dfaSimplifier) split() *dfaSimplifier {

	D.initialSplit() // do the initial split according to the status belongs to term status or non-term status

	for D.flag { // which collection pi is still changing

		D.flag = false

		for idx, ss := range D.pi {

			hasSplit := false

			for _, ch := range D.dfa.chList {
				ssLi, ok := D.isStatusSetSplittable(ss, idx, ch)
				if ok {
					D.pi[idx] = ssLi[0] // since the statusSet is split, we use the status set in ssLi to replace the origin one
					D.pi = append(D.pi, ssLi[1:]...)
					hasSplit = true
					D.flag = true
					break
				}
			}

			if hasSplit { // since we have split a status set in pi, just break it, and redo the scan until pi not changing
				break
			}
		}
	}

	for _, x := range D.pi {
		x.Print()
	}

	return D
}

func (D *dfaSimplifier) rebuildDFA() *SimplifiedDFA {
	simplifiedDFA := &SimplifiedDFA{
		D.dfa,
		len(D.pi),
		D.pi,
		make(map[int]map[byte]int),
	}

	num := simplifiedDFA.SimplifiedStatusNum
	li := simplifiedDFA.SimplifiedStatusList
	trans := simplifiedDFA.SimplifiedTrans

	for j := 1; j <= num; j++ {
		// since each status in pi behaves the same, we can only get the first member to rebuild the transfer map
		m := li[j-1].GetAllStates()[0]

		for _, ch := range D.dfa.chList {
			if toStatus := D.dfa.DTran[m][ch]; D.dfa.DTran[m][ch] > 0 {

				which := D.inWhichStatusSetInPi(toStatus)
				if which < 0 {
					panic(fmt.Sprintf("Fatal error: status %v belongs to no state set", toStatus))
				}

				if which == j-1 {
					continue
				}

				if nil == trans[j] {
					trans[j] = make(map[byte]int)
				}

				trans[j][ch] = which + 1
			}
		}
	}

	return simplifiedDFA
}

// check whether the ch can divide pi, if true, this function will return the status set list after division.
// Depending on the condition of split operation, the return list may >= 2.
func (D *dfaSimplifier) isStatusSetSplittable(ss *StateSet, ssIndex int, ch byte) ([]*StateSet, bool) {
	allStatus := ss.GetAllStates()

	mp := make(map[int][]int)

	for _, s := range allStatus {
		stfOnS := D.dfa.DTran[s]
		if nil == stfOnS {
			continue
		}
		idx := D.inWhichStatusSetInPi(stfOnS[ch])

		// no transfer on ch, then just the use the index of ss in pi.
		if -1 == idx {
			idx = ssIndex
		}

		if nil == mp[idx] {
			mp[idx] = make([]int, 0)
		}

		mp[idx] = append(mp[idx], s)
	}

	if len(mp) < 2 { // not able to split
		return nil, false
	}

	res := make([]*StateSet, 0)
	for _, v := range mp {
		res = append(res, NewStateSet(v...))
	}

	return res, true
}

// inWhichStatusSetInPi: get the index of a statusSet that contains s
func (D *dfaSimplifier) inWhichStatusSetInPi(s int) int {

	// if no transfer, just return a negative num.
	if 0 == s {
		return -1
	}

	for idx, ss := range D.pi {
		if ss.Contains(s) {
			return idx
		}
	}

	// should never run here
	panic("Fatal error: s belongs to no status set")
}

func (D *dfaSimplifier) containsSubset(ss *StateSet) bool {
	for _, st := range D.pi {
		if !st.ContainsSubset(ss) {
			return false
		}
	}
	return true
}

func (D *dfaSimplifier) initialSplit() {
	I2 := NewStateSet(D.dfa.TerminalStatus...) // the terminal statuses
	nonTermStatus := make([]int, 0)            // the non-term statuses

	for i := 1; i < D.dfa.DStatusNum; i++ {
		if !I2.Contains(i) {
			nonTermStatus = append(nonTermStatus, i)
		}
	}
	I1 := NewStateSet(nonTermStatus...)
	D.pi = append(D.pi, I1, I2) // the initial split operation, only split by term of non-term
	D.flag = true
}

func (S *SimplifiedDFA) Print() {
	fmt.Println("--------------original DFA--------------")
	S.DFA.PrintDFA()
	fmt.Println("-------------simplified DFA-------------")

	fmt.Printf("SimplifiedStatusNum: %v \n", S.SimplifiedStatusNum)
	fmt.Printf("SimplifiedStatusList: \n")

	for _, ssLi := range S.SimplifiedStatusList {
		ssLi.Print()
	}
	fmt.Printf("SimplifiedTrans: %v \n", S.SimplifiedTrans)
}
