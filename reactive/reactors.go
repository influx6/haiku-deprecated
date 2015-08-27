package reactive

import (
	"errors"
	"sync"
	"time"

	"github.com/influx6/flux"
)

// ErrTimeCopActive is returned when a TimeReactor has a active timecop
var ErrTimeCopActive = errors.New("TimeCop: TimeReactor is currently being held")

// TimeReactor provides a reactor based on time change
type TimeReactor struct {
	flux.Reactor
	store      ReactorStore
	csd        EventIterator
	locd       Immutable
	currentCop *TimeCop
	ro, ru     sync.Mutex
	paused     bool
}

// TimeTransform returns a time reactor
func TimeTransform(mx flux.Reactor) (t *TimeReactor) {
	tr := TimeReactor{
		Reactor: mx.React(func(r flux.Reactor, err error, signal interface{}) {
			if err != nil {
				r.ReplyError(err)
				return
			}

			t.ro.Lock()
			mux, ok := t.store.Mutate(signal)
			t.ro.Unlock()

			if !ok {
				return
			}

			t.ru.Lock()
			po := t.paused
			t.ru.Unlock()

			if !po {
				t.Reply(mux.Value())
			}

		}, true),
		store: NewListManager(100, nil),
	}

	t = &tr
	return
}

// Resume reconnects back the time-reactor to the change stream
func (t *TimeReactor) resume() {
	t.ru.Lock()
	t.paused = false
	t.currentCop = nil
	t.ru.Unlock()
}

// TimeLord disconnects the time-reactor to the change stream,
func (t *TimeReactor) TimeLord() (*TimeCop, error) {
	return t.TimeLordRange(0, 0)
}

// TimeLordRange disconnects the time-reactor to the change stream, and uses a supplied range
// the provided durations but if both duration are 0,simply performs the same operation as
//  the TimeLord() function by returning the total beginning of all events to the current time(LamportTime)
// if the first duration is not 0 but the second is,it uses the SnapFrom of the EventIterator else uses the SnapRange
// function of the EventIterator
func (t *TimeReactor) TimeLordRange(fo, to time.Duration) (*TimeCop, error) {
	if t.currentCop != nil {
		return nil, ErrTimeCopActive
	}

	//we store the last Immutable before we release a timecop,
	//to ensure we have items to time-travel with
	var ld Immutable

	//lock down the code to ensure we are concurrent when in go-routines
	t.ro.Lock()
	if lr, err := t.store.Last(); err == nil {
		ld = lr
	}
	t.ro.Unlock()

	if ld == nil {
		return nil, ErrEventNotFound
	}

	var mo EventIterator

	t.ro.Lock()
	if fo == time.Duration(0) && to == time.Duration(0) {
		if moe, err := t.store.All(); err == nil {
			mo = moe
		}
	}
	if fo != time.Duration(0) && to == time.Duration(0) {
		if moe, err := t.store.SnapFrom(fo); err == nil {
			mo = moe
		}
	} else {
		if moe, err := t.store.SnapRange(fo, to); err == nil {
			mo = moe
		}
	}
	t.ro.Unlock()

	if mo == nil {
		return nil, ErrEventNotFound
	}

	t.ru.Lock()
	t.paused = true
	t.ru.Unlock()
	t.locd = ld
	//lockdown code again for time copy futuristic time-detro-vento-copulation

	return NewCop(t, mo), nil
}

// TimeCop provides a simple time control for a TimeReactor
type TimeCop struct {
	tm       *TimeReactor
	iterator EventIterator
}

// BuildCop returns a new time cop building a EventIterator from
//  the supplied *MutationRange
func BuildCop(t *TimeReactor, series *MutationRange) *TimeCop {
	return NewCop(t, NewIterator(series))
}

// NewCop returns a new time cop
func NewCop(t *TimeReactor, timeline EventIterator) *TimeCop {
	tc := TimeCop{
		tm:       t,
		iterator: timeline,
	}
	return &tc
}

// ForwardTime moves the timecop into the future one step (weeee-time-travel :) )
func (tc *TimeCop) ForwardTime() {
	if !tc.available() {
		return
	}
	if tc.iterator.IsReversed() {
		// tc.iterator.Next()
		tc.iterator.Reverse()
	}
	if tc.iterator.Next() != nil {
		return
	}
	tc.tm.Reply(tc.iterator.Event().Value())
}

// RestoreTime restores the TimeReactor back to normal operations
func (tc *TimeCop) RestoreTime() {
	if tc.available() {
		tc.tm.resume()
		tc.tm = nil
		tc.iterator = nil
	}
}

// RewindTime moves the timecop into the past one step-backward (noooo-time-travel :) )
func (tc *TimeCop) RewindTime() {
	if !tc.available() {
		return
	}
	if !tc.iterator.IsReversed() {
		tc.iterator.Reverse()
	}
	if tc.iterator.Next() != nil {
		return
	}
	tc.tm.Reply(tc.iterator.Event().Value())
}

func (tc *TimeCop) available() bool {
	if tc.tm == nil {
		return false
	}
	return true
}
