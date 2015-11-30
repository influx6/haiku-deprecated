package types

import "github.com/gopherjs/gopherjs/js"

// EventHandler provides the function type for event callbacks when subscribing
// to events in Haiku.
type EventHandler func(Event)

// Event defines the base interface for browser events and defines the basic interface methods they must provide
type Event interface {
	Bubbles() bool
	Cancelable() bool
	CurrentTarget() *js.Object
	DefaultPrevented() bool
	EventPhase() int
	Target() *js.Object
	Timestamp() int
	Type() string
	PreventDefault()
	StopImmediatePropagation()
	StopPropagation()
	Core() *js.Object
}

// EventMeta provides a basic information about the events and what its targets
type EventMeta struct {
	// Type is the event type to use
	Type string
	//Target is a selector value for matching a event
	Target                   string
	StopPropagation          bool
	StopImmediatePropagation bool
	PreventDefault           bool
	Removed                  bool
}

// // EventManagers defines the Event.EventManager interface type and is used to clean up import path usage and standardize the api
// type EventManagers interface {
//  HasEvent(m string) bool
//  GetEvent(event string)
//  NewEventMeta(*EventMeta) (*EventSub, bool)
//  NewEvent(evtype, evselector string) (*EventSub, bool)
//  AttachManager(EventManagers)
//  DetachManager(EventManagers)
//  HasManager(EventManagers) bool
//  RemoveEvent(event string)
//  AddEvent(*EventSub) bool
//  AddEvents(...*EventSub)
//  EachEvent(fx func(*EventSub))
//  EachManager(fx func(*EventManager))
//  DisconnectRemoved()
//  OffloadDOM()
//  LoadUpEvents()
//  LoadDOM(dom *js.Object) bool
// }
