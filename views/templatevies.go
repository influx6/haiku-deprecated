package views

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/influx6/assets"
	"github.com/influx6/flux"
)

// TemplateView provides a view system based on html.templates
type TemplateView struct {
	Views
	Tmpl *template.Template
}

// NewTemplateView returns a new view specifically created to use go html.template as a rendering system
func NewTemplateView(tag string, t *template.Template, strategy *ViewStrategy, binding interface{}) *TemplateView {
	return &TemplateView{
		Views: NewView(tag, strategy, binding),
		Tmpl:  t,
	}
}

// ExecuteTemplate calls the internal template.Template.ExecuteTemplate and returns the data from the rendering operation. The template is runned with the name but the view as the object/binding
func (v *TemplateView) ExecuteTemplate(name string) ([]byte, error) {
	var buf bytes.Buffer
	err := v.Tmpl.ExecuteTemplate(&buf, name, v.Views)
	return buf.Bytes(), err
}

// Execute calls the internal template.Template.Execute and returns the data from the rendering operation. The template is runned with the name but the view as the object/binding
func (v *TemplateView) Execute() ([]byte, error) {
	var buf bytes.Buffer
	err := v.Tmpl.Execute(&buf, v.Views)
	return buf.Bytes(), err
}

// SilentTemplateStrategy is a simple strategy that when the view is activated calls the View.Execute
func SilentTemplateStrategy() *ViewStrategy {
	return NewViewStrategy(func(v Views) string {
		tv, ok := v.(*TemplateView)

		if !ok {
			return fmt.Sprintf("CustomError(%s): %s", v.Tag(), "Expected type *TemplateView")
		}

		bo, err := tv.Execute()

		if err != nil {
			return fmt.Sprintf("CustomError(%s): %s", v.Tag(), err.Error())
		}

		return string(bo)
	}, func(v Views) string {
		return ""
	})
}

// SilentTemplateNameStrategy is a simple strategy that when the view is activated calls the View.ExecuteTemplate
func SilentTemplateNameStrategy(base string) *ViewStrategy {
	return NewViewStrategy(func(v Views) string {
		tv, ok := v.(*TemplateView)

		if !ok {
			return fmt.Sprintf("CustomError(%s): %s", v.Tag(), "Expected type *TemplateView")
		}

		bo, err := tv.ExecuteTemplate(base)

		if err != nil {
			return fmt.Sprintf("CustomError(%s): %s", v.Tag(), err.Error())
		}

		return string(bo)
	}, func(v Views) string {
		return ""
	})
}

// HiddenTemplateStrategy is a simple strategy that when the view is activated calls the View.Execute
func HiddenTemplateStrategy() *ViewStrategy {
	return NewViewStrategy(func(v Views) string {
		tv, ok := v.(*TemplateView)

		if !ok {
			return fmt.Sprintf("CustomError(%s): %s", v.Tag(), "Expected type *TemplateView")
		}

		bo, err := tv.Execute()

		if err != nil {
			return fmt.Sprintf("CustomError(%s): %s", v.Tag(), err.Error())
		}

		return string(bo)
	}, func(v Views) string {
		tv, ok := v.(*TemplateView)

		if !ok {
			return fmt.Sprintf("CustomError(%s): %s", v.Tag(), "Expected type *TemplateView")
		}

		bo, err := tv.Execute()

		if err != nil {
			return fmt.Sprintf("CustomError(%s): %s", v.Tag(), err.Error())
		}

		return fmt.Sprintf(`<div style="display:none;">\n%s\n</div>`, string(bo))
	})
}

