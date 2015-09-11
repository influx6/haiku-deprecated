package rui

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/influx6/flux"
	"github.com/influx6/prox/reactive"
)

//Reactive data components that are able to react to changes within the given fields they have so an action can be initiated

// ReactorType is the reflect.TypeOf value of the flux.Reactor interface
var ReactorType = reflect.TypeOf((*reactive.Observers)(nil)).Elem()

// DataTreeRegister provides a interface that defines a registering method for datatrees
type DataTreeRegister interface {
	registerObserver(string, reactive.Observers)
}

// DataTrees define a simple datatree type
type DataTrees interface {
	flux.Reactor
	DataTreeRegister
	Track(string) (reactive.Observers, error)
	Tracking(string) bool
	HasTracks() bool
}

// DataTree reprsent a base struct for reactivity of which other structs compose to allow reactive data behaviours
type DataTree struct {
	//Reactor for the tree that emits itself everytime a child Reactor changes
	flux.Reactor `yaml:"-" json:"-"`
	//dirties contain a auto-generated list of field names that have indeed become dirty/received and accepted changes
	trackers map[string]reactive.Observers
	// ro sync.RWMutex
}

// NewDataTree returs a new instance of datatree
func NewDataTree() *DataTree {
	dt := DataTree{
		Reactor:  flux.ReactIdentity(),
		trackers: make(map[string]reactive.Observers),
	}
	return &dt
}

// Track returns the reactor with the fieldname if it exists else return an error
func (b *DataTree) Track(attr string) (reactive.Observers, error) {
	bx, ok := b.trackers[attr]
	if !ok {
		return nil, ErrNotReactor
	}
	return bx, nil
}

// Tracking returns true/false if a field matching the name is being tracked
func (b *DataTree) Tracking(attr string) bool {
	_, ok := b.trackers[attr]
	return ok
}

// HasTracks returns true/false if the tree is being tracked
func (b *DataTree) HasTracks() bool {
	return len(b.trackers) > 0
}

// registerObserver registers a reactor with the tree for change notifications
func (b *DataTree) registerObserver(name string, ob reactive.Observers) {
	if _, ok := b.trackers[name]; ok {
		return
	}

	b.trackers[name] = ob

	ob.React(func(r flux.Reactor, err error, _ interface{}) {
		if err != nil {
			b.SendError(err)
			return
		}
		b.Send(b)
	}, true)
}

// ErrSelfRegister is returned when a tree tries to register itself
var ErrSelfRegister = errors.New("DataTree can not register self")

// ErrNotReactor is returned when a interface is not a reactor
var ErrNotReactor = errors.New("interface is not reactive.Observers type")

// RegisterReflectWith registers the name and reflect.Value if its a flux.Reactor with a DataTree
func RegisterReflectWith(tree DataTreeRegister, name string, rot reflect.Value) error {

	if rot.Interface() == tree {
		return ErrSelfRegister
	}

	// rot := reflect.ValueOf(data)
	if rot.Kind() == reflect.Ptr {
		rot = rot.Elem()
	}

	if !rot.Type().Implements(ReactorType) {
		return ErrNotReactor
	}

	rcfl := rot.Elem().Interface().(reactive.Observers)
	tree.registerObserver(name, rcfl)
	return nil
}

// RegisterStructObservers takes an interface who's type is a struct and searches within if for any flux.Observers and registers them with a DataTreeRegister to enable self reactivity in the tree
func RegisterStructObservers(tree DataTreeRegister, treeable interface{}) error {
	if tree == treeable {
		return ErrSelfRegister
	}

	rot := reflect.ValueOf(treeable)

	if rot.Kind() == reflect.Ptr {
		rot = rot.Elem()
	}

	rotto := rot.Type()
	for i := 0; i < rot.NumField(); i++ {
		//get the field
		fl := rot.Field(i)
		//get the type field from the struct
		flo := rotto.Field(i)

		// since the kind is always indescriminate we cant use it
		// if fl.Kind() != reflect.Struct {
		// 	continue
		// }

		if fl.Elem().Interface() == tree {
			continue
		}

		if !fl.Type().Implements(ReactorType) {
			continue
		}

		rcfl := fl.Elem().Interface().(reactive.Observers)
		tree.registerObserver(flo.Name, rcfl)
	}

	return nil
}

// RegisterListObservers registers a slice/array elements where the elements are flux.Reactors with a DataTree,all indexes are stringed,so if you want 1 do "1"
func RegisterListObservers(tree DataTreeRegister, list []interface{}) error {
	for id, target := range list {
		if target == tree {
			continue
		}

		fl, ok := target.(reactive.Observers)

		if !ok {
			continue
		}

		tree.registerObserver(fmt.Sprintf("%d", id), fl)
	}
	return nil
}

// RegisterMapObservers registers a slice/array elements where the elements are flux.Reactors with a DataTree
func RegisterMapObservers(tree DataTreeRegister, dlist map[string]interface{}) error {
	for id, target := range dlist {
		if target == tree {
			continue
		}

		fl, ok := target.(reactive.Observers)

		if !ok {
			continue
		}

		tree.registerObserver(id, fl)
	}
	return nil
}
