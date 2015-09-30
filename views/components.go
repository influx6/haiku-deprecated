package views

// Components defines the interface member rules for Blueprint instances aka Component
type Components interface {
	ReactiveViews
	Events() *EventManager
}

// component defines the concrete implmentation of a blueprint instance
type component struct {
	ReactiveViews
	events *EventManager
}

// NewComponent returns a new component instance
func NewComponent(v Views) Components {
	co := component{
		ReactiveViews: ReactView(v),
		events:        NewEventManager(),
	}

	return &co
}

// Events returns a EventManager used by the component
func (c *component) Events() *EventManager {
	return c.events
}