// HiddenTemplateNameStrategy is a simple strategy that when the view is activated calls the View.ExecuteTemplate and when hidden wrap it within a div tag laced with a display none style
func HiddenTemplateNameStrategy(base string) *ViewStrategy {
	return NewViewStrategy(func(v Views) string {
		tv, ok := v.(*TemplateView)

		if !ok {
			return fmt.Sprintf("CustomError(%s): %s", v.Tag(), "Expected type *TemplateView")
		}

		bo, err := tv.ExecuteTemplate(base)

		if err != nil {
			return fmt.Sprintf("CustomError(%s): %s", v.Tag(), err.Error())
		}

		return string(bo)
	}, func(v Views) string {
		tv, ok := v.(*TemplateView)

		if !ok {
			return fmt.Sprintf("CustomError(%s): %s", v.Tag(), "Expected type *TemplateView")
		}

		bo, err := tv.ExecuteTemplate(base)

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
func SimpleView(tag string, binding interface{}) (v *TemplateView, err error) {
	return SourceView(tag, ViewLightTemplate, binding)
}

// SimpleTreeView provides a view based on the ViewHTMLTemplate template
func SimpleTreeView(tag string, binding interface{}) (v *TemplateView, err error) {
	return SourceView(tag, ViewHTMLTemplate, binding)
}

// SourceView provides a view that takes the template format of which it will render the view as
func SourceView(tag, tmpl string, binding interface{}) (v *TemplateView, err error) {
	var tl *template.Template

	tl, err = template.New(tag).Parse(tmpl)

	if err != nil {
		return
	}

	v = NewTemplateView(tag, tl, SilentTemplateStrategy(), binding)

	return
}

// AssetView provides a view that takes the template format of which it will render the view as
func AssetView(tag, blockName string, binding interface{}, as *assets.AssetTemplate) (v *TemplateView, err error) {
	v = NewTemplateView(tag, as.Tmpl, SilentTemplateNameStrategy(blockName), binding)
	return
}

// NewReactiveTemplateView provides a decorator function to return a new ReactiveView with the same arguments passed to NewView(...)
func NewReactiveTemplateView(tag string, tl *template.Template, strategy *ViewStrategy, binding interface{}) ReactiveViews {
	return BuildReactiveTemplateView(tag, tl, strategy, binding, true)
}

// BuildReactiveTemplateView provides a decorator function to return a new ReactiveView with the same arguments passed to NewView(...), useRB -> means UseReactiveBinding
func BuildReactiveTemplateView(tag string, tl *template.Template, strategy *ViewStrategy, binding interface{}, useRB bool) ReactiveViews {
	rv := ReactView(NewTemplateView(tag, tl, strategy, binding))
	if useRB {
		BindReactor(rv, binding)
	}
	return rv
}

// ReactiveSourceView provides a view that takes the template format of which it will render the view as
func ReactiveSourceView(tag, tmpl string, binding interface{}, userb bool) (ReactiveViews, error) {
	sv, err := SourceView(tag, tmpl, binding)
	if err != nil {
		return nil, err
	}

	rv := ReactView(sv)

	if userb {
		BindReactor(rv, binding)
	}

	return rv, nil
}

// TemplateBlueprint defines the component blueprint that it generates, like setting the
// building blocks that make up a components behaviour
type TemplateBlueprint struct {
	format  *template.Template
	bluetag string
}

// NewTemplateBlueprint returns a new blueprint instance
func NewTemplateBlueprint(id string, t *template.Template) *TemplateBlueprint {
	bp := TemplateBlueprint{
		format:  t,
		bluetag: id,
	}

	return &bp
}

// Type returns the tagname type of the components generated by this blueprint
func (b *TemplateBlueprint) Type() string {
	return b.bluetag
}

// TemplateView builds up a blueprint with the arguments, the name tag giving to the
// underline view is modded with the blueprint type name + a 5-length random string
// to make it unique in the state machines. All reactive binding are automatically bounded to the view.
func (b *TemplateBlueprint) TemplateView(bind interface{}, vs *ViewStrategy, dobind bool) Components {
	view := BuildReactiveTemplateView(fmt.Sprintf("%s:%s", b.Type(), flux.RandString(5)), b.format, vs, bind, dobind)
	return NewComponent(view)
}

// MixTemplateView creates a new component with a combined template if supplied i.e the parsetree of the
// Blueprint.template adds the parse tree of the supplied template if pressent and if possible else uses the default blueprints template. All reactive binding are automatically bounded to the view.
func (b *TemplateBlueprint) MixTemplateView(bind interface{}, vs *ViewStrategy, subt *template.Template, dobind bool) (Components, error) {
	var sub *template.Template

	if subt != nil {
		so, err := b.format.AddParseTree(subt.Name(), subt.Tree)

		if err != nil {
			return nil, err
		}

		sub = so
	} else {
		sub = b.format
	}

	view := BuildReactiveTemplateView(fmt.Sprintf("%s:%s", b.Type(), flux.RandString(5)), sub, vs, bind, dobind)
	return NewComponent(view), nil
}
