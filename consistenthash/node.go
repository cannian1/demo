package consistenthash

// node
type node struct {
	ip  string
	num uint32
}

type nodeSet []node

func (ns nodeSet) Len() int {
	return len(ns)
}

func (ns nodeSet) Less(i int, j int) bool {
	if ns[i].num < ns[j].num {
		return true
	}
	return false
}

func (ns nodeSet) Swap(i int, j int) {
	ns[j], ns[i] = ns[i], ns[j]
}
