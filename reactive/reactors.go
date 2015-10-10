package reactive

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/influx6/flux"
)

// ErrTimeCopActive is returned when a TimeReactor has a active timecop
var ErrTimeCopActive = errors.New("TimeCop: TimeReactor is currently being held")

// Time defines the interface for a reactive timeable observer
type Time interface {
	flux.Reactor
	DisableTime()
	EnableTime()
	TimeEnabled() bool
	TimeLord() (*TimeCop, error)
	Empty()
	Store() SafeReactorStore
	rstore() ReactorStore
}

// TimeReactor provides a reactor based on time change
type TimeReactor struct {
	flux.Reactor
	ReactorStore
	storage ReactorStore
	csd     EventIterator
	locd    Immutable
	// currentCop *TimeCop
	ro, ru     sync.Mutex
	paused     bool
	enableTime bool
}

// TimeTransformWith returns a time reactor with the specified immutable as origin
func TimeTransformWith(m Immutable) (t *TimeReactor) {
	tr := TimeReactor{
		Reactor: flux.Reactive(func(r flux.Reactor, err error, signal interface{}) {
			if err != nil {
				r.ReplyError(err)
				return
			}

			if t.enableTime {
				t.ro.Lock()
				mux, ok := t.storage.Mutate(signal)
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

				return
			}

			// t.ru.Lock()
			// po := t.paused
			// t.ru.Unlock()

			// if !po {
			t.Reply(signal)
			// }

		}),
		storage: NewListManager(100, m),
	}

	tr.enableTime = true
	t = &tr
	return
}

// TimeTransform returns a time reactor from a pure flux.Reactor, it will receive the values and turn store them appropriate but it will not be closed when the source is closed, so ensure to close the time reactor also
func TimeTransform(mx flux.Reactor) (t *TimeReactor) {
	tm := TimeTransformWith(nil)
	mx.Bind(tm, false)
	return tm
}

// DisableTime disables the stores ability to track changes
func (t *TimeReactor) DisableTime() {
	t.enableTime = false
}

// TimeEnabled returns true/false if the time feature is allowed
func (t *TimeReactor) TimeEnabled() bool {
	return !!t.enableTime
}

// EnableTime enables the stores ability to track changes
func (t *TimeReactor) EnableTime() {
	t.enableTime = true
}

// Store returns the underline time storage system
func (t *TimeReactor) Store() SafeReactorStore {
	return t.storage
}

// Close calls the close method of the internal reactor and empties the store
func (t *TimeReactor) Close() error {
	t.storage.Empty()
	return t.Reactor.Close()
}

// Empty clears all the stored time mutations
func (t *TimeReactor) Empty() {
	t.storage.Empty()
}

// Resume reconnects back the time-reactor to the change stream
func (t *TimeReactor) resume() {
	if !t.enableTime {
		return
	}
	t.ru.Lock()
	t.paused = false
	// t.currentCop = nil
	t.ru.Unlock()
}

// TimeLord disconnects the time-reactor to the change stream,
func (t *TimeReactor) TimeLord() (*TimeCop, error) {
	return t.TimeLordRange(0, 0)
}

