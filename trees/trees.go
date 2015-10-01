package trees

import (
	"fmt"

	"github.com/influx6/flux"
	"github.com/influx6/haiku/reactive"
)

// Mutation defines the capability of an element to state its
// state if mutation occured
type Mutation interface {
	flux.Reactor
	UID() string
	Hash() string
}

// Appliable define the interface specification for applying changes to elements elements in tree
type Appliable interface {
	Apply(*Element)
}

// Markup provide a basic specification type of how a element resolves its content
type Markup interface {
	Appliable
	Mutation
	Name() string
}

// Mutable is a base implementation of the Mutation interface{}
type Mutable struct {
	flux.Reactor
	uid  string
	hash string
}

// NewMutable returns a new mutable instance
func NewMutable() *Mutable {
	m := &Mutable{
		uid:  flux.RandString(8),
		hash: flux.RandString(10),
	}

	m.Reactor = flux.Reactive(func(r flux.Reactor, err error, d interface{}) {
		m.hash = flux.RandString(10)
		if err != nil {
			r.ReplyError(err)
		} else {
			r.Reply(d)
		}
	})

	return m
}

// Hash returns the current hash of the mutable
func (m *Mutable) Hash() string {
	return m.hash
}

// UID returns the current uid of the mutable
func (m *Mutable) UID() string {
	return m.uid
}

// Element represent a concrete implementation of a element node
type Element struct {
	Mutation
	Tagname   string
	Styles    []*Style
	Attrs     []*Attribute
	Children  []Markup
	autoclose bool
}

// NewElement returns a new element instance giving the specificed name
func NewElement(tag string, autoclose bool) *Element {
	return &Element{
		Mutation:  NewMutable(),
		Tagname:   tag,
		Children:  make([]Markup, 0),
		Styles:    make([]*Style, 0),
		Attrs:     make([]*Attribute, 0),
		autoclose: autoclose,
	}
}

// AutoClosed returns true/false if this element uses a </> or a <></> tag convention
func (e *Element) AutoClosed() bool {
	return e.autoclose
}

// Name returns the tag name of the element
func (e *Element) Name() string {
	return e.Tagname
}

// AddChild adds a new markup as the children of this element
func (e *Element) AddChild(em Markup) {
	e.Children = append(e.Children, em)
}

//Apply adds the giving element into the current elements children tree
func (e *Element) Apply(em *Element) {
	em.AddChild(e)
}

// Text represent a text element
type Text struct {
	Mutation
	text string
}

// NewText returns a new Text instance element
func NewText(txt string) *Text {
	mo := NewMutable()

	t := &Text{
		Mutation: mo,
		text:     txt,
	}

	return t
}

// MText returns a text tied to an observable
func MText(m reactive.Observers) *Text {
	t := NewText(fmt.Sprintf("%+v", m.Get()))

	m.React(func(r flux.Reactor, err error, d interface{}) {
		t.Set(fmt.Sprintf("%+v", m.Get()))
	}, true)

	return t
}

// Set sets the value of the text
func (t *Text) Set(tx string) {
	t.text = tx
	t.Send(t)
}

// Get returns the value of the text
func (t *Text) Get() string {
	return t.text
}

// Name returns the tag name of the element
func (t *Text) Name() string {
	return "text"
}

// Apply applies a set change to the giving element children list
func (t *Text) Apply(e *Element) {
	e.Children = append(e.Children, t)
}

// Attribute define the struct  for attributes
type Attribute struct {
	Name  string
	Value string
	Appliable
}

// Apply applies a set change to the giving element attributes list
func (a *Attribute) Apply(e *Element) {
	e.Attrs = append(e.Attrs, a)
}

// Style define the style specification for element styles
type Style struct {
	Name  string
	Value string
	Appliable
}

// Apply applies a set change to the giving element style list
func (s *Style) Apply(e *Element) {
	e.Styles = append(e.Styles, s)
}
