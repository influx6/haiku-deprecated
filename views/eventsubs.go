package views

import (
	"errors"
	"fmt"
	"strings"
	"sync"
)

// EventSubs defines the event interface methods for meta info
type EventSubs interface {
	FlatChains
	ID() string
	EventType() string
	EventSelector() string
	SetStopPropagation(n bool)
	SetStopImmediatePropagation(n bool)
	SetPreventDefault(n bool)
	StopPropagation() bool
	StopImmediatePropagation() bool
	PreventDefault() bool
}

// EventSub represent a single event configuration for dom.Elem objects
// instance which allows chaining of events listeners like middleware
type EventSub struct {
	FlatChains
	// Type is the event type to use
	Type string
	//Target is a selector value for matching a event
	Target                 string
	stopPropagate          bool
	stopImmediatePropagate bool
	preventDefault         bool
}

// NewEventSub returns a new event element config
func NewEventSub(evtype, evtarget string) *EventSub {
	return &EventSub{
		Type:       evtype,
		Target:     evtarget,
		FlatChains: FlatChainIdentity(),
	}
}

// SetStopPropagation sets the value of EventSub.stopPropagate
func (e *EventSub) SetStopPropagation(n bool) {
	e.stopPropagate = n
}

// SetStopImmediatePropagation sets the value of EventSub.stopImmediatePropagate
func (e *EventSub) SetStopImmediatePropagation(n bool) {
	e.stopImmediatePropagate = n
}

// SetPreventDefault sets the value of EventSub.preventDefault
func (e *EventSub) SetPreventDefault(n bool) {
	e.preventDefault = n
}

// StopPropagation returns the value of EventSub.StopPropagation
func (e *EventSub) StopPropagation() bool {
	return e.stopPropagate
}

// StopImmediatePropagation returns the value of EventSub.StopImmediatePropagation
func (e *EventSub) StopImmediatePropagation() bool {
	return e.stopImmediatePropagate
}

// PreventDefault returns the value of EventSub.preventDefault
func (e *EventSub) PreventDefault() bool {
	return e.preventDefault
}

// EventSelector returns the target of the event
func (e *EventSub) EventSelector() string {
	return e.Target
}

// ID returns the event id that EventManager use for this event
func (e *EventSub) ID() string {
	return GetEventID(e)
}

// EventType returns the type of the event
func (e *EventSub) EventType() string {
	return e.Type
}

// ErrEventNotFound is returned when an event is not found
var ErrEventNotFound = errors.New("Event not found")

// EventSubSetup defines a function type for the event setup code
type EventSubSetup func(EventSubs) bool

// EventManager provides a deffered event managing sytem for registery events with
type EventManager struct {
	//events contain the events to be registered on an element
	// it contains the element and the selector type used to match it
	// that is, the value of ElemEvent.Target
	events map[string]EventSubs
	ro     sync.RWMutex
}

// NewEventManager returns a new event manager instance
func NewEventManager() *EventManager {
	em := EventManager{
		events: make(map[string]EventSubs),
	}

	return &em
}

// HasWatch returns true/false if an event target is already marked using the
// format selector#eventType
func (em *EventManager) HasWatch(m string) bool {
	var ok bool
	em.ro.RLock()
	_, ok = em.events[m]
	em.ro.RUnlock()
	return ok
}

// GetEvent returns the event if found by that id
func (em *EventManager) GetEvent(event string) (EventSubs, error) {
	var ed EventSubs
	var ok bool

	em.ro.RLock()
	ed, ok = em.events[event]
	em.ro.RUnlock()

	if !ok {
		return nil, ErrEventNotFound
	}

	return ed, nil
}

// NewEvent allows the adding of event using string values
func (em *EventManager) NewEvent(evtype, evselector string) (EventSubs, bool) {
	eo := BuildEventID(evtype, evselector)

	if em.HasWatch(eo) {
		ed, _ := em.GetEvent(eo)
		return ed, false
	}

	emo := NewEventSub(evtype, evselector)

	em.AddEvent(emo)
	return emo, true
}

//AddEvent adds Event elements into the event manager if a
//event element is already matching that is the combination of selector#eventtype
// then it returns false but if added then true
func (em *EventManager) AddEvent(eo EventSubs) bool {
	id := eo.ID()

	if em.HasWatch(id) {
		return false
	}

	em.ro.Lock()
	em.events[id] = eo
	em.ro.Unlock()
	return true
}

// AddEvents adds a set of ElemEvent into the EventManager
func (em *EventManager) AddEvents(ems ...EventSubs) {
	for _, eo := range ems {
		em.AddEvent(eo)
	}
}

// EachEvent runnings a function over all events
func (em *EventManager) EachEvent(fx func(EventSubs)) {
	em.ro.Lock()
	for _, eo := range em.events {
		fx(eo)
	}
	em.ro.Unlock()
}

// GetEventID returns the id for a ElemEvent object
func GetEventID(m EventSubs) string {
	sel := strings.TrimSpace(m.EventSelector())
	return BuildEventID(sel, m.EventType())
}

// BuildEventID returns the string represent of the values using the select#event format
func BuildEventID(etype, eselect string) string {
	return fmt.Sprintf("%s#%s", eselect, etype)
}
