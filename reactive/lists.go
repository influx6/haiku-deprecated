package reactive

import (
	"errors"
	"sync/atomic"
	"time"
)

var (
	// ErrEmptyList defines when the list is empty
	ErrEmptyList = errors.New("EmptyList")
	// ErrEndIndex defines an error when an iterator can move past a range
	ErrEndIndex = errors.New("End Of Index")
	// ErrEventNotFound defines an error when an event range was not found
	ErrEventNotFound = errors.New("Event Range Impossible")
)

// MutationRange represent a list of immutes
type MutationRange struct {
	root Immutable
	tail Immutable
	size int64
}

// NewMutationRange provides a means of locking the start and end chain between a specific area
func NewMutationRange(r, t Immutable) *MutationRange {
	return &MutationRange{root: r, tail: t}
}

// Allowed calls the tail allowed function
func (m *MutationRange) Allowed(v interface{}) bool {
	return m.tail.Allowed(v)
}

// Tail returns the tail Immutable for this list
func (m *MutationRange) Tail() Immutable {
	return m.tail
}

// Root returns the root Immutable for this list
func (m *MutationRange) Root() Immutable {
	return m.root
}

// Size returns the size of the list
func (m *MutationRange) Size() int {
	return int(m.size)
}

// Stamp returns the time of creation
func (m *MutationRange) Stamp() time.Time {
	return m.tail.Stamp()
}

// Mutate returns a new Immutable of that type
func (m *MutationRange) Mutate(v interface{}) (Immutable, bool) {
	mc, b := m.tail.Mutate(v)
	if b {
		atomic.AddInt64(&m.size, 1)
	}
	m.tail = mc
	return mc, b
}

// ReactorSearch provides an interface for searching a reactor
type ReactorSearch interface {
	Last() (Immutable, error)
	First() (Immutable, error)
	SnapFrom(time.Duration) (EventIterator, error)
	SnapRange(s, e time.Duration) (EventIterator, error)
	All() (EventIterator, error)
}

// Engine returns a new listmanager tagged only with a specific range for searching
func Engine(mf *MutationRange) ReactorSearch {
	ml := NewListManager(20, nil)
	ml.mranges = append(ml.mranges, mf)
	return ml
}

// SafeReactorStore returns a safe interface that presents allowable methods for the outside world
type SafeReactorStore interface {
	ReactorSearch
	AsEngine() ReactorSearch
	Size() int
}

// ReactorStore provides an interface for storing reactor states
type ReactorStore interface {
	ReactorSearch
	Mutate(v interface{}) (Immutable, bool)
	AsEngine() ReactorSearch
	Size() int
	Empty()
	setInitialMutation(mi Immutable)
}

// ListManager defines the managment of mutation changes and provides a simple interface to query the changes over a span of time range
type ListManager struct {
	mranges  []*MutationRange
	maxsplit int
}

// NewListManager creates a new list manager and uses the provided Immutable as the first mutation if it allows linking else it clones that mutation and sets up the necessary settings but ensures to keep restrictions either on or off accordingly with the provided mutation
func NewListManager(maxr int, mf Immutable) *ListManager {
	mc := &ListManager{
		mranges:  make([]*MutationRange, 0),
		maxsplit: maxr,
	}

	if mf != nil {
		if mf.LinkAllowed() {
			mc.setInitialMutation(mf)
		} else {
			var cm Immutable
			if mf.Restricted() {
				cm = StrictAtom(mf.Value(), true)
			} else {
				cm = UnstrictAtom(mf.Value(), true)
			}
			mc.setInitialMutation(cm)
		}
	}

	return mc
}

func (m *ListManager) setInitialMutation(mi Immutable) {
	fmg := NewMutationRange(mi, mi)
	m.mranges = append(m.mranges, fmg)
}

// Mutate creates a new mutation from the list of mutation,registering and collating it within the manager
func (m *ListManager) Mutate(v interface{}) (Immutable, bool) {
	size := m.Size()

	if size <= 0 {
		fm := UnstrictAtom(v, true)
		m.setInitialMutation(fm)
		return fm, true
	}

	last := m.mranges[size-1]

	if last.Size() >= m.maxsplit {
		mutd, ok := last.tail.Mutate(v)

		if !ok {
			return mutd, ok
		}

		mrg := NewMutationRange(last.tail, mutd)
		m.mranges = append(m.mranges, mrg)

		return mrg.tail, true
	}

	return last.Mutate(v)
}

// MaxRange returns the maximum range per mutation list
func (m *ListManager) MaxRange() int {
	return m.maxsplit
}

// AsEngine return the ListManager as a search only interface without a mutation method
func (m *ListManager) AsEngine() ReactorSearch {
	return m
}

// Size returns the total mutation made within the set range
func (m *ListManager) Size() int {
	return len(m.mranges)
}

// Empty empties the list of immutables
func (m *ListManager) Empty() {
	fm, err := m.First()
	if err == nil {
		fm.destroy()
	}
	m.mranges = m.mranges[:0]
}

// First returns the first immutable
func (m *ListManager) First() (Immutable, error) {
	if m.Size() <= 0 {
		return nil, ErrEmptyList
	}

	last := m.mranges[0]
	return last.Root(), nil
}