// TimeLordRange disconnects the time-reactor to the change stream, and uses a supplied range the provided durations but if both duration are 0,simply performs the same operation as the TimeLord() function by returning the total beginning of all events to the current time(LamportTime) if the first duration is not 0 but the second is,it uses the SnapFrom of the EventIterator else uses the SnapRange function of the EventIterator. If the TimeReactor is disabled, this will always return an error
func (t *TimeReactor) TimeLordRange(fo, to time.Duration) (*TimeCop, error) {
	if !t.enableTime {
		return nil, ErrLogicLocked
	}

	// if t.currentCop != nil {
	// 	return nil, ErrTimeCopActive
	// }

	//we store the last Immutable before we release a timecop,
	//to ensure we have items to time-travel with
	var ld Immutable

	//lock down the code to ensure we are concurrent when in go-routines
	t.ro.Lock()
	if lr, err := t.storage.Last(); err == nil {
		ld = lr
	}
	t.ro.Unlock()

	if ld == nil {
		return nil, ErrEventNotFound
	}

	var mo EventIterator

	t.ro.Lock()
	if fo == time.Duration(0) && to == time.Duration(0) {
		if moe, err := t.storage.All(); err == nil {
			mo = moe
		}
	}
	if fo != time.Duration(0) && to == time.Duration(0) {
		if moe, err := t.storage.SnapFrom(fo); err == nil {
			mo = moe
		}
	} else {
		if moe, err := t.storage.SnapRange(fo, to); err == nil {
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

	ct := NewCop(t, mo)
	// t.currentCop = ct
	return ct, nil
}

func (t *TimeReactor) rstore() ReactorStore {
	return t.storage
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

// RestoreTime restores the TimeReactor back to normal operations and disconnects itself from the TimeCop
func (tc *TimeCop) RestoreTime() {
	if tc.available() {
		tc.tm.resume()
		tc.tm = nil
		// tc.tm.currentCop = nil
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

// TimeObject provides the combination of observers and time trackers that both respect the law of mutation and time
type TimeObject struct {
	Time
	origin *atom
}

func newTimeObject(v *atom) (t *TimeObject) {
	return &TimeObject{
		Time:   TimeTransformWith(v),
		origin: v,
	}
}

// StrictTime returns a observer with Time abilities which is set to accept a specific type only
func StrictTime(v interface{}, enableTime bool) (t *TimeObject) {
	nmt := newTimeObject(StrictAtom(v, true))
	if !enableTime {
		nmt.DisableTime()
	}
	return nmt
}

// UnstrictTime returns a observer with Time abilities set to accept any type
func UnstrictTime(v interface{}, enableTime bool) (t *TimeObject) {
	nmt := newTimeObject(UnstrictAtom(v, true))
	if !enableTime {
		nmt.DisableTime()
	}
	return nmt
}

// Equals return true/false if the value equals the data
func (r *TimeObject) Equals(n flux.Equaler) bool {
	return r.Get() == n
}

// MarshalJSON returns the json representation of the observer
func (r *TimeObject) MarshalJSON() ([]byte, error) {
	return []byte(`"` + r.String() + `"`), nil
}

// MarshalYAML returns the yaml representation of the observer
func (r *TimeObject) MarshalYAML() (interface{}, error) {
	return r.Get(), nil
}

//UnmarshalJSON provides a json unmarshaller for observer
func (r *TimeObject) UnmarshalJSON(data []byte) error {
	var newval interface{}

	if err := json.NewDecoder(bytes.NewBuffer(data)).Decode(&newval); err != nil {
		return err
	}

	r.Set(newval)
	return nil
}

//UnmarshalYAML provides a yaml unmarshaller for observer
func (r *TimeObject) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var newval interface{}

	if err := unmarshal(&newval); err != nil {
		return err
	}

	r.Set(newval)
	return nil
}

//Set resets the value of the object to the new value if that value is allowed thereby respecting the rules of type and mutability. If the Time features where disabled using Time.DisableTime() then it sets a new allowed value into the original mutation i.e the first mutation created on the timeObject
func (r *TimeObject) Set(ndata interface{}) {
	if r.TimeEnabled() {
		r.Time.Send(ndata)
	} else {
		if r.origin.Allowed(ndata) {
			r.origin.val = ndata
		}
	}
}

// Empty empties and resets the time mutations infomation to the origin value when the time features where active
func (r *TimeObject) Empty() {
	r.Time.Empty()
	r.Time.rstore().setInitialMutation(r.origin)
}

//Get returns the internal value
func (r *TimeObject) Get() interface{} {
	last, err := r.Time.Store().Last()
	if err != nil {
		return nil
	}
	return last.Value()
}

// Close calls the close method of the internal reactor and empties the store
func (r *TimeObject) Close() error {
	err := r.Time.Close()
	r.Time.rstore().setInitialMutation(r.origin)
	return err
}

// String returns the internal string value of the immutable
func (r *TimeObject) String() string {
	return fmt.Sprintf("%v", r.Get())
}
