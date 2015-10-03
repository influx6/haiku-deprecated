package views

import (
	"html/template"
	"sync/atomic"

	"github.com/influx6/flux"
	"github.com/influx6/haiku/trees"
)

// Viewable defines an interface for any element that can return a rendering of its content out as strings
type Viewable interface {
	Render(...string) trees.Markup
	RenderHTML(...string) template.HTML
}

// StatefulViewable defines an interface for any element that can return a rendering of its content and matches the States interface
type StatefulViewable interface {
	States
	Viewable
}

// Strategy defines the interface method rules for using view strategies
type Strategy interface {
	SwitchActive()
	SwitchInActive()
	Render(Views) trees.Markup
	RenderSource(Views) string
}

// ViewStrategyMux defines a function type that handles and mutates the reply of a view strategy
type ViewStrategyMux func(Views) trees.Markup

// ViewStrategy defines a view behaviour in dealing with a dual state of views i.e views are either active or inactive and ViewStrategy take that and build a custom response provided to the .Render() call
type ViewStrategy struct {
	state    int64
	writer   trees.MarkupWriter
	active   ViewStrategyMux
	inactive ViewStrategyMux
}

// NewViewStrategy returns a new ViewStrategy instance
func NewViewStrategy(w trees.MarkupWriter, active, inactive ViewStrategyMux) *ViewStrategy {
	return &ViewStrategy{
		writer:   w,
		active:   active,
		inactive: inactive,
	}
}

// SwitchActive switches the state of the strategy to active mode
func (v *ViewStrategy) SwitchActive() {
	if atomic.LoadInt64(&v.state) < 1 {
		atomic.StoreInt64(&v.state, 1)
	}
}

// SwitchInActive switches the state of the strategy to inactive mode
func (v *ViewStrategy) SwitchInActive() {
	if atomic.LoadInt64(&v.state) > 0 {
		atomic.StoreInt64(&v.state, 0)
	}
}

// RenderSource renders out the Markup returned by ViewStrategy.Render
func (v *ViewStrategy) RenderSource(vr Views) string {
	source, err := v.writer.Write(v.Render(vr))
	if err != nil {
		return err.Error()
	}
	return source
}

// Render provides the rendering method call for View, it is what is called by a view to implement its show and hide strategy using its internal active and inactive function calls depending on the strategys state
func (v *ViewStrategy) Render(vr Views) trees.Markup {
	var res trees.Markup
	if atomic.LoadInt64(&v.state) < 1 {
		res = v.inactive(vr)
	} else {
		res = v.active(vr)
	}

	return res
}

// Views define an interface for member rules for Views
type Views interface {
	flux.SyncCollectors
	Viewable
	States

	DOM() trees.SearchableMarkup
	Binding() interface{}
	View(string) Viewable
	PartialView() *PartialView
	AddViewable(string, Viewable) error
	AddView(string, string, Viewable) error
	AddStatefulViewable(string, string, StatefulViewable) error
	add(string, Viewable) error
	switchDOM(trees.SearchableMarkup)
	Strategy() Strategy
}

// View provides a base struct for which views can be created with and meets the Views interface
type View struct {
	flux.SyncCollectors
	*State
	//contains the sub-views of the current view
	views *ViewLists
	//strategy defines the view strategy to be used
	strategy Strategy
	binding  interface{}
	dom      trees.SearchableMarkup
}

// NewView returns a new view struct
func NewView(tag string, strategy Strategy, binding interface{}) *View {

	v := &View{
		SyncCollectors: flux.NewSyncCollector(),
		dom:            trees.NewText(""),
		strategy:       strategy,
		State:          NewState(tag),
		views:          NewViewLists(),
		binding:        binding,
	}

	v.State.UseActivator(func(s *StateStat) {
		// v.DOM().Send(true)
		v.strategy.SwitchActive()
	})

	v.State.UseDeactivator(func(s *StateStat) {
		// v.DOM().Send(true)
		v.strategy.SwitchInActive()
	})

	return v
}

// Strategy returns the strategy used by this view
func (v *View) Strategy() Strategy {
	return v.strategy
}

// DOM returns the default text markup over-ride this in subviews
func (v *View) DOM() trees.SearchableMarkup {
	return v.dom
}

// Binding returns the optional binding attached to the view
func (v *View) Binding() interface{} {
	return v.binding
}

// RenderHTML renders the output from .Render() as safe html unescaped
func (v *View) RenderHTML(m ...string) template.HTML {
	var addr string

	if len(m) > 0 {
		addr = m[0]
	}

	if addr != "" {
		v.Engine().All(addr, v.tag)
	}

	return template.HTML(v.strategy.RenderSource(v))
}

// Render calls the internal strategy and renders out the output of that result
func (v *View) Render(m ...string) trees.Markup {
	var addr string

	if len(m) > 0 {
		addr = m[0]
	}

	if addr != "" {
		v.Engine().All(addr, v.tag)
	}

	return v.strategy.Render(v)
}

// PartialView returns a PartialView for this view
func (v *View) PartialView() *PartialView {
	return NewPartialView(v)
}

// View returns the view with the specified tag or nil if not found
func (v *View) View(tag string) Viewable {
	vm := v.views.Get(tag)

	if vm == nil {
		return NewNullRender(tag)
	}

	return vm
}

//Views returns the entire views registered with this view as a list of Viewables
func (v *View) Views() []Viewable {
	return v.views.Views()
}

// AddViewable adds a rendering view which has no state management and will render regardless of state
func (v *View) AddViewable(tag string, vm Viewable) error {
	if err := v.add(tag, vm); err != nil {
		return err
	}
	return nil
}

