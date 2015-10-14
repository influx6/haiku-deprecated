package trees

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"

	"github.com/influx6/assets"
	"github.com/influx6/flux"
	"github.com/influx6/haiku/reactive"
)

// // Renderer provides a interface that defines rendering methods custom renderers
// type Renderer interface {
// 	Render(...string) string
// 	RenderHTML(...string) template.HTML
// 	IsDirty() bool
// }

// AssetRender provides a reactive template for rendering asset templates
type AssetRender struct {
	*TemplateRender
	assets *assets.AssetTemplate
}

// AssetTempler returns a new AssetRender instance using the NameTempler
func AssetTempler(base string, as *assets.AssetTemplate, to interface{}) (*AssetRender, error) {
	if as.Tmpl == nil {
		if err := as.Build(); err != nil {
			return nil, err
		}
	}

	tl, err := NameTempler(as.Tmpl, base, to)

	if err != nil {
		return nil, err
	}

	return &AssetRender{
		TemplateRender: tl,
		assets:         as,
	}, nil
}

// Reload sets up the call to .Build() to reload all filepaths from the given directory again
func (a *AssetRender) Reload() {
	a.assets.Reload()
}

// Build builds/rebuilds the files located within the file map of the AssetTemplate and updates the Templer template pointer
func (a *AssetRender) Build() error {
	if err := a.assets.Build(); err != nil {
		return err
	}

	a.template = a.assets.Tmpl
	return nil
}

// TemplateRender provides a reactive template rendering system for use with DataTrees
type TemplateRender struct {
	reactive.Observers
	template *template.Template
	Tree     DataTrees
	target   interface{}
	name     string
	dirty    bool
	logic    bool
}

// NameTempler returns a TemplateRender set to use the Template.ExecuteTemplate
func NameTempler(t *template.Template, name string, to interface{}) (*TemplateRender, error) {
	ob, err := NewTempler(t, to)
	if err != nil {
		return nil, err
	}
	return ob, UseNameExecLogic(ob, name)
}

// ExecTempler returns a TemplateRender set to use the Template.Execute
func ExecTempler(t *template.Template, to interface{}) (*TemplateRender, error) {
	ob, err := NewTempler(t, to)
	if err != nil {
		return nil, err
	}

	return ob, UseExecLogic(ob)
}

// NewTempler builds a TemplateRender basic skeleton without its update logic
func NewTempler(t *template.Template, to interface{}) (*TemplateRender, error) {
	ob := reactive.ObserveAtom("", false)

	dos, err := StructTree(to)

	if err != nil {
		return nil, err
	}

	tlr := &TemplateRender{
		Observers: ob,
		template:  t,
		Tree:      dos,
		target:    to,
	}

	return tlr, nil
}

// ErrLogicLocked is returned when trying to load a logic block into a TemplateRender
// when it is already loaded with one
var ErrLogicLocked = errors.New("TemplateRender already has logic loaded")

// UseNameExecLogic builds into the TemplateRender the logic of rendering using the template.Template.ExecuteTemplate by tagging with the name of value giving
func UseNameExecLogic(tol *TemplateRender, base string) error {
	if tol.logic {
		return ErrLogicLocked
	}

	tol.name = base
	tol.logic = true

	if err := tol.renderExecTmpl(tol.name, tol.target); err != nil {
		return err
	}

	tol.Tree.React(func(r flux.Reactor, err error, dx interface{}) {
		if err != nil {
			tol.Observers.SendError(err)
			return
		}

		if err := tol.renderExecTmpl(tol.name, dx); err != nil {
			tol.Observers.SendError(err)
		}
	}, true)

	return nil
}

// UseExecLogic builds into the templateRenderer the logic for rendering using the template.Template.Execute function call
func UseExecLogic(tol *TemplateRender) error {
	if tol.logic {
		return ErrLogicLocked
	}

	if err := tol.renderExec(tol.target); err != nil {
		return err
	}

	tol.logic = true
	tol.Tree.React(func(r flux.Reactor, err error, dx interface{}) {
		if err != nil {
			tol.Observers.SendError(err)
			return
		}

		if err = tol.renderExec(dx); err != nil {
			tol.Observers.SendError(err)
		}

	}, true)

	return nil
}

// Set is made empty to ensure only internal observer can handle such an operation
func (t *TemplateRender) Set() {}

// IsDirty returns true/false if the template has being updated
func (t *TemplateRender) IsDirty() bool {
	return !!t.dirty
}

//Get sets the dirty variable as fale and returns the current/last render of the template
func (t *TemplateRender) Get() interface{} {
	t.dirty = false
	return t.Observers.Get()
}

// RenderHTML renders the output from .Render() as safe html unescaped
func (t *TemplateRender) RenderHTML(m ...string) template.HTML {
	return template.HTML(t.Render())
}

// Render renders the template or returns a cache if not yet fully rendered
func (t *TemplateRender) Render(m ...string) string {
	return fmt.Sprintf("%v", t.Get())
}

func (t *TemplateRender) renderExecTmpl(base string, dx interface{}) error {
	var buf bytes.Buffer

	err := t.template.ExecuteTemplate(&buf, base, dx)

	if err != nil {
		// t.Observers.SendError(err)
		return err
	}

	t.dirty = true
	t.Observers.Set(string(buf.Bytes()))
	return nil
}

func (t *TemplateRender) renderExec(dx interface{}) error {
	var buf bytes.Buffer

	err := t.template.Execute(&buf, dx)

	if err != nil {
		// t.Observers.SendError(err)
		return err
	}

	t.dirty = true
	t.Observers.Set(string(buf.Bytes()))
	return nil
}
