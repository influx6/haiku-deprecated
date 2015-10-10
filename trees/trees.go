package trees

import (
	"fmt"
	"strings"

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
	Clone() Markup
	Name() string
	AddChild(Markup) bool
	Augment(...Markup) bool

	GetStyles(f, val string) []*Style
	GetStyle(f string) (*Style, error)
	StyleContains(f, val string) bool
	GetAttrs(f, val string) []*Attribute
	GetAttr(f string) (*Attribute, error)
	AttrContains(f, val string) bool
	ElementsUsingStyle(f, val string) []*Element
	ElementsWithAttr(f, val string) []*Element
	DeepElementsUsingStyle(f, val string, depth int) []*Element
	DeepElementsWithAttr(f, val string, depth int) []*Element
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
func NewElement(tag string, hasNoEndingTag bool) *Element {
	return &Element{
		Mutation:  NewMutable(),
		Tagname:   tag,
		Children:  make([]Markup, 0),
		Styles:    make([]*Style, 0),
		Attrs:     make([]*Attribute, 0),
		autoclose: hasNoEndingTag,
	}
}

// Clone makes a new copy of the markup structure
func (e *Element) Clone() Markup {
	co := NewElement(e.Tagname, e.autoclose)

	//clone the internal styles
	for _, so := range e.Styles {
		so.Clone().Apply(co)
	}

	//clone the internal attribute
	for _, ao := range e.Attrs {
		ao.Clone().Apply(co)
	}

	//clone the internal children
	for _, ch := range e.Children {
		ch.Clone().Apply(co)
	}

	return co
}

// GetStyles returns the styles that contain the specified name and if not empty that contains the specified value also, note that strings
// NOTE: string.Contains is used when checking value parameter if present
func (e *Element) GetStyles(f, val string) []*Style {
	var found []*Style

	for _, as := range e.Styles {
		if as.Name != f {
			continue
		}

		if val != "" {
			if !strings.Contains(as.Value, val) {
				continue
			}
		}

		found = append(found, as)
	}

	return found
}

// GetStyle returns the style with the specified tag name
func (e *Element) GetStyle(f string) (*Style, error) {
	for _, as := range e.Styles {
		if as.Name == f {
			return as, nil
		}
	}
	return nil, ErrNotFound
}

// StyleContains returns the styles that contain the specified name and if the val is not empty then
// that contains the specified value also, note that strings
// NOTE: string.Contains is used
func (e *Element) StyleContains(f, val string) bool {
	for _, as := range e.Styles {
		if !strings.Contains(as.Name, f) {
			continue
		}

		if val != "" {
			if !strings.Contains(as.Value, val) {
				continue
			}
		}

		return true
	}

	return false
}

// GetAttrs returns the attributes that have the specified text within the naming
// convention and if it also contains the set val if not an empty "",
// NOTE: string.Contains is used
func (e *Element) GetAttrs(f, val string) []*Attribute {
	var found []*Attribute

	for _, as := range e.Attrs {
		if as.Name != f {
			continue
		}

		if val != "" {
			if !strings.Contains(as.Value, val) {
				continue
			}
		}

		found = append(found, as)
	}

	return found
}

// AttrContains returns the attributes that have the specified text within the naming
// convention and if it also contains the set val if not an empty "",
// NOTE: string.Contains is used
func (e *Element) AttrContains(f, val string) bool {
	for _, as := range e.Attrs {
		if !strings.Contains(as.Name, f) {
			continue
		}

		if val != "" {
			if !strings.Contains(as.Value, val) {
				continue
			}
		}

		return true
	}

	return false
}

// GetAttr returns the attribute with the specified tag name
func (e *Element) GetAttr(f string) (*Attribute, error) {
	for _, as := range e.Attrs {
		if as.Name == f {
			return as, nil
		}
	}
	return nil, ErrNotFound
}

// ElementsUsingStyle returns the children within the element matching the
// stlye restrictions passed.
// NOTE: is uses Element.StyleContains
func (e *Element) ElementsUsingStyle(f, val string) []*Element {
	return e.DeepElementsUsingStyle(f, val, 1)
}

// ElementsWithAttr returns the children within the element matching the
// stlye restrictions passed.
// NOTE: is uses Element.AttrContains
func (e *Element) ElementsWithAttr(f, val string) []*Element {
	return e.DeepElementsWithAttr(f, val, 1)
}

// DeepElementsUsingStyle returns the children within the element matching the
// style restrictions passed allowing control of search depth
// NOTE: is uses Element.StyleContains
func (e *Element) DeepElementsUsingStyle(f, val string, depth int) []*Element {
	if depth <= 0 {
		return nil
	}

	var found []*Element

	for _, ch := range e.Children {
		if che, ok := ch.(*Element); ok {
			if che.StyleContains(f, val) {
				found = append(found, che)
				cfo := che.DeepElementsUsingStyle(f, val, depth-1)
				if len(cfo) > 0 {
					found = append(found, cfo...)
				}
			}
		}
	}

	return found
}

// DeepElementsWithAttr returns the children within the element matching the
// attributes restrictions passed allowing control of search depth
// NOTE: is uses Element.AttrContains
func (e *Element) DeepElementsWithAttr(f, val string, depth int) []*Element {
	if depth <= 0 {
		return nil
	}

	var found []*Element

	for _, ch := range e.Children {
		if che, ok := ch.(*Element); ok {
			if che.AttrContains(f, val) {
				found = append(found, che)
				cfo := che.DeepElementsWithAttr(f, val, depth-1)
				if len(cfo) > 0 {
					found = append(found, cfo...)
				}
			}
		}
	}

	return found

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
func (e *Element) AddChild(em Markup) bool {
	e.Children = append(e.Children, em)
	return true
}

// Augment provides a generic method for markup addition
func (e *Element) Augment(m ...Markup) bool {
	for _, mo := range m {
		mo.Apply(e)
	}
	return true
}

//Apply adds the giving element into the current elements children tree
func (e *Element) Apply(em *Element) {
	em.AddChild(e)
	e.Bind(em, false)
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

	m.React(func(r flux.Reactor, _ error, _ interface{}) {
		t.Set(fmt.Sprintf("%+v", m.Get()))
	}, true)

	return t
}

// Augment is implemented as a no-op
func (t *Text) Augment(m ...Markup) bool {
	return false
}

// Clone makes a new copy of the markup structure
func (t *Text) Clone() Markup {
	return NewText(t.Get())
}

// AddChild adds a new markup as the children of this element
func (t *Text) AddChild(em Markup) bool {
	return false
}

// GetStyles implement this method as a noop
func (t *Text) GetStyles(f, val string) []*Style {
	return nil
}

// GetStyle implement this method as a noop
func (t *Text) GetStyle(f string) (*Style, error) {
	return nil, ErrNotFound
}

// StyleContains implement this method as a noop
func (t *Text) StyleContains(f, val string) bool {
	return false
}

// GetAttrs implement this method as a noop
func (t *Text) GetAttrs(f, val string) []*Attribute {
	return nil
}

// GetAttr implement this method as a noop
func (t *Text) GetAttr(f string) (*Attribute, error) {
	return nil, ErrNotFound
}

// AttrContains implement this method as a noop
func (t *Text) AttrContains(f, val string) bool {
	return false
}

// ElementsUsingStyle implement this method as a noop
func (t *Text) ElementsUsingStyle(f, val string) []*Element {
	return nil
}

// ElementsWithAttr implement this method as a noop
func (t *Text) ElementsWithAttr(f, val string) []*Element {
	return nil
}

// DeepElementsUsingStyle implement this method as a noop
func (t *Text) DeepElementsUsingStyle(f, val string, depth int) []*Element {
	return nil
}

// DeepElementsWithAttr implement this method as a noop
func (t *Text) DeepElementsWithAttr(f, val string, depth int) []*Element {
	return nil
}

// Set sets the value of the text
func (t *Text) Set(tx string) {
	if t.text == tx {
		return
	}

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
	t.Bind(e, false)
}

// Attribute define the struct  for attributes
type Attribute struct {
	Name  string
	Value string
}

//Clone replicates the attribute into a unique instance
func (a *Attribute) Clone() *Attribute {
	return &Attribute{Name: a.Name, Value: a.Value}
}

// Apply applies a set change to the giving element attributes list
func (a *Attribute) Apply(e *Element) {
	e.Attrs = append(e.Attrs, a)
}

// Style define the style specification for element styles
type Style struct {
	Name  string
	Value string
}

//Clone replicates the style into a unique instance
func (s *Style) Clone() *Style {
	return &Style{Name: s.Name, Value: s.Value}
}

// Apply applies a set change to the giving element style list
func (s *Style) Apply(e *Element) {
	e.Styles = append(e.Styles, s)
}

// Augment adds new markup to an the root if its Element
func Augment(root Markup, m ...Markup) {
	if el, ok := root.(*Element); ok {
		for _, mo := range m {
			mo.Apply(el)
		}
	}
}
