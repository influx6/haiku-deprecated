package views

import (
	"bytes"
	"fmt"
	"html/template"
	"sync/atomic"
)

// Viewable defines an interface for any element that can return a rendering of its content out as strings
type Viewable interface {
	Render() string
	RenderHTML() template.HTML
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
		atomic.StoreInt64(&v.state, 1)
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
		"view": func(data ...interface{}) Viewable {
			if len(data) < 1 {
				return nil
			}

			var tag string
			var ok bool
			tag, ok = data[0].(string)

			if !ok {
				return nil
			}

			return v.View(tag)
		},
		"views": func(data ...interface{}) []Viewable {
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

	v.Tmpl.Funcs(BuildViewTemplateFunctions(&v))

	v.UseActivator(func(_ *StateStat) {
		v.strategy.SwitchActive()
	})

	v.UseDeactivator(func(_ *StateStat) {
		v.strategy.SwitchActive()
	})

	return &v
}

// RenderHTML renders the output from .Render() as safe html unescaped
func (v *View) RenderHTML() template.HTML {
	return template.HTML(v.Render())
}

// Render calls the internal strategy and renders out the output of that result
func (v *View) Render() string {
	v.Engine().All(".", "")
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
	if err := v.views.Add(tag, vm); err != nil {
		return err
	}
	return nil
}

// AddStatefulViewable adds a rendering State into the view lists and allows this to react accordingly the state of the View depending on the views current state address
func (v *View) AddStatefulViewable(tag, addr string, vm StatefulViewable) error {
	if err := v.views.Add(tag, vm); err != nil {
		return err
	}

	return nil
}

// AddView adds a subview into the current view and depending if the view is a StatefulViewable then it adds it with the giving addresspoint else ignores it and adds it as a regular Viewable
func (v *View) AddView(tag, addr string, vm Viewable) error {
	if svm, ok := vm.(StatefulViewable); ok {
		return v.AddStatefulViewable(tag, addr, svm)
	}

	return v.AddViewable(tag, vm)
}

// // UseDeactivator gets overide
// func (v *View) UseDeactivator(StateResponse) {}
//
// // UseActivator gets overide
// func (v *View) UseActivator(StateResponse) {}

// ViewHTMLTemplate simple renders out the internal views of a root View into html like tags
const ViewHTMLTemplate = `
	<view id={{}} name={{}}>
		{{range .Views }}
		  <view>
				.Render()
		  </view>
		{{ end }}
	</view>
`

// ViewLightTemplate simple renders out the internal views of a root View
const ViewLightTemplate = `
	{{range .Views }}
			.Render()
	{{ end }}
`

// SilentStratetgy is a simple strategy that when the view is activated calls the View.Execute
func SilentStratetgy() *ViewStrategy {
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

// SilentNameStratetgy is a simple strategy that when the view is activated calls the View.ExecuteTemplate
func SilentNameStratetgy(base string) *ViewStrategy {
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

// SourceView provides a view that takes the template format of which it will render the view as
func SourceView(tag, tmpl string) (*View, error) {
	tl, err := template.New(tag).Parse(tmpl)
	if err != nil {
		return nil, err
	}
	return NewView(tag, tl, SilentStratetgy()), nil
}
