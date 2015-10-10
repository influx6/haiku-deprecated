package views

import (
	"html/template"
	"strings"

	"github.com/influx6/flux"
	"github.com/influx6/haiku/trees"
	"github.com/influx6/haiku/trees/elems"
)

// Views define a Haiku Component
type Views interface {
	flux.Reactor
	States
	Show()
	Events() *EventManager
	Hide()
	Render(...string) trees.Markup
	RenderHTML(...string) template.HTML
}

// ViewStates defines the two possible behavioral state of a view's markup
type ViewStates interface {
	Render(trees.Markup)
}

// HideView provides a ViewStates for Views inactive state
type HideView struct{}

// Render marks the given markup as display:none
func (v *HideView) Render(m trees.Markup) {
	//if we are allowed to query for styles then check and change display
	if ds, err := m.GetStyle("display"); err == nil {
		if !strings.Contains(ds.Value, "none") {
			ds.Value = "none"
		}
	}
}

// ShowView provides a ViewStates for Views active state
type ShowView struct{}

// Render marks the given markup with a display: block
func (v *ShowView) Render(m trees.Markup) {
	//if we are allowed to query for styles then check and change display
	if ds, err := m.GetStyle("display"); err == nil {
		if strings.Contains(ds.Value, "none") {
			ds.Value = "block"
		}
	}
}

// ViewMux defines a markup generating function for view
type ViewMux func() trees.Markup

// View represent a basic Haiku view
type View struct {
	States
	flux.Reactor
	HideState   ViewStates
	ShowState   ViewStates
	activeState ViewStates
	encoder     trees.MarkupWriter
	events      *EventManager
	fx          ViewMux
}

// NewView returns a basic view
func NewView(fx ViewMux) *View {
	return MakeView(trees.SimpleMarkupWriter, fx)
}

// MakeView returns a Components style
func MakeView(writer trees.MarkupWriter, fx ViewMux) (vm *View) {
	vm = &View{
		Reactor:   flux.FlatAlways(vm),
		States:    NewState(),
		HideState: &HideView{},
		ShowState: &ShowView{},
		events:    NewEventManager(),
		encoder:   writer,
		fx:        fx,
	}

	vm.States.UseActivator(func() {
		vm.Show()
	})

	vm.States.UseDeactivator(func() {
		vm.Hide()
	})

	return
}

// Show activates the view to generate a visible markup
func (v *View) Show() {
	if v.ShowState == nil {
		v.ShowState = &ShowView{}
	}
	v.activeState = v.ShowState
}

// Hide deactivates the view
func (v *View) Hide() {
	if v.HideState == nil {
		v.HideState = &HideView{}
	}
	v.activeState = v.HideState
}

// Events returns the views events manager
func (v *View) Events() *EventManager {
	return v.events
}

// Render renders the generated markup for this view
func (v *View) Render(m ...string) trees.Markup {
	if len(m) <= 0 {
		m = []string{"."}
	}

	v.Engine().All(m[0])

	dom := v.fx()

	if dom == nil {
		return elems.Div()
	}

	return dom
}

// RenderHTML renders out the views markup as a string wrapped with template.HTML
func (v *View) RenderHTML(m ...string) template.HTML {
	ma, _ := v.encoder.Write(v.Render(m...))
	return template.HTML(ma)
}
