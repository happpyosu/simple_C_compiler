package fa

import "fmt"

// StateSet refers to a status set that holds a status collection
type StateSet struct {
	mp map[int]bool // a map used for deduplication
}

func NewStateSet(states ...int) *StateSet {
	set := StateSet{mp: make(map[int]bool)}
	for _, s := range states {
		set.AddStates(s)
	}
	return &set
}

func (S *StateSet) Equals(set *StateSet) bool {
	mp1 := S.mp
	mp2 := set.mp
	for k := range mp1 {
		if mp2[k] == false {
			return false
		}
	}
	for k := range mp2 {
		if mp1[k] == false {
			return false
		}
	}
	return true
}

func (S *StateSet) Contains(state int) bool {
	return S.mp[state]
}

func (S *StateSet) Size() int {
	return len(S.mp)
}

func (S *StateSet) GetAllStates() []int {
	arr := make([]int, 0, 16)
	for k := range S.mp {
		arr = append(arr, k)
	}
	return arr
}

func (S *StateSet) AddStates(states ...int) {
	for _, s := range states {
		S.mp[s] = true
	}
}

func (S *StateSet) RmStates(states ...int) {
	for _, s := range states {
		S.mp[s] = false
	}
}

func (S *StateSet) ContainsSubset(ss *StateSet) bool {
	for k := range ss.mp {
		if !S.Contains(k) {
			return false
		}
	}
	return true
}

func (S *StateSet) ToDStatus() *DStatus {
	return &DStatus{
		StateSet: *S,
		marked:   false,
	}
}

func (S *StateSet) Print() {
	fmt.Println(S.GetAllStates())
}
