package trees

import (
	"errors"
	"fmt"
	"reflect"
	"sync"

	"github.com/influx6/flux"
	"github.com/influx6/haiku/reactive"
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
	trackers map[string]Observers
	rw       sync.RWMutex
}

// StructTree returns a new tree with a struct setup for reactivity
func StructTree(b interface{}) (DataTrees, error) {
	stb := BuildDataTree(b)
	if err := RegisterStructObservers(stb, b); err != nil {
		return nil, err
	}
	return stb, nil
}

//MapTree returns a new datatree for a map[string]interface{}
func MapTree(b map[string]interface{}) (DataTrees, error) {
	stb := BuildDataTree(b)
	if err := RegisterMapObservers(stb, b); err != nil {
		return nil, err
	}
	return stb, nil
}

//ListTree returns a new datatree for a []interface{}
func ListTree(b []interface{}) (DataTrees, error) {
	stb := BuildDataTree(b)
	if err := RegisterListObservers(stb, b); err != nil {
		return nil, err
	}
	return stb, nil
}

// NewDataTree returns a new instance of datatree
func NewDataTree() *DataTree {
	return BuildDataTree(nil)
}

// BuildDataTree returns a new instance of datatree
func BuildDataTree(sub interface{}) (b *DataTree) {
	b = &DataTree{
		trackers: make(map[string]Observers),
	}

	if sub == nil {
		sub = b
	}

	b.Reactor = flux.FlatAlways(sub)
	return
}

// Track returns the reactor with the fieldname if it exists else return an error
func (b *DataTree) Track(attr string) (reactive.Observers, error) {
	b.rw.RLock()
	bx, ok := b.trackers[attr]
	b.rw.RUnlock()
	if !ok {
		return nil, ErrNotReactor
	}
	return bx, nil
}

// Tracking returns true/false if a field matching the name is being tracked
func (b *DataTree) Tracking(attr string) bool {
	b.rw.RLock()
	_, ok := b.trackers[attr]
	b.rw.RUnlock()
	return ok
}

// HasTracks returns true/false if the tree is being tracked
func (b *DataTree) HasTracks() bool {
	b.rw.RLock()
	defer b.rw.RLock()
	return len(b.trackers) > 0
}

// registerObserver registers a reactor with the tree for change notifications
func (b *DataTree) registerObserver(name string, ob reactive.Observers) {
	var ok bool
	b.rw.RLock()
	_, ok = b.trackers[name]
	b.rw.RUnlock()

	if ok {
		return
	}

	b.rw.Lock()
	b.trackers[name] = ob
	b.rw.Unlock()

	// ob.React(flux.IdentityValueMuxer(b), true)
	ob.Bind(b, true)
}

// ErrSelfRegister is returned when a tree tries to register itself
var ErrSelfRegister = errors.New("DataTree can not register self")

// ErrNotStruct is returned when a interface is not a struct
var ErrNotStruct = errors.New("interface is not a struct kind")

// ErrNotReactor is returned when a interface is not a reactor
var ErrNotReactor = errors.New("interface is not reactive.Observers type")

// RegisterStructObservers takes an interface who's type is a struct and searches within if for any flux.Observers and registers them with a DataTreeRegister to enable self reactivity in the tree
func RegisterStructObservers(tree DataTreeRegister, treeable interface{}) error {
	if tree == treeable {
		return ErrSelfRegister
	}

	mo := reflect.TypeOf(treeable)

	if mo.Kind() == reflect.Ptr {
		mo = mo.Elem()
	}

	if mo.Kind() != reflect.Struct {
		return ErrNotStruct
	}

	return RegisterFieldsObservers(tree, treeable)
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

// RegisterFieldsObservers takes an interface who's type is a struct and searches within if for any flux.Observers and registers them with a DataTreeRegister to enable self reactivity in the tree
func RegisterFieldsObservers(tree DataTreeRegister, treeable interface{}) error {
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

		if fl.Kind() != reflect.Interface && fl.Kind() != reflect.Ptr {
			continue
		}

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
