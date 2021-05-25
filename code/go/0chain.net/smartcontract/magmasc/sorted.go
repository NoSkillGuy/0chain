package magmasc

import (
	"sort"
)

// sortedNodes represents slice of Node sorted in alphabetic order by ID.
//
// sortedNodes allows O(logN) access.
type sortedNodes []Node

func (s sortedNodes) getIndex(id string) (i int, ok bool) {
	i = sort.Search(len(s), func(i int) bool {
		return s[i].GetID() >= id
	})
	if i == len(s) {
		return // not found
	}
	if s[i].GetID() == id {
		return i, true // found
	}
	return // not found
}

func (s sortedNodes) get(id string) (b Node, ok bool) {
	var i = sort.Search(len(s), func(i int) bool {
		return s[i].GetID() >= id
	})
	if i == len(s) {
		return // not found
	}
	if s[i].GetID() == id {
		return s[i], true // found
	}
	return // not found
}

func (s *sortedNodes) removeByIndex(i int) {
	(*s) = append((*s)[:i], (*s)[i+1:]...)
}

func (s *sortedNodes) remove(id string) (ok bool) {
	var i int
	if i, ok = s.getIndex(id); !ok {
		return // false
	}
	s.removeByIndex(i)
	return true // removed
}

func (s *sortedNodes) add(n Node) (ok bool) {
	if len(*s) == 0 {
		(*s) = append((*s), n)
		return true // added
	}
	var i = sort.Search(len(*s), func(i int) bool {
		return (*s)[i].GetID() >= n.GetID()
	})
	// out of bounds
	if i == len(*s) {
		(*s) = append((*s), n)
		return true // added
	}
	// the same
	if (*s)[i].GetID() == n.GetID() {
		(*s)[i] = n  // replace
		return false // already have
	}
	// next
	(*s) = append((*s)[:i], append([]Node{n}, (*s)[i:]...)...)
	return true // added
}

// replace if found
func (s *sortedNodes) update(b Node) (ok bool) {
	var i int
	if i, ok = s.getIndex(b.GetID()); !ok {
		return
	}
	(*s)[i] = b // replace
	return
}
