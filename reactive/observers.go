package reactive

import (
	"bytes"
	"encoding/json"
	"fmt"

	"gopkg.in/yaml.v2"

	"github.com/influx6/flux"
)

// Observers define an interface for observable mutators without time features
type Observers interface {
	flux.Reactor
	yaml.Unmarshaler
	yaml.Marshaler
	json.Unmarshaler
	json.Marshaler
	Get() interface{}
	Set(interface{})
	String() string
	Equals(flux.Equaler) bool
}

// Observer defines a basic reactive value
type Observer struct {
	flux.Reactor `yaml:"-" json:"-"`
	// flux.Reactor
	data Immutable
}

//ObserveAtom returns a new observer instance from an interface
func ObserveAtom(m interface{}, chain bool) *Observer {
	// return Reactive(StrictAtom(m, chain))
	return ObserveImmutable(Atom(m, chain))
}

//ObserveImmutable returns a new Observer instance
func ObserveImmutable(m Immutable) *Observer {
	return &Observer{
		Reactor: flux.ReactIdentity(),
		data:    m,
	}
}

func (r *Observer) mutate(ndata interface{}) bool {
	clone, done := r.data.Mutate(ndata)

	//can we make the change(do the type fit if we are type strict and is it a different value from the last)
	if !done {
		return false
	}

	r.data = clone
	return true
}

// Equals return true/false if the value equals the data
func (r *Observer) Equals(n flux.Equaler) bool {
	return r.Get() == n
}

// MarshalJSON returns the json representation of the observer
func (r *Observer) MarshalJSON() ([]byte, error) {
	return []byte(`"` + r.String() + `"`), nil
}

// MarshalYAML returns the yaml representation of the observer
func (r *Observer) MarshalYAML() (interface{}, error) {
	return r.data.Value(), nil
}

//UnmarshalJSON provides a json unmarshaller for observer
func (r *Observer) UnmarshalJSON(data []byte) error {
	var newval interface{}

	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&newval); err != nil {
		return err
	}

	r.Set(newval)
	return nil
}

//UnmarshalYAML provides a yaml unmarshaller for observer
func (r *Observer) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var newval interface{}

	if err := unmarshal(&newval); err != nil {
		return err
	}

	r.Set(newval)
	return nil
}

//Set resets the value of the object
func (r *Observer) Set(ndata interface{}) {
	if r.mutate(ndata) {
		r.Send(r.data.Value())
	}
}

//Get returns the internal value
func (r *Observer) Get() interface{} {
	return r.data.Value()
}

// String returns the internal string value of the immutable
func (r *Observer) String() string {
	return fmt.Sprintf("%v", r.data.Value())
}
