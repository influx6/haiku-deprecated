package react

import "github.com/influx6/flux"

type (

	//Observer defines a basic reactive value
	Observer struct {
		flux.Stacks
		data Immutable
	}
)

//Transform returns a new Reactive instance from an interface
func Transform(m interface{}, chain bool) (*Observer, error) {
	var im Immutable
	var err error

	if im, err = MakeType(m, chain); err != nil {
		return nil, err
	}

	return Reactive(im), nil
}

//Reactive returns a new Reactive instance
func Reactive(m Immutable) *Observer {
	return &Observer{
		Stacks: flux.IdentityStack(),
		data:   m,
	}
}

func (r *Observer) mutate(ndata interface{}) bool {
	clone, done := r.data.Mutate(ndata)

	//can we make the change or his this change proper
	if !done {
		return false
	}

	r.data = clone
	return true
}

//Set resets the value of the object
func (r *Observer) Set(ndata interface{}) {
	if r.mutate(ndata) {
		r.Stacks.Call(r.data.Value())
	}
}

//Identity provides a wrapper over stack.Isolate
func (r *Observer) Identity(ndata interface{}) interface{} {
	if r.mutate(ndata) {
		return r.Stacks.Identity(r.data.Value())
	}
	return r.data.Value()
}

//Isolate provides a wrapper over stack.Isolate
func (r *Observer) Isolate(ndata interface{}) interface{} {
	r.mutate(ndata)
	return r.data.Value()
}

//Apply provides a wrapper over stack.Apply
func (r *Observer) Apply(ndata interface{}) interface{} {
	if r.mutate(ndata) {
		return r.Stacks.Apply(r.data.Value())
		// return r.data.Value()
	}
	return nil
}

//Call provides a wrapper over stack.Call
func (r *Observer) Call(ndata interface{}) interface{} {
	if r.mutate(ndata) {
		r.Stacks.Call(r.data.Value())
		return r.data.Value()
	}
	return nil
}

//Get returns the internal value
func (r *Observer) Get() interface{} {
	return r.data.Value()
}
