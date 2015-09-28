package main

import (
	"fmt"
	"html/template"
	"strconv"

	hdom "github.com/influx6/haiku/dom"
	"github.com/influx6/haiku/reactive"
	"github.com/influx6/haiku/views"
	dom "honnef.co/go/js/dom"
)

type contact struct {
	Name     reactive.Observers
	Location reactive.Observers
	Tel      reactive.Observers
	Age      reactive.Observers
}

//go:generate gopherjs build main.go -o ./app.js

func main() {

	win, ok := dom.GetWindow().(dom.Window)

	if !ok {

		panic("not in a browser-window")
	}

	doc, ok := win.Document().(dom.HTMLDocument)

	if !ok {
		panic("not in a browser-dom")
	}

	container := doc.QuerySelector(".container")

	formtemplate := template.Must(template.New("form").Parse(doc.QuerySelector(".contact-form-template").InnerHTML()))
	viewtemplate := template.Must(template.New("view").Parse(doc.QuerySelector(".contact-view-template").InnerHTML()))

	dude := &contact{
		Name:     reactive.ObserveAtom("John", false),
		Location: reactive.ObserveAtom("Lagos,Nigeria", false),
		Tel:      reactive.ObserveAtom("+01 4234534", false),
		Age:      reactive.ObserveAtom(20, false),
	}

	contactTempler, err := reactive.ExecTempler(viewtemplate, dude)

	if err != nil {
		panic(fmt.Sprintf("Contact template error: %s", err))
	}

	contactView, err := views.ReactiveSourceView("contacts.view", `
			{{ (.View "view").RenderHTML }}
	`)

	contactView.AddView("view", ".", contactTempler)

	formTempler, err := reactive.ExecTempler(formtemplate, dude)

	if err != nil {
		panic(fmt.Sprintf("Contact template error: %s", err))
	}

	// _ = container
	formView, err := views.SourceView("contacts.view", `
			{{ (.View "form").RenderHTML }}
	`)

	if err != nil {
		panic(fmt.Sprintf("Contact view error: %s", err))
	}

	formView.AddView("form", ".", formTempler)

	contacts := hdom.NewElement(container.QuerySelector(".contactView"))

	contacts.AddEvent("keyup", ".contact-form p .input").Next(func(ev dom.Event, next hdom.NextHandler) {
		target := ev.Target()
		name := target.GetAttribute("name")

		var input *dom.HTMLInputElement
		var ok bool

		if input, ok = target.(*dom.HTMLInputElement); !ok {
			return
		}

		switch name {
		case "name":
			dude.Name.Set(input.Value)
		case "age":
			if age, err := strconv.Atoi(input.Value); err == nil {
				dude.Age.Set(age)
			}
		case "location":
			dude.Location.Set(input.Value)
		case "telephone":
			dude.Tel.Set(input.Value)
		}

		next(ev)
	})

	form := hdom.NewViewElement(contacts, formView, ".form")
	view := hdom.NewViewElement(contacts, contactView, ".view")

	form.Sync(".")
	view.Sync(".")
}
