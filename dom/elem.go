package dom

import (
	"errors"
	"html/template"
	"log"
	"strings"
	"sync"

	"github.com/gopherjs/gopherjs/js"
	"github.com/influx6/flux"
	"github.com/influx6/haiku/views"
	hodom "honnef.co/go/js/dom"
)

// JSEventMux represents a js.listener function which is returned when attached
// using AddEventListeners and is used for removals with RemoveEventListeners
type JSEventMux func(*js.Object)

// ElemEvent represent a single event configuration for dom.Elem objects
// instance which allows chaining of events listeners like middleware
type ElemEvent struct {
	// Type is the event type to use
	Type string
	//Target is a selector value for matching a event
	Target                 string
	StopPropagate          bool
	StopImmediatePropagate bool
	PreventDefault         bool
	jslink                 JSEventMux
	FlatChains
}

// NewElemEvent returns a new event element config
func NewElemEvent(evtype, evtarget string) *ElemEvent {
	return &ElemEvent{
		Type:       evtype,
		Target:     evtarget,
		FlatChains: FlatChainIdentity(),
	}
}

// EventSelector returns the target of the event
func (e *ElemEvent) EventSelector() string {
	return e.Target
}

// ID returns the event id that EventManager use for this event
func (e *ElemEvent) ID() string {
	return GetEventID(e)
}

// EventType returns the type of the event
func (e *ElemEvent) EventType() string {
	return e.Type
}

// Matches check if the current event from a specific parent matches this target
func (e *ElemEvent) Matches(h hodom.Event) {
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

	log.Printf("Checking: %s for %s", target, e.ID())

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
		if e.StopImmediatePropagate {
			h.StopImmediatePropagation()
		} else {
			//if we dont want propagation kill it
			if e.StopPropagate {
				h.StopPropagation()
			}
		}

		//we want to PreventDefault then stop default action
		if e.PreventDefault {
			h.PreventDefault()
		}

		e.HandleContext(h)
	}
}

// EventManager provides a deffered event managing sytem for registery events with
type EventManager struct {
	//events contain the events to be registered on an element
	// it contains the element and the selector type used to match it
	// that is, the value of ElemEvent.Target
	dom    hodom.Element
	events map[string]*ElemEvent
	ro     sync.RWMutex
}

