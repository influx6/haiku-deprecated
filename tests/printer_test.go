package tests

import (
	"log"
	"strings"
	"testing"

	"github.com/influx6/flux"
	"github.com/influx6/haiku/trees"
	"github.com/influx6/haiku/trees/elems"
	"github.com/influx6/haiku/trees/styles"
)

var expectedAttr = " data-fax='butler' "

func TestAttrPrinter(t *testing.T) {
	attr := []*trees.Attribute{&Attribute{"data-fax", "butler"}}
	res := trees.SimpleAttrWriter.Write(attr)
	log.Printf("attr -> %s", res)
	flux.Expect(t, expectedAttr, res)
}

func BenchmarkAttrPrinter(t *testing.B) {
	attrs := []*trees.Attribute{}

	for i := 0; i < 10000; i++ {
		attrs = append(attrs, &trees.Attribute{
			fmt.Sprintf("%d", i),
			fmt.Sprintf("%d", i*2),
		})
	}

	for i := 0; i < t.N; i++ {
		trees.SimpleAttrWriter.Write(attrs)
	}
}

var expectedStyle = " width:200px; "

func TestStylePrinter(t *testing.T) {
	su := []*trees.Style{styles.Width(200)}
	res := trees.SimpleStyleWriter.Write(su)
	log.Printf("style -> %s", res)
	flux.Expect(t, expectedStyle, res)
}

func BenchmarkStylePrinter(t *testing.B) {
	su := []*trees.Style{}

	for i := 0; i < 10000; i++ {
		su = append(su, &trees.Style{
			fmt.Sprintf("%d", i),
			fmt.Sprintf("%d", i*2),
		})
	}

	for i := 0; i < t.N; i++ {
		trees.SimpleStyleWriter.Write(su)
	}
}

func TestElementPrinter(t *testing.T) {
	elem := trees.NewElement("bench", true)
	elems.Div(trees.NewText("thunder")).Apply(elem)

	res := trees.SimpleElementWriter.Write(elem)
	log.Printf("elem -> %s", res)

	flux.Truthy(t, "WriteElement", strings.Contains(res, "<bench"))
}

func BenchmarkElementPrinter(t *testing.B) {
	div := elems.Div()

	for i := 0; i < 10000; i++ {
		trees.NewText(fmt.Sprinf("%d", i)).Apply(div)
	}

	for i := 0; i < t.N; i++ {
		trees.SimpleElementWriter.Write(div)
	}
}
