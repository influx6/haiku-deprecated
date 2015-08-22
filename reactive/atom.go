package reactive

import (
	"errors"
	"fmt"
	"reflect"
	"sync/atomic"
	"time"
)

type (

	//CoreImmutable defines an interface method rules for immutables types. All types meeting this rule must be single type values
	CoreImmutable interface {
		Value() interface{}
		Mutate(interface{}) (Immutable, bool)
		Allowed(interface{}) bool
		Stamp() time.Time
	}

	//Immutable defines an interface method rules for immutables types. All types meeting this rule must be single type values
	Immutable interface {
		CoreImmutable
		next() Immutable
		previous() Immutable
	}

	//atom provides a base type for golang base types
	atom struct {
		val    interface{}
		kind   reflect.Kind
		dotype bool
		link   bool
		nxt    Immutable
		prv    Immutable
		stamp  time.Time
	}

	//MList represent a list of immutes
	MList struct {
		root Immutable
		tail Immutable
		size int64
	}

	//MIterator represent an iterator for MList
	MIterator struct {
		from    Immutable
		current Immutable
		end     time.Duration
		started int64
		endless bool
	}

	//TimeRange provides a time duration range store
	TimeRange struct {
		Index    int
		Min, Max time.Time
	}

	//MListManager defines the management of MutationLists. Managers use a maxrange for each MList i.e the total size of a linked Immutable list which decides the proportionality of the largness or smallness of the time ranges,so care must be choosen to choose a good range
	MListManager struct {
		stamps   []*TimeRange
		mranges  []*MList
		maxrange int
	}
)

var (
	//ErrEndIndex defines an error when an iterator can move past a range
	ErrEndIndex = errors.New("End Of Index")
	//ErrEventNotFound defines an error when an event range was not found
	ErrEventNotFound = errors.New("Event Range Impossible")
)

//ImmutableChain provides a means of locking the start and end chain between a specific area
func ImmutableChain(r, t Immutable) *MList {
	return &MList{root: r, tail: t}
}

//Immutables returns a new list with an open-ended chain
func Immutables(m Immutable) *MList {
	return ImmutableChain(m, m)
}

//StrictAtom returns a base type for golang base types(int,string,...etc)
func StrictAtom(data interface{}, link bool) *atom {
	return &atom{
		val:    data,
		kind:   reflect.TypeOf(data).Kind(),
		dotype: true,
		link:   link,
		stamp:  time.Now(),
	}
}

//UnstrictAtom returns a base type for golang base types(int,string,...etc)
func UnstrictAtom(data interface{}, link bool) *atom {
	return &atom{
		val:    data,
		kind:   reflect.TypeOf(data).Kind(),
		dotype: false,
		link:   link,
		stamp:  time.Now(),
	}
}

//Type returns the type of kind
func (a *atom) Kind() reflect.Kind {
	return a.kind
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

//AcceptableKind matches the Kind of a type against a interface supplied
func AcceptableKind(ktype reflect.Kind, m interface{}) bool {
	if GetKind(m) == ktype {
		return true
	}
	return false
}

//Allowed changes the value of the atom
func (a *atom) Allowed(m interface{}) bool {
	if a.dotype {
		if AcceptableKind(a.kind, m) {
			return true
		}
		return false
	}
	return true
}

//String returns the value in fmt.formatted fashion
func (a *atom) String() string {
	return fmt.Sprintf("%q", a.val)
}

//Value returns the data
func (a *atom) Value() interface{} {
	return a.val
}

//Mutate returns a new Immutable set to the value if allowed or returns itself if unchanged
func (a *atom) Mutate(v interface{}) (Immutable, bool) {
	if !a.Allowed(v) {
		return a, false
	}
	mc := a.clone()
	mc.val = v
	return mc, true
}

//Stamp returns the time of creation
func (a *atom) Stamp() time.Time {
	return a.stamp
}

func (a *atom) previous() Immutable {
	return a.prv
}

func (a *atom) next() Immutable {
	return a.nxt
}

//Clone returns a new Immutable
func (a *atom) clone() *atom {
	m := &atom{
		val:    a.val,
		kind:   a.kind,
		dotype: a.dotype,
		link:   a.link,
		stamp:  time.Now(),
	}

	if a.link {
		m.prv = a
		m.nxt = a.nxt
		a.nxt = m
	}

	return m
}

//InfiniteEventIterator returns an instance of MIterator
func InfiniteEventIterator(r Immutable) *MIterator {
	return &MIterator{
		from:    r,
		endless: true,
	}
}

//FiniteEventIterator returns an instance of MIterator
func FiniteEventIterator(r Immutable, end time.Duration) *MIterator {
	return &MIterator{
		from:    r,
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

	if m.started < 0 {
		atomic.StoreInt64(&m.started, 1)
		m.current = m.from
		return nil
	}

	nx := m.current.next()

	if !m.endless {
		if time.Duration(nx.Stamp().Unix()) > m.end {
			atomic.StoreInt64(&m.started, 2)
			return ErrEndIndex
		}
	}

	if nx == nil {
		atomic.StoreInt64(&m.started, 2)
		return ErrEndIndex
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
func (m *MListManager) SnapFrom(e Event) (EventIterator, error) {
	for _, tm := range m.stamps {

		if time.Duration(tm.Min.Unix()) < time.Duration(e) {
			continue
		}

		iml := m.mranges[tm.Index]

		return InfiniteEventIterator(iml.Root()), nil
	}

	return nil, ErrEventNotFound
}

//SnapRange returns the MutationList that matches the range from the specified period
func (m *MListManager) SnapRange(s, e Event) (EventIterator, error) {
	for _, tm := range m.stamps {

		if time.Duration(tm.Min.Unix()) < time.Duration(s) {
			continue
		}

		iml := m.mranges[tm.Index]

		return FiniteEventIterator(iml.Root(), time.Duration(e)), nil
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