// NewEventManager returns a new event manager instance
func NewEventManager(elem hodom.Element) *EventManager {
	em := EventManager{
		events: make(map[string]*ElemEvent),
		dom:    elem,
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

// ErrEventNotFound is returned when an event is not found
var ErrEventNotFound = errors.New("Event not found")

// GetEvent returns the event if found by that id
func (em *EventManager) GetEvent(event string) (*ElemEvent, error) {
	var ed *ElemEvent
	var ok bool

	em.ro.RLock()
	ed, ok = em.events[event]
	em.ro.RUnlock()

	if !ok {
		return nil, ErrEventNotFound
	}

	return ed, nil
}

// AddEvent allows the adding of event using string values
func (em *EventManager) AddEvent(evtype, evselector string) *ElemEvent {
	eo := BuildEventID(evtype, evselector)

	if em.HasWatch(eo) {
		ed, _ := em.GetEvent(eo)
		return ed
	}

	emo := NewElemEvent(evtype, evselector)

	em.WatchEvent(emo)
	return emo
}

//WatchEvent adds Event elements into the event manager if a
//event element is already matching that is the combination of selector#eventtype
// then it returns false but if added then true
func (em *EventManager) WatchEvent(eo *ElemEvent) bool {
	id := eo.ID()

	if !em.setupEvent(eo) {
		return false
	}

	// eo.jslink = em.dom.AddEventListener(eo.Type(), true, eo.Matches)
	em.ro.Lock()
	em.events[id] = eo
	em.ro.Unlock()
	return true
}

// WatchEvents adds a set of ElemEvent into the EventManager
func (em *EventManager) WatchEvents(ems ...*ElemEvent) {
	for _, eo := range ems {
		em.WatchEvent(eo)
	}
}

// UseDOM allows switching the EventManager dom element
func (em *EventManager) UseDOM(dom hodom.Element) {
	em.unRegisterEvents()
	em.dom = dom
	em.ReRegisterEvents()
}

//ReRegisterEvents re-registeres all the event config with dom element
//WARNING: only call this if you switch the event managers dom element
func (em *EventManager) ReRegisterEvents() {

	if em.dom == nil {
		return
	}

	for _, eo := range em.events {
		if eo.jslink != nil {
			em.dom.RemoveEventListener(eo.EventType(), true, eo.jslink)
		}
		em.setupEvent(eo)
	}

}

// unRegisterEvents removes all event bindings from current dom element
func (em *EventManager) unRegisterEvents() {

	if em.dom == nil {
		return
	}

	for _, eo := range em.events {
		if eo.jslink != nil {
			em.dom.RemoveEventListener(eo.EventType(), true, eo.jslink)
			eo.jslink = nil
		}
	}
}

//setupEvent only setsup the event link
func (em *EventManager) setupEvent(eo *ElemEvent) bool {
	id := GetEventID(eo)

	if em.HasWatch(id) {
		return false
	}

	if em.dom != nil {
		eo.jslink = em.dom.AddEventListener(eo.EventType(), true, eo.Matches)
	}

	return true
}

// Elem represent a standard html element
type Elem struct {
	hodom.Element
	*EventManager
}

// NewElement returns a new Element instance for interacting with the dom
func NewElement(e hodom.Element) *Elem {
	return &Elem{
		Element:      e,
		EventManager: NewEventManager(e),
	}
}

//UseDOM resets the elements dom target
//WARNING: use this we care,elements are not generally in need of switching
//they should be concrete and always in the dom until you truly decide to
//delete them but this is provided for the case when there is a need to switch
//dom elements so use carefuly
func (em *Elem) UseDOM(dom hodom.Element) {
	em.Element = dom
	em.EventManager.UseDOM(dom)
}

// Remove removes this element from this parent
func (em *Elem) Remove() {
	pe := em.ParentElement()
	if pe != nil {
		pe.RemoveChild(em.Element)
	}
}

// Html sets value of the inner html
func (em *Elem) Html(el, target string) {
	if target == "" {
		em.SetInnerHTML(el)
		return
	}

	to := em.QuerySelectorAll(target)
	for _, eo := range to {
		eo.SetInnerHTML(el)
	}
}

// Text sets value of the text content or returns the text content value if it receives no arguments
func (em *Elem) Text(el, target string) {
	if target == "" {
		em.SetTextContent(el)
		return
	}

	to := em.QuerySelectorAll(target)
	for _, eo := range to {
		eo.SetTextContent(el)
	}
}

// ViewElement combines the dom.Element and haiku.View package to create a renderable view
type ViewElement struct {
	views.ReactiveViews
	*Elem
	target   string
	lastAddr []string
}

// ViewDOMTemplate creates a new ViewElement from the given set of arguments ready for rendering
func ViewDOMTemplate(viewtag, targetSelector string, dol hodom.Element, tl *template.Template, so *views.ViewStrategy) *ViewElement {
	view := views.NewReactiveView(viewtag, tl, so)
	elem := NewElement(dol)
	return NewViewElement(elem, view, targetSelector)
}

// ViewDOM creates a new ViewElement from the given set of arguments ready for rendering
func ViewDOM(dol hodom.Element, v views.Views, targetSelector string) *ViewElement {
	return NewViewElement(NewElement(dol), v, targetSelector)
}

// NewViewElement returns a new instance of ViewElement which takes a dom *Element,
// a View manager and a optional target string which is used when rendring incase we
// wish to render to a target within the dom.Element and not the element itself has.
// This should be never be mixed into another View because its the last and should be
// the last point of rendering, it can be added to a StateEngine for updates on behaviour and state
// but never for sub-view rendering
func NewViewElement(elem *Elem, v views.Views, target string) *ViewElement {

	ve := &ViewElement{
		ReactiveViews: views.ReactView(v),
		Elem:          elem,
		target:        target,
	}

	ve.React(func(r flux.Reactor, err error, d interface{}) {
		if err != nil {
			r.ReplyError(err)
			return
		}

		ve.Sync(ve.lastAddr...)
		r.Reply(d)
	}, true)

	return ve
}

// Sync calls the views render function and renders out into the attach dom element
// each time it said to sync, its address is cache for when there is an update by
// the internal sub-views
func (v *ViewElement) Sync(addr ...string) {
	v.lastAddr = addr
	v.Elem.Html(v.ReactiveViews.Render(addr...), v.target)
}
