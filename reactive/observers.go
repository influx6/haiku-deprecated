package reactive

import "github.com/influx6/flux"

//Reactive returns a new Reactive instance
func Reactive(m Immutable) *Observer {
	return &Observer{
		Reactors: flux.ReactIdentity(),
		data:     m,
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
		r.Send(r.data.Value())
	}
}

//Get returns the internal value
func (r *Observer) Get() interface{} {
	return r.data.Value()
}
