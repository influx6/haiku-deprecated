package dom

import (
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
	views.EventSubs
	jslink JSEventMux
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

// Events provides a deffered event managing sytem for registery events with
type Events struct {
	*views.EventManager
	dom   hodom.Element
	ro    sync.RWMutex
	fresh bool
}

// NewEvents returns a new event manager instance
func NewEvents(elem hodom.Element) *Events {
	em := Events{
		EventManager: views.NewEventManager(),
		dom:          elem,
	}

	return &em
}

// BuildEvents builds a event manager untop of a exisiting views.EventManager and dom element
func BuildEvents(elem hodom.Element, events *views.EventManager) *Events {
	eo := Events{
		EventManager: events,
		dom:          elem,
	}
	return &eo
}

// AddEvent allows the adding of event using string values
func (em *Events) AddEvent(evtype, evselector string) views.EventSubs {
	eo := views.BuildEventID(evtype, evselector)

	if em.HasWatch(eo) {
		ed, _ := em.GetEvent(eo)
		return ed
	}

	emo := views.NewEventSub(evtype, evselector)
	em.WatchEvent(&ElemEvent{EventSubs: emo})
	return emo
}

//WatchEvent adds Event elements into the event manager if a
//event element is already matching that is the combination of selector#eventtype
// then it returns false but if added then true
func (em *Events) WatchEvent(eo views.EventSubs) bool {
	if em.HasWatch(eo.ID()) {
		return false
	}

	var eom *ElemEvent

	if eod, ok := eo.(*ElemEvent); ok {
		eom = eod
	} else {
		eom = &ElemEvent{EventSubs: eo}
	}

	em.EventManager.AddEvent(eom)
	em.setupEvent(eom)

	return true
}

// WatchEvents adds a set of ElemEvent into the EventManager
func (em *Events) WatchEvents(ems ...views.EventSubs) {
	for _, eo := range ems {
		em.WatchEvent(eo)
	}
}

// UseDOM allows switching the EventManager dom element
func (em *Events) UseDOM(dom hodom.Element) {
	em.unRegisterEvents()
	em.dom = dom
	em.ReRegisterEvents()
}

//ReRegisterEvents re-registeres all the event config with dom element
//WARNING: only call this if you switch the event managers dom element
func (em *Events) ReRegisterEvents() {
	if em.dom == nil {
		return
	}

	if !em.fresh {
		em.unRegisterEvents()
	}

	em.fresh = false
	em.EventManager.EachEvent(func(evo views.EventSubs) {
		if eo, ok := evo.(*ElemEvent); ok {
			em.setupEvent(eo)
		}
	})

}

// unRegisterEvents removes all event bindings from current dom element
func (em *Events) unRegisterEvents() {
	if em.dom == nil {
		return
	}

	em.EventManager.EachEvent(func(evo views.EventSubs) {
		if eo, ok := evo.(*ElemEvent); ok {
			if eo.jslink != nil {
				em.dom.RemoveEventListener(eo.EventType(), true, eo.jslink)
				eo.jslink = nil
			}
		}
	})

	em.fresh = true
}

//setupEvent only setsup the event link
func (em *Events) setupEvent(eo *ElemEvent) {
	if em.dom != nil {
		eo.jslink = em.dom.AddEventListener(eo.EventType(), true, eo.Matches)
	}
}

// Elem represent a standard html element
type Elem struct {
	hodom.Element
	*Events
}

// NewElement returns a new Element instance for interacting with the dom
func NewElement(e hodom.Element) *Elem {
	return &Elem{
		Element: e,
		Events:  NewEvents(e),
	}
}

// UseElement returns a new Element instance for interacting with the dom
func UseElement(e hodom.Element, ex *views.EventManager) *Elem {
	m := &Elem{
		Element: e,
		Events:  BuildEvents(e, ex),
	}

	m.fresh = true
	m.ReRegisterEvents()
	return m
}

//UseDOM overwrites the Events.UseDOM to ensure no issues
// TODO: find out if this can cause race-conditons
func (em *Elem) UseDOM(dom hodom.Element) {
	em.Element = dom
	em.Events.UseDOM(dom)
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

// ViewComponent creates a new ViewElement with a view component for rendering
func ViewComponent(dol hodom.Element, v views.Components, targetSelector string) *ViewElement {
	elem := UseElement(dol, v.Events())
	return NewViewElement(elem, v, targetSelector)
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

	ve.Sync(".")

	return ve
}

// Sync calls the views render function and renders out into the attach dom element
// each time it said to sync, its address is cache for when there is an update by
// the internal sub-views
func (v *ViewElement) Sync(addr ...string) {
	v.lastAddr = addr
	v.Elem.Html(v.ReactiveViews.Render(addr...), v.target)
}
