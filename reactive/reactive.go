package reactive

import (
	"fmt"
	"reflect"
	"time"

	"github.com/influx6/flux"
)

type (

	//Event provides a detail of time
	Event time.Duration

	//EventIterator provides an iterator for Immutable event lists
	EventIterator interface {
		Reset()
		Next() error
		Previous() error
		Event() Immutable
	}

	//ReactorStore provides an interface for storing reactor states
	ReactorStore interface {
		SnapFrom(Event) (EventIterator, error)
		SnapRange(Event, Event) (EventIterator, error)
		All() (EventIterator, error)
		Last() (Immutable, error)
		Mutate(v interface{}) (Immutable, bool)
		ForceSave()
	}

	//ReactorObserver provides an interface for storing reactor states
	ReactorObserver interface {
		ReactorStore
		Get() interface{}
		Set(interface{})
	}

	//Reactor represent the struct type for time-travel
	Reactor struct {
		flux.Stacks
		ReactorStore
	}
)

const (
	//ErrUnacceptedTypeMessage defines the message for types that are not part of the basic units/types in go
	ErrUnacceptedTypeMessage = "Type %s is not acceptable"
)

//MakeType validates accepted types and returns the (Immutable, error)
func MakeType(val interface{}, chain bool) (Immutable, error) {
	switch reflect.TypeOf(val).Kind() {
	case reflect.Struct:
		return nil, fmt.Errorf(ErrUnacceptedTypeMessage, "struct")
	case reflect.Map:
		return nil, fmt.Errorf(ErrUnacceptedTypeMessage, "map")
	case reflect.Array:
		return nil, fmt.Errorf(ErrUnacceptedTypeMessage, "array")
	case reflect.Slice:
		return nil, fmt.Errorf(ErrUnacceptedTypeMessage, "slice")
	}

	return StrictAtom(val, chain), nil
}

//OnlyImmutable returns a stack that vets all data within it is a mutation
func OnlyImmutable(m flux.Stacks) flux.Stacks {
	return m.Stack(func(_ flux.Stacks, data flux.Signal) interface{} {
		mo, ok := data.(Immutable)
		if !ok {
			return nil
		}
		return mo
	}, true)
}

//Identity provides a wrapper over stack.Isolate
func (r *Reactor) Identity(ndata interface{}) interface{} {
	var m Immutable
	var ok bool

	if m, ok = r.ReactorStore.Mutate(ndata); ok {
		return r.Stacks.Identity(m.Value())
	}
	return m.Value()
}

//Isolate provides a wrapper over stack.Isolate
func (r *Reactor) Isolate(ndata interface{}) interface{} {
	m, _ := r.Mutate(ndata)
	return m.Value()
}

//Apply provides a wrapper over stack.Apply
func (r *Reactor) Apply(ndata interface{}) interface{} {
	if m, ok := r.ReactorStore.Mutate(ndata); ok {
		return r.Stacks.Apply(m.Value())
	}
	return nil
}

//Mutate calls the internal .Call function
func (r *Reactor) Mutate(ndata interface{}) (Immutable, bool) {
	data := r.Call(ndata)
	l, _ := r.Last()

	if data == nil {
		return l, false
	}

	return l, true
}

//Call provides a wrapper over stack.Call
func (r *Reactor) Call(ndata interface{}) interface{} {
	if m, ok := r.ReactorStore.Mutate(ndata); ok {
		return r.Stacks.Call(m.Value())
	}
	return nil
}

//Get returns the value of the object
func (r *Reactor) Get() interface{} {
	m, _ := r.Last()
	return m.Value()
}

//Set resets the value of the object
func (r *Reactor) Set(ndata interface{}) {
	_ = r.Call(ndata)
}

//BaseReactor returns a reactor instance
func BaseReactor(max int, m Immutable) *Reactor {
	return &Reactor{
		Stacks:       flux.IdentityStack(),
		ReactorStore: NewManager(m, max),
	}
}

//TypeReactor returns a reactor instance
func TypeReactor(max int, ktype interface{}) *Reactor {
	var rm Immutable
	if ktype == nil {
		rm = UnstrictAtom("", true)
	} else {
		rm = StrictAtom(ktype, true)
	}

	return BaseReactor(max, rm)
}
