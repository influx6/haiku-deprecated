package views

import (
	"html/template"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/influx6/flux"
	"github.com/influx6/haiku/events"
	"github.com/influx6/haiku/trees"
	"github.com/influx6/haiku/trees/elems"
)

// Views define a Haiku Component
type Views interface {
	flux.Reactor
	States
	Show()
	Events() *events.EventManager
	Hide()
	Render(...string) trees.Markup
	RenderHTML(...string) template.HTML
	Mount(*js.Object)
	UseMux(ViewMux)
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
	if ds, err := trees.GetStyle(m, "display"); err == nil {
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
	if ds, err := trees.GetStyle(m, "display"); err == nil {
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
	events      *events.EventManager
	fx          ViewMux
	dom         *js.Object
	//liveMarkup represent the current rendered markup
	liveMarkup trees.Markup
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
		events:    events.NewEventManager(),
		encoder:   writer,
		fx:        fx,
	}

	//set up the reaction chain, if we have node attach then render to it
	vm.React(func(r flux.Reactor, _ error, _ interface{}) {
		//if we are not domless then patch
		// log.Printf("will render")
		if vm.dom != nil {
			html := vm.RenderHTML()
			// log.Printf("will render markup: %s \n-------------------", html)
			// log.Printf("Sending to fragment: -> \n %s", html)
			Patch(CreateFragment(string(html)), vm.dom)
		}
	}, true)

	vm.States.UseActivator(func() {
		vm.Show()
	})

	vm.States.UseDeactivator(func() {
		vm.Hide()
	})

	return
}

// UseMux lets you switch the markup generator
func (v *View) UseMux(fx ViewMux) {
	v.fx = fx
}

// Mount is to be called in the browser to loadup this view with a dom
func (v *View) Mount(dom *js.Object) {
	v.dom = dom
	v.events.OffloadDOM()
	v.events.LoadDOM(dom)
	v.Send(true)
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
func (v *View) Events() *events.EventManager {
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

	if v.liveMarkup != nil {
		dom.Reconcile(v.liveMarkup)
	}

	dom.UseEventManager(v.events)
	v.events.LoadUpEvents()
	v.liveMarkup = dom

	return dom
}

// RenderHTML renders out the views markup as a string wrapped with template.HTML
func (v *View) RenderHTML(m ...string) template.HTML {
	ma, _ := v.encoder.Write(v.Render(m...))
	return template.HTML(ma)
}