// Last returns the current last immutable
func (m *ListManager) Last() (Immutable, error) {
	if m.Size() <= 0 {
		return nil, ErrEmptyList
	}

	last := m.mranges[m.Size()-1]
	return last.Tail(), nil
}

func (m *ListManager) findFrom(e time.Duration) (*MutationRange, error) {
	if m.Size() <= 0 {
		return nil, ErrEmptyList
	}

	var last Immutable
	var err error

	if last, err = m.Last(); err != nil {
		return nil, err
	}

	cms := last.Stamp().Add(e * -1)

	if last.Stamp().Before(cms) {
		return nil, ErrEventNotFound
	}

	var rn *MutationRange

	func() {
		// loopout:
		for _, tm := range m.mranges {

			if tm.Root().Stamp().Before(cms) && !tm.Tail().Stamp().After(cms) {
				continue
			}

			rt := tm.Tail()
			var nxt Immutable

			nxt = rt.previous()

			func() {
			loopm:
				for {

					if nxt == nil {
						break loopm
					}

					if nxt == tm.Root() {
						rn = NewMutationRange(nxt.previous(), last)
						break loopm
					}

					if nxt.Stamp().Before(cms) {
						rn = NewMutationRange(nxt.next(), last)
						break loopm
					}

					nxt = nxt.previous()
				}
			}()

			break

		}
	}()

	if rn == nil {
		return nil, ErrEventNotFound
	}

	return rn, nil
}

// SnapFrom takes a time.Duration aka time.Duration then takes the last mutation and backtracks the total duration return the mutation list iterator from that point
func (m *ListManager) SnapFrom(s time.Duration) (EventIterator, error) {
	mx, err := m.findFrom(s)

	if err != nil {
		return nil, err
	}

	return NewIterator(mx), nil
}

// SnapRange snaps the mutation from a certain point in time and marks an end time range for the iterator
func (m *ListManager) SnapRange(s, e time.Duration) (EventIterator, error) {
	mx, err := m.findFrom(s)

	if err != nil {
		return nil, err
	}

	lms := mx.Root().Stamp().Add(e)

	if mx.Tail().Stamp().Before(lms) {
		return NewIterator(mx), nil
	}

	mxc := mx.Root().next()

	if mxc.Stamp().After(lms) {
		return NewIterator(mx), nil
	}

	for mxc.next() != nil {
		if mxc.Stamp().After(lms) {
			break
		}
		mxc = mxc.next()
	}

	return NewIterator(mx), nil
}

// All returns an iterator for the total mutation set currently in list at that point in time
func (m *ListManager) All() (EventIterator, error) {
	var fs, ls Immutable

	fs, err := m.First()

	if err != nil {
		return nil, err
	}

	ls, err = m.Last()

	if err != nil {
		return nil, err
	}

	return NewIterator(NewMutationRange(fs, ls)), nil
}

// EventIterator provides an iterator for Immutable event lists
type EventIterator interface {
	Reset()
	Next() error
	Reverse()
	IsReversed() bool
	Event() Immutable
}

// MIterator represent an iterator for MList
type MIterator struct {
	imap    *MutationRange
	current Immutable
	started int64
	endless bool
	reverse bool
	isr     bool
}

// NewIterator returns an instance of MIterator
func NewIterator(r *MutationRange) *MIterator {
	return &MIterator{
		imap: r,
	}
}

// NewReverseIterator returns an MIterator with a reverse inclination to its next call
func NewReverseIterator(r *MutationRange) *MIterator {
	return &MIterator{
		imap:    r,
		reverse: true,
		isr:     true,
	}
}

// Event returns the current mutation state
func (m *MIterator) Event() Immutable {
	return m.current
}

// IsReversed returns true/false if the iterator is in reverse mode
func (m *MIterator) IsReversed() bool {
	return m.reverse
}

// Reverse reverses the operation of the iterator
func (m *MIterator) Reverse() {
	if atomic.LoadInt64(&m.started) >= 2 {
		atomic.StoreInt64(&m.started, 0)
	}
	if m.reverse {
		m.reverse = false
	} else {
		m.reverse = true
	}
}

// Next moves to the next item if possible else returns an error of ErrEndIndex
func (m *MIterator) Next() error {
	if atomic.LoadInt64(&m.started) >= 2 {
		return ErrEmptyList
	}

	if m.imap.Tail() == nil && m.imap.Root() == nil {
		return ErrEmptyList
	}

	if m.started <= 0 {
		if m.reverse {
			m.current = m.imap.Tail()
		} else {
			m.current = m.imap.Root()
		}
		atomic.StoreInt64(&m.started, 1)
		return nil
	}

	var nxt Immutable

	if m.current == nil {
		return ErrEndIndex
	}

	if m.reverse {
		nxt = m.current.previous()

		if nxt == nil {
			return ErrEndIndex
		}

		if m.imap.Root() == nxt {
			atomic.StoreInt64(&m.started, 2)
		}

		m.current = nxt
	} else {
		nxt = m.current.next()

		if nxt == nil {
			return ErrEndIndex
		}

		if m.imap.Tail() == nxt {
			atomic.StoreInt64(&m.started, 2)
		}

		m.current = nxt
	}

	return nil
}

// Reset resets the iterator back to the beginning
func (m *MIterator) Reset() {
	m.current = nil
	m.started = 0
	if !m.isr {
		m.reverse = false
	}
}
