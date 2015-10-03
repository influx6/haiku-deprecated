package views

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/influx6/assets"
	"github.com/influx6/haiku/trees"
)

// TemplateView provides a view system based on html.templates
// WARNING: TemplateView will wrap their outputs in a tmlview
// so never try to wrap a full <html></htmL> sequence with them
type TemplateView struct {
	Views
	Tmpl *template.Template
	hdom trees.SearchableMarkup
	txt  *trees.Text
}

// NewTemplateView returns a new view specifically created to use go html.template as a rendering system
func NewTemplateView(tag string, t *template.Template, strategy Strategy, binding interface{}, dobind bool) *TemplateView {
	hdom := trees.NewElement(tag, false)
	view := ReactView(NewView(tag, strategy, binding), dobind)

	tv := &TemplateView{
		Views: view,
		Tmpl:  t,
		hdom:  hdom,
		txt:   trees.NewText(""),
	}

	tv.Views.switchDOM(hdom)
	tv.txt.Apply(hdom)
	return tv
}

// RenderHTML renders the output from .Render() as safe html unescaped
func (v *TemplateView) RenderHTML(m ...string) template.HTML {
	var addr string

	if len(m) > 0 {
		addr = m[0]
	}

	if addr != "" {
		v.Engine().All(addr, v.Tag())
	}

	return template.HTML(v.Strategy().RenderSource(v))
}

