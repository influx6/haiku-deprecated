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
func TimeTransform(mix flux.ReactiveStacks) (t *TimeReactor) {
	rc := flux.Reactive(func(s flux.ReactiveStacks) {
	loop:
		for {
			select {
			case <-s.Closed():
				break loop
			case <-s.Feed().Closed():
				s.End()
				break loop
			case data := <-s.In():
				if !t.paused {
					t.store.Mutate(data)
					flux.GoDefer("TimeTransform-Delivery", func() {
						s.Out() <- data
					})
				} else {
					flux.GoDefer("TimeTransform-AfterPause-Delivery", func() {
						s.Out() <- data
					})
				}
			case data := <-s.Feed().Out():
				t.store.Mutate(data)
				if !t.paused {
					flux.GoDefer("TimeTransform-Delivery", func() {
						s.Out() <- data
					})
				}
			}
		}
	}, mix)

	t = &TimeReactor{
		ReactiveStacks: rc,
		store:          NewListManager(100, nil),
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
