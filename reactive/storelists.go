package reactive

import (
	"reflect"
	"sync/atomic"
	"time"
)

//ImmutableChain provides a means of locking the start and end chain between a specific area
func ImmutableChain(r, t Immutable) *MList {
	return &MList{root: r, tail: t}
}

//Immutables returns a new list with an open-ended chain
func Immutables(m Immutable) *MList {
	return ImmutableChain(m, m)
}

//Allowed calls the tail allowed function
func (m *MList) Allowed(v interface{}) bool {
	return m.tail.Allowed(v)
}

//Tail returns the tail Immutable for this list
func (m *MList) Tail() Immutable {
	return m.tail
}

//Root returns the root Immutable for this list
func (m *MList) Root() Immutable {
	return m.root
}

//Size returns the size of the list
func (m *MList) Size() int {
	return int(m.size)
}

//Stamp returns the time of creation
func (m *MList) Stamp() time.Time {
	return m.tail.Stamp()
}

//Mutate returns a new Immutable of that type
func (m *MList) Mutate(v interface{}) (Immutable, bool) {
	mc, b := m.tail.Mutate(v)
	if b {
		atomic.AddInt64(&m.size, 1)
	}
	m.tail = mc
	return mc, b
}

//GetKind returns the kind of the value
func GetKind(m interface{}) reflect.Kind {
	return reflect.TypeOf(m).Kind()
}

//InfiniteEventIterator returns an instance of MIterator
func InfiniteEventIterator(r Immutable) *MIterator {
	return &MIterator{
		from:    r,
		endless: true,
	}
}

//FiniteEventIterator returns an instance of MIterator
func FiniteEventIterator(r Immutable, end time.Time) *MIterator {
	return &MIterator{
		from:    r,
		end:     end,
		endless: false,
	}
}

//Previous walks to the next Immutable
func (m *MIterator) Previous() error {
	if m.started <= 0 && !m.endless {
		return ErrEndIndex
	}

	if !m.endless && m.current == m.from {
		return ErrEndIndex
	}

	atomic.StoreInt64(&m.started, 1)
	prv := m.current.previous()

	if prv == nil {
		// atomic.StoreInt64(&m.started, 2)
		return ErrEndIndex
	}

	m.current = prv
	return nil
}

//Reset resets the iterator for reuse
func (m *MIterator) Reset() {
	m.current = nil
	m.started = 0
}

//Next walks to the next Immutable
func (m *MIterator) Next() error {
	if int(m.started) > 1 {
		return ErrEndIndex
	}

	var nx Immutable

	if m.started <= 0 {
		atomic.StoreInt64(&m.started, 1)
		m.current = m.from
		nx = m.from
	} else {
		nx = m.current.next()
	}

	if nx == nil {
		atomic.StoreInt64(&m.started, 2)
		return ErrEndIndex
	}

	if !m.endless {

		diff := m.end.UTC().Sub(nx.Stamp().UTC())

		// log.Printf("diff: %s", diff)

		if diff <= 0 {
			atomic.StoreInt64(&m.started, 2)
			return ErrEndIndex
		}
	}

	m.current = nx
	return nil
}

//Event walks to the next Immutable
func (m *MIterator) Event() Immutable {
	if int(m.started) > 1 {
		return nil
	}
	return m.current
}

//NewManager returns a new MListManager instance
func NewManager(m Immutable, mcap int) *MListManager {
	return &MListManager{
		maxrange: mcap,
		mranges:  []*MList{Immutables(m)},
	}
}

//ForceSave saves the current list into the datastore timestamps
func (m *MListManager) ForceSave() {
	sz := len(m.mranges)
	lm := m.mranges[sz-1]

	//get the timestamps so we can store the range
	so, eo := lm.Root().Stamp(), lm.Tail().Stamp()

	tm := &TimeRange{Index: sz, Min: so, Max: eo}

	m.stamps = append(m.stamps, tm)

	//make another mlist,set the tail to be the root of it
	nm := Immutables(lm.Tail())

	m.mranges = append(m.mranges, nm)
}

//Mutate builds on the capability of storing the mutations
func (m *MListManager) Mutate(v interface{}) (Immutable, bool) {
	sz := len(m.mranges)
	lm := m.mranges[sz-1]

	if lm.Size() < m.maxrange {
		return lm.Mutate(v)
	}

	// //get the timestamps so we can store the range
	// so, eo := lm.Root().Stamp(), lm.Tail().Stamp()
	//
	// tm := &TimeRange{Index: sz, Min: so, Max: eo}
	//
	// m.stamps = append(m.stamps, tm)
	//
	// //make another mlist,set the tail to be the root of it
	// nm := Immutables(lm.Tail())
	m.ForceSave()

	return m.mranges[sz].Mutate(v)
}

//SnapFrom returns the MutationList that matches the range from the specified period
func (m *MListManager) SnapFrom(e time.Time) (EventIterator, error) {

	for _, tm := range m.stamps {

		if e.UTC().Sub(tm.Min.UTC()) < 0 {
			continue
		}

		iml := m.mranges[tm.Index]

		return InfiniteEventIterator(iml.Root()), nil
	}

	return nil, ErrEventNotFound
}

//SnapRange returns the MutationList that matches the range from the specified period
func (m *MListManager) SnapRange(s, e time.Time) (EventIterator, error) {
	for _, tm := range m.stamps {

		if s.UTC().Sub(tm.Min.UTC()) < 0 {
			continue
		}

		iml := m.mranges[tm.Index]

		return FiniteEventIterator(iml.Root(), e), nil
	}

	return nil, ErrEventNotFound
}

//All returns an iterator for all immutables or an error
func (m *MListManager) All() (EventIterator, error) {
	if len(m.mranges) <= 0 {
		return nil, ErrEndIndex
	}
	return InfiniteEventIterator((m.mranges[0]).Root()), nil
}

//Last return the last Immutable
func (m *MListManager) Last() (Immutable, error) {
	sz := len(m.mranges)

	if sz <= 0 {
		return nil, ErrEndIndex
	}

	lm := m.mranges[sz-1]
	return lm.Tail(), nil
}
