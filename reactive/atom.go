package reactive

import (
	"fmt"
	"reflect"
	"time"
)

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
	if a.nxt != nil {
		return a.nxt.Mutate(v)
	}

	if !a.Allowed(v) {
		return a, false
	}
	mc := a.clone()
	mc.val = v
	return mc, true
}

//Seq returns the particular sequence of the mutation
func (a *atom) Seq() int {
	return a.stamp.Seq
}

//Stamp returns the time of creation
func (a *atom) Stamp() time.Time {
	return a.stamp.Stamp
}

func (a *atom) previous() Immutable {
	return a.prv
}

func (a *atom) next() Immutable {
	return a.nxt
}

//Clone returns a new Immutable
func (a *atom) clone() *atom {
	m := makeAtom(a.val, a.dotype, a.link, a.timer)

	if a.link {
		m.prv = a
		m.nxt = a.nxt
		a.nxt = m
	}

	return m
}

//Type returns the type of kind
func (a *atom) Kind() reflect.Kind {
	return a.kind
}

func makeAtom(data interface{}, dt, link bool, tm Timer) *atom {

	if tm == nil {
		tm = NewLamport(0)
	}

	return &atom{
		val:    data,
		kind:   reflect.TypeOf(data).Kind(),
		dotype: dt,
		link:   link,
		stamp:  tm.GetTime(),
		timer:  tm,
	}
}

//StrictAtom returns a base type for golang base types(int,string,...etc)
func StrictAtom(data interface{}, link bool, m Timer) Immutable {
	return makeAtom(data, true, link, m)
}

//UnstrictAtom returns a base type for golang base types(int,string,...etc)
func UnstrictAtom(data interface{}, link bool, m Timer) Immutable {
	return makeAtom(data, false, link, m)
}