// Render calls the internal strategy and renders out the output of that result
func (v *TemplateView) Render(m ...string) trees.Markup {
	var addr string

	if len(m) > 0 {
		addr = m[0]
	}

	if addr != "" {
		v.Engine().All(addr, v.Tag())
	}

	return v.Strategy().Render(v)
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

// SilentTemplateStrategy is a simple strategy that when the view is activated calls the View.Execute and returns an empty "" string when deactivated
func SilentTemplateStrategy(w trees.MarkupWriter) Strategy {
	return NewViewStrategy(w, func(v Views) trees.Markup {
		tv, ok := v.(*TemplateView)

		if !ok {
			return trees.NewText(fmt.Sprintf("CustomError(%s): %s", v.Tag(), "Expected type *TemplateView"))
		}

		bo, err := tv.Execute()

		if err != nil {
			tv.txt.Set(fmt.Sprintf("CustomError(%s): %s", v.Tag(), err.Error()))
			return tv.hdom
		}

		tv.txt.Set(string(bo))
		return tv.hdom
	}, func(v Views) trees.Markup {
		tv, ok := v.(*TemplateView)

		if !ok {
			return trees.NewText(fmt.Sprintf("CustomError(%s): %s", v.Tag(), "Expected type *TemplateView"))
		}

		tv.txt.Set("")
		return tv.hdom
	})
}

// SilentTemplateNameStrategy is a simple strategy that when the view is activated calls the View.ExecuteTemplate and returns an empty "" string when deactivated
func SilentTemplateNameStrategy(base string, w trees.MarkupWriter) Strategy {
	return NewViewStrategy(w, func(v Views) trees.Markup {
		tv, ok := v.(*TemplateView)

		if !ok {
			return trees.NewText(fmt.Sprintf("CustomError(%s): %s", v.Tag(), "Expected type *TemplateView"))
		}

		bo, err := tv.ExecuteTemplate(base)

		if err != nil {
			tv.txt.Set(fmt.Sprintf("CustomError(%s): %s", v.Tag(), err.Error()))
			return tv.hdom
		}

		tv.txt.Set(string(bo))
		return tv.hdom
	}, func(v Views) trees.Markup {
		tv, ok := v.(*TemplateView)

		if !ok {
			return trees.NewText(fmt.Sprintf("CustomError(%s): %s", v.Tag(), "Expected type *TemplateView"))
		}

		tv.txt.Set("")
		return tv.hdom
	})
}

// HiddenTemplateStrategy is a simple strategy that when the view is activated calls the View.Execute and if in deactive state returns the original content wrapped in a div with display:none
func HiddenTemplateStrategy(w trees.MarkupWriter) Strategy {
	return NewViewStrategy(w, func(v Views) trees.Markup {
		tv, ok := v.(*TemplateView)

		if !ok {
			return trees.NewText(fmt.Sprintf("CustomError(%s): %s", v.Tag(), "Expected type *TemplateView"))
		}

		dom := v.DOM()

		if ds, err := dom.GetStyle("display"); err == nil {
			if strings.Contains(ds.Value, "none") {
				ds.Value = "block"
			}
		}

		bo, err := tv.Execute()

		if err != nil {
			tv.txt.Set(fmt.Sprintf("CustomError(%s): %s", v.Tag(), err.Error()))
			return tv.hdom
		}

		tv.txt.Set(string(bo))
		return tv.hdom
	}, func(v Views) trees.Markup {
		dom := v.DOM()

		if ds, err := dom.GetStyle("display"); err == nil {
			if strings.Contains(ds.Value, "none") {
				ds.Value = "block"
			}
		}

		return dom
	})
}

// HiddenTemplateNameStrategy is a simple strategy that when the view is activated calls the View.ExecuteTemplate and when deactive wrap it within a div tag laced with a display none style
func HiddenTemplateNameStrategy(base string, w trees.MarkupWriter) Strategy {

	return NewViewStrategy(w, func(v Views) trees.Markup {
		tv, ok := v.(*TemplateView)

		if !ok {
			return trees.NewText(fmt.Sprintf("CustomError(%s): %s", v.Tag(), "Expected type *TemplateView"))
		}

		dom := v.DOM()

		if ds, err := dom.GetStyle("display"); err == nil {
			ds.Value = "block"
		}

		bo, err := tv.ExecuteTemplate(base)

		if err != nil {
			tv.txt.Set(fmt.Sprintf("CustomError(%s): %s", v.Tag(), err.Error()))
			return tv.hdom
		}

		tv.txt.Set(string(bo))
		return tv.hdom
	}, func(v Views) trees.Markup {
		dom := v.DOM()

		if ds, err := dom.GetStyle("display"); err == nil {
			ds.Value = "block"
		}

		return dom
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
func SimpleView(tag string, binding interface{}, dobind bool) (v *TemplateView, err error) {
	return SourceView(tag, ViewLightTemplate, binding, dobind)
}

// SimpleTreeView provides a view based on the ViewHTMLTemplate template
func SimpleTreeView(tag string, binding interface{}, dobind bool) (v *TemplateView, err error) {
	return SourceView(tag, ViewHTMLTemplate, binding, dobind)
}

// SourceView provides a view that takes the template format of which it will render the view as
func SourceView(tag, tmpl string, binding interface{}, dobind bool) (v *TemplateView, err error) {
	var tl *template.Template

	tl, err = template.New(tag).Parse(tmpl)

	if err != nil {
		return
	}

	v = NewTemplateView(tag, tl, SilentTemplateStrategy(trees.SimpleMarkupWriter), binding, dobind)

	return
}

// AssetView provides a view that takes the template format of which it will render the view as
func AssetView(tag, blockName string, binding interface{}, as *assets.AssetTemplate, dobind bool) (v *TemplateView, err error) {
	v = NewTemplateView(tag, as.Tmpl, SilentTemplateNameStrategy(blockName, trees.SimpleMarkupWriter), binding, dobind)
	return
}

// TemplateBlueprint defines the component blueprint that it generates using the TemplateView
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

// CreateComponent builds up a blueprint with the arguments, the name tag giving to the underline view is modded with the blueprint type name + a 5-length random string to make it unique in the state machines. All reactive binding are done if dobind is true hence boudning the binding change notification to the view.
func (b *TemplateBlueprint) CreateComponent(bind interface{}, vs Strategy, dobind bool) Components {
	view := NewTemplateView(MakeBlueprintName(b), b.format, vs, bind, dobind)
	return NewComponent(view, false)
}

// Create returns a new Component using the default HiddenTemplateStrategy
func (b *TemplateBlueprint) Create(bind interface{}, dobind bool) Components {
	return b.CreateComponent(bind, HiddenTemplateStrategy(trees.SimpleMarkupWriter), dobind)
}
