package trees

import (
	"fmt"
	"strings"
)

// This contains printers for the tree dom definition structures

// AttrPrinter defines a printer interface for writing out a Attribute objects into a string form
type AttrPrinter interface {
	Print([]*Attribute) string
}

// AttrWriter provides a concrete struct that meets the AttrPrinter interface
type AttrWriter struct{}

// SimpleAttrWriter provides a basic attribute writer
var SimpleAttrWriter = &AttrWriter{}

const attrformt = " %s='%s' "

// Print returns a stringed repesentation of the attribute object
func (m *AttrWriter) Print(a []*Attribute) string {
	attrs := []string{}

	for _, ar := range a {
		attrs = append(attrs, fmt.Sprintf(attrformt, ar.Name, ar.Value))
	}

	return strings.Join(attrs, " ")
}

// StylePrinter defines a printer interface for writing out a style objects into a string form
type StylePrinter interface {
	Print([]*Style) string
}

// StyleWriter provides a concrete struct that meets the AttrPrinter interface
type StyleWriter struct{}

// SimpleStyleWriter provides a basic style writer
var SimpleStyleWriter = &StyleWriter{}

const styleformt = " %s:%s; "

// Print returns a stringed repesentation of the style object
func (m *StyleWriter) Print(s []*Style) string {
	css := []string{}

	for _, cs := range s {
		css = append(css, fmt.Sprintf(styleformt, cs.Name, cs.Value))
	}

	return strings.Join(css, " ")
}

// TextWriter writes out the text element/node for the vdom into a string
type TextWriter struct{}

// SimpleTextWriter provides a basic text writer
var SimpleTextWriter = &TextWriter{}

// Write returns the string representation of the text object
func (m *TextWriter) Write(t *Text) string {
	return t.Get()
}

// ElementWriter writes out the element out as a string matching the html tag rules
type ElementWriter struct {
	attrWriter  AttrPrinter
	styleWriter StylePrinter
	text        *TextWriter
}

// SimpleElementWriter provides a default writer using the basic attribute and style writers
var SimpleElementWriter = NewElementWriter(SimpleAttrWriter, SimpleStyleWriter, SimpleTextWriter)

// NewElementWriter returns a new writer for Element objects
func NewElementWriter(aw AttrPrinter, sw StylePrinter, tw *TextWriter) *ElementWriter {
	return &ElementWriter{
		attrWriter:  aw,
		styleWriter: sw,
		text:        tw,
	}
}

// Write returns the string representation of the element
func (m *ElementWriter) Write(e *Element) string {

	//collect uid and hash of the element so we can write them along
	hash := &Attribute{"hash", e.Hash()}
	uid := &Attribute{"uid", e.UID()}

	//write out the hash and uid as attributes
	hashes := m.attrWriter.Print([]*Attribute{hash, uid})

	//write out the elements attributes using the AttrWriter
	attrs := m.attrWriter.Print(e.Attrs)

	//write out the elements inline-styles using the StyleWriter
	style := m.styleWriter.Print(e.Styles)

	var closer string

	if e.AutoClosed() {
		closer = "/>"
	} else {
		closer = fmt.Sprintf("</%s>", e.Tagname)
	}

	var children = []string{}

	for _, ch := range e.Children {
		if tch, ok := ch.(*Text); ok {
			children = append(children, m.text.Write(tch))
		}
		if ech, ok := ch.(*Element); ok {
			if ech == e {
				continue
			}
			children = append(children, m.Write(ech))
		}
	}

	//lets create the elements markup now
	return strings.Join([]string{
		fmt.Sprintf("<%s ", e.Tagname),
		hashes,
		attrs,
		style,
		">",
		strings.Join(children, "\n"),
		closer,
	}, "")
}

// MarkupPrinter defines a printer interface for writing out a markup object into a string form
type MarkupPrinter interface {
	Print(Markup) (string, error)
}

// MarkupWriter provides the concrete struct that meets the MarkupPrinter interface
type MarkupWriter struct {
	*ElementWriter
}

// SimpleMarkUpWriter provides a basic markup writer for handling the different markup elements
var SimpleMarkUpWriter = NewMarkupWriter(SimpleElementWriter)

// NewMarkupWriter returns a new markup instance
func NewMarkupWriter(em *ElementWriter) *MarkupWriter {
	return &MarkupWriter{em}
}

// Print returns a stringed repesentation of the markup object
func (m *MarkupWriter) Print(ma Markup) (string, error) {
	if tmr, ok := ma.(*Text); ok {
		return m.ElementWriter.text.Write(tmr), nil
	}

	if emr, ok := ma.(*Element); ok {
		return m.ElementWriter.Write(emr), nil
	}

	return "", ErrNotMarkup
}
