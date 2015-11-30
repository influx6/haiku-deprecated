package tests

import (
	"fmt"
	"strings"
	"testing"

	"github.com/influx6/haiku/tests"
	"github.com/influx6/haiku/trees"
	"github.com/influx6/haiku/trees/attrs"
	"github.com/influx6/haiku/trees/elems"
	"github.com/influx6/haiku/trees/styles"
)

var expectedAttr = " data-fax='butler' "

func TestAttrPrinter(t *testing.T) {
	attr := []*trees.Attribute{&trees.Attribute{Name: "data-fax", Value: "butler"}}
	res := trees.SimpleAttrWriter.Print(attr)
	tests.StrictExpect(t, expectedAttr, res)
}

func BenchmarkAttrPrinter(t *testing.B) {
	attrs := []*trees.Attribute{}

	for i := 0; i < 10000; i++ {
		attrs = append(attrs, &trees.Attribute{
			Name:  fmt.Sprintf("%d", i),
			Value: fmt.Sprintf("%d", i*2),
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
	tests.StrictExpect(t, expectedStyle, res)
}

func BenchmarkStylePrinter(t *testing.B) {
	su := []*trees.Style{}

	for i := 0; i < 10000; i++ {
		su = append(su, &trees.Style{
			Name:  fmt.Sprintf("%d", i),
			Value: fmt.Sprintf("%d", i*2),
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

	classes := &trees.ClassList{}
	classes.Add("x-icon")
	classes.Add("x-lock")

	classes.Apply(elem)

	res := trees.SimpleElementWriter.Print(elem)

	tests.Truthy(t, "Contains '<bench'", strings.Contains(res, "<bench"))
	tests.Truthy(t, "contains '</bench>'", strings.Contains(res, "</bench>"))
	tests.Truthy(t, "contains 'hash='", strings.Contains(res, "hash="))
	tests.Truthy(t, "contains 'uid='", strings.Contains(res, "uid="))
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
