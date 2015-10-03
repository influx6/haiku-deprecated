package dom

import (
	"html/template"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/gopherjs/gopherjs/js"
	"github.com/influx6/flux"
	"github.com/influx6/haiku/trees"
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
func BuildEvents(elem hodom.Element, events *views.EventManager, build bool) *Events {
	eo := &Events{
		EventManager: events,
		dom:          elem,
	}

	if build {
		eo.ReRegisterEvents()
	}

	return eo
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

// ViewComponent combines the dom.Element and haiku.View package to create a renderable view
type ViewComponent struct {
	views.Components
	dom      hodom.Element
	events   *Events
	lastAddr []string
	throttle time.Timer
}

// BasicView creates a new ViewElement with a view component for rendering
func BasicView(v views.Components) *ViewComponent {
	vc := &ViewComponent{
		Components: v,
		events:     BuildEvents(nil, v.Events(), false),
		lastAddr:   []string{"."},
	}

	vc.React(func(r flux.Reactor, _ error, _ interface{}) {
		log.Printf("reacting for change: %s", vc.dom)
		//if we are not domless then patch
		if vc.dom != nil {
			log.Printf("building html for rendering", v)
			html := vc.RenderHTML(vc.lastAddr...)
			log.Printf("html for %s -> %s", v.Tag(), html)
			Patch(CreateFragment(string(html)), vc.dom)
		}
	}, true)

	return vc
}

// Mount assigns the dom element to use with the view
func (v *ViewComponent) Mount(dom hodom.Element) *ViewComponent {
	v.events.unRegisterEvents()
	v.dom = dom
	v.events.dom = dom
	v.events.ReRegisterEvents()
	v.Send(true)
	return v
}

// Render overrides Component.Render
func (v *ViewComponent) Render(n ...string) trees.Markup {
	v.lastAddr = n
	return v.Components.Render(n...)
}

// RenderHTML overrides Component.Render
func (v *ViewComponent) RenderHTML(n ...string) template.HTML {
	v.lastAddr = n
	return v.Components.RenderHTML(n...)
}
