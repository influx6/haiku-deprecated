package reactive

import (
	"errors"
	"reflect"
	"time"

	"github.com/influx6/flux"
)

type (

	//Immutable defines an interface method rules for immutables types. All types meeting this rule must be single type values
	Immutable interface {
		next() Immutable
		previous() Immutable
		Value() interface{}
		Mutate(interface{}) (Immutable, bool)
		Allowed(interface{}) bool
		LinkAllowed() bool
		Restricted() bool
		Stamp() time.Time
		AdjustFuture(time.Time)
		Seq() int
		destroy()
	}

	//atom provides a base type for golang base types
	atom struct {
		val    interface{}
		kind   reflect.Kind
		dotype bool
		link   bool
		nxt    Immutable
		prv    Immutable
		stamp  *TimeStamp
		timer  Timer
	}

	//Observer defines a basic reactive value
	Observer struct {
		flux.Reactors
		data Immutable
	}

	//MIterator represent an iterator for MList
	MIterator struct {
		imap    *MutationRange
		current Immutable
		started int64
		endless bool
		reverse bool
	}

	//MutationRange represent a list of immutes
	MutationRange struct {
		root Immutable
		tail Immutable
		size int64
	}

	//ListManager defines the managment of mutation changes and provides a simple interface to query the changes over a span of time range
	ListManager struct {
		mranges  []*MutationRange
		maxsplit int
	}

	//EventIterator provides an iterator for Immutable event lists
	EventIterator interface {
		Reset()
		Next() error
		Reverse() EventIterator
		Event() Immutable
	}

	//ReactorStore provides an interface for storing reactor states
	ReactorStore interface {
		Last() (Immutable, error)
		First() (Immutable, error)
		SnapFrom(time.Duration) (EventIterator, error)
		SnapRange(s, e time.Duration) (EventIterator, error)
		All() (EventIterator, error)
		Mutate(v interface{}) (Immutable, bool)
	}

	//TimeReactor provides a reactor based on time change
	TimeReactor struct {
		flux.Reactors
		store       ReactorStore
		transformer flux.Reactors
		paused      bool
	}
)

var (
	//ErrEmptyList defines when the list is empty
	ErrEmptyList = errors.New("EmptyList")
	//ErrEndIndex defines an error when an iterator can move past a range
	ErrEndIndex = errors.New("End Of Index")
	//ErrEventNotFound defines an error when an event range was not found
	ErrEventNotFound = errors.New("Event Range Impossible")
)

const (
	//ErrUnacceptedTypeMessage defines the message for types that are not part of the basic units/types in go
	ErrUnacceptedTypeMessage = "Type %s is not acceptable"
	sixty                    = 1 * time.Minute
)
