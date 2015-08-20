package reactive

import (
	"fmt"
	"reflect"
)

type (
	//atom provides a base type for golang base types
	atom struct {
		val  interface{}
		kind reflect.Kind
	}
)

//Atom returns a base type for golang base types(int,string,...etc)
func Atom(data interface{}) *atom {
	return &atom{
		val:  data,
		kind: reflect.TypeOf(data).Kind(),
	}
}

//Type returns the type of kind
func (a *atom) Type() reflect.Kind {
	return a.kind
}

//set changes the value of the atom
func (a *atom) set(m interface{}) bool {
	if reflect.TypeOf(m).Kind() == a.kind {
		a.val = m
		return true
	}
	return false
}

//String returns the value in fmt.formatted fashion
func (a *atom) String() string {
	return fmt.Sprintf("%q", a.val)
}

//Value returns the data
func (a *atom) Value() interface{} {
	return a.val
}

//Clone returns a new Immutable
func (a *atom) Clone() Immutable {
	return &atom{
		val:  a.val,
		kind: a.kind,
	}
}
