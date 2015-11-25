package tests

import (
	"strings"
	"testing"

	"github.com/influx6/haiku/tests"
	"github.com/influx6/haiku/trees"
	"github.com/influx6/haiku/trees/attrs"
	"github.com/influx6/haiku/trees/elems"
)

func TestMarkup(t *testing.T) {
	div := elems.Div(
		elems.Text("20"),
	)

	if len(div.Children()) <= 0 {
		tests.FatalFailed(t, "Inaccurate size of children, expected %d for %s", 1, len(div.Children()))
	}

	firstRender := trees.SimpleElementWriter.Print(div)
	secondRender := trees.SimpleElementWriter.Print(div)

	if firstRender != secondRender {
		tests.FatalFailed(t, "Renders produced unequal results between \n %s and \n %s", firstRender, secondRender)
	}

	tests.LogPassed(t, "Successfully asserted proper markup operation!")
}

var normalRender = `<div hash="901EZEzzkP"  uid="3exhgHR9" style="">20</div>`
var removedRender = `<div hash="aPFZXtl2eW"  uid="2pml1sB0" style="">20</div>`
var cleanRender = `<div hash="aPFZXtl2eW"  uid="2pml1sB0" style=""></div>`

func TestMarkupRemoveRender(t *testing.T) {
	div := elems.Div(
		elems.Text("20"),
	)

	if len(div.Children()) <= 0 {
		tests.FatalFailed(t, "Inaccurate size of children, expected %d for %s", 1, len(div.Children()))
	}

	divCl := div.Clone().(*trees.Element)

	if len(divCl.Children()) <= 0 {
		tests.FatalFailed(t, "Inaccurate size of clone's children, expected %d for %s", 1, len(div.Children()))
	}

	printer := trees.NewElementWriter(trees.SimpleAttrWriter, trees.SimpleStyleWriter, trees.SimpleTextWriter)
	printer.AllowRemoved()

	if ds := printer.Print(div); len(ds) != len(normalRender) {
		tests.FatalFailed(t, "1 Renders produced unequal results between \n %s and \n %s", ds, normalRender)
	}

	trees.ElementsWithTag(divCl, "text")[0].Remove()

	if dl, dcl := len(div.Children()), len(divCl.Children()); dl != dcl {
		tests.FatalFailed(t, "Clone children size is inaccurate, expected %d but got %d", dl, dcl)
	}

	if ds := printer.Print(divCl); len(ds) != len(removedRender) {
		tests.FatalFailed(t, "2 Renders produced unequal results between \n %s and \n %s", ds, removedRender)
	}

	printer.DisallowRemoved()

	if ds := printer.Print(divCl); len(ds) != len(cleanRender) {
		tests.FatalFailed(t, "3 Renders produced unequal results between \n %s and \n %s", ds, cleanRender)
	}

	divCl.CleanRemoved()

	if dl, dcl := len(div.Children()), len(divCl.Children()); dl == dcl {
		tests.FatalFailed(t, "Clone children size is inaccurate, expected %d but got %d", dl, dcl)
	}

	tests.LogPassed(t, "Successfully asserted proper markup operation with .Remove()!")
}

func TestMarkupReconciliation(t *testing.T) {
	div := elems.Div(
		elems.Span(elems.Text("30")),
		elems.Text("20"),
	)

	divCl := elems.Div(
		elems.Span(elems.Text("30")),
		elems.Text("20"),
		elems.Text("400"),
	)

	//lets remove the span with its text child and the parents text child
	trees.ElementsWithTag(trees.ElementsWithTag(divCl, "span")[0], "text")[0].Remove()
	trees.ElementsWithTag(divCl, "text")[0].Remove()
	divCl.CleanRemoved()

	printer := trees.NewElementWriter(trees.SimpleAttrWriter, trees.SimpleStyleWriter, trees.SimpleTextWriter)
	printer.AllowRemoved()

	nrender := printer.Print(div)
	crender := printer.Print(divCl)

	if !strings.Contains(nrender, ">20") && !strings.Contains(nrender, ">30") {
		tests.FatalFailed(t, "Inaccurate rendering occured, has no '>20' or '>30' set", nrender)
	}

	if strings.Contains(crender, ">20") && strings.Contains(crender, ">30") {
		tests.FatalFailed(t, "Inaccurate rendering occured, has '>20' or '>30' and  set", crender)
	}

	//reconcile with the original div
	divCl.Reconcile(div)

	rcrender := printer.Print(divCl)

	if strings.Contains(rcrender, ">20") && strings.Contains(rcrender, ">30") && !strings.Contains(rcrender, ">400") {
		tests.FatalFailed(t, "Inaccurate rendering occured, has '>20' or '>30' and  set", rcrender)
	}

	tests.LogPassed(t, "Successfully reconciled dom markup!")
}

func TestMarkupReconciliation2(t *testing.T) {
	var menuItem = []string{"shops", "janitor", "booky", "drummer"}

	div := elems.Form(elems.Paragraph(elems.Label(elems.Text("name")), elems.Input(attrs.Type("text"))))

	var so = elems.Select()
	for _, mi := range menuItem {
		elems.Anchor(
			// events.Click(clickMe, "").PreventDefault(),
			attrs.Href("#"+mi),
			elems.Text(mi)).Apply(div)
		so.Augment(elems.Option(attrs.Name(mi), elems.Text(mi)))
	}

	div.Augment(elems.Paragraph(so))

	menuItem = append(menuItem, "border", "section", "chief")
	divCl := elems.Form(elems.Paragraph(elems.Label(elems.Text("name")), elems.Input(attrs.Type("text"))))

	var sol = elems.Select()
	for _, mi := range menuItem {
		elems.Anchor(
			// events.Click(clickMe, "").PreventDefault(),
			attrs.Href("#"+mi),
			elems.Text(mi)).Apply(divCl)
		sol.Augment(elems.Option(attrs.Name(mi), elems.Text(mi)))
	}

	divCl.Augment(elems.Paragraph(sol))

	printer := trees.NewElementWriter(trees.SimpleAttrWriter, trees.SimpleStyleWriter, trees.SimpleTextWriter)
	printer.AllowRemoved()

	nrender := printer.Print(div)
	crender := printer.Print(divCl)

	// fmt.Printf("%s\n\n", nrender)
	// fmt.Printf("%s\n\n", crender)

	//reconcile with the original div
	divCl.Reconcile(div)

	rcrender := printer.Print(divCl)

	if rcrender == crender {
		tests.FatalFailed(t, "1. Renders produced wrong results between \n %s and \n %s", rcrender, crender)
	}

	if rcrender == nrender {
		tests.FatalFailed(t, "2. Renders produced wrong results between \n %s and \n %s", rcrender, nrender)
	}

	// fmt.Printf("%s\n\n", rcrender)
	tests.LogPassed(t, "Successfully reconciled dom markup!")
}
