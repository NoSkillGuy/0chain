package magmasc

import (
	"sort"
)

// sortedConsumers represents slice of Consumer sorted in alphabetic order by Consumer.ID.
//
// sortedConsumers allows O(logN) access.
type sortedConsumers []*Consumer

func (s sortedConsumers) getIndex(id string) (i int, ok bool) {
	i = sort.Search(len(s), func(i int) bool {
		return s[i].ID >= id
	})
	if i == len(s) {
		return // not found
	}
	if s[i].ID == id {
		return i, true // found
	}
	return // not found
}

func (s sortedConsumers) get(id string) (b *Consumer, ok bool) {
	var i = sort.Search(len(s), func(i int) bool {
		return s[i].ID >= id
	})
	if i == len(s) {
		return // not found
	}
	if s[i].ID == id {
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

func (s *sortedConsumers) add(b *Consumer) (ok bool) {
	if len(*s) == 0 {
		(*s) = append((*s), b)
		return true // added
	}
	var i = sort.Search(len(*s), func(i int) bool {
		return (*s)[i].ID >= b.ID
	})
	// out of bounds
	if i == len(*s) {
		(*s) = append((*s), b)
		return true // added
	}
	// the same
	if (*s)[i].ID == b.ID {
		(*s)[i] = b  // replace
		return false // already have
	}
	// next
	(*s) = append((*s)[:i], append([]*Consumer{b}, (*s)[i:]...)...)
	return true // added
}

// replace if found
func (s *sortedConsumers) update(b *Consumer) (ok bool) {
	var i int
	if i, ok = s.getIndex(b.ID); !ok {
		return
	}
	(*s)[i] = b // replace
	return
}

func (s sortedConsumers) copy() (cp []*Consumer) {
	cp = make([]*Consumer, 0, len(s))
	for _, b := range s {
		cp = append(cp, b)
	}
	return
}
