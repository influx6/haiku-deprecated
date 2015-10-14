package reactive

import (
	"fmt"
	// "reflect"
	"time"
)

//Immutable defines an interface method rules for immutables types. All types meeting this rule must be single type values
type Immutable interface {
	Value() interface{}
	Mutate(interface{}) (Immutable, bool)
	// Allowed(interface{}) bool
	LinkAllowed() bool
	// Restricted() bool
	Stamp() time.Time
	AdjustFuture(time.Time)
	Seq() int
	destroy()
	next() Immutable
	previous() Immutable
	unsetNext()
	unsetPrevious()
}

//atom provides a base type for golang base types
type atom struct {
	val interface{}
	// kind   reflect.Kind
	// dotype bool
	link  bool
	nxt   Immutable
	prv   Immutable
	stamp *TimeStamp
	timer Timer
}

//LinkAllowed returns true/false if this allows mutation links
func (a *atom) LinkAllowed() bool {
	return a.link
}

// //Restricted returns true/false if kind check is on
// func (a *atom) Restricted() bool {
// 	return a.dotype
// }

//AdjustFuture allows a atom timer to be adjusted for the next mutation
func (a *atom) AdjustFuture(ms time.Time) {
	a.timer.AdjustTime(ms)
}

// //Allowed changes the value of the atom
// func (a *atom) Allowed(m interface{}) bool {
// 	// if a.dotype {
// 	// 	if AcceptableKind(a.kind, m) {
// 	// 		return true
// 	// 	}
// 	// 	return false
// 	// }
// 	return true
// }

//String returns the value in fmt.formatted fashion
func (a *atom) String() string {
	return fmt.Sprintf("%+v", a.val)
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

	// if !a.Allowed(v) || a.val == v {
	// 	return a, false
	// }
	if a.val == v {
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

func (a *atom) unsetPrevious() {
	a.prv = nil
}

func (a *atom) unsetNext() {
	a.nxt = nil
}

//Clone returns a new Immutable
func (a *atom) clone() *atom {
	m := MakeAtom(a.val, a.link, a.timer)

	if a.link {
		m.prv = a
		m.nxt = a.nxt
		a.nxt = m
	}

	return m
}

func (a *atom) destroy() {
	if a.nxt != nil && a.nxt != a {
		a.nxt.unsetPrevious()
		a.nxt.destroy()
		a.nxt = nil
	}
	if a.prv != nil && a.prv != a {
		a.prv.unsetNext()
		a.prv.destroy()
		a.prv = nil
	}
}

// //Type returns the type of kind
// func (a *atom) Kind() reflect.Kind {
// 	return a.kind
// }

//MakeAtom provides more control on the creation of an atom
func MakeAtom(data interface{}, link bool, tm Timer) *atom {

	if tm == nil {
		tm = NewLamport(0)
	}

	return &atom{
		val: data,
		// kind:   reflect.TypeOf(data).Kind(),
		// dotype: dt,
		link:  link,
		stamp: tm.GetTime(),
		timer: tm,
	}
}

//Atom returns a base type for golang base types(int,string,...etc)
func Atom(data interface{}, link bool) *atom {
	return MakeAtom(data, link, nil)
}

// //UnstrictAtom returns a base type for golang base types(int,string,...etc)
// func UnstrictAtom(data interface{}, link bool) *atom {
// 	return MakeAtom(data, false, link, nil)
// }
