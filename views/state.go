package views

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/influx6/flux"
)

// excessStops is a regexp for matching more than one fullstops in the state address which then gets replaced into a single fullstop
var excessStops = regexp.MustCompile(`.\+`)

// ErrStateNotFound is returned when the state address is inaccurate and a state was not found in the path
var ErrStateNotFound = errors.New("State Not Found")

// ErrInvalidStateAddr is returned when a wrong length or format state address is found
var ErrInvalidStateAddr = errors.New("State Not Found")

// StateStat defines a struct that provides the current state address points and data
type StateStat struct {
	Traj string
	Data interface{}
}

// StateResponse defines the function type used by state in response to a state call
type StateResponse func(*StateStat)

//StateValidator defines a function type used in the state validator process
type StateValidator func(string, *StateStat) bool

// States represent the interface defining a state type
type States interface {
	Active() bool
	Tag() string
	Addr() string
	Engine() *StateEngine
	Activate(*StateStat)
	Acceptable(string, *StateStat) bool
	Deactivate(*StateStat)
	UseActivator(StateResponse) States
	UseDeactivator(StateResponse) States
	OverrideValidator(StateValidator) States
}

// State represents a single state of with a specific tag and address
// where the address is a single piece 'home' item in the '.home.files' style state address
type State struct {
	// atomic bit used to indicate active state or inactive state
	active int64

	// tag represent the identifier key used in a super-state
	tag string

	// stateAddr represent the concrent address point for this state e.g home
	stateAddr string

	// stateTraj represent the state trajectory from its root eg {root.stateTraj}.home where {root.stateTraj} is the roots trajectory
	stateTraj string

	// activator and deactivator provide actions to occur when the state is set to be active
	activator   StateResponse
	deactivator StateResponse

	// validator represents an option argument which also takes part in the validation process of a state in validating its rightness/wrongness from a giving state address
	optionalValidator StateValidator

	//internal engine that allows sub-states from a root state
	engine *StateEngine

	// the parent state this is connected to
	// parent States

	vo, ro, do sync.Mutex
}

// NewState builds a new state with a tag and single address point .eg home or files ..etc
func NewState(tag, addrPoint string) *State {
	ns := State{
		tag:       tag,
		stateAddr: addrPoint,
	}

	ns.engine = BuildStateEngine(&ns)

	return &ns
}

// Tag returns the custom tag of this state unique to it
func (s *State) Tag() string {
	return s.tag
}

// Active returns true/false if this state is active
func (s *State) Active() bool {
	return atomic.LoadInt64(&s.active) == 1
}

// Addr returns the state address point eg home
func (s *State) Addr() string {
	return s.stateAddr
}

// Engine returns the internal nested StateEngine
func (s *State) Engine() *StateEngine {
	return s.engine
}

// UseDeactivator assigns the state a new deactivate respone handler
func (s *State) UseDeactivator(so StateResponse) States {
	s.do.Lock()
	s.deactivator = so
	s.do.Unlock()
	return s
}

// UseActivator assigns the state a new active respone handler
func (s *State) UseActivator(so StateResponse) States {
	s.ro.Lock()
	s.activator = so
	s.ro.Unlock()
	return s
}

// OverrideValidator assigns an validator to perform custom matching of the state
func (s *State) OverrideValidator(so StateValidator) States {
	s.vo.Lock()
	s.optionalValidator = so
	s.vo.Unlock()
	return s
}

// Acceptable checks if the state matches the current point
func (s *State) Acceptable(addr string, so *StateStat) bool {
	fxo := s.optionalValidator

	if fxo == nil {
		if addr == s.stateAddr {
			return true
		}
		return false
	}

	s.vo.Lock()
	state := fxo(addr, so)
	s.vo.Unlock()
	return state
}

// Activate activates the state
func (s *State) Activate(so *StateStat) {
	if s.active > 1 {
		return
	}
	atomic.StoreInt64(&s.active, 1)

	subs := s.engine.diffOnlySubs()

	//activate all the subroot states first so they can
	//do any population they want
	for _, ko := range subs {
		ko.Activate(so)
	}

	s.ro.Lock()
	if s.activator != nil {
		s.activator(so)
	}
	s.ro.Unlock()
}

// Deactivate deactivates the state
func (s *State) Deactivate(so *StateStat) {
	if s.active < 1 {
		return
	}
	atomic.StoreInt64(&s.active, 0)
	s.do.Lock()
	if s.activator != nil {
		s.deactivator(so)
	}
	s.do.Unlock()
}

// StateEngine represents the engine that handles the state machine based operations for state-address based states
type StateEngine struct {
	states *flux.SyncCollector
	owner  States
	curr   States
}

// NewStateEngine returns a new engine with a default empty state
func NewStateEngine() *StateEngine {
	// sa := NewState("", ".")
	return BuildStateEngine(nil)
}

// BuildStateEngine returns a new StateEngine instance set with a particular state as its owner
func BuildStateEngine(s States) *StateEngine {
	es := StateEngine{
		states: flux.NewSyncCollector(),
		owner:  s,
	}
	return &es
}

