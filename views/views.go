package views

import (
	"html/template"
	"strings"

	"github.com/gopherjs/gopherjs/js"
	"github.com/influx6/haiku/events"
	"github.com/influx6/haiku/pub"
	"github.com/influx6/haiku/trees"
	"github.com/influx6/haiku/trees/elems"
)

// Renderable provides a interface for a renderable type.
type Renderable interface {
	Render(...string) trees.Markup
}

// ReactiveRenderable provides a interface for a reactive renderable type.
type ReactiveRenderable interface {
	pub.Publisher
	Renderable
}

// Behaviour provides a state changers for haiku.
type Behaviour interface {
	Hide()
	Show()
}

// Views define a Haiku Component
type Views interface {
	pub.Publisher
	States
	Behaviour

	Events() *events.EventManager
	Mount(*js.Object)
	BindView(Views)
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

// View represent a basic Haiku view
type View struct {
	States
	pub.Publisher
	HideState   ViewStates
	ShowState   ViewStates
	activeState ViewStates
	encoder     trees.MarkupWriter
	events      *events.EventManager
	dom         *js.Object
	rview       Renderable
	//liveMarkup represent the current rendered markup
	liveMarkup trees.Markup
}

// NewView returns a basic view
func NewView(view Renderable) *View {
	return MakeView(trees.SimpleMarkupWriter, view)
}

// MakeView returns a Components style
func MakeView(writer trees.MarkupWriter, vw Renderable) (vm *View) {
	vm = &View{
		Publisher: pub.Always(vm),
		States:    NewState(),
		HideState: &HideView{},
		ShowState: &ShowView{},
		events:    events.NewEventManager(),
		encoder:   writer,
		rview:     vw,
	}

	// If its a ReactiveRenderable type then bind the view
	if rxv, ok := vw.(ReactiveRenderable); ok {
		rxv.Bind(vm, true)
	}

	//set up the reaction chain, if we have node attach then render to it
	vm.React(func(r pub.Publisher, _ error, _ interface{}) {
		//if we are not domless then patch
		if vm.dom != nil {
			html := vm.RenderHTML()
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

// BindView binds the given views together,were the view provided as argument will notify this view of change and to act according
func (v *View) BindView(vs Views) {
	vs.Bind(v, true)
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

	dom := v.rview.Render(m...)

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
