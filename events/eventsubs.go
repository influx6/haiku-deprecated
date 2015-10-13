package events

import (
	"errors"
	"fmt"
	"strings"
	"sync"

	"github.com/gopherjs/gopherjs/js"
	// hodom "honnef.co/go/js/dom"
	hodom "github.com/influx6/dom"
)

// EventSubs defines the event interface methods for meta info
type EventSubs interface {
	FlatChains
	ID() string
	EventType() string
	EventSelector() string
	Offload()
	SetStopPropagation(n bool)
	SetStopImmediatePropagation(n bool)
	SetPreventDefault(n bool)
	StopPropagation() bool
	StopImmediatePropagation() bool
	PreventDefault() bool
	DOM(hodom.Element)
	TriggerMatch(hodom.Event)
	Trigger(hodom.Event)
}

// JSEventMux represents a js.listener function which is returned when attached
// using AddEventListeners and is used for removals with RemoveEventListeners
type JSEventMux func(*js.Object)

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
	jslink                 JSEventMux
	dom                    hodom.Element
}

// NewEventSub returns a new event element config
func NewEventSub(evtype, evtarget string) *EventSub {
	return &EventSub{
		Type:       evtype,
		Target:     evtarget,
		FlatChains: FlatChainIdentity(),
	}
}

// DOM sets up the event subs for listening
func (e *EventSub) DOM(dom hodom.Element) {
	e.Offload()
	e.dom = dom
	e.jslink = e.dom.AddEventListener(e.EventType(), true, e.TriggerMatch)
}

// Offload removes all event bindings from current dom element
func (e *EventSub) Offload() {
	if e.dom == nil {
		return
	}

	if e.jslink != nil {
		e.dom.RemoveEventListener(e.EventType(), true, e.jslink)
		e.jslink = nil
	}
}

// Trigger provides bypass for triggering this event sub by passing down an event
// directly without matching target or selector
func (e *EventSub) Trigger(h hodom.Event) {
	e.HandleContext(h)
}

// TriggerMatch check if the current event from a specific parent matches the
// eventarget by using the eventsub selector,if the target is within the results for
// that selector then it triggers the event subscribers
func (e *EventSub) TriggerMatch(h hodom.Event) {
	// if e.dom != nil
	if strings.ToLower(h.Type()) != strings.ToLower(e.EventType()) {
		return
	}

	//get the targets parent
	parent := h.Target().ParentElement()

	var match bool

	//get all possible matches of this query
	posis := parent.QuerySelectorAll(e.EventSelector())

	//get the current event target
	target := h.Target()

	// log.Printf("Checking: %s for %s", target, e.ID())

	//is our target part of those that match the selector
	for _, item := range posis {
		if item.Underlying() != target.Underlying() {
			continue
		}
		match = true
		break
	}

	//if we match then run the listeners registered
	if match {
		//if we dont want immediatepropagation kill it else check propagation also
		if e.StopImmediatePropagation() {
			h.StopImmediatePropagation()
		} else {
			//if we dont want propagation kill it
			if e.StopPropagation() {
				h.StopPropagation()
			}
		}

		//we want to PreventDefault then stop default action
		if e.PreventDefault() {
			h.PreventDefault()
		}

		e.HandleContext(h)
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
	events   map[string]EventSubs
	attaches map[*EventManager]bool
	ro       sync.RWMutex
	wo       sync.RWMutex
	//current dom node attached
	dom hodom.Element
}

// NewEventManager returns a new event manager instance
func NewEventManager() *EventManager {
	em := EventManager{
		events:   make(map[string]EventSubs),
		attaches: make(map[*EventManager]bool),
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

// AttachManager allows a manager to get attached to another manajor to receive a dom binding
// when this receives one
func (em *EventManager) AttachManager(esm *EventManager) {
	//incase of stupid loops, are we attached to the supplied manager? if so duck this
	if esm.HasManager(em) {
		return
	}

	//if we have it already we skip
	if em.HasManager(esm) {
		return
	}

	//its not found so we add it
	em.wo.Lock()
	em.attaches[esm] = true
	em.wo.Unlock()

	//do we already have a dom attached?, then notify this manager immediately
	if em.dom != nil {
		esm.LoadDOM(em.dom)
	}
}

// DetachManager detaches the manager if attached already
func (em *EventManager) DetachManager(esm *EventManager) {
	//if we dont have it attached then skip
	if !em.HasManager(esm) {
		return
	}

	//we got one so we kill it
	em.wo.Lock()
	delete(em.attaches, esm)
	em.wo.Unlock()
}

// HasManager returns true/false if a manager was already attached
func (em *EventManager) HasManager(esm *EventManager) bool {
	//have we already attached before
	em.wo.RLock()
	ok := em.attaches[esm]
	em.wo.RUnlock()
	return ok
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

// EachManager runnings a function over all attached managers
func (em *EventManager) EachManager(fx func(*EventManager)) {
	em.wo.Lock()
	for eo := range em.attaches {
		fx(eo)
	}
	em.wo.Unlock()
}

// LoadDOM passes down the dom element to all EventSub to initialize and listen for their respective events
func (em *EventManager) LoadDOM(dom hodom.Element) {
	//replace the current dom node be used
	em.dom = dom

	// send the dom out to all registered event subs for loadup
	em.EachEvent(func(es EventSubs) {
		es.DOM(dom)
	})

	// send out to all other attach eventmanagers for loadup
	em.EachManager(func(ems *EventManager) {
		ems.LoadDOM(dom)
	})
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