// AddState adds a new state into the engine with the tag used to identify the state, if the address is a empty string then the address recieves the tag as its value, remember the address is a single address point .eg home or files and not the length of the extend address eg .root.home.files
func (se *StateEngine) AddState(tag, addr string) States {
	if addr == "" {
		addr = tag
	}

	sa := NewState(tag, addr)
	se.states.Set(tag, sa)

	return sa
}

// UseState adds a state into the StateEngine with a specific tag, the state address point is still used in matching against it
func (se *StateEngine) UseState(tag string, s States) States {
	se.states.Set(tag, s)
	return s
}

// ShallowState returns the current state of the engine and not the final state i.e with a state address of '.home.files' from its root, it will return State(:home) object
func (se *StateEngine) ShallowState() States {
	if se.curr == nil {
		return nil
	}

	return se.curr
}

// State returns the current last state of the engine with respect to any nested state that is with the state address of '.home.files', it will return State(:files) object
func (se *StateEngine) State() States {
	co := se.curr

	if co == nil {
		// return se.owner
		return nil
	}

	return co.Engine().State()
}

// Partial renders the partial of the last state of the state address
func (se *StateEngine) Partial(addr string, data interface{}) error {
	points, err := se.prepare(addr)

	if err != nil {
		return err
	}

	stat := &StateStat{
		Traj: fmt.Sprintf(".%s", strings.Join(points, ".")),
		Data: data,
	}

	return se.trajectory(points, stat, true)
}

// All renders the partial of the last state of the state address
func (se *StateEngine) All(addr string, data interface{}) error {
	points, err := se.prepare(addr)

	if err != nil {
		return err
	}

	stat := &StateStat{
		Traj: fmt.Sprintf(".%s", strings.Join(points, ".")),
		Data: data,
	}

	return se.trajectory(points, stat, false)
}

// preparePoints prepares the state address into a list of walk points
func (se *StateEngine) prepare(addr string) ([]string, error) {
	addr = excessStops.ReplaceAllString(addr, ".")
	points := strings.Split(addr, ".")
	polen := len(points)

	//check if the length is below 1 then return appropriately
	if polen < 1 {
		return nil, ErrInvalidStateAddr
	}

	//if the first is an empty string, meaning the '.' root was supplied, then we shift so we just start from the first state point else we ignore and use the list as-is.
	if points[0] == "" {
		points = points[1:]
	}

	return points, nil
}

func (se *StateEngine) diffOnlySubs() []States {
	var subs []States

	se.states.Each(eachState(func(so States, tag string, _ func()) {
		if so.Addr() == "." {
			subs = append(subs, so)
		}
	}))

	return subs
}

func (se *StateEngine) diffOnlyNotSubs() []States {
	var subs []States

	se.states.Each(eachState(func(so States, tag string, _ func()) {
		if so.Addr() != "." {
			subs = append(subs, so)
		}
	}))

	return subs
}

func (se *StateEngine) diffSubs() ([]States, []States) {
	var nosubs, subs []States

	se.states.Each(eachState(func(so States, tag string, _ func()) {
		if so.Addr() == "." {
			subs = append(subs, so)
		} else {
			nosubs = append(nosubs, so)
		}
	}))

	return subs, nosubs
}

// DeactivateAll deactivates all states connected to this engine
func (se *StateEngine) DeactivateAll(dso *StateStat) {
	se.states.Each(eachState(func(so States, tag string, _ func()) {
		so.Deactivate(dso)
	}))
}

// trajectory is the real engine which checks the path and passes down the StateStat to the sub-states and determines wether its a full view or partial view
func (se *StateEngine) trajectory(points []string, so *StateStat, partial bool) error {

	nosubs := se.diffOnlyNotSubs()

	//are we out of points to walk through? if so then fire acive and tell others to be inactive
	if len(points) < 1 {
		//deactivate all the non-subroot states
		for _, ko := range nosubs {
			ko.Deactivate(so)
		}

		// //activate all the subroot states first so they can
		// //do any population they want
		// for _, ko := range subs {
		// 	ko.Activate(so)
		// }

		//then active the root itself,so it gets the result
		if se.owner != nil {
			se.owner.Activate(so)
		}
		return nil
	}

	//cache the first point so we dont loose it
	point := points[0]

	if !se.states.Has(point) {
		return ErrStateNotFound
	}

	state := se.states.Get(point).(States)

	if state == nil {
		for _, ko := range nosubs {
			if sko, ok := ko.(States); ok {
				if sko.Acceptable(point, so) {
					state = sko
					break
				}
			}
		}

		if state == nil {
			return ErrStateNotFound
		}
	}

	//set this state as the current active state
	se.curr = state

	//shift the list one more bit for the points
	points = points[1:]

	//we pass down the points since that will handle the loadup downwards
	err := state.Engine().trajectory(points, so, partial)

	if err != nil {
		return err
	}

	if !partial {
		// //activate all the subroot states first so they can
		// //do any population they want
		// for _, ko := range subs {
		// 	ko.Activate(so)
		// }

		if se.owner != nil {
			se.owner.Activate(so)
		}
	}

	return nil
}

// eachState provides a function decorator for calling ech on a flux.SyncCollector
func eachState(fx func(States, string, func())) flux.StringEachfunc {
	return func(item interface{}, key string, stop func()) {
		if so, ok := item.(States); ok {
			fx(so, key, stop)
		}
	}
}