package magmasc

import (
	"sort"
)

// sortedConsumers represents slice of Consumer sorted in alphabetic order by ID.
//
// sortedConsumers allows O(logN) access.
type sortedConsumers []*Consumer

func (s sortedConsumers) getIndex(id string) (i int, ok bool) {
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

func (s sortedConsumers) get(id string) (b *Consumer, ok bool) {
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

func (s *sortedConsumers) removeByIndex(i int) {
	(*s) = append((*s)[:i], (*s)[i+1:]...)
}

func (s *sortedConsumers) remove(id string) (ok bool) {
	var i int
	if i, ok = s.getIndex(id); !ok {
		return // false
	}
	s.removeByIndex(i)
	return true // removed
}

func (s *sortedConsumers) add(n *Consumer) (ok bool) {
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
	(*s) = append((*s)[:i], append([]*Consumer{n}, (*s)[i:]...)...)
	return true // added
}

// replace if found
func (s *sortedConsumers) update(b *Consumer) (ok bool) {
	var i int
	if i, ok = s.getIndex(b.GetID()); !ok {
		return
	}
	(*s)[i] = b // replace
	return
}

// sortedProviders represents slice of Provider sorted in alphabetic order by ID.
//
// sortedProviders allows O(logN) access.
type sortedProviders []*Provider

func (s sortedProviders) getIndex(id string) (i int, ok bool) {
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

func (s sortedProviders) get(id string) (b *Provider, ok bool) {
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

func (s *sortedProviders) removeByIndex(i int) {
	(*s) = append((*s)[:i], (*s)[i+1:]...)
}

func (s *sortedProviders) remove(id string) (ok bool) {
	var i int
	if i, ok = s.getIndex(id); !ok {
		return // false
	}
	s.removeByIndex(i)
	return true // removed
}

func (s *sortedProviders) add(n *Provider) (ok bool) {
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
	(*s) = append((*s)[:i], append([]*Provider{n}, (*s)[i:]...)...)
	return true // added
}

// replace if found
func (s *sortedProviders) update(b *Provider) (ok bool) {
	var i int
	if i, ok = s.getIndex(b.GetID()); !ok {
		return
	}
	(*s)[i] = b // replace
	return
}
