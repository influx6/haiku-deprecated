package reactive

import (
	"errors"
	"reflect"
	"time"

	"github.com/influx6/flux"
)

type (

	// Event provides a detail of time
	Event time.Duration

	//EventIterator provides an iterator for Immutable event lists
	EventIterator interface {
		Reset()
		Next() error
		Previous() error
		Event() Immutable
	}

	//Observer defines a basic reactive value
	Observer struct {
		flux.ReactiveStacks
		data Immutable
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
	}

	//CoreImmutable defines an interface method rules for immutables types. All types meeting this rule must be single type values
	CoreImmutable interface {
		Value() interface{}
		Mutate(interface{}) (Immutable, bool)
		Allowed(interface{}) bool
		Stamp() time.Time
		Seq() int
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
		stamp  *TimeStamp
		timer  Timer
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
		end     time.Time
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

const (
	//ErrUnacceptedTypeMessage defines the message for types that are not part of the basic units/types in go
	ErrUnacceptedTypeMessage = "Type %s is not acceptable"
)