// AddStatefulViewable adds a rendering State into the view lists and allows this to react accordingly the state of the View depending on the views current state address
func (v *View) AddStatefulViewable(tag, addr string, vm StatefulViewable) error {
	if err := v.add(tag, vm); err != nil {
		return err
	}

	v.Engine().UseState(addr, vm)
	return nil
}

// AddView adds a subview into the current view and depending if the view is a StatefulViewable then it adds it with the giving addresspoint else ignores it and adds it as a regular Viewable
func (v *View) AddView(tag, addr string, vm Viewable) error {
	if svm, ok := vm.(StatefulViewable); ok {
		return v.AddStatefulViewable(tag, addr, svm)
	}

	return v.AddViewable(tag, vm)
}

// String simply calls the internal .Render() function
func (v *View) String() string {
	return string(v.RenderHTML())
}

// switchDOM lets you switch out the dom returned by the view
func (v *View) switchDOM(dom trees.SearchableMarkup) {
	v.dom = dom
}

// add internally adds a view with the tag into the views list
func (v *View) add(tag string, vm Viewable) error {
	return v.views.Add(tag, vm)
}

// PartialView provides a wrapper around View that enforces only partial rendering of the internal state by call .Partial() instead of a .All() view on all .Render() calls
type PartialView struct {
	Views
}

// NewPartialView returns a new PartialView instance
func NewPartialView(v Views) *PartialView {
	return &PartialView{Views: v}
}

// Render calls the internal strategy and renders out the output of that result
func (pv *PartialView) Render(m ...string) trees.Markup {
	var addr string

	if len(m) > 0 {
		addr = m[0]
	}

	if addr != "" {
		pv.Engine().Partial(addr, pv.Tag())
	}

	return pv.Strategy().Render(pv.Views)
}

// PartialView returns a PartialView for this view
func (pv *PartialView) PartialView() *PartialView {
	return pv
}

// RenderHTML renders the output from .Render() as safe html unescaped
func (pv *PartialView) RenderHTML(m ...string) template.HTML {
	var addr string

	if len(m) > 0 {
		addr = m[0]
	}

	if addr != "" {
		pv.Engine().Partial(addr, pv.Tag())
	}

	return template.HTML(pv.Strategy().RenderSource(pv.Views))
}

// ObserveViewable defines an interface for any element that can return a rendering of its content out as strings
type ObserveViewable interface {
	flux.Reactor
	Viewable
}

// ObserveStatefulViewable defines an interface for any element that can return a rendering of its content and matches the States interface
type ObserveStatefulViewable interface {
	flux.Reactor
	StatefulViewable
}

// ReactiveViews provides the interface type for ReactiveView
type ReactiveViews interface {
	Views
	flux.Reactor
}

// ReactiveView defines a struct that handles the addition of views that react to change, meaning it can deal with the standard Viewable and StatefulViewable types and the combination of Observers by letting the providing listen for change in the main view or subviews to take appropriate action i.e it does nothing than the normal views but only to signal a change reaction from subviews upward to anyone who wishes to listen and react to that
type ReactiveView struct {
	Views
	flux.Reactor
	dombinder flux.Reactor
}

// BindReactor binds the reactorView with a binding value if that value is a flux.Rector type
func BindReactor(v ReactiveViews, b interface{}) {
	if dok, ok := b.(flux.Reactor); ok {
		dok.Bind(v, false)
	}
}

// ReactView returns a new ReactiveView instance using a Views type as a composition, thereby turning a simple view into a reactable view. if the `dobind` bool is true and if the binding from the view is reactive then the binding is made to the new ReactiveView
func ReactView(v Views, dobind bool) ReactiveViews {
	rve, ok := v.(ReactiveViews)

	if ok {
		return rve
	}

	vr := &ReactiveView{
		Reactor: flux.ReactIdentity(),
		Views:   v,
	}

	if dobind {
		BindReactor(vr, v.Binding())
	}

	return vr
}

// AddViewable adds a rendering view which has no state management and will render regardless of state
func (v *ReactiveView) AddViewable(tag string, vm Viewable) error {
	if err := v.add(tag, vm); err != nil {
		return err
	}

	if osm, ok := vm.(ObserveViewable); ok {
		osm.Bind(v, false)
	}

	return nil
}

// AddStatefulViewable adds a rendering State into the view lists and allows this to react accordingly the state of the View depending on the views current state address
func (v *ReactiveView) AddStatefulViewable(tag, addr string, vm StatefulViewable) error {
	if err := v.add(tag, vm); err != nil {
		return err
	}

	if osm, ok := vm.(ObserveStatefulViewable); ok {
		osm.Bind(v, false)
	}

	v.Engine().UseState(addr, vm)

	return nil
}

// AddView adds a subview into the current view and depending if the view is a StatefulViewable then it adds it with the giving addresspoint else ignores it and adds it as a regular Viewable
func (v *ReactiveView) AddView(tag, addr string, vm Viewable) error {
	if svm, ok := vm.(StatefulViewable); ok {
		return v.AddStatefulViewable(tag, addr, svm)
	}

	return v.AddViewable(tag, vm)
}

// switchDOM lets you switch out the dom returned by the view
func (v *ReactiveView) switchDOM(dom trees.SearchableMarkup) {
	if v.dombinder != nil {
		v.dombinder.Close()
	}

	binder := dom.React(func(r flux.Reactor, _ error, d interface{}) {
		v.Reactor.Send(d)
	}, true)

	v.Views.switchDOM(dom)

	v.dombinder = binder
}

// Blueprint defines a interface type for blueprints
type Blueprint interface {
	Type() string
}
