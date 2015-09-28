package views

import (
	"bytes"
	"fmt"
	"html/template"
	"sync/atomic"

	"github.com/influx6/assets"
	"github.com/influx6/flux"
)

// Viewable defines an interface for any element that can return a rendering of its content out as strings
type Viewable interface {
	Render(...string) string
	RenderHTML(...string) template.HTML
}

// StatefulViewable defines an interface for any element that can return a rendering of its content and matches the States interface
type StatefulViewable interface {
	States
	Viewable
}

// ViewStrategyMux defines a function type that handles and mutates the reply of a view strategy
type ViewStrategyMux func(*View) string

// ViewStrategy defines a view behaviour in dealing with a dual state of views i.e views are either active or inactive and ViewStrategy take that and build a custom response provided to the .Render() call
type ViewStrategy struct {
	state    int64
	active   ViewStrategyMux
	inactive ViewStrategyMux
}

// NewViewStrategy returns a new ViewStrategy instance
func NewViewStrategy(active, inactive ViewStrategyMux) *ViewStrategy {
	return &ViewStrategy{
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

// Render provides the rendering method call for View, it is what is called by a view to implement its show and hide strategy using its internal active and inactive function calls depending on the strategys state
func (v *ViewStrategy) Render(vr *View) string {
	if atomic.LoadInt64(&v.state) < 1 {
		return v.inactive(vr)
	}
	return v.active(vr)
}

// View provides a base struct for which views can be created with and meets the Views interface
type View struct {
	*State
	Tmpl *template.Template
	//contains the sub-views of the current view
	views *ViewLists
	//strategy defines the view strategy to be used
	strategy *ViewStrategy
}

// BuildViewTemplateFunctions returns a template.FuncMap that contains the template functions (View and Views) for use with a template. Ensure to pass this to the root template so you can acccess it down
func BuildViewTemplateFunctions(v *View) template.FuncMap {
	return template.FuncMap{
		"view": func(tag string) Viewable {
			return v.View(tag)
		},
		"views": func() []Viewable {
			return v.Views()
		},
	}
}

// NewView returns a new view struct
func NewView(tag string, tl *template.Template, strategy *ViewStrategy) *View {

	v := View{
		Tmpl:     tl,
		strategy: strategy,
		State:    NewState(tag),
		views:    NewViewLists(),
	}

	// bx := BuildViewTemplateFunctions(&v)
	// log.Printf("%s", bx)
	// v.Tmpl = tl.Funcs(bx).

	v.State.UseActivator(func(s *StateStat) {
		v.strategy.SwitchActive()
	})

	v.State.UseDeactivator(func(s *StateStat) {
		v.strategy.SwitchInActive()
	})

	return &v
}

// RenderHTML renders the output from .Render() as safe html unescaped
func (v *View) RenderHTML(m ...string) template.HTML {
	return template.HTML(v.Render(m...))
}

// Render calls the internal strategy and renders out the output of that result
func (v *View) Render(m ...string) string {
	var addr string

	if len(m) > 0 {
		addr = m[0]
	}

	if addr != "" {
		v.Engine().All(addr, v.tag)
	}

	return v.strategy.Render(v)
}

// String simply calls the internal .Render() function
func (v *View) String() string {
	return v.Render()
}

// ExecuteTemplate calls the internal template.Template.ExecuteTemplate and returns the data from the rendering operation. The template is runned with the name but the view as the object/binding
func (v *View) ExecuteTemplate(name string) ([]byte, error) {
	var buf bytes.Buffer
	err := v.Tmpl.ExecuteTemplate(&buf, name, v)
	return buf.Bytes(), err
}

// Execute calls the internal template.Template.Execute and returns the data from the rendering operation. The template is runned with the name but the view as the object/binding
func (v *View) Execute() ([]byte, error) {
	var buf bytes.Buffer
	err := v.Tmpl.Execute(&buf, v)
	return buf.Bytes(), err
}

// Views define an interface for member rules for Views
type Views interface {
	Viewable
	States
	// NameTag() string
	// Stategy() *ViewStrategy
	View(string) Viewable
	PartialView() *PartialView
	AddViewable(string, Viewable) error
	AddView(string, string, Viewable) error
	AddStatefulViewable(string, string, StatefulViewable) error
	add(string, Viewable) error
}

// // Stategy returns the views assigned tag
// func (v *View) Stategy() *ViewStrategy {
// 	return v.strategy
// }

// // NameTag returns the views assigned tag
// func (v *View) NameTag() string {
// 	return v.tag
// }

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

func (v *View) add(tag string, vm Viewable) error {
	return v.views.Add(tag, vm)
}

// PartialView provides a wrapper around View that enforces only partial rendering of the internal state by call .Partial() instead of a .All() view on all .Render() calls
type PartialView struct {
	*View
}

// NewPartialView returns a new PartialView instance
func NewPartialView(v *View) *PartialView {
	return &PartialView{View: v}
}

// Render calls the internal strategy and renders out the output of that result
func (pv *PartialView) Render(m ...string) string {
	var addr string

	if len(m) > 0 {
		addr = m[0]
	}

	if addr != "" {
		pv.Engine().Partial(addr, pv.Tag())
	}

	return pv.strategy.Render(pv.View)
}

// PartialView returns a PartialView for this view
func (pv *PartialView) PartialView() *PartialView {
	return pv
}

// RenderHTML renders the output from .Render() as safe html unescaped
func (pv *PartialView) RenderHTML(m ...string) template.HTML {
	return template.HTML(pv.Render(m...))
}

// SilentStrategy is a simple strategy that when the view is activated calls the View.Execute
func SilentStrategy() *ViewStrategy {
	return NewViewStrategy(func(v *View) string {
		bo, err := v.Execute()

		if err != nil {
			return fmt.Sprintf("CustomError(%s): %s", v.Tag(), err.Error())
		}

		return string(bo)
	}, func(v *View) string {
		return ""
	})
}

// SilentNameStrategy is a simple strategy that when the view is activated calls the View.ExecuteTemplate
func SilentNameStrategy(base string) *ViewStrategy {
	return NewViewStrategy(func(v *View) string {
		bo, err := v.ExecuteTemplate(base)

		if err != nil {
			return fmt.Sprintf("CustomError(%s): %s", v.Tag(), err.Error())
		}

		return string(bo)
	}, func(v *View) string {
		return ""
	})
}

// HiddenStrategy is a simple strategy that when the view is activated calls the View.Execute
func HiddenStrategy() *ViewStrategy {
	return NewViewStrategy(func(v *View) string {
		bo, err := v.Execute()

		if err != nil {
			return fmt.Sprintf("CustomError(%s): %s", v.Tag(), err.Error())
		}

		return string(bo)
	}, func(v *View) string {
		bo, err := v.Execute()

		if err != nil {
			return fmt.Sprintf("CustomError(%s): %s", v.Tag(), err.Error())
		}

		return fmt.Sprintf(`<div style="display:none;">\n%s\n</div>`, string(bo))
	})
}

// HiddenNameStrategy is a simple strategy that when the view is activated calls the View.ExecuteTemplate and when hidden wrap it within a div tag laced with a display none style
func HiddenNameStrategy(base string) *ViewStrategy {
	return NewViewStrategy(func(v *View) string {
		bo, err := v.ExecuteTemplate(base)

		if err != nil {
			return fmt.Sprintf("CustomError(%s): %s", v.Tag(), err.Error())
		}

		return string(bo)
	}, func(v *View) string {
		bo, err := v.ExecuteTemplate(base)

		if err != nil {
			return fmt.Sprintf("CustomError(%s): %s", v.Tag(), err.Error())
		}

		return fmt.Sprintf(`<div style="display:none;">\n%s\n</div>`, string(bo))
	})
}

// ViewHTMLTemplate simple renders out the internal views of a root View into html like tags
const ViewHTMLTemplate = `
	<masterview>
		{{range .Views }}
		  <view>
				{{ .RenderHTML }}
		  </view>
		{{ end }}
	</masterview>
`

// ViewLightTemplate simple renders out the internal views of a root View
const ViewLightTemplate = `
	{{range .Views }}
			{{ .RenderHTML }}
	{{ end }}
`

// SimpleView provides a view based on the ViewLightTemplate template
func SimpleView(tag string) (v *View, err error) {
	return SourceView(tag, ViewLightTemplate)
}

// SimpleTreeView provides a view based on the ViewHTMLTemplate template
func SimpleTreeView(tag string) (v *View, err error) {
	return SourceView(tag, ViewHTMLTemplate)
}

// SourceView provides a view that takes the template format of which it will render the view as
func SourceView(tag, tmpl string) (v *View, err error) {
	var tl *template.Template

	tl, err = template.New(tag).Parse(tmpl)

	if err != nil {
		return
	}

	v = NewView(tag, tl, SilentStrategy())

	return
}

// AssetView provides a view that takes the template format of which it will render the view as
func AssetView(tag, blockName string, as *assets.AssetTemplate) (v *View, err error) {
	v = NewView(tag, as.Tmpl, SilentNameStrategy(blockName))
	return
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
	flux.Reactor
	Views
}

// ReactiveView defines a struct that handles the addition of views that react to change, meaning it can deal with the standard Viewable and StatefulViewable types and the combination of Observers by letting the providing listen for change in the main view or subviews to take appropriate action i.e it does nothing than the normal views but only to signal a change reaction from subviews upward to anyone who wishes to listen and react to that
type ReactiveView struct {
	flux.Reactor
	Views
}

// NewReactiveView provides a decorator function to return a new ReactiveView with the same arguments passed to NewView(...)
func NewReactiveView(tag string, tl *template.Template, strategy *ViewStrategy) ReactiveViews {
	return ReactView(NewView(tag, tl, strategy))
}

// ReactView returns a new ReactiveView instance using a Views type as a composition, thereby turning
// a simple view into a reactable view
func ReactView(v Views) ReactiveViews {
	if rve, ok := v.(ReactiveViews); ok {
		return rve
	}

	return &ReactiveView{
		Reactor: flux.ReactIdentity(),
		Views:   v,
	}
}

// ReactiveSourceView provides a view that takes the template format of which it will render the view as
func ReactiveSourceView(tag, tmpl string) (ReactiveViews, error) {
	sv, err := SourceView(tag, tmpl)
	if err != nil {
		return nil, err
	}

	return ReactView(sv), nil
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
