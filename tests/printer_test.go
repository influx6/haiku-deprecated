package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/influx6/flux"
	"github.com/influx6/haiku/trees"
	"github.com/influx6/haiku/trees/attrs"
	"github.com/influx6/haiku/trees/elems"
	"github.com/influx6/haiku/trees/styles"
)

var expectedAttr = " data-fax='butler' "

func TestAttrPrinter(t *testing.T) {
	attr := []*trees.Attribute{&trees.Attribute{"data-fax", "butler"}}
	res := trees.SimpleAttrWriter.Print(attr)
	flux.StrictExpect(t, expectedAttr, res)
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
		trees.SimpleAttrWriter.Print(attrs)
	}
}

var expectedStyle = " width:200px; "

func TestStylePrinter(t *testing.T) {
	su := []*trees.Style{styles.Width(styles.Px(200))}
	res := trees.SimpleStyleWriter.Print(su)
	flux.StrictExpect(t, expectedStyle, res)
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
		trees.SimpleStyleWriter.Print(su)
	}
}

func TestElementPrinter(t *testing.T) {
	elem := trees.NewElement("bench", false)
	attrs.ClassName("grid col1").Apply(elem)
	elems.Div(trees.NewText("thunder")).Apply(elem)

	res := trees.SimpleElementWriter.Print(elem)

	flux.Truthy(t, "Contains '<bench'", strings.Contains(res, "<bench"))
	flux.Truthy(t, "contains '</bench>'", strings.Contains(res, "</bench>"))
	flux.Truthy(t, "contains 'hash='", strings.Contains(res, "hash="))
	flux.Truthy(t, "contains 'uid='", strings.Contains(res, "uid="))
}

func BenchmarkElementPrinter(t *testing.B) {
	div := elems.Div()

	for i := 0; i < 10000; i++ {
		trees.NewText(fmt.Sprintf("%d", i)).Apply(div)
	}

	for i := 0; i < t.N; i++ {
		trees.SimpleElementWriter.Print(div)
	}
}
