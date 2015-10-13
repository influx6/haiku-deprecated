package views

import (
	"bytes"
	"html/template"
)

// //NullRender provides a null rendering that returns a string of error stating which render tag failed
// type NullRender struct {
// 	tag string
// }
//
// // NewNullRender returns a new NullRender instance
// func NewNullRender(tag string) *NullRender {
// 	return &NullRender{tag: tag}
// }
//
// // RenderHTML renders the output from .Render() as safe html unescaped
// func (n *NullRender) RenderHTML(m ...string) template.HTML {
// 	return template.HTML(n.Render(m...).(*trees.Text).Get())
// }
//
// // Render returns the error message
// func (n *NullRender) Render(_ ...string) trees.Markup {
// 	return trees.NewText(fmt.Sprintf(`Render.Error: "%s" view not found!`, n.tag))
// }
//

// Renderables are objects capable of rendering out themselves as strings using any variety of methods but most importantly do cache that last render for returning when desired to do so

// TemplateRenderable defines a basic example of a Renderable
type TemplateRenderable struct {
	tmpl  *template.Template
	cache *bytes.Buffer
}

// NewTemplateRenderable returns a new Renderable
func NewTemplateRenderable(content string) (*TemplateRenderable, error) {
	tl, err := template.New("").Parse(content)
	if err != nil {
		return nil, err
	}

	tr := TemplateRenderable{
		tmpl:  tl,
		cache: bytes.NewBuffer([]byte{}),
	}

	return &tr, nil
}

// Execute effects the inner template with the supplied data
func (t *TemplateRenderable) Execute(v interface{}) error {
	t.cache.Reset()
	err := t.tmpl.Execute(t.cache, v)
	return err
}

// Render renders out the internal cache
func (t *TemplateRenderable) Render(_ ...string) string {
	return string(t.cache.Bytes())
}

// RenderHTML renders the output from .Render() as safe html unescaped
func (t *TemplateRenderable) RenderHTML(_ ...string) template.HTML {
	return template.HTML(t.Render())
}

// String calls the render function
func (t *TemplateRenderable) String() string {
	return t.Render()
}
