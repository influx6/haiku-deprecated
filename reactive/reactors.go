package reactive

import "github.com/influx6/flux"

//ObserveTransform returns a new Reactive instance from an interface
func ObserveTransform(m interface{}, chain bool) (*Observer, error) {
	var im Immutable
	var err error

	if im, err = MakeType(m, chain); err != nil {
		return nil, err
	}

	return Reactive(im), nil
}

//TimeTransform returns a time reactor
func TimeTransform(mx flux.Reactors) (t *TimeReactor) {

	proc := func(signal flux.Signal) flux.Signal {
		mux, ok := t.store.Mutate(signal)
		if ok && !t.paused {
			return mux
		}
		return nil
	}

	stream := mx.React(flux.DataReactProcessor(proc))

	t = &TimeReactor{
		Reactors:    flux.ReactIdentity(),
		store:       NewListManager(100, nil),
		transformer: stream,
	}

	return t
}

//Resume reconnects back the time-reactor to the change stream
func (t *TimeReactor) Resume() {
	t.paused = false
}

//Pause disconnects the time-reactor to the change stream
func (t *TimeReactor) Pause() {
	t.paused = true
}
