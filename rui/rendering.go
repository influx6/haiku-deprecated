package rui

import (
	"bytes"
	"fmt"
	"html/template"

	"github.com/influx6/flux"
	"github.com/influx6/prox/reactive"
)

// Renderer provides a interface that defines rendering methods custom renderers
type Renderer interface {
	Render() []byte
	Dirty() bool
}

// Templator returns a generator which it returns a new instance  of TemplateRender for rendering templates
type Templator struct {
	template *template.Template
}

// FileTemplator returns a Templator whoes template is received from a file
func FileTemplator(file ...string) (*Templator, error) {
	tml, err := template.ParseFiles(file...)

	if err != nil {
		return nil, err
	}

	return &Templator{
		template: tml,
	}, nil
}

// SourceTemplator returns a Templator whoes template is received from a file
func SourceTemplator(name, src string) (*Templator, error) {
	tml, err := template.New(name).Parse(src)

	if err != nil {
		return nil, err
	}

	return &Templator{
		template: tml,
	}, nil
}

// TemplateRender provides a reactive template rendering system for use with DataTrees
type TemplateRender struct {
	reactive.Observers
	template *template.Template
	Tree     DataTrees
	target   interface{}
}

func (t *Templator) buildTempler(to interface{}) (*TemplateRender, error) {
	ob, err := reactive.ObserveTransform("", false)

	if err != nil {
		return nil, err
	}

	dos, err := StructTree(to)

	if err != nil {
		return nil, err
	}

	cotl := (*t.template)
	tlr := &TemplateRender{
		Observers: ob,
		template:  &cotl,
		Tree:      dos,
		target:    to,
	}

	return tlr, nil
}

// BuildName returns a new TemplateRender set to use the current template and the given reactive object
func (t *Templator) BuildName(to interface{}, base string) (*TemplateRender, error) {
	tol, err := t.buildTempler(to)

	if err != nil {
		return nil, err
	}

	if err = tol.renderExecTmpl(base, to); err != nil {
		return nil, err
	}

	tol.Tree.React(func(r flux.Reactor, err error, dx interface{}) {
		if err != nil {
			tol.Observers.SendError(err)
			return
		}

		if err := tol.renderExecTmpl(base, dx); err != nil {
			tol.Observers.SendError(err)
		}
	}, true)

	return tol, nil
}

// Build returns a new TemplateRender set to use the current template and the given reactive object
func (t *Templator) Build(to interface{}) (*TemplateRender, error) {
	tol, err := t.buildTempler(to)

	if err != nil {
		return nil, err
	}

	if err = tol.renderExec(to); err != nil {
		return nil, err
	}

	tol.Tree.React(func(r flux.Reactor, err error, dx interface{}) {
		if err != nil {
			tol.Observers.SendError(err)
			return
		}

		if err = tol.renderExec(dx); err != nil {
			tol.Observers.SendError(err)
		}

	}, true)

	return tol, nil
}

// Set is made empty to ensure only internal observer can handle such an operation
func (t *TemplateRender) Set() {}

func (t *TemplateRender) renderExecTmpl(base string, dx interface{}) error {
	var buf bytes.Buffer

	err := t.template.ExecuteTemplate(&buf, base, dx)

	if err != nil {
		// t.Observers.SendError(err)
		return err
	}

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

	t.Observers.Set(string(buf.Bytes()))
	return nil
}

// Render renders the template or returns a cache if not yet fully rendered
func (t *TemplateRender) Render() string {
	return fmt.Sprintf("%v", t.Get())
}
